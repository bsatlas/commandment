# Help and Localization

Help should stay correct because it is a projection of the CLI you declared, not because a second usage description was kept in sync by hand.

Commandment's built-in Help can project commands, aliases, shortcuts, argument forms, flag spellings, finite value domains, lifecycle state, inherited flags, examples, and intrinsic actions from the accepted application model.

## Start from the Help your declarations already imply

A realistic declaration can produce output such as this without another command registry.

```text
Navigate expedition craft across nearby worlds.

Usage:
  alienctl [flags] <command>

Navigation:
  navigate  Navigate the expedition to a destination.
  travel    Navigate the expedition to a destination.
            [alias]
  home      Navigate directly to Earth.
            [shortcut]

Commands:
  help        Show help.
  completion  Generate shell completion source

Flags:
  --verbose  Show navigation detail.
  --help     Show help.
  -h         Show help.
  --version  Show version.
```

The important property is not the exact formatting. The important property is that the same command entries, spellings, descriptions, lifecycle metadata, and domains used by execution are available to Help.

See [presentation/help](../examples/presentation/help/) for a complete rich application.

## Write descriptions where the semantic subject is declared

Commands, flags, arguments, values, and examples should carry the presentation text that belongs to them instead of relying on a later renderer to rediscover the meaning.

```go
scanCmd := cli.Command[scanConfig]{
	Name: "scan",
	Description: cli.Description{
		Short: "Scan a nearby world.",
	},
}
```

Finite value domains can carry descriptions and lifecycle information for their entries. Aliases and shortcuts retain their relationship to the canonical command. Deprecated spellings retain the lifecycle and replacement information that Help may present.

That is why a Help renderer should not need a second application-specific table of commands and options.

## Groups organize presentation without redefining grammar

Presentation groups are useful when a large root command has several families of commands or flags. Grouping changes how the page is organized. It does not create new command ownership or routing semantics.

Use groups when they make a real Help page easier to scan. Do not create one group per declaration merely to mirror source structure.

## Help can show compatibility relationships instead of hiding them

Lifecycle and replacement metadata can appear directly beside the syntax users need to migrate.

```text
--velocity VALUE [deprecated 2.0.0: velocity was renamed; use --speed],
--velocity=VALUE [deprecated 2.0.0: velocity was renamed; use --speed],
--speed VALUE, --speed=VALUE
  Set approach speed.
```

The Help renderer did not invent the deprecation. It projected lifecycle and replacement declared on the spelling.

[Evolving a CLI](08-evolving-a-cli.md) owns the complete release and migration model.

## Sensitive values should not leak because a renderer forgot the rule

Sensitivity belongs to the semantic value or producer that owns the secret. Commandment-owned Help, diagnostics, inspection, and related presentation can then honor the same redaction policy.

A custom renderer should consume the structured `HelpPage` it receives rather than reaching around it to reconstruct raw sensitive values from unrelated state.

[Values and Configuration](04-values-and-configuration.md) explains destination and producer-specific sensitivity.

## Replace Help rendering without replacing the Help model

Applications that need another textual layout can provide a renderer.

```go
func renderHelp(_ context.Context, page cli.HelpPage) (string, error) {
	text := fmt.Sprintf("%s %s\n%s\n", strings.ToUpper(page.ApplicationName), page.Name, page.Description.Short)
	return text, nil
}

var definition = cli.Definition[struct{}]{
	Name: "alienctl",
	Help: cli.HelpConfig{
		Renderer: renderHelp,
	},
}
```

The renderer receives structured Help data. It does not need to parse argv, walk private declarations, or maintain its own command catalog.

That boundary is what makes custom presentation safe. The application changes formatting while Commandment continues to own semantic preparation.

See [presentation/custom-help](../examples/presentation/custom-help/).

## Localize application text without duplicating the CLI

Application-authored presentation text can be translated through `Localization`.

```go
const scanDescription = "scan nearby worlds"

var definition = cli.Definition[rootConfig]{
	Name: "alienctl",
	Localization: cli.Localization{
		DefaultLocale: "en",
		Bundles: []cli.TranslationBundle{
			{
				Locale: "fr",
				Translations: cli.Translations{
					scanDescription: "analyser les mondes proches",
				},
			},
		},
	},
}
```

The source text remains the authored identity for the application string. The translation bundle changes presentation for another locale without creating another command or description registry.

```text
alienctl help scan
scan nearby worlds
```

```text
alienctl --locale fr help scan
analyser les mondes proches
```

See [presentation/localization](../examples/presentation/localization/).

## Let the invocation select its effective locale explicitly

A value can be marked as the locale selector.

```go
locale := cb.Values.String(&cb.Config.Locale)
cb.Flags.ValueFor(locale).LocaleSelector().Occurrences(cli.Count.Optional()).Long("locale")
```

`LocaleSelector()` tells Commandment that the resolved value selects the operation's effective locale. Help and other Commandment-owned localized presentation can then use that locale consistently.

The effective locale belongs to the operation. Avoid package-global mutable locale state that would make concurrent operations influence each other.

## Application and framework text have different localization owners

Application descriptions and examples are authored by the application. Commandment also owns framework phrases such as built-in Help labels and intrinsic-action wording.

Those ownership domains remain distinct. Application translation bundles should not become an undocumented override table for framework internals, and framework translations should not own application-specific text.

When custom provider output contributes localized presentation, it follows the provider and localization contracts for that operation rather than mutating the static translation catalog.

Exact translation-bundle and framework-translation contracts live in GoDoc.

## Help is not the static tooling API

Help is a presentation surface for humans. A generator that needs structural static data should normally use `Schema` or `Reference` rather than scraping the rendered Help text.

Use `Reference` when the result is reader-oriented static reference data. Use `Schema` when the tool needs structural identities and relationships. [Inspection and Tooling](11-inspection-and-tooling.md) develops that distinction.

Continue with [Completion](07-completion.md) to carry the same accepted command and value model into incomplete editor input and real shell integrations.
