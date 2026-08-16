package code

import (
	cmethod "github.com/actorgo-game/actorgo/net/method"
	cproto "github.com/actorgo-game/actorgo/net/proto"
)

func NewInvokeError(value int32) error {
	return &cmethod.InvokeError{
		Code:    cproto.StatusCode(value),
		Message: GetMessage(value),
	}
}
