package main

import (
	"github.com/actorgo-game/actorgo"
	cherryGORM "github.com/actorgo-game/actorgo/components/gorm"
)

func main() {
	app := actorgo.Configure(
		"../config/demo-gorm.json", // 使用环境的配置
		"game-1",                   // 使用game-1 的节点id
		false,
		actorgo.Standalone,
	)

	// 注册gorm组件，数据库具体配置请查看 config/demo-gorm.json文件
	app.Register(cherryGORM.NewComponent())

	app.AddActors(
		&ActorDB{},
	)

	app.Startup()
}
