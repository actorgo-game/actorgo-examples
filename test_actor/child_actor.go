package main

import (
	"fmt"

	cfacade "github.com/actorgo-game/actorgo/facade"
	cactor "github.com/actorgo-game/actorgo/net/actor"
)

type childActor struct {
	cactor.Base
}

func (p *childActor) OnInit() {
	fmt.Println("[childActor] Execute OnInit()")

	p.Methods().Register(childHelloMethodID, p.hello)
}

func (p *childActor) hello(_ *cfacade.RequestContext, _ *helloRequest) (*helloResponse, error) {
	text := "[childActor] Call hello()"
	fmt.Println(text)
	return &helloResponse{Text: text}, nil
}

func (*childActor) OnStop() {
}
