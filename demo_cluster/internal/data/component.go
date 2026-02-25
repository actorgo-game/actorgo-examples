package data

import (
	cdataconfig "github.com/actorgo-game/actorgo/components/data-config"
	cherryMapStructure "github.com/actorgo-game/actorgo/extend/mapstructure"
	"github.com/actorgo-game/examples/demo_cluster/internal/types"
)

var (
	AreaConfig       = &areaConfig{}
	AreaGroupConfig  = &areaGroupConfig{}
	AreaServerConfig = &areaServerConfig{}
	SdkConfig        = &sdkConfig{}
	CodeConfig       = &codeConfig{}
	PlayerInitConfig = &playerInitConfig{}
)

func New() *cdataconfig.Component {
	dataConfig := cdataconfig.New()
	dataConfig.Register(
		AreaConfig,
		AreaGroupConfig,
		AreaServerConfig,
		SdkConfig,
		CodeConfig,
		PlayerInitConfig,
	)
	return dataConfig
}

func DecodeData(input interface{}, output interface{}) error {
	return cherryMapStructure.HookDecode(
		input,
		output,
		"json",
		types.GetDecodeHooks(),
	)
}
