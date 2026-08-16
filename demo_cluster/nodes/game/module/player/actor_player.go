package player

import (
	cfacade "github.com/actorgo-game/actorgo/facade"
	cactor "github.com/actorgo-game/actorgo/net/actor"
	cproto "github.com/actorgo-game/actorgo/net/proto"
	"github.com/actorgo-game/examples/demo_cluster/internal/code"
	"github.com/actorgo-game/examples/demo_cluster/internal/data"
	"github.com/actorgo-game/examples/demo_cluster/internal/event"
	"github.com/actorgo-game/examples/demo_cluster/internal/methodid"
	"github.com/actorgo-game/examples/demo_cluster/internal/pb"
	sessionKey "github.com/actorgo-game/examples/demo_cluster/internal/session_key"
	"github.com/actorgo-game/examples/demo_cluster/nodes/game/db"
	"github.com/actorgo-game/examples/demo_cluster/nodes/game/module/online"
)

type actorPlayer struct {
	cactor.Base
	isOnline bool
	playerID int64
	uid      int64
}

func (p *actorPlayer) OnInit() {
	p.Methods().Register(methodid.GameSessionClose, p.sessionClose)
	p.Methods().Register(methodid.GamePlayerSelect, p.playerSelect)
	p.Methods().Register(methodid.GamePlayerCreate, p.playerCreate)
	p.Methods().Register(methodid.GamePlayerEnter, p.playerEnter)
}

func (p *actorPlayer) OnStop() {
}

func (p *actorPlayer) sessionClose(ctx *cfacade.RequestContext, req *pb.Int64) error {
	uid := req.Value
	if uid == 0 && ctx != nil && ctx.Session != nil {
		uid = ctx.Session.Uid
	}
	online.UnBindPlayer(uid)
	p.isOnline = false
	p.Exit()
	return nil
}

func (p *actorPlayer) playerSelect(ctx *cfacade.RequestContext, _ *pb.None) (*pb.PlayerSelectResponse, error) {
	session, err := requestSession(ctx)
	if err != nil {
		return nil, err
	}
	response := &pb.PlayerSelectResponse{}
	if playerID := db.GetPlayerIdWithUID(session.Uid); playerID > 0 {
		if playerTable, found := db.GetPlayerTable(playerID); found {
			playerInfo := buildPBPlayer(playerTable)
			response.List = append(response.List, &playerInfo)
		}
	}
	return response, nil
}

func (p *actorPlayer) playerCreate(ctx *cfacade.RequestContext, req *pb.PlayerCreateRequest) (*pb.PlayerCreateResponse, error) {
	session, err := requestSession(ctx)
	if err != nil {
		return nil, err
	}
	if req.Gender > 1 || len(req.PlayerName) < 1 {
		return nil, code.NewInvokeError(code.PlayerCreateFail)
	}
	if db.GetPlayerIdWithUID(session.Uid) > 0 {
		return nil, code.NewInvokeError(code.PlayerCreateFail)
	}
	playerInitRow, found := data.PlayerInitConfig.Get(req.Gender)
	if !found {
		return nil, code.NewInvokeError(code.PlayerCreateFail)
	}

	serverID := session.GetInt32(sessionKey.ServerID)
	newPlayerTable, errCode := db.CreatePlayer(session, req.PlayerName, serverID, playerInitRow)
	if code.IsFail(errCode) {
		return nil, code.NewInvokeError(errCode)
	}

	playerCreateEvent := event.NewPlayerCreate(newPlayerTable.PlayerId, req.PlayerName, req.Gender)
	p.PostEvent(&playerCreateEvent)
	playerInfo := buildPBPlayer(newPlayerTable)
	return &pb.PlayerCreateResponse{Player: &playerInfo}, nil
}

func (p *actorPlayer) playerEnter(ctx *cfacade.RequestContext, req *pb.Int64) (*pb.PlayerEnterResponse, error) {
	session, err := requestSession(ctx)
	if err != nil {
		return nil, err
	}
	playerID := req.Value
	if playerID < 1 {
		return nil, code.NewInvokeError(code.PlayerIDError)
	}
	playerTable, found := db.GetPlayerTable(playerID)
	if !found {
		return nil, code.NewInvokeError(code.PlayerIDError)
	}

	online.BindPlayer(playerID, playerTable.UID, session.Sid)
	p.uid = playerTable.UID
	p.playerID = playerTable.PlayerId
	p.isOnline = true

	response := &pb.PlayerEnterResponse{GuideMaps: map[int32]int32{}}
	loginEvent := event.NewPlayerLogin(p.ActorID(), playerID)
	p.PostEvent(&loginEvent)
	return response, nil
}

func requestSession(ctx *cfacade.RequestContext) (*cproto.Session, error) {
	if ctx == nil || ctx.Session == nil || ctx.Session.Uid < 1 {
		return nil, code.NewInvokeError(code.PlayerDenyLogin)
	}
	return ctx.Session, nil
}

func buildPBPlayer(playerTable *db.PlayerTable) pb.Player {
	return pb.Player{
		PlayerId:   playerTable.PlayerId,
		PlayerName: playerTable.Name,
		Level:      playerTable.Level,
		CreateTime: playerTable.CreateTime,
		Exp:        playerTable.Exp,
		Gender:     playerTable.Gender,
	}
}
