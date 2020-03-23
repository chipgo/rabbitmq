package publisher

import (
	"strconv"
	"time"

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

	ticker := time.NewTicker(2 * time.Second)
	counter := 0
	go func() {
		for {
			select {
			case <-ticker.C:
				counter++
				body := "message 0" + strconv.Itoa(counter)
				err = ch.Publish(
					"",
					q.Name,
					false,
					false,
					amqp.Publishing{
						// we need to mark our messages as persistent
						// by using the amqp.Persistent option amqp.Publishing takes.
						DeliveryMode: amqp.Persistent,
						ContentType:  "text/plain",
						Body:         []byte(body),
					},
				)
				util.FailOnError(err, "Failed to publish message")
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()
}
