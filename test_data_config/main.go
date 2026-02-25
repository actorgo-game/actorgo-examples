package main

import (
	"time"

	"github.com/actorgo-game/actorgo"
	cdataconfig "github.com/actorgo-game/actorgo/components/data-config"
	cfacade "github.com/actorgo-game/actorgo/facade"
	clog "github.com/actorgo-game/actorgo/logger"
)

func main() {
	testApp := actorgo.NewApp(
		"../config/test.json",
		"game-1",
		false,
		actorgo.Standalone,
	)

	dataConfig := cdataconfig.New()
	dataConfig.Register(&DropList, &DropOne)
	testApp.Register(dataConfig)

	go func(testApp *actorgo.Application) {
		//120秒后退出应用
		getDropConfig(testApp)
		testApp.Shutdown()
	}(testApp)

	testApp.Startup()
}

func getDropConfig(_ cfacade.IApplication) {
	time.Sleep(5 * time.Second)

	for {
		clog.Info("DropOneConfig %p, %v", &DropOne, DropOne)

		x1 := DropList.Get(1011)
		clog.Info("DropConfig %p, %v", x1, x1)

		itemTypeList := DropList.GetItemTypeList(3)
		clog.Info("DropConfig %p, %v", itemTypeList, itemTypeList)

		time.Sleep(500 * time.Millisecond)
	}
}
