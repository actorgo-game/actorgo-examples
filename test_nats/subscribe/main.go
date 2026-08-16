package main

import (
	"context"
	"encoding/binary"
	"flag"
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

	connection, err := config.Connect("actorgo-example-subscribe")
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	for _, subscriberName := range []string{"subscriber-1", "subscriber-2"} {
		name := subscriberName
		if _, err = connection.Subscribe(config.Subject, func(message *nats.Msg) {
			if len(message.Data) != 8 {
				log.Printf("%s received invalid %d-byte message", name, len(message.Data))
				return
			}
			sentAt := int64(binary.BigEndian.Uint64(message.Data))
			log.Printf("%s received message, latency=%dµs", name, time.Now().UnixMicro()-sentAt)
		}); err != nil {
			log.Fatal(err)
		}
	}
	if err = connection.Flush(); err != nil {
		log.Fatal(err)
	}
	log.Printf("subscribed to %s; press Ctrl+C to stop", config.Subject)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	if err = connection.Drain(); err != nil {
		log.Printf("drain connection: %v", err)
	}
}
