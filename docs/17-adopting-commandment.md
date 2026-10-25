# Adopting Commandment

Replacing the frontend of an established CLI is a compatibility project before it is a framework project.

Do not begin by translating old parser calls into Commandment declarations. Begin by deciding which observable contracts users, scripts, shells, operators, and application code actually depend on.

The migration then has a concrete shape.

```text
observe existing contract
    -> decide preserve or change
    -> assign the Commandment owner
    -> choose the proof boundary
    -> implement
    -> compare the result
```

This works for a small home-grown parser, a CLI built on another framework, or a large kubectl-scale application.

## Build a compatibility ledger before changing the frontend

Record behavior by responsibility rather than by old implementation type.

| Existing contract | Decision | Commandment owner | Proof boundary |
| --- | --- | --- | --- |
| `fleet launch <ship>` command path | preserve | commands and arguments | `clitest` plus selected process checks |
| `--config` before or after the command | preserve or deliberately change | grammar profile / dialect | semantic grammar tests |
| flag beats environment | preserve | value producers and resolution | `clitest` Inspection |
| invalid region prints a replacement hint | preserve | rejection, Diagnostic, repair | `clitest` plus rendered error check |
| Bash completes dynamic outposts | preserve | completion and shell provider | semantic completion plus `shelltest` |
| SIGINT begins graceful shutdown | preserve | signal policy | semantic signal test plus subprocess signal test |
| exit code 7 for application rejection | preserve | application error mapping | subprocess |
| startup under a release-specific threshold | preserve only if it is a real requirement | application/process performance owner | benchmark at that boundary |

A ledger forces every migration difference to be deliberate. It also stops old source structure from masquerading as a public contract.

## Inventory the language users type

Capture the external grammar first.

Record at least the behaviors that matter for the application.

- root and nested command paths;
- aliases and shortcuts;
- default commands;
- long and short flag spellings;
- flag counting and implicit values;
- positional argument forms and optionality;
- pass-through arguments;
- end-of-options behavior;
- flag placement and matching rules;
- command matching and normalization;
- response-file recognition;
- completion-visible grammar.

Then model the intended language with Commandment commands, flags, arguments, grammar profiles, dialect, token languages, and command entries.

Do not preserve an old parser object model merely because it happened to produce the grammar. Preserve the grammar itself when that is the contract.

If the product should improve an awkward legacy grammar, classify that as an intentional CLI change and prove it separately from the framework migration.

[Commands, Flags, and Arguments](03-commands-flags-and-arguments.md) owns the target authoring model.

## Inventory configuration behavior independently from grammar

A mature CLI often has years of accumulated configuration rules hidden across parser setup, environment loading, config files, globals, and handler defaults.

For each meaningful value, record

- its Go type;
- every accepted source;
- precedence;
- static and host defaults;
- repeated-input behavior;
- legal domain or validation;
- dynamic dependencies;
- provenance requirements;
- sensitivity requirements.

Then map those responsibilities onto one Commandment destination and the producers, resolution rules, domains, derivations, and policy that actually own the final value.

Do not copy an old merge function into every handler. The migration should converge duplicated precedence logic into the value model.

When precedence matters operationally, prove both the effective value and the contributing producer. A test that sees `4` is weaker than a test that proves `4` came from the flag rather than the environment or a default.

[Values and Configuration](04-values-and-configuration.md) develops that model.

## Inventory invalid-input behavior and diagnostics

Do not reduce compatibility to successful invocations.

Record important failures too.

- unknown commands and flags;
- unavailable conditional syntax;
- missing required input;
- parse failures;
- rule violations;
- replacement suggestions;
- advisory warnings;
- exit status;
- structured error information consumed by tools.

Commandment may let you improve weak legacy failures. Decide whether exact text, semantic diagnostic identity, repair guidance, exit behavior, or only the broad rejection is the actual contract.

Prefer semantic compatibility where exact presentation has never been promised. Preserve exact text only when users or automation genuinely depend on it.

[Conditions, Validation, and Diagnostics](05-conditions-validation-and-diagnostics.md) explains the target failure model.

## Inventory Help, localization, version, and completion as product behavior

Many mature CLIs accidentally maintain separate command models for parsing, Help, completion, and documentation.

Record what users actually depend on.

- Help organization and descriptions;
- examples;
- hidden or visible declarations;
- lifecycle and replacement guidance;
- localization and locale selection;
- version output;
- semantic completion behavior;
- filesystem or other native delegation;
- generated shell integrations;
- supported shell families.

The target should let Commandment project these behaviors from the accepted semantic application rather than preserve a second legacy registry solely for presentation.

Use [Help and Localization](06-help-and-localization.md), [Completion](07-completion.md), and [Evolving a CLI](08-evolving-a-cli.md) for the target model.

## Inventory the host and process boundary

Record how the current executable interacts with its environment.

- process arguments and environment;
- working directory;
- standard streams;
- response files;
- facts and host-provided defaults;
- application services;
- context cancellation;
- signals;
- application errors and exit codes;
- locale and terminal capabilities.

Keep the resulting ownership honest.

User-selected meaning belongs in Commandment values. Trusted invocation context can belong in facts. Overridable host fallback can belong in invocation defaults. Long-lived clients, stores, and domain services remain ordinary Go dependencies. Process adaptation belongs near the executable boundary.

Do not turn every existing global into a CLI value simply because the migration is an opportunity to remove globals.

See [Runtime and Host Integration](09-runtime-and-host-integration.md) and [Signals](10-signals.md).

## Inventory lifecycle and compatibility syntax

A mature CLI usually has syntax that exists because history exists.

Record which commands, flag spellings, aliases, shortcuts, and other declarations are

```text
current
newly introduced
deprecated
scheduled for removal
replacement syntax
legacy behavior that is not actually a contract
```

Then encode real lifecycle and replacement relationships explicitly.

Do not keep an internal spelling forever merely because it once existed. Do not remove widely used syntax merely because the new declaration is cleaner. Make the compatibility decision at the product boundary and let Help, diagnostics, reference tooling, and tests project the same decision.

[Evolving a CLI](08-evolving-a-cli.md) owns this workflow.

## Preserve only real performance contracts

A migration is not automatically required to reproduce every benchmark number from the old implementation.

Keep a performance constraint when the application genuinely treats it as a release requirement. Name the exact boundary being measured.

```text
Commandment operation latency
process startup latency
shell completion round trip
memory at a representative workload
application service latency
```

Use `clibench` when the question is exactly one supported successful Commandment operation over a reusable workload. Use a wider benchmark when the contract includes process startup, shell execution, network work, or another owner outside that operation.

[Benchmarking with clibench](14-benchmarking-with-clibench.md) explains the narrow helper boundary.

## Design the target ownership from the contracts, not the legacy packages

Once the compatibility ledger is explicit, design the new application around semantic ownership.

A common target has one root assembly owner, command/domain packages that own their typed configuration and CLI semantics, ordinary Go packages for application services, and thin host/process integration.

Do not preserve old package boundaries automatically. An old `flags` package or global command registry may exist only because the previous frontend required it.

Do not flatten everything into one Commandment package either. The goal is not replacing one framework-shaped architecture with another framework-shaped architecture.

[Building Large Applications](16-building-large-applications.md) covers the target composition in detail.

## Build a compatibility corpus with one proof owner per contract

Choose the smallest legitimate proof for each ledger row.

Use `clitest` for Commandment-owned semantic facts such as

- command selection;
- value resolution and provenance;
- conditions and availability;
- rules and diagnostics;
- repairs;
- lifecycle and replacement;
- completion semantics;
- Schema, Reference, and Inspection relationships.

Use ordinary Go tests for handlers, domain services, and application-owned computation.

Use subprocess tests when the contract includes process exit, standard streams, process environment, or another executable boundary.

Use `shelltest` when the contract includes a real shell or editor integration.

A single giant golden test is rarely the strongest migration proof. It makes unrelated changes look equivalent and makes failures harder to attribute.

## Migrate in coherent ownership units

The right increment depends on the existing application. It might be a command subtree, a value/configuration family, a process adapter, or another coherent user-visible boundary.

Each increment should leave one clear owner for the behavior it migrates.

Do not invent a dual-runtime architecture merely to claim incremental migration. If the legacy and Commandment frontends cannot safely coexist in one process, perform the work behind source and test boundaries and switch the executable once the compatibility corpus for the chosen cutover is green.

The supported migration unit is determined by the application and current product semantics, not by a desire to keep two parsers alive indefinitely.

## Classify every difference before changing the expectation

When old and new behavior differ, use a small decision vocabulary.

| Classification | Meaning |
| --- | --- |
| Bug | Intended compatibility changed accidentally |
| Intentional redesign | Product behavior is deliberately changing |
| Old implementation detail | The observed difference was never a supported external contract |
| Unsupported legacy behavior | The target cannot currently represent an important contract and a product decision is required |

Do not automatically update a test because the new implementation is different. Do not automatically preserve an old snapshot because it exists.

The classification comes first.

## Re-prove external behavior at the boundary that owns it

Before cutover, rerun the compatibility corpus against the exact target candidate.

A useful final pass asks

```text
Does the same command select?
Does the same value win for the same reason?
Does invalid input fail with the intended semantic guidance?
Does Help expose the intended compatibility state?
Does semantic completion agree?
Do supported shells still behave correctly?
Do application errors map to the intended exit behavior?
Do signals still produce the intended process outcome?
Do real performance requirements still pass?
```

Then inspect the compiled application with `Schema`, `Reference`, and representative `Inspect` requests. Those observations often expose duplicate declarations, missing lifecycle, unexpected precedence, or presentation drift more clearly than source review alone.

## A successful adoption should simplify the frontend

The target should have one compiled semantic model rather than independent parser, Help, completion, diagnostics, and test descriptions.

Commands should own typed configuration. Resolution should be explicit rather than repeated as handler-side merge logic. Dynamic namespaces should use providers only when they are genuinely dynamic. Host integration should remain distinct from CLI semantics. Compatibility and lifecycle should be explicit. Tests should prove semantic, process, and shell contracts at their natural boundaries.

That is a stronger endpoint than making every old command run. It preserves the behavior worth preserving while removing frontend architecture that no longer earns its existence.

For a specific migration failure, continue with [Troubleshooting](15-troubleshooting.md). For exact declaration contracts, use GoDoc after the relevant guide has established the mental model.
