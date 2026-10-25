# Runtime and Host Integration

After `Compile`, a host has two jobs. It supplies invocation capabilities and it chooses which application operation to run. Commandment should not require a global process singleton or a general-purpose dependency-injection container to do either job.

## Use `RunFromProcess` for an ordinary executable

A conventional `main` can stay thin.

```go
func main() {
	app, err := cli.Compile(definition)
	if err != nil {
		_ = cli.WriteError(os.Stderr, err)
		os.Exit(1)
	}
	if err := app.RunFromProcess(context.Background()); err != nil {
		_ = cli.WriteError(os.Stderr, err)
		os.Exit(cli.ExitCode(err))
	}
}
```

`RunFromProcess` is the process adapter. It uses the process invocation and supported process capabilities rather than making application packages read `os.Args` and `os.Environ` themselves.

## Use `Run` when another host owns the invocation explicitly

An embedded application, daemon, alternate process adapter, or test can construct a `RunRequest`.

```go
request := cli.RunRequest{
	Args:          args,
	Env:           env,
	Defaults:      defaults,
	Facts:         facts,
	Notify:        notify,
	ResponseFiles: responseFiles,
}
err := app.Run(ctx, request)
```

Provide only the capabilities the host actually owns. An explicit request makes host dependencies visible and keeps concurrent invocations independent.

## Supply invocation defaults when the host owns a fallback

An invocation default contributes to an ordinary destination at runtime while still allowing stronger producers to win according to that destination's resolution policy.

The [host-context/defaults](../examples/host-context/defaults/) scenario captures the destination while building the definition, then builds a provider around that durable handle.

```go
supplyDefault := func(_ context.Context, supply *cli.DefaultSupply[statusConfig]) error {
	supply.Set(ship, "scout-7")
	return nil
}
defaults := cli.DefaultsFor[statusConfig](supplyDefault)
```

The request supplies the provider.

```go
request := cli.RunRequest{
	Args:     os.Args[1:],
	Defaults: []cli.DefaultProvider{defaults},
}
```

A user can still override the fallback through an explicit producer such as `--ship` when the value model allows it.

Use a factory when it binds a real declaration handle to a provider with a lifetime relationship. Do not add one merely to create a test-access wrapper.

## Supply facts when the host owns trusted context outside CLI input

Facts are typed host-supplied context rather than ordinary grammar input.

```go
ship = cb.Facts.Require(&cb.Config.Ship)
```

The host then supplies that fact through a provider.

```go
supplyFacts := func(_ context.Context, supply *cli.FactSupply[contactConfig]) error {
	supply.Set(ship, "mothership")
	return nil
}
facts := cli.FactsFor[contactConfig](supplyFacts)
```

In

```text
alienctl contact martian
```

`martian` came from the user. `mothership` came from the host. Inspection and tests can preserve that distinction rather than pretending both values were parser input.

See [host-context/facts](../examples/host-context/facts/).

## Give handlers the standard resources they actually need

Execution access can place supported runtime resources into selected-command configuration.

```go
type scanConfig struct {
	World  string
	Output io.Writer
}

cmd := cli.Command[scanConfig]{
	Name: "scan",
	Declare: func(cb *cli.CommandBuilder[scanConfig]) {
		cb.Arguments.String(&cb.Config.World, "world")
		cb.Access.Stdout(&cb.Config.Output)
	},
}
```

A handler can then write to the configured output without reaching for package-global `os.Stdout`.

Execution access also supports selected resolved values and other declared resources where the public API provides them. Use those capabilities for execution concerns that Commandment actually owns. Ordinary application services remain ordinary Go dependencies constructed by the application.

## Use middleware for admitted execution concerns

Middleware wraps application execution after Commandment has completed semantic preparation and admitted the invocation.

```go
func scanMiddleware(ctx context.Context, cfg *scanConfig, next cli.MiddlewareNext) error {
	if _, err := fmt.Fprintln(cfg.Output, "preparing sensors"); err != nil {
		return fmt.Errorf("write middleware output: %w", err)
	}
	if err := next(ctx); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(cfg.Output, "recording scan"); err != nil {
		return fmt.Errorf("write middleware output: %w", err)
	}
	return nil
}
```

Use middleware for cross-cutting execution behavior such as tracing, execution-scoped resource management, or application concerns that genuinely wrap handlers.

Do not use middleware to reparse flags, resolve values, decide command availability, or rebuild validation. Those belong before execution admission.

See [execution/middleware](../examples/execution/middleware/).

## Keep application errors as ordinary Go errors

A handler can return a domain error directly.

```go
var errLandingDenied = errors.New("landing denied: earth is protected")

func runLand(_ context.Context, cfg *landConfig) error {
	if cfg.World == "earth" {
		return cli.WithExitCode(errLandingDenied, 7)
	}
	return nil
}
```

`cli.WithExitCode` adds process exit policy while preserving normal `errors.Is` and `errors.As` identity for the application error.

Do not manufacture a Commandment diagnostic for an application failure merely to make every error look alike. Diagnostics are for failures Commandment semantically owns.

See [execution/errors](../examples/execution/errors/).

## Keep response-file grammar separate from host loading

A command can enable response-file recognition.

```go
cb.Dialect.ResponseFiles(true)
```

The ordinary process host can then load a filesystem source such as `@crew.args` through the default process behavior.

See [response-files/filesystem](../examples/response-files/filesystem/).

An embedded host can provide another loader.

```go
loadCrew := func(_ context.Context, request cli.ResponseFileRequest) (cli.ResponseFile, error) {
	if request.Path != "crew" {
		return cli.ResponseFile{}, errUnknownResponseSource
	}
	response := cli.ResponseFile{
		Key:      "crew",
		Location: "crew",
		Content:  []byte("Zorblax"),
	}
	return response, nil
}
```

The included text re-enters normal Commandment grammar with source history. It does not become an opaque application string.

See [response-files/custom-loader](../examples/response-files/custom-loader/).

## Cancellation has an owner

The context supplied to an operation belongs to the caller. Commandment may also create cancellation relationships as part of declared signal policy or operation handling.

Do not spawn an unrelated application cancellation tree merely to mirror Commandment's invocation lifecycle. Decide whether cancellation came from the host context, a signal policy, application logic, or another explicit owner and preserve that distinction.

[Signals](10-signals.md) develops signal-driven cancellation and escalation separately because it has a richer model than ordinary host context.

## Keep process construction outside the CLI when it is not CLI configuration

Database pools, service clients, caches, repositories, and long-lived domain state should be constructed in ordinary Go packages under application ownership.

Pass them into application behavior through ordinary constructors, closures, command configuration that truthfully owns them, or another explicit application boundary. Do not force every dependency through Commandment merely because the executable uses Commandment.

[Building Large Applications](16-building-large-applications.md) shows how to keep root assembly visible without turning the root declaration into a service locator.

## Test the boundary that owns the fact

Use `clitest` when the question belongs to Commandment's semantic model. Use ordinary Go tests for handler or service behavior. Use subprocess tests for real exit codes, process streams, or process lifecycle. Use `shelltest` for real shell/editor completion. Use real process signals when signal delivery itself is the claim.

A narrower semantic test should not be mislabeled as proof of a wider host boundary.

Continue with [Signals](10-signals.md) for declared process and application-defined signal behavior.
