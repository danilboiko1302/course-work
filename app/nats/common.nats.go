package nats

import (
	"os"

	"github.com/nats-io/nats.go"
)

var connection *nats.Conn
var err error

func Init() error {
	url := os.Getenv("NATS_URL")

	connection, err = nats.Connect(url)

	if err != nil {
		return err
	}

	return nil
}

func Close() {
	connection.Drain()
	connection.Close()
}
