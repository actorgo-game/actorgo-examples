package main

import (
	"strings"
	"unicode/utf8"

	cfacade "github.com/actorgo-game/actorgo/facade"
	clog "github.com/actorgo-game/actorgo/logger"
	cactor "github.com/actorgo-game/actorgo/net/actor"
)

type HTTPChildActor struct {
	cactor.Base
}

func (actor *HTTPChildActor) OnInit() {
	actor.Methods().Register(methodChildHello, actor.hello)
	clog.Info("[HTTPChildActor] method registered. [actorPath = %s, methodID = %d]", actor.PathString(), methodChildHello)
}

func (actor *HTTPChildActor) hello(ctx *cfacade.RequestContext, request *childHelloRequest) (*ChildHelloResponse, error) {
	message := strings.TrimSpace(request.Message)
	if message == "" {
		clog.Warn("[HTTPChildActor.hello] invalid request. [requestID = %d, actorPath = %s, err = message is required]", ctx.RequestID, actor.PathString())
		return nil, invalidArgument("message is required")
	}
	clog.Info(
		"[HTTPChildActor.hello] request handled. [methodID = %d, requestID = %d, actorPath = %s, messageLength = %d]",
		methodChildHello,
		ctx.RequestID,
		actor.PathString(),
		utf8.RuneCountInString(message),
	)
	return &ChildHelloResponse{
		ChildID:   actor.ActorID(),
		ActorPath: actor.PathString(),
		Message:   "child received: " + message,
	}, nil
}
