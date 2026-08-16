package main

import (
	"context"

	"github.com/actorgo-game/actorgo"
	cgin "github.com/actorgo-game/actorgo/components/gin"
	cfacade "github.com/actorgo-game/actorgo/facade"
	clog "github.com/actorgo-game/actorgo/logger"
	cconnector "github.com/actorgo-game/actorgo/net/connector"
	"github.com/actorgo-game/actorgo/net/parser"
)

// main 启动单节点聊天室。
func main() {
	app := actorgo.Configure(
		"../../config/demo-chat.json",
		"0.0.1.1",
		actorgo.Standalone,
	)
	app.SetDefaultBodyCodec(cfacade.CodecJSON)

	agpServer := parser.New(
		"chat",
		[]cfacade.IConnector{cconnector.NewWS(app.Address())},
		parser.WithOnDisconnect(func(connection *parser.Connection) {
			session := connection.Session()
			if session == nil || session.Uid < 1 {
				return
			}
			ctx := cfacade.NewRequestContext(context.Background())
			ctx.Codec = cfacade.CodecJSON
			ctx.Session = session
			result := app.ActorSystem().Notify(ctx, MethodExit, &Int64{Value: session.Uid})
			if !result.OK() {
				clog.Warn("notify room exit failed. uid=%d code=%d message=%s", session.Uid, result.Code, result.Message)
			}
		}),
	)

	app.Register(agpServer)
	app.AddActors(newActorRoom(agpServer))
	registerHTTP(app)
	app.Startup()
}

// registerHTTP 部署 H5 静态文件。
func registerHTTP(app *actorgo.AppBuilder) {
	httpComponent := cgin.New("web", ":8081")
	httpComponent.Use(cgin.RecoveryWithZap(true))
	httpComponent.Static("/", "../static/")
	app.Register(httpComponent)
}
