package main

import (
	"github.com/actorgo-game/actorgo"
	ccluster "github.com/actorgo-game/actorgo/net/cluster"
)

func main() {
	app := actorgo.NewApp(
		"../config/test-discovery.json",
		"game-1",
		false,
		actorgo.Cluster,
	)
	app.Register(ccluster.New())

	app.Startup()
}
