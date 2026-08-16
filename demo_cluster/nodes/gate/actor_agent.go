package gate

import (
	"fmt"
	"maps"

	cstring "github.com/actorgo-game/actorgo/extend/string"
	cfacade "github.com/actorgo-game/actorgo/facade"
	clog "github.com/actorgo-game/actorgo/logger"
	cactor "github.com/actorgo-game/actorgo/net/actor"
	cmethod "github.com/actorgo-game/actorgo/net/method"
	"github.com/actorgo-game/actorgo/net/parser"
	cproto "github.com/actorgo-game/actorgo/net/proto"
	"github.com/actorgo-game/examples/demo_cluster/internal/code"
	"github.com/actorgo-game/examples/demo_cluster/internal/data"
	"github.com/actorgo-game/examples/demo_cluster/internal/methodid"
	"github.com/actorgo-game/examples/demo_cluster/internal/pb"
	rpcCenter "github.com/actorgo-game/examples/demo_cluster/internal/rpc/center"
	sessionKey "github.com/actorgo-game/examples/demo_cluster/internal/session_key"
	"github.com/actorgo-game/examples/demo_cluster/internal/token"
)

type ActorGate struct {
	cactor.Base
	server *parser.Component
}

func NewActorGate(server *parser.Component) *ActorGate {
	return &ActorGate{server: server}
}

func (p *ActorGate) AliasID() string {
	return "user"
}

func (p *ActorGate) OnInit() {
	p.Methods().Register(methodid.GateLogin, p.login)
	p.Methods().Register(methodid.GamePlayerSelect, p.playerSelect)
	p.Methods().Register(methodid.GamePlayerCreate, p.playerCreate)
	p.Methods().Register(methodid.GamePlayerEnter, p.playerEnter)
}

func (p *ActorGate) login(ctx *cfacade.RequestContext, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if ctx == nil || ctx.Session == nil || ctx.Session.Sid == "" {
		return nil, code.NewInvokeError(code.LoginError)
	}
	userToken, errCode := p.validateToken(req.Token)
	if code.IsFail(errCode) {
		return nil, code.NewInvokeError(errCode)
	}
	sdkRow := data.SdkConfig.Get(userToken.PID)
	if sdkRow == nil {
		return nil, code.NewInvokeError(code.PIDError)
	}
	serverRow, found := data.AreaServerConfig.Get(req.ServerId)
	if !found || serverRow.NodeID == "" || serverRow.Status != 1 {
		return nil, code.NewInvokeError(code.LoginError)
	}

	uid, errCode := rpcCenter.GetUIDContext(ctx, p.App(), sdkRow.SdkId, userToken.PID, userToken.OpenID)
	if uid == 0 || code.IsFail(errCode) {
		return nil, code.NewInvokeError(code.AccountBindFail)
	}
	if oldConnection, ok := p.server.Connections().GetUID(uid); ok && oldConnection.ID() != ctx.Session.Sid {
		_ = oldConnection.Kick(code.PlayerDuplicateLogin, code.GetMessage(code.PlayerDuplicateLogin), true)
	}

	sessionData := map[string]string{
		sessionKey.ServerID:   cstring.ToString(req.ServerId),
		sessionKey.GameNodeID: serverRow.NodeID,
		sessionKey.PID:        cstring.ToString(userToken.PID),
		sessionKey.OpenID:     userToken.OpenID,
	}
	if err := p.server.Bind(ctx.Session.Sid, uid, sessionData); err != nil {
		clog.Warn("bind gate session fail. uid = %d, err = %v", uid, err)
		return nil, code.NewInvokeError(code.AccountBindFail)
	}

	return &pb.LoginResponse{Uid: uid, Pid: userToken.PID, OpenId: userToken.OpenID}, nil
}

func (p *ActorGate) playerSelect(ctx *cfacade.RequestContext, req *pb.None) (*pb.PlayerSelectResponse, error) {
	rsp := &pb.PlayerSelectResponse{}
	if err := p.invokeGame(ctx, methodid.GamePlayerSelect, req, rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}

func (p *ActorGate) playerCreate(ctx *cfacade.RequestContext, req *pb.PlayerCreateRequest) (*pb.PlayerCreateResponse, error) {
	rsp := &pb.PlayerCreateResponse{}
	if err := p.invokeGame(ctx, methodid.GamePlayerCreate, req, rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}

func (p *ActorGate) playerEnter(ctx *cfacade.RequestContext, req *pb.Int64) (*pb.PlayerEnterResponse, error) {
	rsp := &pb.PlayerEnterResponse{}
	if err := p.invokeGame(ctx, methodid.GamePlayerEnter, req, rsp); err != nil {
		return nil, err
	}

	sessionData := maps.Clone(ctx.Session.Data)
	if sessionData == nil {
		sessionData = make(map[string]string)
	}
	sessionData[sessionKey.PlayerID] = cstring.ToString(req.Value)
	if err := p.server.Bind(ctx.Session.Sid, ctx.Session.Uid, sessionData); err != nil {
		return nil, fmt.Errorf("update gate session: %w", err)
	}
	return rsp, nil
}

func (p *ActorGate) invokeGame(ctx *cfacade.RequestContext, methodID uint32, request, response any) error {
	if ctx == nil || ctx.Session == nil || ctx.Session.Uid < 1 {
		return code.NewInvokeError(code.PlayerDenyLogin)
	}
	gameNodeID := ctx.Session.GetString(sessionKey.GameNodeID)
	if gameNodeID == "" {
		return code.NewInvokeError(code.LoginError)
	}
	target := cfacade.NewChildPath(gameNodeID, "player", ctx.Session.Uid)
	result := p.App().ActorSystem().InvokeTarget(ctx, target, methodID, request)
	if result == nil {
		return fmt.Errorf("game actor returned nil result")
	}
	if !result.OK() {
		return &cmethod.InvokeError{Code: cproto.StatusCode(result.Code), Message: result.Message}
	}
	return result.Decode(p.App().BodyCodecs(), cfacade.CodecProtobuf, response)
}

func (p *ActorGate) validateToken(base64Token string) (*token.Token, int32) {
	userToken, ok := token.DecodeToken(base64Token)
	if !ok {
		return nil, code.AccountTokenValidateFail
	}
	platformRow := data.SdkConfig.Get(userToken.PID)
	if platformRow == nil {
		return nil, code.PIDError
	}
	statusCode, ok := token.Validate(userToken, platformRow.Salt)
	if !ok {
		return nil, statusCode
	}
	return userToken, code.OK
}
