package pubsub

type Subscriber interface {
	Unsubscribe() error
	SetLoop(loop func() error)
	Loop() error
	Channel() <-chan []byte
}
