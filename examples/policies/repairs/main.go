package main

import (
	"context"
	"fmt"
	"os"

	"github.com/abstractionsinc/commandment/cli"
)

type dockConfig struct {
	Code string
}

func parseDockCode(_ context.Context, tokens []string) (string, error) {
	code := tokens[0]
	if code == "A1" {
		return code, nil
	}
	rejection := cli.Rejection.Message(
		"dock-code",
		"invalid dock code",
		cli.Rejection.At(0, 0, len(code)),
		cli.Rejection.ReplaceWith("A1"),
	)
	return "", rejection
}

func formatDockCode(code string) []string {
	return []string{code}
}

var definition = cli.Definition[struct{}]{
	Name: "alienctl",
	Declare: func(cb *cli.CommandBuilder[struct{}]) {
		cmd := cli.Command[dockConfig]{
			Name: "dock",
			Declare: func(cb *cli.CommandBuilder[dockConfig]) {
				language := cli.TokenLanguage[string]{
					Name:   "dock-code",
					Count:  cli.Count.Exactly(1),
					Parse:  parseDockCode,
					Format: formatDockCode,
				}
				code := cb.Values.String(&cb.Config.Code)
				cb.Flags.Tokens(code, language).Long("code")
				cb.Rule.Require(code)
			},
			Run: func(_ context.Context, cfg *dockConfig) error {
				fmt.Printf("docking at %s\n", cfg.Code)
				return nil
			},
		}
		cb.Subcommands.Command(cmd)
	},
}

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
