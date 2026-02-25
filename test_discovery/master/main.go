package main

import (
	"github.com/actorgo-game/actorgo"
	ccluster "github.com/actorgo-game/actorgo/net/cluster"
)

func main() {
	app := actorgo.NewApp(
		"../config/test-discovery.json",
		"master-1",
		true,
		actorgo.Cluster,
	)

	app.Register(ccluster.New())
	app.Startup()
}
