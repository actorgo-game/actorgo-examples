package main

import (
	"github.com/actorgo-game/actorgo"
	cgin "github.com/actorgo-game/actorgo/components/gin"
)

func main() {
	app := actorgo.NewApp(
		"../config/test.json",
		"web-1",
		false,
		actorgo.Standalone,
	)

	httpServer := cgin.NewHttp("web_1", app.Address())
	httpServer.Use(cgin.Cors(), cgin.MaxConnect(2))
	httpServer.Register(new(Test1Controller))

	app.Register(httpServer)
	app.Startup()
}
