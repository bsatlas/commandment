# Evolving a CLI

A shipped CLI accumulates users, scripts, documentation, shell history, and automation. Renaming a flag or removing a command is therefore a compatibility change, not merely a source edit.

Commandment can keep release identity, declaration lifecycle, replacement, Help, Reference, completion, and tests connected to one model.

## Give the application one release authority

A definition can declare both its current release and the comparison scheme used for lifecycle decisions.

```go
var definition = cli.Definition[struct{}]{
	Name:    "alienctl",
	Release: cli.Release.SemVer("2.0.0"),
	Version: cli.VersionConfig{},
}
```

With semantic-version authority, the intrinsic version action can project the same release.

```text
alienctl --version
2.0.0
```

See [evolution/version](../examples/evolution/version/).

`Release.SemVer` is appropriate when the application's releases follow semantic-version ordering. Commandment also supports an explicit ordered release scheme for applications whose release labels have another deliberate order. The release scheme is comparison authority for lifecycle evaluation. It is not a generic package-version parser for unrelated application data.

## Describe when a declaration entered, became deprecated, or stopped existing

A lifecycle records release boundaries.

```go
legacy := speed.Long("velocity").Release(cli.ReleaseLifecycle{
	IntroducedIn: "1.0.0",
	DeprecatedIn: "2.0.0",
	Reason:       "velocity was renamed",
})
```

`IntroducedIn`, `DeprecatedIn`, and `RemovedIn` answer different questions.

- Before `IntroducedIn`, the declaration is not yet part of that release's effective surface.
- From `DeprecatedIn`, the declaration remains present but is deprecated.
- At or after `RemovedIn`, the declaration no longer participates as a current/deprecated declaration for that release.

The exact lifecycle states projected by Schema, Reference, Help, completion, and inspection come from this shared evaluation rather than independent version checks in each frontend.

## Use replacement to tell users where to move

Lifecycle says that a declaration is deprecated or removed. Replacement says which declaration should take its place.

```go
current := speed.Long("speed")
legacy.ReplacedBy(current)
```

The relationship is deliberately separate. A declaration can be deprecated without a replacement, and a replacement relationship should not be inferred merely because two declarations look similar.

The [evolution/deprecation](../examples/evolution/deprecation/) scenario keeps `--velocity` available at release `2.0.0` while directing users to `--speed`.

## Let Help carry migration information at the point of use

At release `2.0.0`, the same declaration model can render this.

```text
--velocity VALUE [deprecated 2.0.0: velocity was renamed; use --speed],
--velocity=VALUE [deprecated 2.0.0: velocity was renamed; use --speed],
--speed VALUE, --speed=VALUE
  Set approach speed.
```

That is more useful than a separate migration table that can drift away from the syntax itself.

Help presentation is not the lifecycle authority. It is one projection. A tooling client can inspect lifecycle and replacement structurally through Schema or Reference, and completion can carry lifecycle state without parsing Help text.

## Distinguish aliases and compatibility syntax from declaration replacement

An alias is another accepted name for the same command. A deprecated spelling can remain accepted for compatibility. A replacement relates one declaration to the declaration users should adopt.

Those are different facts.

Do not use aliases merely to avoid declaring lifecycle when the compatibility surface is genuinely deprecated. Do not declare replacement merely because an alias exists. Preserve the relationship that actually describes the public CLI.

## Removal should be deliberate and testable

Before adding `RemovedIn`, decide whether the application truly intends to stop accepting the old declaration at that release.

The migration sequence for a typical rename is often

```text
release N
    old spelling current

release N+1
    old spelling deprecated
    replacement -> new spelling
    both accepted

later release
    old spelling removed
    new spelling current
```

That sequence is a product compatibility decision. Commandment can project it consistently, but it should not decide the application's removal policy for you.

## Test lifecycle semantically instead of only snapshotting text

A semantic test can locate the old and new spellings and assert their lifecycle and relationship directly.

```go
flag := clitest.Command("approach").Flag("--speed")
legacy := flag.Spelling("--velocity")
current := flag.Spelling("--speed")
clitest.Schema(app).Require(t,
	legacy.Lifecycle(cli.LifecycleStateDeprecated),
	legacy.ReplacedBy(current),
	current.Lifecycle(cli.LifecycleStateCurrent),
)
```

Use a focused Help assertion as well when the exact migration presentation is important to the product experience. Semantic lifecycle tests and presentation tests answer different questions.

## Keep version presentation derived from release authority

Do not maintain one version string for `--version`, another for lifecycle comparison, and a third for generated docs unless the application truly has distinct concepts.

The definition's release authority should normally be the source from which built-in version presentation and lifecycle evaluation proceed.

Application build metadata can still have its own ownership where the API supports it. The key rule is to avoid shadow release registries that frontends must keep synchronized manually.

## Provider-produced declarations follow the same release model

Dynamic commands or other provider-produced declarations may carry lifecycle metadata where the provider contract supports it. Invocation-time origin does not create a second lifecycle system.

The provider still publishes declarations into the operation's semantic model, and release validation applies under the application's release scheme.

## Treat lifecycle as a public contract during adoption

When migrating an existing CLI to Commandment, inventory which names, spellings, commands, and deprecation behaviors are intentionally supported today. Preserve the ones you mean to preserve and deliberately change the ones you do not.

[Adopting Commandment](17-adopting-commandment.md) turns those observations into a migration ledger and proof plan.

Continue with [Runtime and Host Integration](09-runtime-and-host-integration.md) for process execution and explicit host capabilities.
