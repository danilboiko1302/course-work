package types

type Publisher interface {
	Pub(topic string, data []byte) error
}

type Subscriber interface {
	Sub(topic string, cb func(data []byte)) (unsub func() error, err error)
}

type PubSub interface {
	Publisher
	Subscriber
}
