package main

import (
	"testing"

	"github.com/abstractionsinc/commandment/cli"
	"github.com/abstractionsinc/commandment/clitest"
)

func TestDerivedSummaryUsesResolvedInputs(t *testing.T) {
	t.Parallel()

	app := clitest.Compile(t, definition)
	command := clitest.Command("abduct")
	target := command.ArgumentFormAt(0).Term("target").Value()
	beamStrength := command.Flag("--beam-strength").Value()
	summary := command.ValueAt(2)
	request := cli.InspectionRequest{
		Args: []string{"abduct", "--beam-strength", "4", "mars"},
	}
	clitest.Inspect(t.Context(), app, request).Require(t,
		clitest.InspectionReady(),
		beamStrength.Text("4"),
		summary.Resolved(),
		summary.ContributedBy(cli.ValueProducerKindDerived),
	)
	clitest.Schema(app).Require(t,
		summary.DependenciesExactly(target, beamStrength),
	)
}
