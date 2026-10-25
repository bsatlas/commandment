package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/abstractionsinc/commandment/cli"
)

func renderHelp(_ context.Context, page cli.HelpPage) (string, error) {
	text := fmt.Sprintf("%s %s\n%s\n", strings.ToUpper(page.ApplicationName), page.Name, page.Description.Short)
	return text, nil
}

var definition = cli.Definition[struct{}]{
	Name: "alienctl",
	Help: cli.HelpConfig{
		Renderer: renderHelp,
	},
	Declare: func(cb *cli.CommandBuilder[struct{}]) {
		scanCmd := cli.Command[struct{}]{
			Name: "scan",
			Description: cli.Description{
				Short: "scan a nearby world",
			},
		}
		cb.Subcommands.Command(scanCmd)
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
