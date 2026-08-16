package main

import (
	"encoding/binary"
	"flag"
	"log"
	"time"

	"github.com/actorgo-game/examples/test_nats/internal/natsutil"
)

func int64ToBytes(value int64) []byte {
	buffer := make([]byte, 8)
	binary.BigEndian.PutUint64(buffer, uint64(value))
	return buffer
}

func main() {
	config := natsutil.BindFlags()
	count := flag.Int("count", 10, "number of messages to publish")
	interval := flag.Duration("interval", time.Second, "interval between messages")
	flag.Parse()
	if *count < 1 {
		log.Fatal("count must be positive")
	}

	connection, err := config.Connect("actorgo-example-publish")
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	for index := 0; index < *count; index++ {
		if err = connection.Publish(config.Subject, int64ToBytes(time.Now().UnixMicro())); err != nil {
			log.Fatal(err)
		}
		if index+1 < *count {
			time.Sleep(*interval)
		}
	}
	if err = connection.Flush(); err != nil {
		log.Fatal(err)
	}
	log.Printf("published %d message(s) to %s", *count, config.Subject)
	if err = connection.Drain(); err != nil {
		log.Fatal(err)
	}
}
