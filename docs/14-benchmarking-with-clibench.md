# Benchmarking with clibench

`clibench` is a small benchmark frontend for successful Commandment public operations. It removes repeated benchmark-loop plumbing without deciding what your application should measure.

A normal inspection benchmark keeps the compiled application and request outside the measured loop.

```go
func BenchmarkInspect(b *testing.B) {
	request := cli.InspectionRequest{
		Args: []string{"abduct", "--beam-strength", "3", "cow"},
	}
	clibench.Inspect(b, app, request)
}
```

When compilation itself is the subject, benchmark the definition instead. `clibench.Compile` recompiles it on every benchmark iteration.

```go
func BenchmarkCompile(b *testing.B) {
	clibench.Compile(b, definition)
}
```

## Choose the operation before the helper

Start with the performance question.

```text
How expensive is a successful inspection of this invocation?
How expensive is completion for this partial input?
How expensive is this supported run path?
```

Then choose the corresponding `clibench` helper.

Do not benchmark whichever helper is easiest to call and then infer a different product property from the number.

## Keep the application and request representative

Compile the application outside the measured loop when compilation is not part of the question. Construct stable request data outside the loop when request allocation is not part of the question.

A representative benchmark should exercise the same public operation and semantic workload whose cost you care about.

For example, a completion benchmark should use a partial input shape representative of the application rather than a trivial root completion if the product concern is a deep dynamic namespace.

## Let clibench keep the loop narrow

The helper should repeatedly invoke the selected successful public operation and report failures through `testing.B` rather than requiring every benchmark to duplicate the same loop and error handling.

Use the helper when that is the complete benchmark lifecycle you want.

Use `testing.B` directly when setup, teardown, mutation, concurrency, process behavior, shell integration, or another lifecycle element is itself part of the measurement.

## Interpret the result narrowly

An operation benchmark measures the workload you supplied under the Go benchmark environment that ran it.

It does not automatically establish a public latency SLO, end-to-end process startup time, shell completion latency, memory ceiling, or performance guarantee for unrelated application shapes.

Record the exact operation and relevant input shape when the result matters to an engineering decision.

[Adopting Commandment](17-adopting-commandment.md) explains why performance should be preserved during migration only when the existing CLI actually owns a concrete performance contract worth preserving.

Use GoDoc for the exact helper signatures and supported operation request types.
