package publisher

import (
	"strconv"
	"time"

	"github.com/streadway/amqp"

	"github.com/vdntruong/rabbitmq/util"
)

func DefaultPublisher(ch *amqp.Channel, stop chan bool) {
	// you can redefine queue with same params
	q, err := ch.QueueDeclare("queue", false, false, false, false, nil)
	util.FailOnError(err, "Failed to declare a queue")

	ticker := time.NewTicker(2 * time.Second)
	counter := 0
	go func() {
		for {
			select {
			case <-ticker.C:
				// body := util.BodyFrom(os.Args)
				counter++
				body := "message 0" + strconv.Itoa(counter)
				err = ch.Publish(
					"",
					q.Name, // routing key | directly to queue have same name
					false,
					false,
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(body),
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
