// publisher should be able to publish messages to a topic.
package pubsub

type Publisher interface {
	Publish(message []byte)
	Topic() string
	Close()
}
