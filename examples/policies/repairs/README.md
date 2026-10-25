# Structured repairs

`alienctl dock` rejects an invalid dock code as normal CLI input while keeping a machine-readable replacement for tools and tests.

Build the CLI from this directory:

```sh
go build -o alienctl .
```

Try an invalid code:

```sh
./alienctl dock --code ZZ
```

```text
alienctl: invalid dock code
  input: ZZ
```

Use the accepted code:

```sh
./alienctl dock --code A1
```

```text
docking at A1
```

The parser attaches `A1` as a structured repair rather than applying it silently. See [Conditions, validation, and diagnostics](../../../docs/05-conditions-validation-and-diagnostics.md).
