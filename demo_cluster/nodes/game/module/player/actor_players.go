package player

import (
	"time"

	cfacade "github.com/actorgo-game/actorgo/facade"
	clog "github.com/actorgo-game/actorgo/logger"
	"github.com/actorgo-game/actorgo/net/parser/pomelo"
	"github.com/actorgo-game/examples/demo_cluster/internal/event"
	"github.com/actorgo-game/examples/demo_cluster/nodes/game/module/online"
)

type (
	// ActorPlayers 玩家总管理actor
	ActorPlayers struct {
		pomelo.ActorBase
		childExitTime time.Duration
	}
)

func (p *ActorPlayers) AliasID() string {
	return "player"
}

func (p *ActorPlayers) OnInit() {
	p.childExitTime = time.Minute * 30

	// 注册角色登陆事件
	p.Event().Register(event.PlayerLoginKey, p.onLoginEvent)
	p.Event().Register(event.PlayerLogoutKey, p.onLogoutEvent)
	p.Event().Register(event.PlayerCreateKey, p.onPlayerCreateEvent)
}

func (p *ActorPlayers) OnFindChild(msg *cfacade.Message) (cfacade.IActor, bool) {
	// 动态创建 player child actor
	childID := msg.TargetPath().ChildID
	childActor, err := p.Child().Create(childID, &actorPlayer{
		isOnline: false,
	})

	if err != nil {
		return nil, false
	}

	return childActor, true
}

// onLoginEvent 玩家登陆事件处理
func (p *ActorPlayers) onLoginEvent(e cfacade.IEventData) {
	evt, ok := e.(*event.PlayerLogin)
	if ok == false {
		return
	}

	clog.Info("[PlayerLoginEvent] [playerId = %d, onlineCount = %d]",
		evt.PlayerId,
		online.Count(),
	)
}

// onLoginEvent 玩家登出事件处理
func (p *ActorPlayers) onLogoutEvent(e cfacade.IEventData) {
	evt, ok := e.(*event.PlayerLogout)
	if !ok {
		return
	}

	clog.Info("[PlayerLogoutEvent] [playerId = %d, onlineCount = %d]",
		evt.PlayerId,
		online.Count(),
	)
}

// onPlayerCreateEvent 玩家创建事件
func (p *ActorPlayers) onPlayerCreateEvent(e cfacade.IEventData) {
	evt, ok := e.(*event.PlayerCreate)
	if !ok {
		return
	}

	clog.Info("[PlayerCreateEvent] [%+v]", evt)
}

func (p *ActorPlayers) OnStop() {
	clog.Info("onlineCount = %d", online.Count())
}
