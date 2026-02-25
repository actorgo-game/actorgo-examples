package rpcGame

import (
	"fmt"

	cfacade "github.com/actorgo-game/actorgo/facade"
	clog "github.com/actorgo-game/actorgo/logger"
	cproto "github.com/actorgo-game/actorgo/net/proto"
	"github.com/actorgo-game/examples/demo_cluster/internal/pb"
	sessionKey "github.com/actorgo-game/examples/demo_cluster/internal/session_key"
)

const (
	playerActor = "player"
)

const (
	sessionClose = "sessionClose"
)

// SessionClose 如果session已登录，则调用rpcGame.SessionClose() 告知游戏服
func SessionClose(app cfacade.IApplication, session *cproto.Session) {
	nodeID := session.GetString(sessionKey.ServerID)
	if nodeID == "" {
		clog.Warn("Get server id fail. session = %s", session.Sid)
		return
	}

	targetPath := fmt.Sprintf("%s.%s.%s", nodeID, playerActor, session.Sid)
	app.ActorSystem().Call("", targetPath, sessionClose, &pb.Int64{
		Value: session.Uid,
	})

	//clog.Info("send close session to game node. [node = %s, uid = %d]", nodeID, session.Uid)
}
