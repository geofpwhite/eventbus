package pubsub

type Subscriber[T any] interface {
	Unsubscribe() error
	SetLoop(loop func() error)
	Loop() error
	Channel() <-chan T
}
