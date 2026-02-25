package game

import (
	"github.com/actorgo-game/actorgo"
	ccron "github.com/actorgo-game/actorgo/components/cron"
	cherryGops "github.com/actorgo-game/actorgo/components/gops"
	cherrySnowflake "github.com/actorgo-game/actorgo/extend/snowflake"
	cstring "github.com/actorgo-game/actorgo/extend/string"
	cherryUtils "github.com/actorgo-game/actorgo/extend/utils"
	checkCenter "github.com/actorgo-game/examples/demo_cluster/internal/component/check_center"
	"github.com/actorgo-game/examples/demo_cluster/internal/data"
	"github.com/actorgo-game/examples/demo_cluster/nodes/game/db"
	"github.com/actorgo-game/examples/demo_cluster/nodes/game/module/player"
)

func Run(profileFilePath, nodeID string) {
	if !cherryUtils.IsNumeric(nodeID) {
		panic("node parameter must is number.")
	}

	// snowflake global id
	serverId, _ := cstring.ToInt64(nodeID)
	cherrySnowflake.SetDefaultNode(serverId)

	// 配置cherry引擎
	app := actorgo.Configure(profileFilePath, nodeID, false, actorgo.Cluster)

	// diagnose
	app.Register(cherryGops.New())
	// 注册调度组件
	app.Register(ccron.New())
	// 注册数据配置组件
	app.Register(data.New())
	// 注册检测中心节点组件，确认中心节点启动后，再启动当前节点
	app.Register(checkCenter.New())
	// 注册db组件
	app.Register(db.New())

	app.AddActors(
		&player.ActorPlayers{},
	)

	app.Startup()
}
