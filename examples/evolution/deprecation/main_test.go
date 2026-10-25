package main

import (
	"testing"

	"github.com/abstractionsinc/commandment/cli"
	"github.com/abstractionsinc/commandment/clitest"
)

func TestLegacyVelocitySpellingPointsToSpeed(t *testing.T) {
	t.Parallel()

	app := clitest.Compile(t, definition)
	flag := clitest.Command("approach").Flag("--speed")
	legacy := flag.Spelling("--velocity")
	current := flag.Spelling("--speed")
	clitest.Schema(app).Require(t,
		legacy.Lifecycle(cli.LifecycleStateDeprecated),
		legacy.ReplacedBy(current),
		current.Lifecycle(cli.LifecycleStateCurrent),
	)
}
