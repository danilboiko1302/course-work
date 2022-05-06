package nats

import (
	"os"

	"github.com/nats-io/nats.go"
)

type PubSub struct {
	Conn *nats.Conn
}

var Connection *PubSub

func Init() error {
	url := os.Getenv("NATS_URL")

	connection, err := nats.Connect(url)

	if err != nil {
		return err
	}

	Connection = &PubSub{Conn: connection}

	return nil
}

func Close() {
	Connection.Conn.Drain()
	Connection.Conn.Close()
}

func (ps *PubSub) Pub(topic string, data []byte) error {
	return ps.Conn.Publish(topic, data)
}

func (ps *PubSub) Sub(topic string, cb func(data []byte)) (unsub func() error, err error) {
	s, err := ps.Conn.Subscribe(topic, func(msg *nats.Msg) {
		cb(msg.Data)
	})

	if err != nil {
		return nil, err
	}

	return s.Unsubscribe, nil
}
