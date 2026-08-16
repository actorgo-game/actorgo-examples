package main

import (
	"flag"

	"github.com/actorgo-game/actorgo"
	cgorm "github.com/actorgo-game/actorgo/components/gorm"
)

func main() {
	profilePath := flag.String("path", "../config/demo-gorm.json", "profile config path")
	nodeID := flag.String("node", "0.0.5.1", "ActorGo node ID")
	flag.Parse()

	app := actorgo.Configure(
		*profilePath,
		*nodeID,
		actorgo.Standalone,
	)

	// 注册 GORM 组件，数据库配置见 config/demo-gorm.json。
	app.Register(cgorm.NewComponent())

	app.AddActors(
		&ActorDB{},
	)

	app.Startup()
}
