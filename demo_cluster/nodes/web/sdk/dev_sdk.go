package sdk

import (
	cgin "github.com/actorgo-game/actorgo/components/gin"
	cerror "github.com/actorgo-game/actorgo/error"
	cstring "github.com/actorgo-game/actorgo/extend/string"
	cfacade "github.com/actorgo-game/actorgo/facade"
	"github.com/actorgo-game/examples/demo_cluster/internal/code"
	"github.com/actorgo-game/examples/demo_cluster/internal/data"
	rpcCenter "github.com/actorgo-game/examples/demo_cluster/internal/rpc/center"
)

type devSdk struct {
	app cfacade.IApplication
}

func (devSdk) SdkId() int32 {
	return DevMode
}

func (p devSdk) Login(_ *data.SdkRow, params Params, callback Callback) {
	accountName, _ := params.GetString("account")
	password, _ := params.GetString("password")

	if accountName == "" || password == "" {
		err := cerror.Error("account or password params is empty.")
		callback(code.LoginError, nil, err)
		return
	}

	accountId := rpcCenter.GetDevAccount(p.app, accountName, password)
	if accountId < 1 {
		callback(code.LoginError, nil)
		return
	}

	callback(code.OK, map[string]string{
		"open_id": cstring.ToString(accountId),
	})
}

func (devSdk) PayCallback(_ *data.SdkRow, _ *cgin.Context) {
}
