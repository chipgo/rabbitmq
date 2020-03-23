package consumer

import (
	"log"

	"github.com/streadway/amqp"
	"github.com/vdntruong/rabbitmq/util"
)

func WorkerQueue(ch *amqp.Channel, stop chan bool) {
	q, err := ch.QueueDeclare("queue01", false, false, false, false, nil)
	util.FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(q.Name, "01 On Queue01", true, false, false, false, nil)
	util.FailOnError(err, "Failed to register a consumer")

	msgs1, err := ch.Consume(q.Name, "02 On Queue01", true, false, false, false, nil)
	util.FailOnError(err, "Failed to register a consumer")

	msgs2, err := ch.Consume(q.Name, "03 On Queue01", true, false, false, false, nil)
	util.FailOnError(err, "Failed to register a consumer")

	qq, err := ch.QueueDeclare("queue01", false, false, false, false, nil)
	util.FailOnError(err, "Failed to declare a queue")

	msgs3, err := ch.Consume(qq.Name, "04 On Other Queue01", true, false, false, false, nil)
	util.FailOnError(err, "Failed to register a consumer")

	msgs4, err := ch.Consume(qq.Name, "05 On Other Queue01", true, false, false, false, nil)
	util.FailOnError(err, "Failed to register a consumer")

	go func() {
		for {
			select {
			case d := <-msgs:
				log.Printf("Consumer %s received a message: %s", d.ConsumerTag, d.Body)
			case d := <-msgs1:
				log.Printf("Consumer %s received a message: %s", d.ConsumerTag, d.Body)
			case d := <-msgs2:
				log.Printf("Consumer %s received a message: %s", d.ConsumerTag, d.Body)
			case d := <-msgs3:
				log.Printf("Consumer %s received a message: %s", d.ConsumerTag, d.Body)
			case d := <-msgs4:
				log.Printf("Consumer %s received a message: %s", d.ConsumerTag, d.Body)
			case <-stop:
				return
			}
		}
	}()

}
