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
	count := flag.Int("count", 10, "number of requests")
	interval := flag.Duration("interval", time.Second, "interval between requests")
	timeout := flag.Duration("timeout", 2*time.Second, "request timeout")
	flag.Parse()
	if *count < 1 {
		log.Fatal("count must be positive")
	}

	connection, err := config.Connect("actorgo-example-request")
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	for index := 0; index < *count; index++ {
		response, requestErr := connection.Request(
			config.Subject,
			int64ToBytes(time.Now().UnixMicro()),
			*timeout,
		)
		if requestErr != nil {
			log.Printf("request %d failed: %v", index+1, requestErr)
		} else {
			log.Printf("request %d response: %s", index+1, response.Data)
		}
		if index+1 < *count {
			time.Sleep(*interval)
		}
	}
	if err = connection.Drain(); err != nil {
		log.Fatal(err)
	}
}
