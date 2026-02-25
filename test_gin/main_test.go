package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	chttp "github.com/actorgo-game/actorgo/extend/http"
)

func TestControllerMaxConnect(t *testing.T) {

	for i := 0; i < 100; i++ {
		go func(i int) {
			result, rsp, _ := chttp.GET("http://127.0.0.1:10820")

			if rsp != nil && rsp.StatusCode != http.StatusOK {
				fmt.Printf("index = %d, result = %s, code = %v\n", i, result, rsp.StatusCode)
			}
		}(i)
	}

	time.Sleep(1 * time.Hour)
}
