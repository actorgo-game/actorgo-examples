package rpcGame

import (
	cfacade "github.com/actorgo-game/actorgo/facade"
	clog "github.com/actorgo-game/actorgo/logger"
	cproto "github.com/actorgo-game/actorgo/net/proto"
	"github.com/actorgo-game/examples/demo_cluster/internal/methodid"
	"github.com/actorgo-game/examples/demo_cluster/internal/pb"
	sessionKey "github.com/actorgo-game/examples/demo_cluster/internal/session_key"
)

func SessionClose(app cfacade.IApplication, session *cproto.Session) {
	nodeID := session.GetString(sessionKey.GameNodeID)
	if nodeID == "" || session.Uid == 0 {
		clog.Warn("Get game node id fail. session = %s", session.Sid)
		return
	}

	ctx := cfacade.NewRequestContext(nil)
	ctx.Codec = cfacade.CodecProtobuf
	ctx.Session = session
	targetPath := cfacade.NewChildPath(nodeID, "player", session.Uid)
	result := app.ActorSystem().NotifyTarget(ctx, targetPath, methodid.GameSessionClose, &pb.Int64{Value: session.Uid})
	if result != nil && !result.OK() {
		clog.Warn("send close session to game node fail. node = %s, uid = %d, err = %s", nodeID, session.Uid, result.Message)
	}
}
