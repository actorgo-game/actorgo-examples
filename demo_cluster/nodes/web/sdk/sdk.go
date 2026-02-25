package sdk

import (
	cgin "github.com/actorgo-game/actorgo/components/gin"
	cerror "github.com/actorgo-game/actorgo/error"
	cstring "github.com/actorgo-game/actorgo/extend/string"
	cfacade "github.com/actorgo-game/actorgo/facade"
	"github.com/actorgo-game/examples/demo_cluster/internal/data"
)

// sdk平台类型
const (
	DevMode  int32 = 1 // 开发模式，注册开发帐号登陆(开发时使用)
	QuickSDK int32 = 2 // quick sdk
)

var (
	invokeMaps = make(map[int32]Invoke)
)

type (
	Invoke interface {
		SdkId() int32                                                // sdk id
		Login(config *data.SdkRow, params Params, callback Callback) // Login 登录验证接口
		PayCallback(config *data.SdkRow, c *cgin.Context)            // PayCallback 支付回调接口
	}

	Params map[string]string

	Callback func(code int32, result Params, error ...error)
)

func (p Params) GetInt(key string, defaultValue ...int) int {
	defVal := 0
	if len(defaultValue) > 0 {
		defVal = defaultValue[0]
	}

	val, found := p[key]
	if !found {
		return defVal
	}

	intVal, ok := cstring.ToInt(val)
	if ok {
		return intVal
	}

	return defVal
}

func (p Params) GetString(key string) (string, bool) {
	v, ok := p[key]
	return v, ok
}

func register(invoke Invoke) {
	invokeMaps[invoke.SdkId()] = invoke
}

func GetInvoke(sdkId int32) (invoke Invoke, error error) {
	invoke, found := invokeMaps[sdkId]
	if found == false {
		return nil, cerror.Errorf("[sdkId = %d] not found.", sdkId)
	}

	return invoke, nil
}

func Init(app cfacade.IApplication) {
	register(devSdk{app})
	register(quickSdk{})
}
