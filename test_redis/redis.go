package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	clogger "github.com/actorgo-game/actorgo/logger"
	"github.com/go-redis/redis/v8"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()

	data := struct {
		UserName string `json:"user_name"`
		Address  string `json:"address"`
	}{
		UserName: "tom",
		Address:  "china shenzhen",
	}

	jsonData, err := json.Marshal(&data)
	if err != nil {
		return
	}

	cmd := rdb.Set(ctx, "data_config:test", jsonData, 10*time.Hour)
	clogger.Debug(cmd.Val(), cmd.Err())

	keysVal := rdb.Keys(ctx, "node_list*")
	clogger.Debug("val[%v] err[%v]", keysVal.Val(), keysVal.Err())

	scanVal := rdb.Scan(ctx, 0, "node_list*", 0)
	clogger.Debug("%v", scanVal)

	go func() {
		subscribe := rdb.Subscribe(ctx, "aaaa")
		defer subscribe.Close()

		ch := subscribe.Channel()
		for msg := range ch {
			fmt.Println(msg.Channel, msg.Payload)
		}
	}()

	time.Sleep(1 * time.Hour)

}
