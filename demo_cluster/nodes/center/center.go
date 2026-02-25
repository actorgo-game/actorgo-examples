package center

import (
	"github.com/actorgo-game/actorgo"
	ccron "github.com/actorgo-game/actorgo/components/cron"
	"github.com/actorgo-game/examples/demo_cluster/internal/data"
	"github.com/actorgo-game/examples/demo_cluster/nodes/center/db"
	"github.com/actorgo-game/examples/demo_cluster/nodes/center/module/account"
	"github.com/actorgo-game/examples/demo_cluster/nodes/center/module/ops"
)

func Run(profileFilePath, nodeID string) {
	app := actorgo.Configure(
		profileFilePath,
		nodeID,
		false,
		actorgo.Cluster,
	)

	app.Register(ccron.New())
	app.Register(data.New())
	app.Register(db.New())

	app.AddActors(
		&account.ActorAccount{},
		&ops.ActorOps{},
	)

	app.Startup()
}
