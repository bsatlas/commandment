# Inspection and Tooling

A compiled Commandment application is useful to more than the executable that runs commands. Developer tools, documentation generators, configuration explainers, IDE integrations, and diagnostics can consume public detached observations of the same semantic model.

The key is to choose the operation that owns the question.

## Explain one invocation with `Inspect`

Suppose a support tool needs to answer

```text
Which command did this input select?
Why is range 5?
Which world was parsed?
Did the handler run?
```

`Application.Inspect` analyzes one invocation without running application middleware or handlers.

```go
request := cli.InspectionRequest{
	Args: []string{"scan", "--range", "5", "mars"},
}
inspection, err := app.Inspect(context.Background(), request)
```

Inspection can expose command selection, resolved values, producers and candidates, provenance, dynamic provider/domain outcomes, conditions, rule results, diagnostics, execution admission, and other invocation observations.

See [tooling/inspect](../examples/tooling/inspect/).

## Use Schema to locate stable static subjects before reading invocation state

The inspection utility needs to know which static command and value identities correspond to `scan`, `--range`, and `world`. `Application.Schema()` provides the accepted static structure.

The example uses Schema to locate those identities, then uses Inspection to read the invocation-specific representations.

That pairing demonstrates the boundary well.

```text
Schema
    static accepted identity and relationships

Inspection
    one invocation's effective semantic state
```

Schema is not an invocation trace. Provider-produced batches and dynamic-domain contents that depend on one invocation belong to operation observations rather than to the static program.

Schema is most useful when it enables a real tool such as inspection, code generation, static analysis, or declaration mapping. Reading static structure is usually one step inside that larger task rather than the task by itself.

## Follow a value back to its contributing producer

Inspection retains value, candidate, producer, contributor, and origin relationships.

A low-level tool can use those identities directly.

```go
value := inspection.Values[rangeID]
candidate := inspection.ValueCandidates[value.Contributors[0]]
producer := value.Producers[candidate.ProducerIndex]

fmt.Printf("value: %s\n", value.Representation.Text)
if producer.Kind == cli.ValueProducerKindFlag {
	fmt.Println("source: flag")
}
```

Higher-level tests can express the same question more directly with `clitest`.

```go
rangeValue := clitest.Command("scan").Flag("--range").Value()
clitest.Inspect(t.Context(), app, request).Require(t,
	rangeValue.Text("5"),
	rangeValue.ContributedBy(cli.ValueProducerKindFlag),
)
```

Use direct public artifacts when building a tool. Use `clitest` when the goal is to assert a semantic fact in a test.

## Use `Reference` for reader-oriented static data

A generator that needs a reader-facing command reference should not scrape rendered Help and should not rebuild presentation from private declaration state.

```go
reference, err := app.Reference(cli.ReferenceOptions{})
```

`Reference` is a detached static projection that carries reader-oriented information such as command presentation, descriptions, lifecycle state, and other supported reference relationships.

The [tooling/reference](../examples/tooling/reference/) utility prints static command descriptions from that public projection.

## Schema and Reference are intentionally different

The distinction is practical.

| Need | Use |
| --- | --- |
| Structural identity, dependency, declaration, and static semantic relationships | `Schema` |
| Reader-facing static reference data | `Reference` |
| Effective state for one invocation | `Inspect` |
| Partial editor input | `Complete` |
| Structured meaning of a failure | `Diagnostics` |

A documentation generator may use Reference. A source-mapping or architecture tool may need Schema. A configuration explainer needs Inspection. Using the narrowest operation avoids coupling the tool to details it does not need.

## Use structured Diagnostics instead of parsing error text

A tooling client can call `cli.Diagnostics(err)` and inspect structured failures. The complete model for rules, rejections, diagnostics, and repairs belongs in [Conditions, Validation, and Diagnostics](05-conditions-validation-and-diagnostics.md).

For tooling, the practical rule is simple. Do not parse human error text when a structured diagnostic exists.

## Completion is another observation, but partial input is its own problem

A developer tool may call `Application.Complete` directly. Completion has its own semantics for partial input, matching, ranking, native delegation, and shell adapters. [Completion](07-completion.md) develops those semantics.

Use completion when the input is incomplete. Do not use Inspect and then guess what token the user intended to type next.

## Detached artifacts belong to the caller

Schema, Reference, Inspection, Completion, and diagnostic reports are detached public artifacts. Mutating the returned Go data does not mutate the compiled application.

That makes them safe inputs to downstream tooling without creating a back door into application state.

The detachment rule also means a tool should not expect changes to its local copy to influence later operations. To change the CLI, change the definition and compile a new application.

## Respect completeness when a result is partial

An operation can sometimes return valid partial evidence alongside a failure or truncation condition.

Presence in a partial artifact is evidence. Exact absence requires the observation to be complete enough to prove it.

A support tool should therefore say "this result contains X" when that is all it knows rather than upgrading incomplete evidence into "X is the only possible result."

`clitest` exactness expectations enforce this distinction in tests.

## Build tools from public artifacts rather than implementation internals

The public operations are designed so useful tooling does not need to reach into unexported compiler state.

If a tool requires private implementation data merely to answer a normal consumer question, first ask whether the public Schema, Reference, Inspection, Completion, or Diagnostics artifact already owns the fact. If none can express an important supported consumer question, that is an API/product-design issue rather than a reason for the tool to couple to internals.

Continue with [Testing with clitest](12-testing-with-clitest.md) when the same observations need to become durable semantic assertions.
