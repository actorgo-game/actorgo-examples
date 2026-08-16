package ops

import (
	cfacade "github.com/actorgo-game/actorgo/facade"
	cactor "github.com/actorgo-game/actorgo/net/actor"
	"github.com/actorgo-game/examples/demo_cluster/internal/methodid"
	"github.com/actorgo-game/examples/demo_cluster/internal/pb"
)

var (
	pingReturn = &pb.Bool{Value: true}
)

type (
	ActorOps struct {
		cactor.Base
	}
)

func (p *ActorOps) AliasID() string {
	return "ops"
}

// OnInit 注册remote函数
func (p *ActorOps) OnInit() {
	p.Methods().Register(methodid.CenterPing, p.ping)
}

// ping 请求center是否响应
func (p *ActorOps) ping(_ *cfacade.RequestContext, _ *pb.None) (*pb.Bool, error) {
	return pingReturn, nil
}
