package bus

import (
	"errors"
)

var ErrSendFailed = errors.New("failed to send message")

type publisher struct {
	closeChannel chan struct{}
	channel      chan<- []byte
	topic        string
}

func (p *publisher) Publish(message []byte) {
	p.channel <- message
}

func (p *publisher) Topic() string {
	return p.topic
}

func (p *publisher) Close() {
	p.closeChannel <- struct{}{}
}
