# Reference

`alienctl-reference` prints Commandment's reader-facing command reference.

Build the utility from this directory:

```sh
go build -o alienctl-reference .
```

Print the command reference:

```sh
./alienctl-reference
```

```text
scan       scan nearby worlds
status     show fleet status
completion Generate shell completion source
```

See [Inspection and tooling](../../../docs/11-inspection-and-tooling.md) for the difference between Schema, Reference, Inspection, and Completion.
