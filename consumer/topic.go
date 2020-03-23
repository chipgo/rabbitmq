package consumer

import (
	"log"

	"github.com/streadway/amqp"
	"github.com/vdntruong/rabbitmq/util"
)

func Topic(ch *amqp.Channel, stop chan bool) {
	err := ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	util.FailOnError(err, "Failed to declare topic exchange")

	q, err := ch.QueueDeclare(
		"", false, false,
		true, // exclusive
		false, nil)
	util.FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil,
	)
	util.FailOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, "", true, false, false, false, nil)
	util.FailOnError(err, "Failed to register a consumer")

	q02, err := ch.QueueDeclare(
		"", false, false,
		true, // exclusive
		false, nil)
	util.FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q02.Name, // queue name
		"",       // routing key
		"logs",   // exchange
		false,
		nil,
	)
	util.FailOnError(err, "Failed to bind a queue")

	msgs02, err := ch.Consume(
		q02.Name, "", true, false, false, false, nil)
	util.FailOnError(err, "Failed to register a consumer")

	go func() {
		for {
			select {
			case d := <-msgs:
				log.Printf("Consumer %s received a message: %s", d.ConsumerTag, d.Body)
			case d := <-msgs02:
				log.Printf("Consumer %s received a message: %s", d.ConsumerTag, d.Body)
			case <-stop:
				return
			}
		}
	}()
}
