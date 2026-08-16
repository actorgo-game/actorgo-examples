package rpcCenter

import (
	cfacade "github.com/actorgo-game/actorgo/facade"
	clog "github.com/actorgo-game/actorgo/logger"
	"github.com/actorgo-game/examples/demo_cluster/internal/code"
	"github.com/actorgo-game/examples/demo_cluster/internal/constant"
	"github.com/actorgo-game/examples/demo_cluster/internal/methodid"
	"github.com/actorgo-game/examples/demo_cluster/internal/pb"
)

func Ping(app cfacade.IApplication) bool {
	nodeID := GetCenterNodeID(app)
	if nodeID == "" {
		return false
	}

	rsp := &pb.Bool{}
	result := app.ActorSystem().InvokeNode(newContext(), nodeID, methodid.CenterPing, &pb.None{})
	return decodeResult(app, result, rsp) == code.OK && rsp.Value
}

func RegisterDevAccount(app cfacade.IApplication, accountName, password, ip string) int32 {
	req := &pb.DevRegister{AccountName: accountName, Password: password, Ip: ip}
	rsp := &pb.Int32{}
	result := app.ActorSystem().InvokeNode(newContext(), GetCenterNodeID(app), methodid.CenterRegisterDevAccount, req)
	if errCode := decodeResult(app, result, rsp); code.IsFail(errCode) {
		clog.Warn("[RegisterDevAccount] accountName = %s, errCode = %v", accountName, errCode)
		return errCode
	}
	return rsp.Value
}

func GetDevAccount(app cfacade.IApplication, accountName, password string) int64 {
	req := &pb.DevRegister{AccountName: accountName, Password: password}
	rsp := &pb.Int64{}
	result := app.ActorSystem().InvokeNode(newContext(), GetCenterNodeID(app), methodid.CenterGetDevAccount, req)
	if errCode := decodeResult(app, result, rsp); code.IsFail(errCode) {
		clog.Warn("[GetDevAccount] accountName = %s, errCode = %v", accountName, errCode)
		return 0
	}
	return rsp.Value
}

func GetUID(app cfacade.IApplication, sdkID, pid int32, openID string) (cfacade.UID, int32) {
	return GetUIDContext(newContext(), app, sdkID, pid, openID)
}

func GetUIDContext(ctx *cfacade.RequestContext, app cfacade.IApplication, sdkID, pid int32, openID string) (cfacade.UID, int32) {
	req := &pb.User{SdkId: sdkID, Pid: pid, OpenId: openID}
	rsp := &pb.Int64{}
	result := app.ActorSystem().InvokeNode(ctx, GetCenterNodeID(app), methodid.CenterGetUID, req)
	if errCode := decodeResult(app, result, rsp); code.IsFail(errCode) {
		clog.Warn("[GetUID] errCode = %v", errCode)
		return 0, errCode
	}
	return rsp.Value, code.OK
}

func GetCenterNodeID(app cfacade.IApplication) string {
	list := app.Discovery().ListByType(constant.CenterNodeType)
	if len(list) > 0 {
		return list[0].GetNodeID()
	}
	return ""
}

func newContext() *cfacade.RequestContext {
	ctx := cfacade.NewRequestContext(nil)
	ctx.Codec = cfacade.CodecProtobuf
	return ctx
}

func decodeResult(app cfacade.IApplication, result *cfacade.InvokeResult, response any) int32 {
	if result == nil {
		return code.Error
	}
	if !result.OK() {
		return result.Code
	}
	if err := result.Decode(app.BodyCodecs(), cfacade.CodecProtobuf, response); err != nil {
		clog.Warn("decode actor response fail. err = %v", err)
		return code.Error
	}
	return code.OK
}
