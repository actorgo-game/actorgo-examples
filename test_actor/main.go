package main

import (
	"fmt"

	"github.com/actorgo-game/actorgo"
)

func main() {
	fmt.Println("test actor &  child actor")

	app := actorgo.Configure(
		"../config/test.json", // 使用环境的配置
		"0.0.5.1",             // 使用 game 节点的 NodeID
		actorgo.Standalone,
	)

	app.AddActors(&actor{})
	app.Startup()
}
