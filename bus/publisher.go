package bus

import (
	"errors"
)

var ErrSendFailed = errors.New("failed to send message")

type publisher[T any] struct {
	closeChannel chan struct{}
	channel      chan<- T
	topic        string
}

func (p *publisher[T]) Publish(message T) {
	p.channel <- message
}

func (p *publisher[T]) Topic() string {
	return p.topic
}

func (p *publisher[T]) Close() {
	p.closeChannel <- struct{}{}
}
