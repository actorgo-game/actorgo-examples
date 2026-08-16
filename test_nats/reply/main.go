package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/actorgo-game/examples/test_nats/internal/natsutil"
	"github.com/nats-io/nats.go"
)

func main() {
	config := natsutil.BindFlags()
	flag.Parse()

	connection, err := config.Connect("actorgo-example-reply")
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	if _, err = connection.Subscribe(config.Subject, func(message *nats.Msg) {
		if len(message.Data) != 8 {
			log.Printf("received invalid %d-byte request", len(message.Data))
			return
		}
		sentAt := int64(binary.BigEndian.Uint64(message.Data))
		latency := time.Now().UnixMicro() - sentAt
		if respondErr := message.Respond([]byte(fmt.Sprintf("latency=%dµs", latency))); respondErr != nil {
			log.Printf("respond: %v", respondErr)
		}
	}); err != nil {
		log.Fatal(err)
	}
	if err = connection.Flush(); err != nil {
		log.Fatal(err)
	}
	log.Printf("reply service subscribed to %s; press Ctrl+C to stop", config.Subject)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	if err = connection.Drain(); err != nil {
		log.Printf("drain connection: %v", err)
	}
}
