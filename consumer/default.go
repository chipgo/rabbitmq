package consumer

import (
	"log"

	"github.com/streadway/amqp"

	"github.com/vdntruong/rabbitmq/util"
)

func DefaultConsumer(ch *amqp.Channel, stop chan bool) {
	q, err := ch.QueueDeclare("queue", false, false, false, false, nil)
	util.FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(q.Name, "System Events", true, false, false, false, nil)
	util.FailOnError(err, "Failed to register a consumer")

	go func() {
		for {
			select {
			case d := <-msgs:
				log.Printf("Consumer %s received a message: %s", d.ConsumerTag, d.Body)
			case <-stop:
				return
			}
		}
	}()
}
