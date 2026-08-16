package main

import (
	"fmt"
	"sync/atomic"

	cfacade "github.com/actorgo-game/actorgo/facade"
	clog "github.com/actorgo-game/actorgo/logger"
	cactor "github.com/actorgo-game/actorgo/net/actor"
	"github.com/actorgo-game/actorgo/net/parser"
)

var nextUID atomic.Int64

type (
	actorRoom struct {
		cactor.Base
		server  *parser.Component
		userMap map[int64]*User
	}

	User struct {
		uid      int64
		nickname string
		balance  int64
		message  int
	}
)

func newActorRoom(server *parser.Component) *actorRoom {
	return &actorRoom{server: server}
}

func (*actorRoom) AliasID() string { return "room" }

func (p *actorRoom) OnInit() {
	p.userMap = make(map[int64]*User)
	p.Methods().Register(MethodLogin, p.login)
	p.Methods().Register(MethodSyncMessage, p.syncMessage)
	p.Methods().Register(MethodExit, p.exit)
}

func (p *actorRoom) login(ctx *cfacade.RequestContext, req *LoginRequest) (*LoginResponse, error) {
	if ctx == nil || ctx.Session == nil || ctx.Session.Sid == "" {
		return nil, fmt.Errorf("chat session is unavailable")
	}
	if ctx.Session.Uid > 0 {
		return &LoginResponse{Code: 0}, nil
	}

	uid := nextUID.Add(1)
	if err := p.server.Bind(ctx.Session.Sid, uid, map[string]string{"nickname": req.Nickname}); err != nil {
		return nil, fmt.Errorf("bind chat session: %w", err)
	}
	p.userMap[uid] = &User{uid: uid, nickname: req.Nickname, balance: 1000}

	clog.Debug("new chat session. sid=%s uid=%d nickname=%s", ctx.Session.Sid, uid, req.Nickname)
	p.broadcast(MethodNewUser, &NewUserBroadcast{Content: fmt.Sprintf("user join: %s", req.Nickname)})
	return &LoginResponse{Code: 0}, nil
}

func (p *actorRoom) syncMessage(ctx *cfacade.RequestContext, req *SyncMessage) error {
	if ctx == nil || ctx.Session == nil || ctx.Session.Uid < 1 {
		return fmt.Errorf("chat session is not logged in")
	}
	user, found := p.userMap[ctx.Session.Uid]
	if !found {
		return fmt.Errorf("chat user %d not found", ctx.Session.Uid)
	}

	user.message++
	user.balance--
	p.broadcast(MethodMessage, req)
	if err := p.server.NotifyUID(user.uid, MethodBalance, &UserBalanceResponse{CurrentBalance: user.balance}); err != nil {
		clog.Warn("notify balance failed. uid=%d err=%v", user.uid, err)
	}
	return nil
}

func (p *actorRoom) exit(_ *cfacade.RequestContext, req *Int64) error {
	if req.Value < 1 {
		return nil
	}
	clog.Debug("chat user exit. uid=%d", req.Value)
	delete(p.userMap, req.Value)
	return nil
}

func (p *actorRoom) broadcast(methodID uint32, payload any) {
	for uid := range p.userMap {
		if err := p.server.NotifyUID(uid, methodID, payload); err != nil {
			clog.Warn("chat broadcast failed. uid=%d methodID=%d err=%v", uid, methodID, err)
		}
	}
}
