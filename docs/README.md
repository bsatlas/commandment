# Commandment documentation

The [root README](../README.md) shows what Commandment can do. This manual teaches how to build with it.

Start with [Getting Started](01-getting-started.md) if you have not built a Commandment application yet. Read [How Commandment Works](02-how-commandment-works.md) when you want the mental model that makes the focused guides predictable.

Use the rest of the manual by the problem you are solving.

## Build the CLI users type

| You want to | Read |
| --- | --- |
| Define commands, arguments, flags, aliases, shortcuts, defaults, grammar, or dynamic command namespaces | [Commands, Flags, and Arguments](03-commands-flags-and-arguments.md) |
| Resolve typed configuration from flags, environment, defaults, host input, custom Go types, domains, or derived values | [Values and Configuration](04-values-and-configuration.md) |
| Make syntax conditional, validate input, produce structured diagnostics, offer repairs, or emit advisories | [Conditions, Validation, and Diagnostics](05-conditions-validation-and-diagnostics.md) |
| Control Help, descriptions, groups, custom rendering, localization, or locale selection | [Help and Localization](06-help-and-localization.md) |
| Complete commands and values, delegate native completion, generate shell integration, or add a shell provider | [Completion](07-completion.md) |
| Introduce, deprecate, replace, or remove CLI syntax across releases and present application version information | [Evolving a CLI](08-evolving-a-cli.md) |

## Connect the CLI to the application and host

| You want to | Read |
| --- | --- |
| Run the application, supply facts or invocation defaults, use streams, middleware, response files, application errors, or exit behavior | [Runtime and Host Integration](09-runtime-and-host-integration.md) |
| Handle operating-system or application-defined signals, cancellation, typed events, forwarding, escalation, or reraising | [Signals](10-signals.md) |
| Build configuration explainers, reference generators, static analyzers, or other tools from the compiled application | [Inspection and Tooling](11-inspection-and-tooling.md) |

## Prove behavior

| You want to | Read |
| --- | --- |
| Assert Commandment-owned semantic facts without recreating parser behavior in tests | [Testing with clitest](12-testing-with-clitest.md) |
| Prove generated completion against a real supported shell or editor boundary | [Testing Shell Integrations with shelltest](13-testing-shell-integrations-with-shelltest.md) |
| Benchmark one supported Commandment operation over a reusable workload | [Benchmarking with clibench](14-benchmarking-with-clibench.md) |
| Diagnose a wrong command, value, diagnostic, Help result, completion, process result, or signal behavior | [Troubleshooting](15-troubleshooting.md) |

## Grow or replace a real CLI

[Building Large Applications](16-building-large-applications.md) shows how to preserve semantic ownership as the command tree, host integration, tooling, and tests spread across packages.

[Adopting Commandment](17-adopting-commandment.md) starts from the observable contracts of an existing CLI and turns them into an explicit compatibility and migration plan.

## Runnable examples and exact API contracts

The [examples corpus](../examples/README.md) contains focused executable `alienctl` scenarios. Use it when you want complete source worth copying after a guide has established the concept.

Use GoDoc when you already understand the concept and need an exact local contract such as a signature, zero or nil behavior, callback lifetime, ordering, snapshotting, concurrency, or failure semantics.
