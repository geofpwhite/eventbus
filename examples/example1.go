package examples

import (
	"fmt"
	"time"

	"github.com/geofpwhite/eventbus/bus"
)

func Example() {
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
