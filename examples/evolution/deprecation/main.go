package main

import (
	"context"
	"fmt"
	"os"

	"github.com/abstractionsinc/commandment/cli"
)

type approachConfig struct {
	Speed int
}

var definition = cli.Definition[struct{}]{
	Name:    "alienctl",
	Release: cli.Release.SemVer("2.0.0"),
	Version: cli.VersionConfig{},
	Declare: func(cb *cli.CommandBuilder[struct{}]) {
		cmd := cli.Command[approachConfig]{
			Name: "approach",
			Declare: func(cb *cli.CommandBuilder[approachConfig]) {
				speed := cb.Flags.Int(&cb.Config.Speed).Description("Set approach speed.")
				legacy := speed.Long("velocity").Release(cli.ReleaseLifecycle{
					IntroducedIn: "1.0.0",
					DeprecatedIn: "2.0.0",
					Reason:       "velocity was renamed",
				})
				current := speed.Long("speed")
				legacy.ReplacedBy(current)
			},
			Run: func(_ context.Context, cfg *approachConfig) error {
				fmt.Printf("approaching at speed %d\n", cfg.Speed)
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
