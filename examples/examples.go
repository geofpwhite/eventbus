package examples

import (
	"fmt"
	"time"

	"github.com/geofpwhite/eventbus/bus"
)

func Example1pub1sub() {
	// create a new event bus
	eb := bus.New()

	pub, err := eb.NewPublisher("example-topic")
	if err != nil {
		panic(err)
	}

	sub, err := eb.Subscribe("example-topic")
	if err != nil {
		panic(err)
	}

	sub.SetLoop(func() error {
		i := 0
		for msg := range sub.Channel() {
			fmt.Println(string(msg))
			i++
			if i == 2 {
				return nil
			}
		}
		return nil
	})

	go sub.Loop()
	pub.Publish([]byte("hi friend"))

	pub.Publish([]byte("hopefully this works"))
	pub.Close()
	time.Sleep(5 * time.Second)
}

func Example2pub2sub() {
	// create a new event bus
	eb := bus.New()
	pub, err := eb.NewPublisher("example-topic")
	if err != nil {
		panic(err)
	}
	sub1, err := eb.Subscribe("example-topic")
	if err != nil {
		panic(err)
	}
	sub2, err := eb.Subscribe("example-topic")
	if err != nil {
		panic(err)
	}
	pub2, err := eb.NewPublisher("example-topic")
	if err != nil {
		panic(err)
	}
	sub1.SetLoop(func() error {
		i := 0
		for msg := range sub1.Channel() {
			fmt.Printf("sub1: %s\n", string(msg))
			i++
			if i == 4 {
				return nil
			}
		}
		return nil
	})
	sub2.SetLoop(func() error {
		i := 0
		for msg := range sub2.Channel() {
			fmt.Printf("sub2: %s\n", string(msg))
			i++
			if i == 4 {
				return nil
			}
		}
		return nil
	})
	go sub1.Loop()
	go sub2.Loop()
	pub.Publish([]byte("hi friend"))
	pub.Publish([]byte("hopefully this works"))
	pub2.Publish([]byte("this is the second publisher"))
	pub2.Publish([]byte("it should work too"))
	pub.Close()
	pub2.Close()
	time.Sleep(5 * time.Second)
}
