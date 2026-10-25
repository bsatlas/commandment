# Values

Model application configuration as typed values instead of parser-specific fields.

- [Value choices](choices/) - Restrict a value to a finite set of meaningful choices.
- [Dynamic choices](dynamic-choices/) - Choose which values are valid from another resolved value.
- [Custom value types](custom-types/) - Parse an application-defined Go type through `ValueCodec[T]`.
- [Derived configuration](derived/) - Compute one configuration value from already-resolved inputs.
