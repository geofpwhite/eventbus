package bus

import (
	"context"
	"errors"
	"sync"

	"github.com/geofpwhite/eventbus/pubsub"
)

var (
	ErrTopicDoesNotExist      = errors.New("topic does not exist")
	ErrTopicAlreadyExists     = errors.New("topic already exists")
	ErrSubscriberDoesNotExist = errors.New("subscriber does not exist")
)

type EventBus struct {
	topics map[string]*TopicBus
}

type TopicBus struct {
	ctx                   context.Context
	done                  context.CancelFunc // this is used to stop the topic bus loop when the topic is deleted
	numPublishers         int
	publisherCloseChannel chan struct{}
	topic                 string
	subscribers           map[int]chan<- []byte // map exists so subcribers can be removed without changing ID of other subscribers
	publishers            chan []byte
	subscriberMut         sync.Mutex
	curID                 int
}

func (tb *TopicBus) Loop() {
	for {
		select {
		case <-tb.ctx.Done():
			tb.subscriberMut.Lock()
			for _, sub := range tb.subscribers {
				close(sub)
			}
			tb.subscriberMut.Unlock()
			return
		case msg := <-tb.publishers:
			tb.subscriberMut.Lock()
			for _, sub := range tb.subscribers {
				sub <- msg
			}
			tb.subscriberMut.Unlock()
		case <-tb.publisherCloseChannel:
			tb.numPublishers--
			if tb.numPublishers == 0 {
				tb.done()
			}
		}
	}
}

func (tb *TopicBus) AddSubscriber() (int, <-chan []byte) {
	tb.subscriberMut.Lock()
	defer tb.subscriberMut.Unlock()
	channel := make(chan []byte)
	tb.subscribers[tb.curID] = channel
	tb.curID++
	return tb.curID - 1, channel
}

func (tb *TopicBus) NewPublisher() pubsub.Publisher {
	tb.numPublishers++
	p := publisher{
		channel: tb.publishers,
		topic:   tb.topic,
	}
	return &p
}

func (eb *EventBus) Unsubscribe(topic string, id int) error {
	if topic, ok := eb.topics[topic]; ok {
		topic.subscriberMut.Lock()
		defer topic.subscriberMut.Unlock()
		if sub, ok := topic.subscribers[id]; ok {
			close(sub)
			delete(topic.subscribers, id)
			return nil
		}
		return ErrSubscriberDoesNotExist
	}
	return ErrTopicDoesNotExist
}

func (eb *EventBus) Subscribe(topic string) (pubsub.Subscriber, error) {
	if topic, ok := eb.topics[topic]; ok {
		id, channel := topic.AddSubscriber()
		s := subscriber{
			id:      id,
			eb:      eb,
			channel: channel,
			topic:   topic.topic,
		}
		return &s, nil
	}
	return nil, ErrTopicDoesNotExist
}

func (eb *EventBus) NewTopic(topic string) error {
	if _, ok := eb.topics[topic]; ok {
		return ErrTopicAlreadyExists
	}
	ctx, done := context.WithCancel(context.Background())
	eb.topics[topic] = &TopicBus{
		ctx:                   ctx,
		done:                  done,
		publisherCloseChannel: make(chan struct{}),
		publishers:            make(chan []byte),
		subscribers:           make(map[int]chan<- []byte),
	}
	go eb.topics[topic].Loop()

	return nil
}

func (eb *EventBus) NewPublisher(topic string) (pubsub.Publisher, error) {
	if topic, ok := eb.topics[topic]; ok {
		return topic.NewPublisher(), nil
	}
	err := eb.NewTopic(topic)
	if err != nil {
		return nil, err
	}
	topicBus := eb.topics[topic]
	return topicBus.NewPublisher(), nil
}

func New() *EventBus {
	return &EventBus{
		topics: make(map[string]*TopicBus),
	}
}
