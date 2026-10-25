package main

import (
	"testing"

	"github.com/abstractionsinc/commandment/cli"
	"github.com/abstractionsinc/commandment/clitest"
)

func TestAddressUsesTheApplicationCodec(t *testing.T) {
	t.Parallel()

	app := clitest.Compile(t, definition)
	address := clitest.Command("relay").Flag("--address").Value()
	request := cli.InspectionRequest{
		Args: []string{"relay", "--address", "203.0.113.7"},
	}
	clitest.Inspect(t.Context(), app, request).Require(t,
		clitest.InspectionReady(),
		address.Text("203.0.113.7"),
	)
}
