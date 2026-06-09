package bus

import "errors"

var ErrSubscriberLoopNotSet = errors.New("subscriber loop not set")

type subscriber[T any] struct {
	id               int
	eb               *EventBus[T]
	channel          <-chan T
	topic            string
	subscriptionLoop func() error
}

func (s *subscriber[T]) Unsubscribe() error {
	return s.eb.Unsubscribe(s.topic, s.id)
}

func (s *subscriber[T]) SetLoop(loop func() error) {
	s.subscriptionLoop = loop
}

func (s *subscriber[T]) Loop() error {
	if s.subscriptionLoop == nil {
		return ErrSubscriberLoopNotSet
	}
	return s.subscriptionLoop()
}

func (s *subscriber[T]) Channel() <-chan T {
	return s.channel
}
