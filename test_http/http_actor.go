package main

import (
	"strings"
	"unicode/utf8"

	cfacade "github.com/actorgo-game/actorgo/facade"
	clog "github.com/actorgo-game/actorgo/logger"
	cactor "github.com/actorgo-game/actorgo/net/actor"
	cmethod "github.com/actorgo-game/actorgo/net/method"
	cproto "github.com/actorgo-game/actorgo/net/proto"
)

type HTTPActor struct {
	cactor.Base
}

func (actor *HTTPActor) AliasID() string {
	return "http-example"
}

func (actor *HTTPActor) OnInit() {
	actor.Methods().Register(MethodHTTPHealth, actor.health)
	actor.Methods().Register(MethodHTTPHello, actor.hello)
	actor.Methods().Register(MethodHTTPEcho, actor.echo)
	actor.Methods().Register(MethodHTTPChildHello, actor.childHello)
	clog.Info(
		"[HTTPActor] methods registered. [actorPath = %s, methodIDs = %d,%d,%d,%d]",
		actor.PathString(),
		MethodHTTPHealth,
		MethodHTTPHello,
		MethodHTTPEcho,
		MethodHTTPChildHello,
	)
}

func (actor *HTTPActor) OnFindChild(message *cfacade.Message) (cfacade.IActor, bool) {
	childID := message.TargetPath().ChildID
	child, err := actor.Child().Create(childID, new(HTTPChildActor))
	if err != nil {
		clog.Warn("[HTTPActor] create child failed. [childID = %s, methodID = %d, err = %v]", childID, message.MethodID, err)
		return nil, false
	}
	clog.Info("[HTTPActor] child is ready. [childID = %s, methodID = %d]", childID, message.MethodID)
	return child, true
}

func (actor *HTTPActor) health(ctx *cfacade.RequestContext, _ *HealthRequest) (*HealthResponse, error) {
	clog.Debug("[HTTPActor.health] request handled. [methodID = %d, requestID = %d]", MethodHTTPHealth, ctx.RequestID)
	return &HealthResponse{Status: "up"}, nil
}

func (actor *HTTPActor) hello(ctx *cfacade.RequestContext, request *HelloRequest) (*HelloResponse, error) {
	name := strings.TrimSpace(request.Name)
	if name == "" {
		name = "ActorGo"
	}
	clog.Info("[HTTPActor.hello] request handled. [methodID = %d, requestID = %d, name = %s]", MethodHTTPHello, ctx.RequestID, name)
	return &HelloResponse{
		Name:     name,
		Greeting: "Hello, " + name + "!",
	}, nil
}

func (actor *HTTPActor) echo(ctx *cfacade.RequestContext, request *EchoRequest) (*EchoResponse, error) {
	message := strings.TrimSpace(request.Message)
	if message == "" {
		clog.Warn("[HTTPActor.echo] invalid request. [methodID = %d, requestID = %d, err = message is required]", MethodHTTPEcho, ctx.RequestID)
		return nil, invalidArgument("message is required")
	}
	messageLength := utf8.RuneCountInString(message)
	clog.Info("[HTTPActor.echo] request handled. [methodID = %d, requestID = %d, messageLength = %d]", MethodHTTPEcho, ctx.RequestID, messageLength)
	return &EchoResponse{
		Message: message,
		Length:  messageLength,
	}, nil
}

func (actor *HTTPActor) childHello(ctx *cfacade.RequestContext, request *ChildHelloRequest) (*ChildHelloResponse, error) {
	childID := strings.TrimSpace(request.ChildID)
	if childID == "" {
		clog.Warn("[HTTPActor.childHello] invalid request. [methodID = %d, requestID = %d, err = child_id is required]", MethodHTTPChildHello, ctx.RequestID)
		return nil, invalidArgument("child_id is required")
	}

	clog.Info("[HTTPActor.childHello] invoking child. [methodID = %d, requestID = %d, childID = %s]", MethodHTTPChildHello, ctx.RequestID, childID)
	result := actor.InvokeChild(ctx, childID, methodChildHello, &childHelloRequest{
		Message: request.Message,
	})
	if !result.OK() {
		clog.Warn("[HTTPActor.childHello] child invocation failed. [requestID = %d, childID = %s, code = %d, message = %s]", ctx.RequestID, childID, result.Code, result.Message)
		return nil, &cmethod.InvokeError{
			Code:    cproto.StatusCode(result.Code),
			Message: result.Message,
		}
	}
	response, ok := result.Payload.(*ChildHelloResponse)
	if !ok {
		clog.Error("[HTTPActor.childHello] unexpected child response. [requestID = %d, childID = %s, payloadType = %T]", ctx.RequestID, childID, result.Payload)
		return nil, &cmethod.InvokeError{
			Code:    cproto.StatusCode_STATUS_INTERNAL,
			Message: "unexpected child response",
		}
	}
	clog.Info("[HTTPActor.childHello] child invocation succeeded. [requestID = %d, childID = %s, actorPath = %s]", ctx.RequestID, childID, response.ActorPath)
	return response, nil
}
