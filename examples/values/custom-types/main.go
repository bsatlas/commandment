package main

import (
	"context"
	"fmt"
	"net/netip"
	"os"

	"github.com/abstractionsinc/commandment/cli"
)

type relayConfig struct {
	Address netip.Addr
}

var definition = cli.Definition[struct{}]{
	Name: "alienctl",
	Declare: func(cb *cli.CommandBuilder[struct{}]) {
		cmd := cli.Command[relayConfig]{
			Name: "relay",
			Declare: func(cb *cli.CommandBuilder[relayConfig]) {
				addressCodec := cli.ValueCodec[netip.Addr]{
					Name:     "ip-address",
					Parse:    netip.ParseAddr,
					Format:   netip.Addr.String,
					Snapshot: cli.Value.Assign[netip.Addr](),
				}
				address := cb.Values.Value(&cb.Config.Address, addressCodec)
				cb.Flags.ValueFor(address).Long("address")
				cb.Rule.Require(address)
			},
			Run: func(_ context.Context, cfg *relayConfig) error {
				fmt.Printf("relaying through %s\n", cfg.Address)
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
