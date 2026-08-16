package natsutil

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

const DefaultSubject = "actorgo.examples.game.10001"

// Config contains the connection flags shared by all NATS examples.
type Config struct {
	URL      string
	User     string
	Password string
	Subject  string
}

// BindFlags registers the common command-line flags. Environment variables
// make it possible to override credentials without changing the examples.
func BindFlags() *Config {
	config := &Config{}
	flag.StringVar(&config.URL, "url", envOrDefault("NATS_URL", nats.DefaultURL), "NATS server URL")
	flag.StringVar(&config.User, "user", envOrDefault("NATS_USER", "actorgo"), "NATS user")
	flag.StringVar(&config.Password, "password", envOrDefault("NATS_PASSWORD", "actorgo2026"), "NATS password")
	flag.StringVar(&config.Subject, "subject", envOrDefault("NATS_SUBJECT", DefaultSubject), "NATS subject")
	return config
}

// Connect creates an authenticated connection with consistent reconnect
// behavior for the four examples.
func (config *Config) Connect(clientName string) (*nats.Conn, error) {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second
	opts := []nats.Option{
		nats.Name(clientName),
		nats.Timeout(5 * time.Second),
		nats.ReconnectWait(reconnectDelay),
		nats.MaxReconnects(int(totalWait / reconnectDelay)),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			if err != nil {
				log.Printf("NATS disconnected (%v); reconnecting for up to %.0fm", err, totalWait.Minutes())
			}
		}),
		nats.ReconnectHandler(func(connection *nats.Conn) {
			log.Printf("NATS reconnected to %s", connection.ConnectedUrl())
		}),
		nats.ClosedHandler(func(connection *nats.Conn) {
			if err := connection.LastError(); err != nil {
				log.Printf("NATS connection closed: %v", err)
			}
		}),
	}
	if config.User != "" {
		opts = append(opts, nats.UserInfo(config.User, config.Password))
	}
	return nats.Connect(config.URL, opts...)
}

func envOrDefault(name, fallback string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}
	return fallback
}
