package publisher

import (
	"strconv"
	"time"

	"github.com/streadway/amqp"
	"github.com/vdntruong/rabbitmq/util"
)

func WorkerQueue(ch *amqp.Channel, stop chan bool) {
	q, err := ch.QueueDeclare("queue", false, false, false, false, nil)
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
					q.Name, // routing key | directly to queue have same name
					false,  // mandatory
					false,
					amqp.Publishing{
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
