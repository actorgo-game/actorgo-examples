package account

import (
	"strings"

	cfacade "github.com/actorgo-game/actorgo/facade"
	cactor "github.com/actorgo-game/actorgo/net/actor"
	"github.com/actorgo-game/examples/demo_cluster/internal/code"
	"github.com/actorgo-game/examples/demo_cluster/internal/methodid"
	"github.com/actorgo-game/examples/demo_cluster/internal/pb"
	"github.com/actorgo-game/examples/demo_cluster/nodes/center/db"
)

type (
	ActorAccount struct {
		cactor.Base
	}
)

func (p *ActorAccount) AliasID() string {
	return "account"
}

// OnInit center为后端节点，不直接与客户端通信，所以了一些remote函数，供RPC调用
func (p *ActorAccount) OnInit() {
	p.Methods().Register(methodid.CenterRegisterDevAccount, p.registerDevAccount)
	p.Methods().Register(methodid.CenterGetDevAccount, p.getDevAccount)
	p.Methods().Register(methodid.CenterGetUID, p.getUID)
}

// registerDevAccount 注册开发者帐号
func (p *ActorAccount) registerDevAccount(_ *cfacade.RequestContext, req *pb.DevRegister) (*pb.Int32, error) {
	accountName := req.AccountName
	password := req.Password

	if strings.TrimSpace(accountName) == "" || strings.TrimSpace(password) == "" {
		return &pb.Int32{Value: code.LoginError}, nil
	}

	if len(accountName) < 3 || len(accountName) > 18 {
		return &pb.Int32{Value: code.LoginError}, nil
	}

	if len(password) < 3 || len(password) > 18 {
		return &pb.Int32{Value: code.LoginError}, nil
	}

	return &pb.Int32{Value: db.DevAccountRegister(accountName, password, req.Ip)}, nil
}

// getDevAccount 根据帐号名获取开发者帐号表
func (p *ActorAccount) getDevAccount(_ *cfacade.RequestContext, req *pb.DevRegister) (*pb.Int64, error) {
	accountName := req.AccountName
	password := req.Password

	devAccount, _ := db.DevAccountWithName(accountName)
	if devAccount == nil || devAccount.Password != password {
		return nil, code.NewInvokeError(code.AccountAuthFail)
	}

	return &pb.Int64{Value: devAccount.AccountId}, nil
}

// getUID 获取uid
func (p *ActorAccount) getUID(_ *cfacade.RequestContext, req *pb.User) (*pb.Int64, error) {
	uid, ok := db.BindUID(req.SdkId, req.Pid, req.OpenId)
	if uid == 0 || ok == false {
		return nil, code.NewInvokeError(code.AccountBindFail)
	}

	return &pb.Int64{Value: uid}, nil
}
