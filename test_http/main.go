package main

import (
	"flag"

	"github.com/actorgo-game/actorgo"
	clog "github.com/actorgo-game/actorgo/logger"
	httpactor "github.com/actorgo-game/actorgo/net/httpactor"
)

func main() {
	profilePath := flag.String("path", "../config/test-http.json", "profile config path")
	nodeID := flag.String("node", "0.0.3.1", "ActorGo node ID")
	flag.Parse()

	app := actorgo.Configure(*profilePath, *nodeID, actorgo.Standalone)
	app.AddActors(new(HTTPActor))
	app.Register(httpactor.NewComponent("http_example", app.Address()))
	clog.Info("[test_http] starting MethodID HTTP service. [address = %s, nodeID = %s]", app.Address(), app.NodeID())
	app.Startup()
}
