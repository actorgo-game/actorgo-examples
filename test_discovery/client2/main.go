package main

import (
	"github.com/actorgo-game/actorgo"
)

func main() {
	app := actorgo.Configure(
		"../../config/test-discovery.json",
		"0.0.5.2",
		actorgo.Cluster,
	)
	app.Startup()
}
