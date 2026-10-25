package main

import (
	"context"
	"fmt"
	"os"

	"github.com/abstractionsinc/commandment/cli"
)

type abductConfig struct {
	Target       string
	BeamStrength int
	Summary      string
}

var definition = cli.Definition[struct{}]{
	Name: "alienctl",
	Declare: func(cb *cli.CommandBuilder[struct{}]) {
		cmd := cli.Command[abductConfig]{
			Name: "abduct",
			Declare: func(cb *cli.CommandBuilder[abductConfig]) {
				target := cb.Arguments.String(&cb.Config.Target, "target")
				beamStrength := cb.Values.Int(&cb.Config.BeamStrength)
				cb.Values.Default(beamStrength, 1)
				cb.Flags.ValueFor(beamStrength).Long("beam-strength")

				deriveSummary := func(_ context.Context, _ *abductConfig, dependencies cli.Dependencies) (string, error) {
					resolvedTarget, _ := dependencies.Get(target)
					resolvedBeamStrength, _ := dependencies.Get(beamStrength)
					summary := fmt.Sprintf("abducting %s with beam strength %d", resolvedTarget, resolvedBeamStrength)
					return summary, nil
				}
				summary := cli.Value.ComputedBy(deriveSummary).From(target, beamStrength)
				cb.Values.Derive(&cb.Config.Summary, summary)
			},
			Run: func(_ context.Context, cfg *abductConfig) error {
				fmt.Println(cfg.Summary)
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
