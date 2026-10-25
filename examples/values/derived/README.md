# Derived configuration

`alienctl abduct` computes its final summary from already-resolved input. The handler receives the finished configuration instead of rebuilding the dependency in application code.

Build the CLI from this directory:

```sh
go build -o alienctl .
```

Run an abduction:

```sh
./alienctl abduct --beam-strength 4 mars
```

```text
abducting mars with beam strength 4
```

See [Values and configuration](../../../docs/04-values-and-configuration.md) for dependencies, derivation, precedence, and provenance.
