# How Commandment Works

[Getting Started](01-getting-started.md) leaves you with a working command, typed configuration, validation, and one semantic test. The model underneath that program keeps later features predictable as the CLI grows.

The useful starting point is not the type taxonomy. It is the path one program takes.

```text
Go declaration
    -> cli.Compile
    -> immutable Application
    -> operation request
    -> invocation analysis when needed
    -> detached observation and/or execution
```

The same accepted application drives execution, Help, completion, inspection, tooling, and semantic tests. Commandment does not ask each frontend to rediscover the CLI from raw strings.

## Declare the CLI where its ownership is visible

A `cli.Definition[C]` is the authored root of an application. Its declaration callback receives a `CommandBuilder[C]`, and nested command callbacks receive builders for their own configuration types.

```go
type rootConfig struct{}

type scanConfig struct {
	World string
}

var definition = cli.Definition[rootConfig]{
	Name: "alienctl",
	Declare: func(cb *cli.CommandBuilder[rootConfig]) {
		cmd := cli.Command[scanConfig]{
			Name: "scan",
			Declare: func(cb *cli.CommandBuilder[scanConfig]) {
				cb.Arguments.String(&cb.Config.World, "world")
			},
		}
		cb.Subcommands.Command(cmd)
	},
}
```

The lexical nesting is useful information. The `scan` command owns its argument and its `scanConfig`. A reader does not have to jump through `declareRoot`, `declareScan`, and other one-use helpers to reconstruct that ownership.

Extract a helper when it names behavior or a reusable declaration boundary that independently earns a name. Dynamic command providers, completion algorithms, host providers, substantial renderers, handlers, and reusable declaration fragments can all earn extraction. Shortening a callback is not enough by itself.

## Declaration capabilities are temporary

The builder and the declaration capabilities reached through it exist for the active declaration scope. They are not durable application objects.

Use them synchronously while the callback is active. Do not retain a builder or one of its callback-scoped capabilities in package state, use it concurrently, or reuse it in another declaration scope.

Durable declaration values are different. Returned handles such as destinations, declarations, conditions, and accepted relationships can participate in the supported relationships their APIs define. The important distinction is between a callback-scoped authoring capability and a durable semantic value produced through that capability.

This is why a helper that needs to return something should normally return the durable semantic value it creates rather than a builder or capability used to create it.

## Some authoring expressions are inert until you use them

Commandment lets declarations describe relationships before those relationships affect the program. A condition is a common example.

```go
stealth := cb.Flags.Bool(&cb.Config.Stealth).Occurrences(cli.Count.Optional())
stealth.Long("stealth")

cloakLevel := cb.Flags.Int(&cb.Config.CloakLevel)
cloakLevel.AvailableWhen(cb.When.Equal(stealth, true)).Long("cloak-level")
```

`cb.When.Equal(...)` describes a condition. The condition gains semantic effect when a declaration owner such as `AvailableWhen` accepts it.

This pattern keeps authoring declarative. Compose the relationship inside the declaration scope, apply it through the owning declaration operation, and let compilation turn the accepted declarations into the durable application.

## Compile accepts one immutable program

`cli.Compile` validates and accepts the static declaration program.

```go
app, err := cli.Compile(definition)
```

A successful `Application` is immutable. Operations do not reopen the original definition as a mutable registry. Repeated and concurrent operations own their own operation state.

That matters when the application contains invocation-dependent structure. Dynamic command providers, dynamic value domains, host-supplied defaults, and facts can contribute information for one operation without mutating the static application for another operation.

A provider-produced command becomes an ordinary command after its batch is accepted for that operation. A dynamic value-domain entry becomes ordinary domain data for the operation that materialized it. Their origin and lifetime differ from static declarations, but they enter the same semantic model.

## Commands own configuration that handlers actually need

A large Commandment application should not grow one root configuration struct containing every descendant field.

```go
type rootConfig struct {
	Locale string
}

type deployConfig struct {
	Region string
	DryRun bool
}
```

The deploy command should own `deployConfig`. Its handler receives that typed configuration after Commandment has selected the command, resolved values, applied policy, and admitted execution.

Cross-scope relationships should be explicit. Execution access can move an ancestor-owned resolved value into selected-command configuration when the child legitimately needs it. Dependencies can name values required by derivations or dynamic domains. Those relationships are preferable to hidden package globals or a giant shared config object.

[Building Large Applications](16-building-large-applications.md) develops this ownership model across packages.

## An invocation creates operation-owned semantic state

Operations that analyze invocation input produce an effective invocation model for that operation. The model can include command selection, token ownership, resolved values, dynamic provider output, policy outcomes, diagnostics, completion state, and execution admission.

That state is not the immutable static program and it is not application business state.

A typical complete invocation goes through work conceptually like this.

```text
select commands
    -> establish token ownership
    -> materialize demanded dynamic structure
    -> resolve typed values
    -> evaluate conditions and policy
    -> determine execution admission
    -> execute middleware and handler when admitted
```

The exact implementation is not a public pipeline contract. The useful public rule is that semantic preparation belongs to Commandment before application execution begins. Handlers should receive admitted typed configuration rather than reparse `os.Args`, reapply precedence, or re-run validation.

## Pick the operation that owns the question

A compiled application supports several public questions.

| Question | Public operation |
| --- | --- |
| Execute an accepted invocation | `Run` or `RunFromProcess` |
| Explain one invocation without running its handler | `Inspect` |
| Reason about incomplete editor input | `Complete` |
| Read accepted static structure | `Schema` |
| Build reader-oriented static reference data | `Reference` |

Do not stretch one operation into another's job. `Schema` cannot answer which provider-produced command appeared for a particular invocation. `Inspect` is unnecessary when the question is entirely static. Completion clients should not scrape Help text. A real shell integration needs the completion adapter and shell boundary, not just a successful `Application.Complete` call.

The artifacts returned by these operations are detached snapshots. Mutating a returned `Schema`, `Reference`, `Inspection`, `Completion`, or diagnostic report does not mutate the application.

[Inspection and Tooling](11-inspection-and-tooling.md) covers the non-executing tooling operations. [Completion](07-completion.md) owns partial-input analysis and shell integration.

## Partial observations prove only what they actually contain

Some operations can preserve useful evidence even when the whole question cannot be completed. A partial or truncated result can prove presence without proving absence.

If a partial completion contains `europa`, the presence of `europa` is real evidence. If `ganymede` is missing from a truncated result, that does not prove the application has no such candidate.

The same principle matters for provider namespaces, inspection state, diagnostics, and exact collection assertions. `clitest` distinguishes a fact that is disproved from a fact that cannot be proved from incomplete evidence.

## Presentation is a projection of accepted semantics

Help, localization, version output, completion presentation, diagnostics, notices, Schema, Reference, and Inspection should project the accepted application model rather than create parallel registries.

A finite domain declared for a value can appear in Help and completion because both are reading the same value model. Lifecycle declared on a flag spelling can appear in Help and reference output. A rule violation can render a human error while also carrying a structured diagnostic. The feature is not that several frontends happen to agree. The feature is that they have one semantic owner to agree about.

## Ordinary application code stays ordinary Go

Commandment owns CLI semantics. It does not need to become the application's dependency-injection container, domain model, persistence layer, or service framework.

Construct application services in ordinary Go. Keep long-lived application state under application ownership. Use Commandment declarations for command language, typed configuration, policy, presentation, completion, host integration, and execution boundaries that genuinely belong to the CLI.

After execution admission, handlers and middleware are normal Go functions receiving the configuration and resources declared for them.

Continue with [Commands, Flags, and Arguments](03-commands-flags-and-arguments.md) to define the language users type, or [Values and Configuration](04-values-and-configuration.md) to model how typed configuration gets its final meaning.
