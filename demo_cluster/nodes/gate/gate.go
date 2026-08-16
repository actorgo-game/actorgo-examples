package gate

import (
	"github.com/actorgo-game/actorgo"
	cherryGops "github.com/actorgo-game/actorgo/components/gops"
	cfacade "github.com/actorgo-game/actorgo/facade"
	cconnector "github.com/actorgo-game/actorgo/net/connector"
	"github.com/actorgo-game/actorgo/net/parser"
	checkCenter "github.com/actorgo-game/examples/demo_cluster/internal/component/check_center"
	"github.com/actorgo-game/examples/demo_cluster/internal/data"
	rpcGame "github.com/actorgo-game/examples/demo_cluster/internal/rpc/game"
)

func Run(profileFilePath, nodeID string) {
	app := actorgo.Configure(profileFilePath, nodeID, actorgo.Cluster)
	agpServer := parser.New(
		"gate",
		[]cfacade.IConnector{
			cconnector.NewTCP(":10011"),
			cconnector.NewWS(app.Address()),
		},
		parser.WithOnDisconnect(func(connection *parser.Connection) {
			rpcGame.SessionClose(app, connection.Session())
		}),
	)

	app.Register(cherryGops.New())
	app.Register(checkCenter.New())
	app.Register(data.New())
	app.Register(agpServer)
	app.AddActors(NewActorGate(agpServer))
	app.Startup()
}
