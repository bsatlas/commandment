# Examples

Commandment examples are grouped by the CLI problem they solve. Each scenario is an independent fictional application or developer utility unless its README says otherwise.

If you are building your first Commandment application, start with [Getting Started](../docs/01-getting-started.md). For a focused example, choose the collection that matches what you are trying to build.

- [Arguments](arguments/) - positional input, alternate forms, optional terms, and pass-through command lines.
- [Commands](commands/) - nesting, aliases, shortcuts, defaults, and runtime-generated command namespaces.
- [Flags](flags/) - alternate spellings and repeated counting flags.
- [Values](values/) - typed values, finite choices, dynamic choices, application-defined types, and derived configuration.
- [Host context](host-context/) - trusted facts and overridable invocation defaults supplied by the host.
- [Policies](policies/) - conditional syntax, blocking or advisory policy, notices, and structured repairs.
- [Completion](completion/) - semantic candidates, native delegation, and real-shell integration.
- [Execution](execution/) - middleware and application-owned error behavior.
- [Presentation](presentation/) - built-in Help, custom Help, and localization.
- [Evolution](evolution/) - version presentation, deprecation, and replacement across releases.
- [Tooling](tooling/) - explain invocations and build static reader-facing reference data.
- [Response files](response-files/) - filesystem and host-provided response sources.
- [Signals](signals/) - process interrupts and application-defined typed signal sources.

Every scenario README gives commands relative to that scenario directory. The complete program and focused test are beside it in `main.go` and `main_test.go`.
