package consumer

import (
	"log"

	"github.com/streadway/amqp"
	"github.com/vdntruong/rabbitmq/util"
)

func Duration(ch *amqp.Channel, stop chan bool) {
	q, err := ch.QueueDeclare(
		"ourqueue",
		true, // durable
		false, false, false, nil,
	)
	util.FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(q.Name, "cs_duration", true, false, false, false, nil)
	util.FailOnError(err, "Failed to consum queue")

	for {
		select {
		case d := <-msgs:
			log.Printf("Consumer %s received a message: %s", d.ConsumerTag, d.Body)
		case <-stop:
			return
		}
	}
}
