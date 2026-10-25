# Commands, Flags, and Arguments

A Commandment grammar starts from a user question rather than a parser object.

> What should a person be able to type, and what should each token mean?

Commands, flags, positional arguments, shortcuts, defaults, profiles, and recognition policy describe that language before a handler runs.

## Build the command path users think in

Nested commands should reflect the nouns and actions users naturally navigate.

```text
alienctl fleet status
alienctl fleet launch voyager
```

The declaration mirrors that structure.

```go
fleetCmd := cli.Command[struct{}]{
	Name: "fleet",
	Declare: func(cb *cli.CommandBuilder[struct{}]) {
		statusCmd := cli.Command[struct{}]{Name: "status"}
		cb.Subcommands.Command(statusCmd)

		launchCmd := cli.Command[launchConfig]{
			Name: "launch",
			Declare: func(cb *cli.CommandBuilder[launchConfig]) {
				cb.Arguments.String(&cb.Config.Ship, "ship")
			},
		}
		cb.Subcommands.Command(launchCmd)
	},
}
```

See [commands/nested](../examples/commands/nested/) for the complete application.

Each command should own the configuration its handler needs. Do not flatten descendant arguments and flags into a root config simply because they share one executable.

## Give an existing command another name with an alias

An alias is another immediate lexical entry for the same command.

```text
alienctl scan mars
alienctl survey mars
```

Use an alias when the alternative is genuinely another name for the same command rather than a macro that contributes extra input.

```go
cmd := cli.Command[scanConfig]{
	Name:    "scan",
	Aliases: cli.Aliases("survey"),
}
```

The canonical command keeps one semantic identity even though the namespace recognizes another spelling. Help, Reference, completion, lifecycle, and tests can distinguish the canonical command from the alias entry when that matters.

See [commands/aliases](../examples/commands/aliases/).

## Use a shortcut when the alternate entry contributes bounded input

A shortcut can select a canonical command and contribute arguments or other supported bounded input.

```text
alienctl home
```

can mean

```text
alienctl navigate earth
```

without creating a second navigation handler.

```go
shortcut := cli.CommandShortcut{
	Name: "home",
	Args: []string{"earth"},
}
cb.Shortcut(shortcut)
```

An alias changes how a command is named. A shortcut changes the entry path and can contribute input. Keeping those concepts distinct lets diagnostics, provenance, Help, and compatibility relationships describe what actually happened.

See [commands/shortcuts](../examples/commands/shortcuts/).

## Let a parent select a default child when no token is consumed

A default child is useful when one child is the natural behavior of the parent itself.

```text
alienctl fleet
```

can select the same status command as

```text
alienctl fleet status
```

Declare the relationship at the parent.

```go
statusCmd := cli.Command[struct{}]{Name: "status"}
cb.Subcommands.Command(statusCmd).Default()
```

A default relationship is not an alias and it does not consume a hidden command token. That distinction matters for selection, provenance, Help, and completion.

See [commands/default](../examples/commands/default/).

## Model positional input in the order users read it

The simplest positional declaration binds one argument to one typed destination.

```go
cb.Arguments.String(&cb.Config.World, "world")
```

Optional arguments should be optional because the CLI naturally has a useful meaning without them, not because the application handler can tolerate an empty string.

```go
cb.Arguments.String(&cb.Config.Ship, "ship").Optional()
```

See [arguments/optional](../examples/arguments/optional/).

### Use forms for genuinely different positional languages

Some commands have several legitimate positional shapes.

```text
alienctl navigate europa
alienctl navigate 12 4 9
```

Declare each accepted form instead of parsing the tail manually in the handler.

```go
cb.Arguments.Form(func(form *cli.ArgumentFormBuilder[navigateConfig]) {
	form.String(&form.Config.World, "world")
})
cb.Arguments.Form(func(form *cli.ArgumentFormBuilder[navigateConfig]) {
	form.Int(&form.Config.X, "x")
	form.Int(&form.Config.Y, "y")
	form.Int(&form.Config.Z, "z")
})
```

Argument-form callbacks stay lexical because the form itself is the declaration being taught. See [arguments/forms](../examples/arguments/forms/).

### Use pass-through when the remaining tokens belong to another program

A relay command may own a rover name while intentionally passing the rest of the command line to that rover.

```text
alienctl relay rover-7 -- reboot --force
```

The pass-through destination makes that ownership explicit instead of teaching the handler to guess where Commandment parsing should stop.

See [arguments/pass-through](../examples/arguments/pass-through/).

## A flag owns a value while its spellings own syntax

One flag can accept more than one spelling without creating more than one destination.

```go
scanRange := cb.Flags.Int(&cb.Config.Range)
scanRange.Long("range")
scanRange.Short('r')
```

Both spellings contribute to the same typed value.

```text
alienctl scan --range 5 mars
alienctl scan -r 5 mars
```

That separation becomes especially useful when one spelling is deprecated, sensitive, or replaced while another remains current.

See [flags/spellings](../examples/flags/spellings/).

### Counting flags model repeated occurrences directly

For conventional verbosity-style input, declare a count rather than parse repeated booleans yourself.

```go
cb.Flags.Count(&cb.Config.Verbosity).Short('v')
```

```text
alienctl scan -v -v mars
```

The current example intentionally teaches repeated `-v` occurrences. Do not infer clustered `-vvv` support unless the configured dialect accepts it. See [flags/counting](../examples/flags/counting/).

## Put conditional syntax in the grammar instead of rejecting it later

If `--cloak-level` only exists in stealth mode, make availability part of the declaration.

```go
stealth := cb.Flags.Bool(&cb.Config.Stealth).Occurrences(cli.Count.Optional())
stealth.Long("stealth")

cloakLevel := cb.Flags.Int(&cb.Config.CloakLevel)
cloakLevel.AvailableWhen(cb.When.Equal(stealth, true)).Long("cloak-level")
```

Without `--stealth`, the spelling is unavailable. The handler never needs a second `if !stealth && cloakLevel != 0` policy.

[Conditions, Validation, and Diagnostics](05-conditions-validation-and-diagnostics.md) owns the complete availability and validation model. The runnable case is [policies/conditions](../examples/policies/conditions/).

## Use a grammar profile when one command truly has alternative languages

A grammar profile represents mutually exclusive command grammars rather than optional fragments that can be mixed arbitrarily.

Suppose `alienctl transmit` has a normal brief mode and an explicitly selected verbose grammar that adds `--verbose`.

```go
type transmitConfig struct {
	Profile string
	Verbose bool
}

declareProfiles := func(values *cli.ValueDomainBuilder[string]) {
	values.Value("brief", "brief")
	values.Value("verbose", "verbose")
}

profiles := cli.Value.Domain(cli.Value.Assign[string](), declareProfiles)
profile := cb.Values.Value(&cb.Config.Profile, profiles)
cb.Values.Default(profile, "brief")
cb.Flags.ValueFor(profile).Long("profile")

cb.Profiles.Switch(profile, func(cases *cli.GrammarProfileCases[transmitConfig, string]) {
	cases.Case("brief", func(*cli.CommandBuilder[transmitConfig]) {})
	cases.Case("verbose", func(cb *cli.CommandBuilder[transmitConfig]) {
		cb.Flags.Bool(&cb.Config.Verbose).
			Occurrences(cli.Count.Optional()).
			Long("verbose")
	})
})
```

The resulting language has one selected grammar at a time.

```text
alienctl transmit --profile=brief europa
alienctl transmit --profile=verbose --verbose europa
```

Use a profile when the selector chooses one coherent syntax family such as a legacy compatibility form versus a new form, or two established invocation styles that should not combine.

Do not reach for profiles merely because several arguments are optional. Ordinary arguments, conditions, defaults, and subcommands usually express those cases more clearly. Use separate commands when the alternatives are actually different operations.

The selected profile remains part of Schema, Reference, inspection, completion, and testing. The profile declaration stays lexical because the profile cases themselves are the grammar being taught. GoDoc owns the exact profile methods and lifecycle contracts.

## Change dialect only to express a real compatibility or language requirement

Commandment exposes recognition policy for cases where a CLI deliberately differs from the default language. Dialect controls include flag ordering, short-flag clustering, flag and command matching, normalization, dashless flags, end-of-options behavior, response-file recognition, and default flag placement.

A dialect change should answer a concrete user-language requirement.

For example, repeated short flags do not cluster by default. If an established CLI accepts `-vvv`, enable that recognition rule deliberately.

```go
cb.Flags.Count(&cb.Config.Verbosity).Short('v')
cb.Dialect.ShortFlagClustering(cli.ShortFlagClusteringEnabled)
```

The command can then accept the clustered spelling.

```text
alienctl scan -vvv mars
```

That is different from merely accepting three separate `-v` occurrences. The dialect changed what one token means.

An existing mature CLI may instead require unique-prefix command matching, a particular flag-placement rule, dashless short flags, or another recognition contract. Configure only the rule the product actually owns and prove the observable grammar after changing it.

Keep the default dialect when the application has no compatibility or language reason to diverge. Dialect is not a style customization layer. It changes recognition semantics.

[Adopting Commandment](17-adopting-commandment.md) shows how to inventory these contracts before migration. GoDoc owns the complete set of recognition options.

## Use a token language when one semantic value genuinely owns several tokens

`ValueCodec[T]` parses one scalar text value. A `TokenLanguage[T]` owns a bounded sequence of tokens that together represent one semantic value.

For example, a position can own exactly two positional tokens while the handler receives one typed value.

```go
type position struct {
	X int
	Y int
}

type locateConfig struct {
	Position position
}

parsePosition := func(_ context.Context, tokens []string) (position, error) {
	x, err := strconv.Atoi(tokens[0])
	if err != nil {
		return position{}, err
	}
	y, err := strconv.Atoi(tokens[1])
	if err != nil {
		return position{}, err
	}
	return position{X: x, Y: y}, nil
}

formatPosition := func(value position) []string {
	tokens := []string{
		strconv.Itoa(value.X),
		strconv.Itoa(value.Y),
	}
	return tokens
}

language := cli.TokenLanguage[position]{
	Name:   "position",
	Count:  cli.Count.Exactly(2),
	Parse:  parsePosition,
	Format: formatPosition,
}

destination := cb.Values.Destination(&cb.Config.Position, cli.Value.Assign[position]())
cb.Arguments.Tokens(destination, "position", language)
```

The user types two tokens.

```text
alienctl locate 12 4
```

The handler receives one `position` value with `X` set to `12` and `Y` set to `4`.

Use a token language when the token sequence itself is one semantic value. Do not hide subcommand routing or an arbitrary mini-parser inside it. Ordinary arguments, forms, and flags remain clearer when the tokens are independent pieces of CLI grammar.

Token languages can also reject a particular token range and attach structured repairs. [policies/repairs](../examples/policies/repairs/) demonstrates that path. GoDoc owns exact count, parse, format, and snapshot contracts.

## Generate commands from runtime data only when the namespace is actually dynamic

A provider is appropriate when the available command namespace depends on invocation-time input or host data.

The [commands/dynamic](../examples/commands/dynamic/) scenario makes `alienctl outpost` choose its available child command from `--system`.

```go
provideOutposts := func(_ context.Context, batch *cli.CommandBatch[outpostConfig]) error {
	switch batch.Config.System {
	case "jovian":
		batch.Command(cli.Command[struct{}]{Name: "europa"})
	case "saturnian":
		batch.Command(cli.Command[struct{}]{Name: "titan"})
	}
	return nil
}
cb.Subcommands.CommandsFrom(provideOutposts, system)
```

Use static lexical declaration for static topology. A command provider should not become an indirection layer around commands known at compile time.

## Keep response-file recognition separate from loading

Recognition of `@file` syntax belongs to grammar.

```go
cb.Dialect.ResponseFiles(true)
```

How that source is loaded belongs to the host boundary. A normal process host can read the filesystem. An embedded host can provide another `ResponseFileAccess` implementation.

[Runtime and Host Integration](09-runtime-and-host-integration.md) owns loading and host capability wiring. [response-files](../examples/response-files/) contains both filesystem and custom-loader scenarios.

## Verify grammar at the boundary that owns the question

Use `Schema` or `clitest.Schema` for static accepted grammar. Use `Inspect` when conditions, dynamic providers, lifecycle, or resolved input affect the effective invocation. Use `Complete` for partial editor input. Use process tests only when process behavior itself is the contract.

Continue with [Values and Configuration](04-values-and-configuration.md) when the language is established and you need to decide how typed configuration gets its final value.
