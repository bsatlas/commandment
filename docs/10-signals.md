# Signals

Commandment can make external or application-defined control events part of command execution policy instead of leaving every command to build its own goroutine and `os.Signal` loop.

The important model is broader than operating-system signals.

```text
signal source
    -> occurrence
    -> match and phase
    -> selection or aggregation
    -> policy
    -> action sequence
    -> execution effect
```

## Start with the familiar case of Ctrl-C

A command can route an interrupt into a typed endpoint that its handler receives through configuration.

```go
type listenConfig struct {
	Events <-chan cli.SignalEvent
}

cmd := cli.Command[listenConfig]{
	Name: "listen",
	Declare: func(cb *cli.CommandBuilder[listenConfig]) {
		events := cb.Signals.Endpoint(&cb.Config.Events)
		interrupt := cb.Signals.Match.
			Interrupt().
			During(cli.SignalPhaseExecution).
			Each().
			Hint("interrupt")

		cb.Signals.Policy(interrupt, cb.Signals.Action.Emit(events))
	},
}
```

The handler can wait on the endpoint just like any other typed Go channel.

```go
select {
case <-cfg.Events:
	fmt.Println("signal received")
	return nil
case <-ctx.Done():
	return context.Cause(ctx)
}
```

The [signals/interrupt](../examples/signals/interrupt/) scenario proves the complete process boundary by starting the executable, waiting for readiness, delivering a real interrupt, and observing successful shutdown.

## A signal source owns how occurrences enter Commandment

Operating-system delivery is one source family. Applications can also register typed sources for events that have signal semantics in their own host environment.

```go
descriptor := cli.SignalDescriptor[string]{
	Signal:   "alarm",
	Fidelity: cli.SignalSourceDeliveryFidelityExact,
}

source := cli.SignalSource[string]{
	Signals: []cli.SignalDescriptor[string]{descriptor},
	Open:    openAlarm,
}

registered := cb.Signals.Source(source)
alarm := registered.Signal("alarm")
```

A custom source may represent a supervisory control channel, runtime notification, host event, reload request, or another application-defined control plane.

See [signals/custom-source](../examples/signals/custom-source/).

## Match the semantic signal and the phase where it matters

A policy should say which signal it responds to and when that response is meaningful.

```go
match := cb.Signals.Match.
	On(alarm).
	During(cli.SignalPhaseExecution).
	Each().
	Hint("alarm")
```

Phase selection prevents a control event intended for active command execution from accidentally acquiring the same meaning during unrelated preparation or teardown.

Matching is semantic policy. It should not be replaced with a handler that reads a generic channel and then reimplements source identity and phase rules in application code.

## Selection and aggregation decide which occurrences drive policy

A policy can operate on occurrences individually or according to the supported selector and aggregation semantics declared through the signal API.

Use the narrowest behavior that represents the application requirement. A single graceful-shutdown action and a "first interrupt begins shutdown, second interrupt escalates" policy are different control semantics and should be declared differently.

Exact selector and aggregation constructors live in GoDoc because they are reference-level details. The design question is whether your application needs individual occurrences, a selected occurrence, or an aggregate control decision.

## Endpoints deliver structured events to application behavior

An endpoint is useful when the handler or another execution component needs to observe the matched signal as application data.

```go
events := cb.Signals.Endpoint(&cb.Config.Events)
cb.Signals.Policy(match, cb.Signals.Action.Emit(events))
```

The resulting `SignalEvent` retains signal-model identity instead of collapsing the occurrence into an untyped boolean.

An endpoint does not require every policy to use event delivery. Some policies should act directly.

## Choose the action that owns the control effect

Signal actions can express several kinds of policy under the public API, including ignoring a match, cancelling primary execution, invoking application handling, forwarding or sending a signal to a declared target, emitting a structured event, aborting with an exit code, or reraising to the host.

Choose based on ownership.

- Emit when application code should receive structured signal data.
- Cancel when the signal means the operation's primary execution should be cancelled.
- Handle when a declared application callback owns the immediate reaction.
- Forward or send when another declared target should receive the signal.
- Abort when the policy owns an explicit termination result.
- Reraise when host-level signal behavior should resume after Commandment's declared handling.

Do not emulate one action with another merely because they can produce superficially similar process output.

## Sequence actions when the policy genuinely has stages

Actions can be sequenced.

```go
plan := cb.Signals.Action.
	Emit(events).
	ThenCancel()
```

A policy can also express delayed escalation where supported.

```go
plan := cb.Signals.Action.
	Handle(beginGracefulShutdown).
	After(5 * time.Second).
	Abort(1)
```

This is useful for policies such as graceful shutdown followed by forced termination. The sequence belongs to signal semantics rather than to an unrelated timer goroutine hidden in `main`.

## Source fidelity and sender information matter when the host can provide them

Not every source observes or preserves the same information. A source descriptor can express delivery fidelity, and occurrences may carry sender or source metadata according to the source contract.

Do not invent precision the source does not have. If the host only knows that a signal-like event occurred, it should not claim exact sender identity. If a source promises exact delivery identity, downstream policy and tooling can rely on that stronger fact.

Exact fidelity and sender fields are GoDoc-level contracts, but the application design should preserve the source's actual guarantees.

## Source failure is different from a matched signal

Opening, reading, or maintaining a source can fail independently of receiving a valid signal occurrence. Treat source failure according to the signal-source contract rather than translating it into an arbitrary synthetic signal.

Likewise, overflow or delivery loss is not automatically the same thing as application cancellation. Preserve the distinction exposed by the public source and occurrence model.

## Signals and context cancellation are related without being identical

A signal policy may cancel execution, but a caller can also cancel the operation context directly. Application logic may have additional cancellation relationships.

When debugging shutdown behavior, identify which owner caused cancellation.

```text
host context cancellation
signal policy cancellation
application-owned cancellation
```

Do not infer a signal merely because `ctx.Done()` closed.

## Test semantic policy and process delivery at different layers

Use `clitest` and its signal fixtures for deterministic checks of declared signal semantics when the question is about matching, actions, endpoints, or operation state.

Use a real subprocess and real signal delivery when the claim is that the operating system, process adapter, and command cooperate correctly.

The [signals/interrupt](../examples/signals/interrupt/) example deliberately uses the wider proof. The [signals/custom-source](../examples/signals/custom-source/) example proves an application-defined source without pretending it came from the operating system.

Continue with [Inspection and Tooling](11-inspection-and-tooling.md) to observe the compiled application and invocation state without executing handlers.
