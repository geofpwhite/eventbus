// publisher should be able to publish messages to a topic.
package pubsub

type Publisher[T any] interface {
	Publish(message T)
	Topic() string
	Close()
}
