package main

import (
	"io"
	"testing"

	"github.com/abstractionsinc/commandment/cli"
	"github.com/abstractionsinc/commandment/clitest"
)

func TestInvalidDockCodeCarriesARepair(t *testing.T) {
	t.Parallel()

	app := clitest.Compile(t, definition)
	request := cli.RunRequest{
		Args:   []string{"dock", "--code", "ZZ"},
		Stdout: io.Discard,
	}
	err := app.Run(t.Context(), request)
	if err == nil {
		t.Fatal("Run() error = nil, want invalid dock code")
	}
	diagnostic := clitest.Diagnostic(cli.DiagnosticInvalidValue)
	clitest.Diagnostics(err).Require(t,
		clitest.DiagnosticsExactly(diagnostic),
		diagnostic.RepairCount(1),
		diagnostic.RepairAt(0).Present(),
	)
}
