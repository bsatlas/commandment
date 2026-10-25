# Deprecation and replacement

`alienctl` is at release `2.0.0`. The old `--velocity` spelling still works, but Help marks it deprecated and directs users to `--speed`.

Build the CLI from this directory:

```sh
go build -o alienctl .
```

Inspect the command Help:

```sh
./alienctl approach --help
```

```text
Usage:
  alienctl approach [flags]

Flags:
  --velocity VALUE [deprecated 2.0.0: velocity was renamed; use --speed],
  --velocity=VALUE [deprecated 2.0.0: velocity was renamed; use --speed],
  --speed VALUE, --speed=VALUE
    Set approach speed.

Global flags:
  --help  Show help.
  -h      Show help.
```

See [Evolving a CLI](../../../docs/08-evolving-a-cli.md) for release schemes, lifecycle boundaries, replacement, Help projection, and version presentation.
