package main

import (
	"context"
	"fmt"

	cfacade "github.com/actorgo-game/actorgo/facade"
	cactor "github.com/actorgo-game/actorgo/net/actor"
)

const childHelloMethodID uint32 = 1001

type helloRequest struct{}

type helloResponse struct {
	Text string `json:"text"`
}

type actor struct {
	cactor.Base
}

func (*actor) AliasID() string {
	return "parentActor"
}

func (p *actor) OnInit() {
	fmt.Println("[actor] Execute OnInit()")

	childActorID := "1"
	if _, err := p.Child().Create(childActorID, &childActor{}); err != nil {
		panic(err)
	}

	ctx := cfacade.NewRequestContext(context.Background())
	result := p.InvokeChild(ctx, childActorID, childHelloMethodID, &helloRequest{})
	if !result.OK() {
		panic(fmt.Sprintf("invoke child failed: code=%d message=%s", result.Code, result.Message))
	}
	reply, ok := result.Payload.(*helloResponse)
	if !ok {
		panic(fmt.Sprintf("unexpected child response: %T", result.Payload))
	}
	fmt.Println(reply.Text)
}

func (*actor) OnStop() {
}
