# Troubleshooting

Start from the symptom the user or developer can actually observe. Then choose the public operation that owns the questionable fact.

Do not begin by tracing private compiler code unless the public observations have already shown a product defect that requires implementation work.

## My command is not selected

Use `Inspect` on the exact invocation.

Check the selected command path, command-entry outcome, default relationship, alias or shortcut contribution, dynamic provider output, and any availability or lifecycle state affecting the entry.

If the command is absent from static Schema, the declaration itself is the first place to look. If it exists statically but disappears for one invocation, inspect effective availability, provider materialization, or lifecycle state rather than assuming routing is broken.

The [tooling/inspect](../examples/tooling/inspect/) utility demonstrates this split between Schema and Inspection.

## My flag is rejected or unavailable

First ask whether the spelling is part of the accepted grammar.

- Use Schema for the static spelling and lifecycle.
- Use Inspect for invocation-dependent availability.
- Use Diagnostics for the structured failure when the invocation is rejected.

A condition can make a spelling unavailable even though it exists statically. A release can make a spelling removed even though older releases accepted it. A dialect rule can alter recognition. Those are different causes and should not be collapsed into "flag parsing failed."

See [Conditions, Validation, and Diagnostics](05-conditions-validation-and-diagnostics.md).

## The wrong value won

Use `Inspect` and follow the destination through its candidates, contributors, producers, resolution tier, and provenance.

Questions worth asking include

- Which candidates were produced?
- Which one became effective?
- Which producer kind created it?
- Which input source did that producer read?
- Did the destination use the resolution policy you intended?
- Was a dynamic dependency blocked or failed?

A handler-output assertion can confirm the final value but cannot explain the causal chain.

[Values and Configuration](04-values-and-configuration.md) owns the resolution model.

## A value is missing, but zero looks valid

Do not infer absence from a Go zero value. Inspect the destination's resolution state.

`0`, `false`, and `""` can all be valid resolved values. Commandment tracks whether the destination is resolved, failed, blocked, inactive, or otherwise incomplete separately from the Go representation.

## Conditional syntax is behaving strangely

Inspect the condition inputs first. Then inspect the declaration's availability or visibility outcome.

A hidden declaration can remain valid. An unavailable declaration does not participate in the effective grammar. A rule violation is later than availability and means the semantic value existed but failed blocking policy.

Those distinctions prevent debugging one layer by changing another.

## Validation failed, but I need more than the error text

Call `cli.Diagnostics(err)`.

Inspect the diagnostic code, subject, source span, rejection relationship, repair indexes, and causal error as applicable. If the input language attached a repair, inspect the structured `InputRepair` instead of parsing a hint line.

Use `cli.WriteError` for normal user-facing publication. Use Diagnostics for programmatic meaning.

The [policies/repairs](../examples/policies/repairs/) scenario shows both sides of the same failure.

## Help is wrong

Decide whether the semantic model or the renderer is wrong.

Use Schema or Reference to check the command, spelling, domain, lifecycle, and description data independent of rendered Help. If those public projections are correct but the final Help text is wrong, the Help presentation path is the likely owner.

If the underlying lifecycle or domain is wrong, fix the declaration rather than patching the renderer to lie about it.

[Help and Localization](06-help-and-localization.md) owns the presentation model.

## Completion is semantically wrong

Call `Application.Complete` or `clitest.Complete` with the smallest request that reproduces the problem.

Inspect the completion target, candidates, matching, ranking, native requests, diagnostics, and completeness state.

If the semantic result is already wrong, do not start with shell scripts.

[Completion](07-completion.md) owns semantic completion.

## Semantic completion is right, but Bash or another shell is wrong

Now the shell boundary is the correct subject.

Check generated integration, provider selection, runtime version, sourcing/setup, observed editor state, quoting, and real-shell behavior through `shelltest` or an equivalent supported real-shell proof.

Do not treat a successful `Application.Complete` call as proof that Bash invokes the adapter correctly.

[Testing Shell Integrations with shelltest](13-testing-shell-integrations-with-shelltest.md) owns that proof.

## A dynamic command is missing

Confirm the provider was demanded for the invocation and inspect the resolved dependencies it uses. Then inspect provider output or failure state.

A command that is not in static Schema may still be valid provider-produced structure for one invocation. Conversely, a provider that returns no command for the current dependencies has not necessarily failed.

The [commands/dynamic](../examples/commands/dynamic/) scenario provides a small reference implementation.

## The handler never runs

Inspect execution admission before debugging application code.

Selection may have failed. A required value may be unresolved. A rule may have blocked execution. The invocation may have selected an intrinsic action such as Help or version rather than the application handler.

Only after inspection shows execution was admitted should you debug handler or middleware behavior.

## The handler ran and returned an error

Application errors are ordinary Go errors. Use `errors.Is`, `errors.As`, and application-specific context.

If an explicit exit code was attached with `cli.WithExitCode`, verify the process boundary with a subprocess when exit status itself matters.

Do not expect `cli.Diagnostics` to invent a Commandment diagnostic for an application-owned failure.

See [execution/errors](../examples/execution/errors/).

## A signal never reaches my command

Separate source delivery, signal matching, phase selection, policy action, endpoint delivery, and context cancellation.

A closed context does not prove a signal arrived. A source occurrence does not prove the policy matched it. A matched policy does not imply endpoint emission if the action was cancellation or another effect.

Use deterministic signal fixtures for semantic policy and a real subprocess signal when operating-system delivery is the claim.

[Signals](10-signals.md) owns the full model.

## clitest says a negative fact is not provable

Check whether the retained observation is partial, incomplete, or truncated.

Missing evidence in a partial result is not proof of absence. Request or establish the completeness needed for the negative assertion, or weaken the assertion to the positive fact the artifact can actually prove.

## Build a useful bug report from public evidence

A good report should include

- the smallest public declaration that reproduces the behavior;
- the exact operation and request;
- the relevant public Schema, Inspection, Completion, Reference, or Diagnostic facts;
- literal output only when presentation is part of the problem;
- shell/runtime versions when the shell boundary matters;
- process/platform context when process behavior matters.

Avoid private stack archaeology unless it is needed after a public defect has been isolated. Public evidence makes the report useful even if implementation structure changes.

Return to the [documentation index](README.md) after identifying the owning subsystem.
