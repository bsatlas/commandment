# Values and Configuration

Most real CLIs accept configuration from more than one place. A scan range may have a built-in default, an environment override, and an explicit flag. Application code should still receive one typed value.

Commandment models that problem around a destination rather than around the parser surface that happened to produce it.

## Give one application value one destination

```go
beamStrength := cb.Values.Int(&cb.Config.BeamStrength)
cb.Values.Default(beamStrength, 1)
cb.Values.Environment(beamStrength, "ALIEN_BEAM_STRENGTH")
cb.Flags.ValueFor(beamStrength).Long("beam-strength")
```

The destination is the typed semantic value. The default, environment input, and flag are producers that may contribute candidates for it.

The handler reads `cfg.BeamStrength`. It does not ask whether the value came from argv, environment, a default, or another supported producer.

## Let competing sources resolve before the handler runs

Resolution determines which accepted contributions become effective. Commandment separates producer identity from precedence so tools can explain both what contributed and why it won.

Built-in resolution policies include `KeepFirst`, `KeepLast`, `RejectMultipleTiers`, and `CombineTiers` where the destination shape supports them.

```go
cb.Values.ResolveWith(beamStrength, cli.Value.KeepLast())
```

A resolution tier is a precedence stratum, not a synonym for input source. An environment variable is an input source. The environment-input declaration is a producer. A candidate records the producer kind and participates at a resolution tier. Keeping those concepts separate is what lets provenance and inspection answer precise questions later.

Absence is also distinct from a resolved Go zero value. `0`, `false`, and `""` can be legitimate resolved values, so tools should not infer presence from the Go representation alone.

[Getting Started](01-getting-started.md) shows a destination receiving contributions from several ordinary sources.

## Use application-defined Go types with `ValueCodec[T]`

A framework cannot have a built-in parser for every meaningful application type. `ValueCodec[T]` lets the application define one scalar value language without building a parallel configuration system.

```go
addressCodec := cli.ValueCodec[netip.Addr]{
	Name:     "ip-address",
	Parse:    netip.ParseAddr,
	Format:   netip.Addr.String,
	Snapshot: cli.Value.Assign[netip.Addr](),
}

address := cb.Values.Value(&cb.Config.Address, addressCodec)
cb.Flags.ValueFor(address).Long("address")
```

The codec owns three things. `Parse` converts text to the application type. `Format` gives Commandment a stable textual projection when it needs one. `Snapshot` defines how a retained typed value is detached from caller mutation.

The resulting destination participates in arguments, flags, domains, diagnostics, completion, artifacts, and tests just like a built-in scalar language.

Use a codec for one semantic scalar text value. Use a `TokenLanguage[T]` when a value genuinely spans a bounded sequence of CLI tokens. Use keyed value languages for keyed input shapes rather than hiding key parsing inside an unrelated scalar codec.

See [values/custom-types](../examples/values/custom-types/) for a complete `netip.Addr` application.

## Make legal choices part of the value

A finite domain records the accepted values in the value model itself.

```go
declareApproaches := func(values *cli.ValueDomainBuilder[string]) {
	values.Value("orbit", "orbit")
	values.Value("land", "land")
	values.Value("flyby", "flyby")
}

approaches := cli.Value.Domain(cli.Value.Assign[string](), declareApproaches)
approach := cb.Values.Value(&cb.Config.Approach, approaches)
cb.Values.Default(approach, "orbit")
cb.Flags.ValueFor(approach).Long("approach")
```

The domain is not merely a parser check. Help, completion, inspection, diagnostics, Reference, and `clitest` can all project the same accepted entries.

See [values/choices](../examples/values/choices/).

## Let resolved input control a dynamic domain

Some choices depend on values already settled for this invocation. `alienctl visit` can use the selected planetary system to decide which moons are legal.

```go
declareMoons := func(_ context.Context, _ *visitConfig, dependencies cli.Dependencies, domain *cli.ValueDomainBuilder[string]) error {
	resolvedSystem, _ := dependencies.Get(system)
	switch resolvedSystem {
	case "jovian":
		domain.Value("europa", "europa")
		domain.Value("ganymede", "ganymede")
	case "saturnian":
		domain.Value("titan", "titan")
		domain.Value("enceladus", "enceladus")
	}
	return nil
}

moons := cb.Values.Domain(declareMoons, system)
moon := cb.Values.Value(&cb.Config.Moon, moons)
cb.Arguments.ValueFor(moon, "moon")
```

Passing `system` to `Values.Domain` declares the dependency. Commandment can resolve the system before materializing the moon domain and can validate dependency scope and cycles instead of discovering them accidentally inside the callback.

The domain materialization belongs to that operation. It does not mutate the static application for future invocations.

See [values/dynamic-choices](../examples/values/dynamic-choices/).

## Compute configuration from values Commandment already resolved

A derived value is useful when application configuration should be computed from other semantic inputs before execution.

```go
deriveSummary := func(_ context.Context, _ *abductConfig, dependencies cli.Dependencies) (string, error) {
	resolvedTarget, _ := dependencies.Get(target)
	resolvedBeamStrength, _ := dependencies.Get(beamStrength)
	summary := fmt.Sprintf("abducting %s with beam strength %d", resolvedTarget, resolvedBeamStrength)
	return summary, nil
}

summary := cli.Value.ComputedBy(deriveSummary).From(target, beamStrength)
cb.Values.Derive(&cb.Config.Summary, summary)
```

`.From(...)` is not decorative metadata. It declares the destinations the computation requires, which lets Commandment reason about ordering, scope, and dependency validity.

The handler receives the finished `Summary`. It should not repeat the computation after Commandment has already resolved the source values.

See [values/derived](../examples/values/derived/).

## Ask where a value came from with provenance

Inspection preserves the relationship between the effective value, the candidate that contributed it, the semantic producer, and the input source history.

A tool can therefore answer questions such as these without parsing logs.

- Did `--range` win over `ALIEN_RANGE`?
- Was the effective value supplied by a host default or an explicit flag?
- Did a shortcut contribute this argument?
- Which derived value depends on this destination?
- Which source span produced an invalid value?

At the lowest level, `Inspection.Values` exposes contributor identities and producer observations. `clitest` provides higher-level expectations such as `ContributedBy` for common semantic checks.

[tooling/inspect](../examples/tooling/inspect/) demonstrates using `Schema` to locate static subjects and `Inspect` to read one invocation.

## Mark the semantic value sensitive when all of its presentations need protection

```go
token := cb.Values.String(&cb.Config.Token).Sensitive()
cb.Flags.ValueFor(token).Long("token")
```

Destination sensitivity says the semantic value itself is sensitive. Commandment can redact its representation in Commandment-owned presentation and observations.

Sometimes only one producer is sensitive.

```go
token := cb.Values.String(&cb.Config.Token)
cb.Values.Environment(token, "ALIEN_TOKEN").Sensitive()
cb.Flags.ValueFor(token).Long("token")
```

In that case the environment contribution is sensitive without asserting that every possible producer for the destination must be treated identically.

Flag spellings and other supported producer declarations can own sensitivity for the same reason.

Sensitivity is presentation and observation policy. It is not encryption, secret storage, or application authorization.

## Keep host facts distinct from ordinary CLI values

A fact is trusted typed context supplied by the host outside ordinary CLI input. An invocation default is a host-supplied fallback contribution to an ordinary destination.

Those are different responsibilities.

A fact can answer a question like "which spacecraft is this process running aboard?" without pretending the user typed the answer. A host default can provide a default spacecraft while still allowing an explicit `--ship` value to override it under the destination's resolution policy.

See [host-context/facts](../examples/host-context/facts/) and [host-context/defaults](../examples/host-context/defaults/). [Runtime and Host Integration](09-runtime-and-host-integration.md) owns provider construction and `RunRequest` wiring.

## Use resolution state instead of guessing from representations

Inspection distinguishes states such as resolved, failed, blocked, inactive, or otherwise incomplete according to the public resolution model. A tooling client should inspect that state instead of treating a missing textual representation as synonymous with "no value."

Derived destinations are a useful example. A derived value can be resolved even when it has no ordinary scalar textual representation suitable for user presentation. Tests can assert the resolution state and producer kind without inventing a string rendering that the destination never declared.

## Test the semantic value rather than duplicating precedence logic

`clitest` can locate a value through the command grammar and assert its resolved meaning, provenance, domain, sensitivity, and producer relationships.

```go
beamStrength := clitest.Command("abduct").Flag("--beam-strength").Value()
clitest.Inspect(t.Context(), app, request).Require(t,
	clitest.InspectionReady(),
	beamStrength.Text("4"),
	beamStrength.ContributedBy(cli.ValueProducerKindFlag),
)
```

Use an ordinary Go test for application-owned computation after the configuration reaches the handler. Use `clitest` when the fact belongs to Commandment's value model.

Continue with [Conditions, Validation, and Diagnostics](05-conditions-validation-and-diagnostics.md) when the next question is whether an otherwise well-typed invocation is allowed, invalid, repairable, or advisory.
