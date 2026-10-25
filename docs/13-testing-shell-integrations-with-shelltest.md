# Testing Shell Integrations with shelltest

Use `shelltest` when you own a generated shell integration and need to prove that it works through the actual shell and editor boundary.

That need is different from testing semantic completion.

```text
clitest.Complete
    proves what Commandment semantically completes

shelltest
    proves what the real shell/editor does with generated integration
```

This distinction matters most when an application adds a custom shell provider or changes shell-specific integration behavior. Commandment itself uses the same kind of real-shell proof for its built-in integrations.

## Prove semantic completion before blaming the shell

Start with `clitest.Complete` when the question is about the completion target, candidates, ranking, matching, lifecycle, native delegation, or structured completion diagnostics.

If that semantic result is wrong, the shell adapter is not the first owner of the defect.

Move to `shelltest` when quoting, generated source, editor state, shell protocol, runtime version, sourcing, or the shell process itself is part of the claim.

[Completion](07-completion.md) owns the semantic model.

## A `Target` describes one real integration

A `shelltest.Target[T]` binds a Commandment shell provider to the runtime and interpretation needed for one real shell/editor path.

```go
runtimeConfig := shelltest.Runtime{
	Executable:   bash,
	VersionArgs:  []string{"--version"},
	ParseVersion: parseBashVersion,
}

target := shelltest.Target[string]{
	Provider:  cli.Shell.Bash,
	Runtime:   runtimeConfig,
	Prepare:   prepareBash,
	Interpret: interpretBash,
}
```

`Provider` owns Commandment's generated integration. `Runtime` identifies and starts the shell. `Prepare` establishes the callback-scoped test environment. `Interpret` reduces the editor observation to the stable result the test actually cares about.

Those pieces compose directly with editor input and `shelltest.Complete`.

```go
input := shelltest.Input{
	BeforeCursor: "alienctl sc",
}
screen := shelltest.Complete(t, app, target, input)
if !strings.Contains(screen, "scan") {
	t.Fatalf("settled Bash screen does not contain scan: %q", screen)
}
```

A Bash target normally uses `Prepare` to source the generated integration before the interactive shell starts.

```go
prepareBash := func(setup shelltest.Setup) error {
	integration := setup.IntegrationFile("completion.bash")
	rc := setup.File("bashrc", "source "+shellQuote(integration)+"\n")
	setup.Args("--noprofile", "--rcfile", rc, "-i")
	return nil
}
```

`Interpret` should reduce the resulting `shelltest.Observation` to the stable fact the test wants to assert. The complete [completion/shell](../examples/completion/shell/) test contains the version parser, safe Bash quoting, preparation, interpretation, target construction, `Complete` call, and assertion together.

The complete generic contracts live in GoDoc. The design rule is to keep one target focused on one supported integration path.

## Preparation is setup, not a persistent shell-control API

The setup capability passed to `Prepare` belongs to that target preparation. Use it to install generated integration, write required files, configure environment, or supply arguments needed by the proof.

Do not retain preparation capabilities after the callback returns. Do not use them as an alternate way to mutate an already-running shell session.

Keep setup as small as the integration requires. Reconstructing a developer's entire interactive shell configuration usually makes the test less deterministic without proving more about your integration.

## `Complete` drives the real editor path

`shelltest.Complete` starts the configured runtime, loads the generated integration, presents editor state, and observes what the shell actually does.

The [completion/shell](../examples/completion/shell/) scenario verifies that typing a prefix at a real Bash prompt reaches the intrinsic Commandment completion integration and produces the expected command completion.

That is stronger evidence than checking that `alienctl completion bash` returned non-empty source.

## Generated target version and running shell version are different

A provider may generate source for a requested shell version while the test runtime has its own discovered version.

Keep those identities separate.

- Generation target version says what source was requested from the provider.
- Runtime version says which shell actually executed the integration.

A version-sensitive provider test should state both facts rather than assuming they are interchangeable.

## Interpret stable behavior, not shell noise

A real shell can emit prompts, control sequences, redraws, and runtime-specific formatting that are not the semantic result under test.

`Interpret` should reduce that observation to a stable application-relevant fact such as the settled command line, completion replacement, or other bounded result.

Do not normalize away behavior that actually matters. The point is to discard incidental terminal noise without weakening the contract.

## Let shelltest own process lifecycle

The package should start, communicate with, and clean up the shell process required by the target. Tests should not build a second process-management layer around `shelltest` unless they are deliberately testing something outside its supported lifecycle.

Use deterministic readiness and protocol observations instead of timing sleeps when possible.

If a required shell runtime is unavailable in the current environment, report that proof as unavailable or pending. Do not substitute source generation and label it real-shell verification.

## Custom providers should reuse Commandment semantic completion

A custom provider created with `cli.Shell.Custom` should adapt the target shell/editor protocol to `Application.Complete`. The provider should not maintain another command or value registry.

A good application-level proof therefore has two layers.

1. Semantic completion tests prove the application returns the right Commandment completion result.
2. `shelltest` proves the provider carries the shell's real editor state into that model and returns the result correctly.

This keeps a shell-specific transport defect from being confused with an application completion defect.

## Keep clitest and shelltest separate

Use `clitest` for Commandment semantic facts that do not require a real shell. Use `shelltest` only for the wider shell/editor integration boundary.

A test suite that runs every ordinary completion assertion through Bash is slower, more brittle, and less precise. A suite that never enters Bash cannot prove a Bash integration.

Choose the smallest truthful boundary.

Continue with [Benchmarking with clibench](14-benchmarking-with-clibench.md) for performance measurements of supported Commandment operations.
