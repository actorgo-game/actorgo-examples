package data

import (
	cerror "github.com/actorgo-game/actorgo/error"
	clogger "github.com/actorgo-game/actorgo/logger"
)

type (
	AreaServerRow struct {
		ServerId   int32  `json:"serverId"`   // 游戏服id
		ServerName string `json:"serverName"` // 游戏服名称
		AreaId     int32  `json:"areaId"`     // 游戏服所属的区id
		Status     int32  `json:"status"`     // 游戏服状态
	}

	// 游戏区分组
	areaServerConfig struct {
		maps map[int32]*AreaServerRow
	}
)

// Name 根据名称读取 ./config/data/areaServerConfig.json文件
func (p *areaServerConfig) Name() string {
	return "areaServerConfig"
}

func (p *areaServerConfig) Init() {
	p.maps = make(map[int32]*AreaServerRow)
}

func (p *areaServerConfig) OnLoad(maps interface{}, _ bool) (int, error) {
	list, ok := maps.([]interface{})
	if !ok {
		return 0, cerror.Error("maps convert to []interface{} error.")
	}

	loadMaps := make(map[int32]*AreaServerRow)
	for index, data := range list {
		loadConfig := &AreaServerRow{}
		err := DecodeData(data, loadConfig)
		if err != nil {
			clogger.Warn("decode error. [row = %d, %v], err = %s", index+1, loadConfig, err)
			continue
		}

		loadMaps[loadConfig.ServerId] = loadConfig
	}

	p.maps = loadMaps

	return len(list), nil
}

func (p *areaServerConfig) OnAfterLoad(_ bool) {
}

func (p *areaServerConfig) Get(pk int32) (*AreaServerRow, bool) {
	i, found := p.maps[pk]
	return i, found
}

func (p *areaServerConfig) Contain(pk int32) bool {
	_, found := p.Get(pk)
	return found
}

func (p *areaServerConfig) ListWithAreaId(areaId int32) []*AreaServerRow {
	var list []*AreaServerRow

	for _, row := range p.maps {
		if row.AreaId == areaId {
			list = append(list, row)
		}
	}

	return list
}
