# Conditions, Validation, and Diagnostics

A well-designed CLI should make invalid states difficult to express and should reject the remaining invalid input before application code has to reason about it.

Commandment keeps several responsibilities separate so one declaration can drive grammar, validation, human error output, structured diagnostics, repairs, and tests without collapsing them into one callback.

## Make syntax exist only when it makes sense

Suppose `--cloak-level` is meaningful only when stealth mode is active.

```go
stealth := cb.Flags.Bool(&cb.Config.Stealth).Occurrences(cli.Count.Optional())
stealth.Long("stealth")

cloakLevel := cb.Flags.Int(&cb.Config.CloakLevel)
cloakLevel.AvailableWhen(cb.When.Equal(stealth, true)).Long("cloak-level")
```

With stealth enabled, the spelling participates in the effective CLI.

```text
alienctl approach --stealth --cloak-level 3 earth
```

Without stealth, the flag is unavailable rather than silently accepted and rejected later by the handler.

```text
alienctl approach --cloak-level 3 earth
alienctl: flag --cloak-level is unavailable
  suggestion: --cloak-level
```

Availability is structural participation. Help, completion, inspection, diagnostics, and execution all see the same decision.

See [policies/conditions](../examples/policies/conditions/).

## Keep availability and visibility separate

A declaration may remain valid while being omitted from ordinary presentation. That is a visibility decision, not an availability decision.

Use availability when the declaration should not participate in the effective CLI. Use visibility when the declaration remains accepted but should not ordinarily be shown.

Keeping those states separate prevents a hidden compatibility spelling from accidentally becoming invalid or an unavailable option from appearing merely because a renderer forgot the condition.

## Put blocking semantic validation in rules

Rules evaluate settled semantic values before the selected handler runs.

```go
speed := cb.Flags.Int(&cb.Config.Speed)
speed.Long("speed")
cb.Rule.Require(speed)
cb.Rule.Between(speed, 1, 5)
```

The user gets CLI-specific failure at the CLI boundary.

```text
alienctl approach earth
alienctl: --speed is required
```

```text
alienctl approach --speed 9 earth
alienctl: --speed must be between 1 and 5
```

The handler can assume admitted configuration satisfied the declared rule instead of rebuilding the same validation in application code.

Use ordinary handler validation for business rules that belong to the application rather than the CLI model. A database constraint, remote authorization decision, or domain transaction failure should remain ordinary application behavior unless the CLI itself owns the rule.

## Let `WriteError` own normal Commandment error presentation

A process entry point normally reports a Commandment failure with `cli.WriteError` and exits with the error's policy.

```go
if err := app.RunFromProcess(context.Background()); err != nil {
	_ = cli.WriteError(os.Stderr, err)
	os.Exit(cli.ExitCode(err))
}
```

Do not parse or reconstruct Commandment's rendered error in application code. The error retains structured meaning separately from its terminal presentation.

Application-owned handler errors remain ordinary Go errors and may carry an explicit exit code through `cli.WithExitCode`. They do not become Commandment diagnostics merely because they occurred during a CLI invocation.

[execution/errors](../examples/execution/errors/) shows that boundary.

## Read structured diagnostics when a program needs more than the message

When Commandment owns a failure, `cli.Diagnostics(err)` returns a detached diagnostic report.

```go
request := cli.RunRequest{
	Args: []string{"approach", "--speed", "9", "earth"},
}
err := app.Run(context.Background(), request)
report := cli.Diagnostics(err)
```

A diagnostic can expose stable machine-readable meaning such as its code, semantic subject, producer relationship, source span, input origin, rejection details, repair references, and causal error information. Tools should use those structured fields instead of scraping `Message`, `Reason`, or `Hint` text.

If an application handler returns an ordinary error with no Commandment-owned diagnostic, the diagnostic report is empty. That is expected and preserves the ownership boundary.

## A parser can reject input and attach a repair without silently changing the input

Token languages can return a structured rejection when parsing fails.

```go
func parseDockCode(_ context.Context, tokens []string) (string, error) {
	code := tokens[0]
	if code == "A1" {
		return code, nil
	}
	rejection := cli.Rejection.Message(
		"dock-code",
		"invalid dock code",
		cli.Rejection.At(0, 0, len(code)),
		cli.Rejection.ReplaceWith("A1"),
	)
	return "", rejection
}
```

The ordinary user still gets a normal CLI error.

```text
alienctl dock --code ZZ
alienctl: invalid dock code
  input: ZZ
```

The same failure also carries an `InputRepair` describing the replacement. Commandment does not apply the repair automatically. The repair is detached structured data that an editor, interactive frontend, diagnostic tool, or test may choose to present or apply under its own policy.

This separation is important. A suggestion is not a hidden mutation of the user's invocation.

See [policies/repairs](../examples/policies/repairs/).

## Rejections belong to the input language that discovered them

A `Rejection` can describe parser-local identity and exact mapped spans. This is useful when a value language or token language knows more about its invalid input than a generic rule could know.

Use a rule when already-parsed semantic values violate a CLI rule. Use a rejection when the input language itself cannot produce the semantic value and can identify why.

Both can become Commandment diagnostics, but they answer different questions.

## Advisories describe nonfatal semantic outcomes

Not every semantic concern should block execution.

`alienctl approach --reckless earth` can warn while still running.

```go
advisory := cli.Advisory{
	Code:     "reckless-approach",
	Severity: cli.NoticeSeverityWarning,
	Message:  "reckless approach enabled",
}
cb.Advisory.When(cb.When.Equal(reckless, true), advisory).Subject(reckless)
```

The advisory is semantic truth. Notice delivery is how that truth reaches a caller or user.

See [policies/advisories](../examples/policies/advisories/).

## Suppress notice delivery without erasing the advisory

A quiet mode can suppress a matching notice.

```go
cb.Notice.SuppressWhen(cb.When.Equal(quiet, true), cli.NoticeMatch{
	AdvisoryCode: advisory.Code,
})
```

The advisory outcome still exists in inspection. The notice is simply suppressed from delivery.

That distinction matters for tooling, testing, and auditability. "The advisory did not exist" and "the advisory existed but its notice was suppressed" are different facts.

Notice delivery is host-owned. A host that uses `Application.Run` supplies `RunRequest.Notify` when it wants Notices delivered. `RunFromProcess` does not automatically publish Notices. Suppression can prevent a matching Notice from reaching `Notify` without erasing the underlying advisory semantic result. Tests can inspect advisory and suppression state without depending on terminal text.

## Execution admission is later than parsing and validation

Several states near execution are intentionally distinct.

A command can be reachable in the command graph yet unavailable for the current invocation. A declaration can be available but hidden. An invocation can select a command but fail blocking rules. An application action can be invocable yet not admitted to the selected handler because analysis found an error.

You do not need to memorize every state before building a CLI. The practical rule is to ask the public observation that owns the question rather than infer one state from another.

- Use `Schema` for static declared structure.
- Use `Inspect` for effective invocation state.
- Use `Diagnostics` for structured failure meaning.
- Use `Reference` or Help for presentation state.
- Use `Run` when application execution itself is the subject.

[Inspection and Tooling](11-inspection-and-tooling.md) shows how those observations fit together.

## Test the same semantic fact that the user experiences

A good semantic test should not duplicate the rule implementation.

```go
cloak := clitest.Command("approach").Flag("--cloak-level")
clitest.Inspect(t.Context(), app, ordinaryRequest).Require(t,
	clitest.InspectionReady(),
	cloak.Unavailable(),
)
```

For a repair, test the structured diagnostic rather than only a rendered hint.

```go
diagnostic := clitest.Diagnostic(cli.DiagnosticInvalidValue)
clitest.Diagnostics(err).Require(t,
	clitest.DiagnosticsExactly(diagnostic),
	diagnostic.RepairCount(1),
	diagnostic.RepairAt(0).Present(),
)
```

Use focused golden assertions when literal presentation is itself the contract under test. Do not make every semantic test depend on exact terminal wording.

Continue with [Help and Localization](06-help-and-localization.md) for presentation derived from the same declarations, or [Troubleshooting](15-troubleshooting.md) when you need to diagnose a CLI that is not behaving as expected.
