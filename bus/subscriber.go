package bus

import "errors"

var ErrSubscriberLoopNotSet = errors.New("subscriber loop not set")

type subscriber struct {
	id               int
	eb               *EventBus
	channel          <-chan []byte
	topic            string
	subscriptionLoop func() error
}

func (s *subscriber) Unsubscribe() error {
	return s.eb.Unsubscribe(s.topic, s.id)
}

func (s *subscriber) SetLoop(loop func() error) {
	s.subscriptionLoop = loop
}

func (s *subscriber) Loop() error {
	if s.subscriptionLoop == nil {
		return ErrSubscriberLoopNotSet
	}
	return s.subscriptionLoop()
}

func (s *subscriber) Channel() <-chan []byte {
	return s.channel
}
