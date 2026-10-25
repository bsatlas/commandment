# Building Large Applications

A small Commandment program can keep its entire declaration in one file. A large application should not keep stretching that shape until one callback owns every command, service, policy, and integration.

The scaling rule is simple. Keep Commandment semantics with the part of the CLI that owns them, keep ordinary application dependencies in ordinary Go, and keep one visible place where the application is assembled.

A healthy large application usually has this flow.

```text
process entry
    -> construct ordinary application dependencies
    -> assemble one Commandment definition
    -> compile one immutable application
    -> run or observe that application

command owners
    -> own command configuration
    -> own grammar and values
    -> own command-local policy
    -> bind the handler to ordinary application services

host boundary
    -> supplies process context, facts, defaults, streams, files, and signals
```

The package tree is an implementation choice. Those ownership relationships are the part worth preserving.

## Keep one obvious application assembly owner

A large application needs one place where a reader can answer which top-level capabilities become `alienctl`.

One workable shape is this.

```text
main.go

internal/
    program/
        program.go
        program_test.go
    fleet/
        command.go
    worlds/
        command.go
    services/
        services.go
    tooling/
        reference.go
```

`main.go` owns process adaptation. `internal/program` owns root Commandment assembly. Domain packages own the commands and configuration that belong to them. Ordinary services remain ordinary Go values.

Do not treat this exact tree as a framework convention. A smaller application may not need `internal/program`. A larger application may have many domain packages. The useful property is that root composition is easy to find and command semantics are not centralized merely because `Compile` starts at one definition.

## Bind real dependencies at real package boundaries

A command package earns a constructor when it binds something that genuinely varies by application instance, such as a service, client, durable policy, catalog, or provider state.

```go
type FleetService interface {
	Status(context.Context) (FleetStatus, error)
}

func Command(service FleetService) cli.Command[struct{}] {
	statusCmd := cli.Command[statusConfig]{
		Name: "status",
		Run: func(ctx context.Context, cfg *statusConfig) error {
			status, err := service.Status(ctx)
			if err != nil {
				return err
			}
			return writeStatus(cfg.Output, status)
		},
	}
	fleetCmd := cli.Command[struct{}]{
		Name: "fleet",
		Declare: func(cb *cli.CommandBuilder[struct{}]) {
			cb.Subcommands.Command(statusCmd)
		},
	}
	return fleetCmd
}
```

The constructor is not valuable because the source file is long. It is valuable because the `fleet` package owns a command and the service dependency its handler needs.

When a command has no real dependency or reusable boundary, keep static declaration lexical rather than manufacturing a factory solely to move lines into another file.

## Let command packages own command configuration

Configuration should follow the command that consumes it.

```go
type launchConfig struct {
	Ship  string
	Speed int
}
```

A `fleet launch` package or owner should declare the argument, speed value, policy, lifecycle, middleware, and handler associated with that configuration. Do not grow a root struct containing every descendant field merely so one package can reach another package's values.

Cross-scope relationships should be explicit. If a descendant legitimately depends on an ancestor value, use the supported dependency or execution-access relationship rather than a package global or hidden pointer into another command's configuration.

## Do not move declaration capabilities across lifetimes

Builders and callback-scoped capabilities are temporary authoring capabilities. A package boundary does not make them durable.

A helper may receive one synchronously when the package boundary genuinely owns that declaration work. It must finish using the capability before returning.

```text
callback-scoped builder or capability
    -> use synchronously
    -> never retain
    -> never use concurrently

durable semantic value
    -> return or reuse only when its public contract permits it
```

This distinction matters more as the repository grows. A package layout that requires storing builders in structs or global registries is fighting Commandment's ownership model.

## Keep static topology static

If the command tree is known when the application is built, declare that topology statically.

A command provider is appropriate when the namespace itself depends on invocation-time or external data. For example, an `outpost` command may expose children discovered from a catalog only after `--system` is known.

```go
provideOutposts := func(_ context.Context, batch *cli.CommandBatch[outpostConfig]) error {
	for _, name := range catalog.Outposts(batch.Config.System) {
		outpost := name
		cmd := cli.Command[struct{}]{
			Name: outpost,
			Run: func(context.Context, *struct{}) error {
				return contactOutpost(outpost)
			},
		}
		batch.Command(cmd)
	}
	return nil
}
cb.Subcommands.CommandsFrom(provideOutposts, system)
```

Do not use providers as a source-organization technique. A provider that always reconstructs a statically known subtree adds invocation-time machinery without adding application meaning.

See [Dynamic commands](../examples/commands/dynamic/).

## Keep shared policy narrow and named for the application

Reuse is valuable when the meaning repeats.

A durable value language for `PlanetID`, a release scheme shared by the whole CLI, or a domain policy used by several related commands can deserve a shared owner. A `common` package full of every repeated flag, helper, and rule usually does not.

Prefer abstractions with domain meaning.

```text
navigation.PlanetDomain
fleet.ReleasePolicy
host.CredentialSource
```

Avoid abstractions whose main purpose is hiding Commandment calls.

```text
common.AddStringFlag
helpers.RegisterCommand
cliutil.BuildRule
```

The first group explains why the abstraction exists. The second group often hides ownership.

## Keep host integration at the boundary where the information enters

Not every dependency belongs in command configuration.

| Information | Normal owner |
| --- | --- |
| User-selected CLI meaning | values and command configuration |
| Trusted invocation context | facts |
| Overridable host fallback | invocation defaults |
| Commandment execution resources | execution access |
| Long-lived services and business dependencies | ordinary Go |
| Process arguments, environment, streams, signals | executable or host adapter |

A host fact may require coordination between the command that declares its destination and the root program that supplies the fact provider. That is a real semantic relationship. Keep it explicit rather than turning all services into facts or reaching into process globals from the handler.

[Runtime and Host Integration](09-runtime-and-host-integration.md) covers those boundaries in detail.

## Let tooling consume the compiled application

A large CLI often needs more than one executable path around its command model. Documentation generation, configuration explainers, IDE tooling, compatibility checks, and release tooling may all need structural information.

Compile the same definition and give the tool the public application or detached observation it needs.

```go
reference, err := app.Reference(cli.ReferenceOptions{})
```

Do not maintain another command registry for tooling. The point of `Schema`, `Reference`, `Inspect`, and `Complete` is that tools can ask the accepted application directly.

[Inspection and Tooling](11-inspection-and-tooling.md) explains which operation owns each question.

## Put cross-cutting behavior in the right layer

Middleware is for admitted execution behavior such as tracing, authorization of an admitted action, or an application transaction boundary.

Do not move these concerns into middleware merely because they affect many commands.

```text
grammar
value precedence
availability
validation rules
diagnostics
Help
completion
```

Those semantics need to exist before handler execution so inspection, completion, Help, diagnostics, and tests can observe the same model.

## Test each contract at its natural boundary

Large applications become easier to maintain when every test does not start a process.

Use `clitest` for Commandment-owned semantic facts such as command selection, value resolution and provenance, policy, diagnostics, lifecycle, completion, Schema, Reference, and Inspection.

Use ordinary Go tests for domain services, handlers, middleware effects, and application-owned computation.

Use subprocess tests for process-visible contracts such as exit status and stream behavior. Use `shelltest` for real shell or editor integration.

A root semantic test can cheaply prove that assembly still contains the expected major subtrees while package-local tests prove the details owned by each command package.

## Judge the architecture by ordinary changes

The best package tree is the one that keeps normal evolution local.

| Change | Files that should normally care |
| --- | --- |
| Add a command inside `fleet` | `fleet` owner and its tests |
| Add a command-local flag or rule | that command owner and its tests |
| Deprecate one spelling | declaration owner, plus root release authority only if the application did not already have one |
| Add a genuinely dynamic namespace | namespace owner and root assembly where it is attached |
| Add a host fact | semantic destination owner, host/provider owner, and root process assembly |
| Add a reference generator | tooling owner that consumes the existing application |
| Move a command path | old/new semantic owners plus compatibility tests if the old path remains a contract |

If adding one command repeatedly requires edits to a universal configuration struct, unrelated feature packages, a custom registry, and several parallel Help/completion descriptions, the architecture has lost semantic ownership.

## Moving source is not automatically behavior-neutral

Moving a Go file between packages is an internal refactor. Moving a command in the CLI tree changes its canonical command path unless compatibility is preserved deliberately.

If `alienctl fleet status` becomes `alienctl status`, decide whether the old path disappears, remains as compatibility syntax, or enters lifecycle deprecation. Then prove that decision at the semantic and process boundaries that matter.

Do not let a source reorganization silently become a CLI redesign.

## Keep process adaptation thin

Once a large application has a real assembly owner, `main.go` can remain small.

```go
dependencies := program.Dependencies{
	Fleet:   fleetService,
	Catalog: outpostCatalog,
}
assembled := program.Build(dependencies)
app, err := cli.Compile(assembled.Definition)
```

If the program owns host fact/default providers, the same assembly result can expose those providers to an explicit `RunRequest`. If it does not, `RunFromProcess` remains the smallest process adapter.

The executable should construct dependencies, compile the application, adapt the host process, and choose the exit path. It should not become the place where every command's semantics are declared.

## A healthy large application remains understandable from both directions

Starting from `main.go`, a reader should be able to find root assembly and then the semantic owner of any command.

Starting from a command package, a reader should be able to understand that command's configuration, grammar, policy, behavior, and real dependencies without reconstructing a global mutable registry.

Starting from a tool or test, a reader should be able to compile the same definition and query public Commandment observations rather than rebuild the CLI model.

If all three paths stay clear as the application grows, the architecture is doing its job.

For an existing mature CLI, continue with [Adopting Commandment](17-adopting-commandment.md).
