package main

import (
	cmethod "github.com/actorgo-game/actorgo/net/method"
	cproto "github.com/actorgo-game/actorgo/net/proto"
)

const (
	MethodHTTPHealth     uint32 = 5001
	MethodHTTPHello      uint32 = 5002
	MethodHTTPEcho       uint32 = 5003
	MethodHTTPChildHello uint32 = 5004

	methodChildHello uint32 = 5101
)

type HealthRequest struct{}

type HealthResponse struct {
	Status string `json:"status"`
}

type HelloRequest struct {
	Name string `json:"name"`
}

type HelloResponse struct {
	Name     string `json:"name"`
	Greeting string `json:"greeting"`
}

type EchoRequest struct {
	Message string `json:"message"`
}

type EchoResponse struct {
	Message string `json:"message"`
	Length  int    `json:"length"`
}

type ChildHelloRequest struct {
	ChildID string `json:"child_id"`
	Message string `json:"message"`
}

type ChildHelloResponse struct {
	ChildID   string `json:"child_id"`
	ActorPath string `json:"actor_path"`
	Message   string `json:"message"`
}

type childHelloRequest struct {
	Message string `json:"message"`
}

func invalidArgument(message string) error {
	return &cmethod.InvokeError{
		Code:    cproto.StatusCode_STATUS_INVALID_ARGUMENT,
		Message: message,
	}
}
