package nats

import (
	"os"

	"github.com/nats-io/nats.go"
)

var connection *nats.Conn
var err error

func Init() error {
	url := os.Getenv("NATS_URL")
	port := os.Getenv("NATS_PORT")

	connection, err = nats.Connect(url + ":" + port)

	if err != nil {
		return err
	}

	return nil
}

func Close() {
	connection.Drain()
	connection.Close()
}
