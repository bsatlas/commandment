# Custom value types

`alienctl relay` accepts an IP address directly into a `netip.Addr` field. The application teaches Commandment that type once with `ValueCodec[T]`.

Build the CLI from this directory:

```sh
go build -o alienctl .
```

Relay through an address:

```sh
./alienctl relay --address 203.0.113.7
```

```text
relaying through 203.0.113.7
```

See [Values and configuration](../../../docs/04-values-and-configuration.md) for value languages, resolution, domains, derivation, provenance, and sensitivity.
