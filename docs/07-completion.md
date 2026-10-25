# Completion

Completion is not a separate command registry. `Application.Complete` analyzes incomplete input using the same commands, grammar, values, conditions, domains, providers, lifecycle, and presentation metadata that define normal invocation behavior.

The useful mental model is one path from partial input to the client that will present the result.

```text
partial input
    -> Commandment semantic completion
    -> candidates and/or native delegation
    -> finalized completion result
    -> optional shell provider adapter
    -> shell/editor
```

## Complete the semantic value the user is editing

A value can contribute application-specific candidates.

```go
func completeWorld(_ context.Context, _ *visitConfig, request cli.ValueCompletionRequest) (cli.CompletionContribution, error) {
	worlds := []string{"europa", "ganymede", "io"}
	candidates := make([]cli.CompletionCandidate, 0, len(worlds))
	for _, world := range worlds {
		if strings.HasPrefix(world, request.Text) {
			candidate := cli.CompletionCandidate{Value: world}
			candidates = append(candidates, candidate)
		}
	}
	contribution := cli.CompletionContribution{Candidates: candidates}
	return contribution, nil
}
```

Attach the behavior to the destination that owns the value.

```go
world := cb.Arguments.String(&cb.Config.World, "world")
world.CompleteBy(completeWorld)
```

A programmatic client can then ask the application directly.

```go
request := cli.CompletionRequest{
	Args: []string{"visit", "eu"},
	Cursor: cli.CompletionCursor{
		Argument: 1,
		Byte:     2,
	},
}
completion, err := app.Complete(context.Background(), request)
```

The [completion/values](../examples/completion/values/) utility prints `europa` for that request.

## Candidates are proposals before they become client-facing items

A completion callback contributes candidates. Commandment then applies the completion model, including matching, ranking, lifecycle, availability, grouping, and request limits as appropriate, to produce the final result.

Do not assume a candidate is automatically an emitted item. The candidate is a semantic proposal that still belongs to Commandment's completion analysis.

When testing a completion algorithm, distinguish the application contribution from the finalized client-facing result if the difference matters to the behavior you are proving.

## Let a direct client choose matching and ranking when it owns that experience

A programmatic completion client may want a different textual search experience from a shell adapter. The client describes that preference in `CompletionRequest.Client` rather than post-filtering an unrelated result.

For example, an editor that wants case-insensitive substring search over insertion and display text, then lexical ordering, can ask for it directly.

```go
request := cli.CompletionRequest{
	Args: []string{"visit", "eu"},
	Cursor: cli.CompletionCursor{
		Argument: 1,
		Byte:     2,
	},
	Client: cli.CompletionClient{
		Name: "alien-editor",
		Match: cli.CompletionMatchOptions{
			Mode: cli.CompletionMatchModeSubstring,
			Case: cli.CompletionCaseInsensitive,
			Fields: []cli.CompletionMatchField{
				cli.CompletionMatchFieldValue,
				cli.CompletionMatchFieldDisplayText,
			},
		},
		Ranking: cli.CompletionRankingModeLexical,
		Capabilities: cli.CompletionClientCapabilities{
			Descriptions: true,
			Groups:       true,
			MatchRanges:  true,
		},
	},
}
completion, err := app.Complete(context.Background(), request)
```

The matching modes cover exact, prefix, substring, fuzzy, approximate, and unfiltered client behavior. Ranking can use Commandment relevance, preserve semantic source order, use lexical ordering, or request no additional ranking transform. GoDoc owns the exact enum and field contracts.

Use client-selected matching when the completion consumer actually owns a search interaction. Do not change the application's grammar merely to make one editor fuzzy-match command names.

Client matching runs after semantic candidate production. A completion callback that prefilters its own contribution can narrow the candidate set before client matching, so it should not discard candidates that the intended client policy still needs.

Shell adapters intentionally use a smaller matching profile described later in this guide. Direct `Application.Complete` is the expert boundary for arbitrary completion-client matching policy.

## Dynamic domains and providers can participate in completion

A dynamic value domain may need already-resolved input before it knows the legal values. A command provider may need invocation state before it knows the available child namespace.

Completion can materialize those same operation-scoped contributions when the partial input demands them. The completion path should not maintain a second set of dynamic choices merely because the invocation is incomplete.

This is one reason explicit dependencies matter. Commandment can settle the required values before asking the dependent domain or provider for completion-relevant structure.

## Delegate filesystem work to clients that already own the filesystem experience

File completion is often better handled by a shell or editor that already knows local filesystem semantics and presentation.

```go
manifest := cb.Flags.String(&cb.Config.Manifest)
manifest.CompleteWith(cli.CompletionNativeFiles).Long("manifest")
```

A client advertises native capabilities in its completion request. When native file completion is available, Commandment can return a native request rather than generating a fake filesystem candidate list itself.

The [completion/native-files](../examples/completion/native-files/) scenario demonstrates that delegation.

Delegation is still part of the semantic result. Tools can observe that Commandment requested file completion without pretending Commandment enumerated the files.

## Use exact argv when you have exact argv

A direct programmatic caller that owns execution-equivalent arguments should use `CompletionRequest.Args` and its cursor coordinates.

Shell/editor adapters do not always observe execution-equivalent argv. Quoting and editing state can leave the adapter with shell words whose exact decoded argv topology is uncertain. The completion API also provides the observed-input lane for those clients.

Do not manufacture exact argv merely to fit the simpler request shape. Preserve the fidelity the client actually has and let the completion model reason in the correct coordinate space.

Exact observed-input fields and cursor contracts belong in GoDoc because they are adapter-level reference details.

## Every normal application can expose built-in shell integration

The zero `ShellCompletion` configuration enables the built-in provider catalog.

```text
Bash
Zsh
Fish
PowerShell
Elvish
Nushell
```

A normal application gets the intrinsic `completion` command.

```sh
go build -o alienctl .
export PATH="$PWD:$PATH"
source <(alienctl completion bash)
```

At a Bash prompt, the generated integration adapts Bash/editor state to the same `Application.Complete` operation.

```text
alienctl sc<TAB>
```

becomes

```text
alienctl scan
```

See [completion/shell](../examples/completion/shell/) for the real Bash proof.

### Built-in compatibility floors

The built-in generators support these shell runtime floors.

| Provider | Supported runtime |
| --- | --- |
| Bash | 3.2 or newer |
| Zsh | 4.0.1 or newer |
| Fish | 3.x or newer |
| PowerShell | Windows PowerShell 5.1 or PowerShell 7 or newer |
| Elvish | 0.11 or newer |
| Nushell | 0.72 or newer |

Elvish and Nushell have version-sensitive generated forms for older supported bands. Use `ShellCompletionScriptForVersion` when you deliberately need one of those target-specific artifacts. The other built-in providers use canonical generation and reject a nonempty generation target rather than pretending the target changed their source.

A support floor is not evidence that every shell release behaves identically. Commandment keeps shell-specific compatibility behavior inside the provider while preserving one semantic completion model.

## Generation target and running shell version are different facts

A shell provider can generate integration source successfully while an installed shell integration still fails because of quoting, editor protocol, shell version, PATH, sourcing, or runtime behavior.

Keep those proof boundaries separate.

- `Application.Complete` proves semantic partial-input behavior.
- `ShellCompletionScript` proves canonical generation for a provider.
- `ShellCompletionScriptForVersion` proves target-specific generation when that provider supports it.
- `shelltest` proves the real shell/editor boundary.

`ShellCompletionGeneration.TargetVersion` is an artifact syntax target. It is never a discovered runtime shell version. A generated adapter observes its actual shell family and version when it runs and reports those facts in the adapter request.

## Add a custom shell provider without creating another completion engine

A Definition with a nil `ShellCompletion.Shells` slice uses the canonical built-in catalog in this order.

```text
Bash
Zsh
Fish
PowerShell
Elvish
Nushell
```

A nonnil, nonempty `Shells` slice replaces that catalog completely. Commandment does not append the defaults. A nonnil empty slice is invalid. `Disabled` removes the intrinsic completion command and requires `Shells` to remain nil.

Provider names are exact opaque catalog identities and must be unique. A custom provider does not override a same-named built-in that is still present. To replace a built-in provider, omit the built-in from the explicit catalog and supply the custom provider under that exact built-in name. Nil providers and duplicate names make the catalog invalid. `ShellCompletion` and `Shell.Custom` GoDoc own the exact custom-name validity rules.

To extend the default catalog with Xonsh, list the six built-ins together with the custom provider.

```go
func generateXonsh(generation cli.ShellCompletionGeneration) (string, error) {
	return renderXonshAdapter(generation)
}

var xonsh = cli.Shell.Custom("xonsh", generateXonsh)

var definition = cli.Definition[struct{}]{
	Name: "alienctl",
	ShellCompletion: cli.ShellCompletion{
		Shells: []cli.ShellCompletionProvider{
			cli.Shell.Bash,
			cli.Shell.Zsh,
			cli.Shell.Fish,
			cli.Shell.PowerShell,
			cli.Shell.Elvish,
			cli.Shell.Nushell,
			xonsh,
		},
	},
}
```

The intrinsic `alienctl completion <shell>` command uses exactly the Definition-owned catalog. Direct generation is a different boundary. `Application.ShellCompletionScript` and `ShellCompletionScriptForVersion` accept any valid provider passed by the caller even when that provider is not present in the Definition catalog.

`renderXonshAdapter` owns Xonsh registration, observation, quoting, and candidate presentation. The Commandment side of that renderer is fully determined by `ShellCompletionGeneration` plus the adapter protocol documented at the end of this guide.

The generation value contains exactly four facts.

- `ApplicationName` is the executable spelling the adapter invokes.
- `TargetVersion` is the requested artifact target and may be empty.
- `ActivationEnvironment` is the environment variable name that enters the completion adapter lane.
- `ActivationValue` is the opaque value the adapter assigns to that variable.

Embed both activation values verbatim. Do not parse, compare, increment, or negotiate `ActivationValue`. A custom provider that interprets `TargetVersion` owns validation and generation decisions for its own shell.

A generator retained by a compiled application must be deterministic for the same generation value and safe for concurrent calls. Mutable closure state remains application-owned.

## A custom adapter invokes the application through one process boundary

At completion time the generated shell code does four Commandment-facing things.

1. Observe the shell's actual words, cursor, runtime family, and runtime version without claiming more fidelity than the shell provides.
2. Encode one adapter request and write it to the application's stdin.
3. Run `ApplicationName` with `ActivationEnvironment=ActivationValue` and read the adapter response from stdout.
4. Translate returned items, native requests, or delegation into the target shell's completion facilities.

There is no helper executable and no second application command registry. Commandment receives the adapter request before normal CLI parsing, maps it into one `CompletionRequest`, executes semantic completion, and projects the result back into the adapter response.

The shell owns observation and presentation. Commandment owns semantic completion.

## Preserve the truthfulness of shell observation

The adapter request can distinguish whether the shell has exact, prefix-exact, or uncertain argv topology. Each word can be exact, lexical, or expanded. Cursor coordinates are expressed in truthful argument or observed-word byte space.

Do not parse a shell command line yourself and report the result as exact argv when the shell only exposes lexical editor words. Converting source coordinates into truthful observed-word or argument coordinates is shell-specific adapter work.

This fidelity is why generated shell adapters can use one completion engine without pretending every shell exposes the same editor model.

## Shell adapters intentionally use a bounded matching profile

The adapter protocol exposes three matching profiles rather than the complete direct-client matching API.

| Adapter profile | Commandment behavior |
| --- | --- |
| `prefix` | prefix matching with smart case over the default completion fields |
| `value-prefix` | prefix matching with smart case over the candidate value only |
| `client` | no server-side textual filtering over the candidate value |

The default is `prefix`.

Use `client` when the shell itself must perform textual candidate filtering. Use direct `Application.Complete` when a non-shell client needs the richer exact, substring, fuzzy, approximate, field-selection, or ranking controls described earlier.

Shell adapter requests use invoked completion, public visibility, relevance ranking, and no result limit. The adapter does not negotiate those semantic policies.

## Capabilities tell Commandment what the shell can represent

A custom adapter should request only presentation it can actually carry back to the user.

Supported presentation capability tokens are

```text
descriptions
details
groups
candidate-kinds
match-ranges
commit-characters
messages
delegation
```

Supported native capability tokens are

```text
paths
files
directories
executables
users
hosts
environment-variables
```

Unknown optional capability tokens are treated as unsupported rather than corrupting the core request. Commandment emits capability-gated response material only when the adapter requested the corresponding known capability.

Native enumeration stays shell-owned. When the semantic completion result requests a native category the adapter declared it can perform, the response identifies that active native kind and the shell carries out the enumeration. Commandment does not send its internal target or fallback graph across the adapter protocol.

## The response contains only shell-actionable completion material

Every successful protocol response carries the completion failure state, replacement range, result spacing, and ordered items. Items carry insertion values, optional structural replacements, and spacing. Display text is available when present.

Requested capabilities can add descriptions, details, group metadata, candidate kinds, match ranges, commit characters, structured messages, and delegation. Native kinds are returned only for native capabilities the adapter declared.

The response deliberately does not expose Commandment declaration identities, completion targets, ranking scores or algorithms, provenance graphs, Diagnostics, Repairs, provider state, or sensitivity flags. The shell receives what it needs to finish the editor interaction rather than a second serialized Commandment model.

## Custom providers need real-shell proof only when that runtime is available

Source generation proves only that the generator returned an artifact. It does not prove registration, quoting, shell observation, protocol transport, native enumeration, or editor behavior.

A custom provider deserves real-shell proof when the required shell runtime is available. Do not fabricate a custom-shell example around an unavailable runtime merely to demonstrate `Shell.Custom`.

[Testing Shell Integrations with shelltest](13-testing-shell-integrations-with-shelltest.md) owns the real shell/editor proof boundary. The `shelltest` guide shows how a provider, runtime, preparation callback, editor input, and interpreted observation compose into one test.

## Completion can be partial or truncated

A limited client may request fewer results than the application could otherwise produce. A dynamic provider may also leave some structure unavailable when analysis cannot fully complete.

Treat presence in a partial result as evidence. Do not infer exact absence unless the result has enough completeness to support that claim.

`clitest` completion expectations preserve this distinction so tests do not accidentally turn truncation into a semantic assertion.

## Keep completion presentation tied to the semantic subject

Descriptions, lifecycle state, aliases, finite-domain entry metadata, and other presentation should come from the accepted semantic subjects rather than a completion-only annotation registry.

That lets a deprecated command or value-domain entry remain consistently deprecated in Help, Reference, completion, inspection, and tests.

## Custom shell adapter protocol reference

This section is the public wire reference for a provider author. The target shell's registration and editor APIs remain shell-owned.

### Activation and framing

The generated adapter invokes `ApplicationName` with `ActivationEnvironment` set exactly to `ActivationValue`. It writes one request to stdin and reads one response from stdout.

The first input line selects the request codec.

```text
records
```

or

```text
json
```

The response uses the same codec unless the request asks for the other codec through `response_encoding`.

JSON requests are strict UTF-8 and unknown fields are rejected.

### JSON request shape

A minimal JSON request can look like this.

```json
{
  "shell": {
    "family": "xonsh",
    "version": "0.19.3"
  },
  "input": {
    "words": [
      {"value": "alienctl"},
      {"value": "deploy"},
      {"value": "--pla"}
    ],
    "cursor": {
      "word": 2,
      "byte": 5
    }
  }
}
```

A provider that needs explicit input truth, capabilities, native delegation, or a different response codec uses the same shape with those fields present.

```json
{
  "response_encoding": "json",
  "shell": {
    "family": "xonsh",
    "version": "0.19.3"
  },
  "input": {
    "topology": "prefix-exact",
    "exact_argv_prefix": 2,
    "words": [
      {"value": "alienctl", "fidelity": "exact"},
      {"value": "deploy", "fidelity": "exact"},
      {"value": "--pla", "fidelity": "lexical"}
    ],
    "cursor": {
      "word": 2,
      "byte": 5,
      "space": "observed-word-bytes"
    }
  },
  "matching": "value-prefix",
  "capabilities": [
    "descriptions",
    "groups",
    "candidate-kinds"
  ],
  "native_capabilities": [
    "files",
    "directories"
  ]
}
```

Topology is `exact`, `prefix-exact`, or `uncertain`. Word fidelity is `exact`, `lexical`, or `expanded`. Cursor space is `argument-bytes` or `observed-word-bytes`. Matching is `prefix`, `value-prefix`, or `client`.

When omitted, JSON uses uncertain topology, zero exact argv prefix, lexical word fidelity, observed-word byte coordinates, prefix matching, the request codec for the response, and no optional capabilities.

### Request validity

Both request codecs must describe one admissible completion observation.

- `words` is nonempty.
- `cursor.word` indexes `words` and `cursor.byte` is nonnegative.
- Every word uses `exact`, `lexical`, or `expanded` fidelity.
- `exact` topology establishes the entire observed word vector as one-to-one argv topology. `exact_argv_prefix` is not used for that topology.
- `prefix-exact` requires `exact_argv_prefix` between zero and `len(words)`, inclusive.
- `uncertain` requires `exact_argv_prefix` to be zero.
- Within an established exact prefix, exact semantic settlement stops at the first word whose fidelity is not `exact`; later exact-looking words do not regain exact positional meaning.
- `argument-bytes` is valid only when the cursor-owned word has `exact` fidelity and its position lies inside the established exact argv prefix. `observed-word-bytes` remains local to the reported observed word.
- A byte offset that applies to the analyzed word may not exceed that word's UTF-8 byte length.

JSON is deliberately strict. The complete request must be one valid UTF-8 JSON value. Duplicate object members, `null`, wrong JSON types, unknown fields, malformed values, noncanonical number tokens, and trailing non-whitespace are rejected. JSON numbers use the same canonical integer form as Records. They are base-10 integer text only, with no leading zeroes, leading `+`, `-0`, fraction, or exponent. Unknown strings inside `capabilities` and `native_capabilities` are the only forward-compatible extension area; unsupported tokens there are ignored as unavailable capabilities.

### Records request shape

The equivalent minimal request can use Records.

```text
records
shell	xonsh	0.19.3
input	uncertain	0	2	5	observed-word-bytes
word	0	alienctl	lexical
word	1	deploy	lexical
word	2	--pla	lexical
matching	prefix
end
```

The request record vocabulary is fixed.

```text
response-encoding <records|json>
shell             <family> <version>
input             <topology> <exact-prefix> <cursor-word> <cursor-byte> <cursor-space>
word              <index> <value> <fidelity>
matching          <prefix|value-prefix|client>
capability        <token>
native-capability <token>
end
```

`response-encoding`, `matching`, `capability`, and `native-capability` are optional. `shell`, `input`, contiguous indexed `word` records, and `end` are required. Capability records may repeat.

TAB separates fields and LF terminates records. Fields may contain literal bytes only from printable ASCII `0x20` through `0x7e`, except backslash. Encode backslash as `\\`, TAB as `\t`, LF as `\n`, and CR as `\r`. Encode every other nonzero byte as uppercase `\xHH`. NUL is invalid. This means UTF-8 bytes outside printable ASCII are escaped byte by byte rather than written literally.

Every integer or index uses canonical base-10 text. Leading zeroes, a leading `+`, and `-0` are invalid.

### JSON response shape

The required JSON response fields are `failed`, `replace`, `spacing`, and `items`. The fully populated shape below assumes the corresponding presentation and native capabilities were requested and the semantic result contains each optional value.

```json
{
  "failed": false,
  "replace": {
    "argument": 2,
    "start": 0,
    "end": 5
  },
  "spacing": "default",
  "items": [
    {
      "value": "--planet",
      "replacement": ["--planet"],
      "spacing": "append-space",
      "display_text": "--planet",
      "description": "Target planet",
      "detail": "Known navigation target",
      "kind": "flag",
      "group": "flags",
      "match_ranges": [
        {"start": 0, "end": 5}
      ],
      "commit_characters": ["="]
    }
  ],
  "groups": [
    {"id": "flags", "label": "Flags"}
  ],
  "messages": [
    {"kind": "information", "text": "Showing navigation flags"}
  ],
  "native": ["files"],
  "delegation": {
    "args": ["alienctl", "deploy", "--planet"],
    "cursor": {"argument": 2, "byte": 8}
  }
}
```

`replacement` is an optional nonempty string array and is distinct from the top-level replacement span. `match_ranges` is an optional array of byte ranges with `start` and `end`. `commit_characters` is an optional array of one-rune strings. `groups` is an optional array of `id` and `label` objects. `messages` is an optional array whose `kind` is `information`, `warning`, or `error` and whose `text` carries the message. `native` is an optional array using the native capability tokens listed above. `delegation` is an optional object with `args` and a cursor containing `argument` and `byte`.

`kind`, when present, is one of the public completion candidate tokens.

```text
value
command
command-alias
command-shortcut
action
flag
keyword
file
directory
executable
resource
namespace
user
host
environment-variable
url
```

`spacing` is `default`, `append-space`, or `no-space` both at the result level and on each item.

### Records response shape

Records responses carry the same semantics as tagged lines.

```text
result                <failed>
replace               <argument> <start> <end>
spacing               <default|append-space|no-space>
item                  <index> <value> <replacement-present> <spacing>
item-replacement      <item-index> <replacement-index> <value>
item-display          <item-index> <text>
item-description      <item-index> <text>
item-detail           <item-index> <text>
item-kind             <item-index> <kind>
item-group            <item-index> <group-id>
item-match-range      <item-index> <range-index> <start> <end>
item-commit-character <item-index> <character-index> <value>
group                 <index> <id> <label>
message               <index> <information|warning|error> <text>
native                <index> <kind>
delegation            <cursor-argument> <cursor-byte>
delegation-arg        <index> <value>
end
```

`result`, `replace`, and `spacing` occur exactly once. Item indexes are contiguous from zero. Replacement records exist only when the item marks a replacement present. Optional presentation, native, message, and delegation records appear only when the request negotiated the corresponding capability. Records fields use the same escaping and canonical integer rules as requests.

Continue with [Evolving a CLI](08-evolving-a-cli.md) for the release model that completion can project, or [Testing Shell Integrations with shelltest](13-testing-shell-integrations-with-shelltest.md) when the shell/editor boundary is the thing you need to prove.
