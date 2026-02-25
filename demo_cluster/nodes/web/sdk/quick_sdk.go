package sdk

import (
	cgin "github.com/actorgo-game/actorgo/components/gin"
	cerror "github.com/actorgo-game/actorgo/error"
	chttp "github.com/actorgo-game/actorgo/extend/http"
	cstring "github.com/actorgo-game/actorgo/extend/string"
	clogger "github.com/actorgo-game/actorgo/logger"
	"github.com/actorgo-game/examples/demo_cluster/internal/code"
	"github.com/actorgo-game/examples/demo_cluster/internal/data"
	sessionKey "github.com/actorgo-game/examples/demo_cluster/internal/session_key"
)

type (
	quickSdk struct {
	}
)

func (quickSdk) SdkId() int32 {
	return QuickSDK
}

func (quickSdk) Login(config *data.SdkRow, params Params, callback Callback) {
	token, found := params.GetString("token")
	if found == false || cstring.IsBlank(token) {
		err := cerror.Error("token is nil")
		callback(code.LoginError, nil, err)
		return
	}

	uid, found := params.GetString("uid")
	if found == false || cstring.IsBlank(uid) {
		err := cerror.Error("uid is nil")
		callback(code.LoginError, nil, err)
		return
	}

	rspBody, _, err := chttp.GET(config.LoginURL(), map[string]string{
		"token":        token,
		"uid":          uid,
		"product_code": config.GetString("productCode"),
	})

	if err != nil || rspBody == nil {
		callback(code.LoginError, nil, err)
		return
	}

	var result = string(rspBody)
	if result != "1" {
		clogger.Warn("quick sdk login fail. rsp =%s", rspBody)
		callback(code.LoginError, nil, err)
		return
	}

	callback(code.OK, map[string]string{
		sessionKey.OpenID: uid, //返回 quick的uid做为 open id
	})
}

func (s quickSdk) PayCallback(config *data.SdkRow, c *cgin.Context) {
	// TODO 这里实现quick sdk 支付回调的逻辑
	c.RenderHTML("FAIL")
}
