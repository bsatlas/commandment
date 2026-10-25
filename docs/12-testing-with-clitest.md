# Testing with clitest

`clitest` is the semantic testing frontend for Commandment applications. Use it when the test needs to prove what Commandment understood, not merely what a handler happened to print.

A good `clitest` test starts from an application question.

```text
Did this command exist?
Which command was selected?
Which producer supplied this value?
Was this flag available?
Which rule blocked execution?
Did this diagnostic carry a repair?
What completion item was produced?
```

The assertion should name that fact directly.

## Prove why a value won

```go
func TestBeamStrengthFlagWins(t *testing.T) {
	t.Parallel()

	app := clitest.Compile(t, definition)
	beamStrength := clitest.Command("abduct").Flag("--beam-strength").Value()
	request := cli.InspectionRequest{
		Env: cli.Env{{
			Name:  "ALIEN_BEAM_STRENGTH",
			Value: "2",
		}},
		Args: []string{"abduct", "--beam-strength", "4", "cow"},
	}
	clitest.Inspect(t.Context(), app, request).Require(t,
		clitest.InspectionReady(),
		beamStrength.Text("4"),
		beamStrength.ContributedBy(cli.ValueProducerKindFlag),
	)
}
```

A handler-output test could prove that application code received `4`. It would not prove that the explicit flag contribution beat the environment contribution under Commandment resolution. `clitest` asks the public semantic observation directly.

## Choose the observation that owns the fact

| Question | Use |
| --- | --- |
| Should this definition compile? | `clitest.Compile` |
| Why should this definition fail to compile? | `clitest.CompileFails` |
| What does this complete invocation mean before execution? | `clitest.Inspect` |
| What should completion return for partial input? | `clitest.Complete` |
| What static structure did compilation accept? | `clitest.Schema` |
| What reader-facing static reference is exposed? | `clitest.Reference` |
| What Commandment diagnostics are carried by this error? | `clitest.Diagnostics` |

Do not force a fact through a wider boundary simply because that boundary is familiar. A static flag relationship belongs to Schema. Effective availability under invocation conditions belongs to Inspection. Partial editor input belongs to Completion.

## Locate the semantic subject, then state the fact

Locators describe recognizable public coordinates.

```go
command := clitest.Command("approach")
cloakLevel := command.Flag("--cloak-level")
```

Expectations then state what must be true about that subject.

```go
clitest.Inspect(t.Context(), app, request).Require(t,
	clitest.InspectionReady(),
	cloakLevel.Unavailable(),
)
```

Prefer meaningful coordinates such as command paths, flag spellings, argument terms, and public relationships. ID- and order-based locators exist for cases where those are the public coordinates, but a test should not become harder to read than the behavior it proves.

## Make expected failures precise

A test that expects one Commandment rejection should not accept any non-nil error.

```go
diagnostics := clitest.CompileFails(t, badDefinition)
diagnostics.Require(t,
	clitest.Diagnostic(cli.DiagnosticInvalidDeclaration).Present(),
)
```

For fallible analysis operations, use the observation and failure expectations that identify the exact diagnostic under test.

Application-owned handler errors remain ordinary Go errors. Use `errors.Is`, `errors.As`, and application assertions for them rather than trying to force them into the diagnostic frontend.

## Assert repairs as structure, not hint text

The repair scenario exposes a structured replacement for invalid input.

```go
diagnostic := clitest.Diagnostic(cli.DiagnosticInvalidValue)
clitest.Diagnostics(err).Require(t,
	clitest.DiagnosticsExactly(diagnostic),
	diagnostic.RepairCount(1),
	diagnostic.RepairAt(0).Present(),
)
```

When the exact edit matters, descend through the repair source and text-edit locators. When the fact under test is simply that the parser offered one repair, keep the assertion at that level.

This is stronger and more stable than searching the rendered error for `A1`.

See [policies/repairs](../examples/policies/repairs/).

## Reuse one observation for several related assertions

Observation helpers retain the public result rather than rerunning the operation for each `Require` call.

```go
proof := clitest.Inspect(t.Context(), app, request)
proof.Require(t, clitest.InspectionReady())
proof.Require(t, clitest.Command("abduct").Selected())
```

If ordinary Go is clearer for some downstream processing, access the retained detached public artifact through the observation's accessor.

Use the structured expectation surface when it expresses the semantic fact directly. Drop to ordinary Go when you are doing application-specific aggregation over the public data rather than testing another Commandment rule.

## Do not infer absence from incomplete evidence

A candidate found in a partial completion proves presence. A candidate missing from a truncated completion does not prove absence.

The same rule applies to exact collections in partial inspection/provider results.

`clitest` preserves completeness state and can report that an assertion is not provable rather than incorrectly treating missing partial evidence as false.

This is especially important for negative assertions such as "this command does not exist" or "these are exactly the completion items." Prove the required completeness first when exactness matters.

## Use Commandment protocol fixtures only for Commandment-owned protocols

`clitest.ResponseFiles` supplies deterministic data for the public response-file host protocol. `clitest.SignalSource` and related fixtures let tests control Commandment signal-source protocols without starting a real OS signal path.

Those fixtures are not a general mocking framework. Use ordinary Go fakes or real collaborators for application-owned services.

When the process itself is the contract, use a subprocess. When the shell/editor path is the contract, use `shelltest`.

## Share a real semantic contract with `Conformance`

`clitest.Conformance` is useful when the same Commandment behavior must hold across multiple applications or reusable components.

```go
checks := clitest.Conformance{
	clitest.SchemaCase{
		Name: "root command",
		Expectations: []clitest.SchemaExpectation{
			clitest.RootCommand().Name("alienctl"),
		},
	},
	clitest.InspectionCase{
		Name:    "abduct is selectable",
		Request: cli.InspectionRequest{Args: []string{"abduct", "cow"}},
		Expectations: []clitest.InspectionExpectation{
			clitest.InspectionReady(),
			clitest.Command("abduct").Selected(),
		},
	},
}

checks.Require(t, t.Context(), app)
```

Do not use Conformance merely to make a table test look framework-specific. It earns itself when a reusable semantic contract is genuinely shared.

## Keep ordinary Go responsible for ordinary Go facts

Use package `testing` directly for application computation, handler side effects, collaborators, application error identity, panic behavior, subprocess orchestration, races, fuzzing, and other ordinary Go concerns.

Use focused golden or textual assertions when literal presentation itself is the product contract. Use `clitest` when Commandment already exposes the underlying semantic fact structurally.

The examples corpus demonstrates this layering beside complete programs. Start with [examples/values/derived](../examples/values/derived/) for derived semantics, [examples/policies/repairs](../examples/policies/repairs/) for diagnostics, or [examples/signals](../examples/signals/) for signal boundaries.

Continue with [Testing Shell Integrations with shelltest](13-testing-shell-integrations-with-shelltest.md) when semantic completion is correct and the remaining question belongs to a real shell or editor.
