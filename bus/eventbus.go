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

type EventBus[T any] struct {
	topics map[string]*TopicBus[T]
}

type TopicBus[T any] struct {
	ctx                   context.Context
	done                  context.CancelFunc // this is used to stop the topic bus loop when the number of publishers goes to 0
	subscribers           map[int]chan<- T   // map exists so subcribers can be removed without changing ID of other subscribers
	publisherCloseChannel chan struct{}
	publishers            chan T
	subscriberMut         sync.Mutex
	topic                 string
	curID                 int
	numPublishers         int
}

func (tb *TopicBus[T]) Loop() {
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

func (tb *TopicBus[T]) AddSubscriber() (int, <-chan T) {
	tb.subscriberMut.Lock()
	defer tb.subscriberMut.Unlock()
	channel := make(chan T)
	tb.subscribers[tb.curID] = channel
	tb.curID++
	return tb.curID - 1, channel
}

func (tb *TopicBus[T]) NewPublisher() pubsub.Publisher[T] {
	tb.numPublishers++
	p := publisher[T]{
		channel:      tb.publishers,
		topic:        tb.topic,
		closeChannel: tb.publisherCloseChannel,
	}
	return &p
}

func (eb *EventBus[T]) Unsubscribe(topic string, id int) error {
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

func (eb *EventBus[T]) Subscribe(topic string) (pubsub.Subscriber[T], error) {
	if topic, ok := eb.topics[topic]; ok {
		id, channel := topic.AddSubscriber()
		s := subscriber[T]{
			id:      id,
			eb:      eb,
			channel: channel,
			topic:   topic.topic,
		}
		return &s, nil
	}
	return nil, ErrTopicDoesNotExist
}

func (eb *EventBus[T]) NewTopic(topic string) error {
	if _, ok := eb.topics[topic]; ok {
		return ErrTopicAlreadyExists
	}
	ctx, done := context.WithCancel(context.Background())
	eb.topics[topic] = &TopicBus[T]{
		ctx:                   ctx,
		done:                  done,
		publisherCloseChannel: make(chan struct{}),
		publishers:            make(chan T),
		subscribers:           make(map[int]chan<- T),
	}
	go eb.topics[topic].Loop()

	return nil
}

func (eb *EventBus[T]) NewPublisher(topic string) (pubsub.Publisher[T], error) {
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

func New[T any]() *EventBus[T] {
	return &EventBus[T]{
		topics: make(map[string]*TopicBus[T]),
	}
}
