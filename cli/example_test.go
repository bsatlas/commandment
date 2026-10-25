package cli_test

import (
	"context"
	"fmt"

	"github.com/abstractionsinc/commandment/cli"
)

type exampleVerboseConfig struct{ Verbose bool }

func ExampleCompile() {
	run := false
	definition := cli.Definition[struct{}]{
		Name: "alienctl",
		Run: func(context.Context, *struct{}) error {
			run = true
			return nil
		},
	}
	app, err := cli.Compile(definition)
	if err != nil {
		panic(err)
	}
	if err := app.Run(context.Background(), cli.RunRequest{}); err != nil {
		panic(err)
	}

	fmt.Println(run)
	// Output: true
}

func ExampleApplication_Schema() {
	definition := cli.Definition[struct{}]{
		Name: "alienctl",
	}
	app, err := cli.Compile(definition)
	if err != nil {
		panic(err)
	}

	fmt.Println(len(app.Schema().Commands))
	// Output: 2
}

func ExampleCommandBuilder() {
	type config = exampleVerboseConfig
	definition := cli.Definition[config]{
		Name: "alienctl",
		Declare: func(cb *cli.CommandBuilder[config]) {
			cb.Flags.Bool(&cb.Config.Verbose).Long("verbose")
		},
		Run: func(_ context.Context, cfg *config) error {
			fmt.Println(cfg.Verbose)
			return nil
		},
	}
	app, err := cli.Compile(definition)
	if err != nil {
		panic(err)
	}
	request := cli.RunRequest{
		Args: []string{"--verbose"},
	}
	if err := app.Run(context.Background(), request); err != nil {
		panic(err)
	}

	// Output: true
}
