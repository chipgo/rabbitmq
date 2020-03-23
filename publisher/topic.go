package publisher

import (
	"strconv"
	"time"

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

	ticker := time.NewTicker(2 * time.Second)
	counter := 0
	go func() {
		for {
			select {
			case <-ticker.C:
				counter++
				body := "log message 0" + strconv.Itoa(counter)
				err = ch.Publish(
					"logs",
					"", // we don't need routing keys, with fanout exchange
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
