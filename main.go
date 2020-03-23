package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/streadway/amqp"

	"github.com/vdntruong/rabbitmq/consumer"
	"github.com/vdntruong/rabbitmq/publisher"
	"github.com/vdntruong/rabbitmq/util"
)

type (
	RabbitConfig struct {
		User string `envconfig:"RABBIT_USER" default:"admin"`
		Pass string `envconfig:"RABBIT_PASS" default:"123456"`
		Host string `envconfig:"RABBIT_HOST" default:"localhost"`
		Port string `envconfig:"RABBIT_PORT" default:"5672"`
	}
)

func main() {
	var conf RabbitConfig
	util.LoadConfig(&conf)
	strConf := fmt.Sprintf("amqp://%s:%s@%s:%s", conf.User, conf.Pass, conf.Host, conf.Port)

	conn, err := amqp.Dial(strConf)
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open channel")
	defer ch.Close()

	exchangeType := flag.String("type", "default", "type exchange to demo")
	flag.Parse()

	stopSignal := make(chan bool)
	switch *exchangeType {
	case "default":
		consumer.DefaultConsumer(ch, stopSignal)
		publisher.DefaultPublisher(ch, stopSignal)
	case "workerqueue":
		consumer.WorkerQueue(ch, stopSignal)
		publisher.WorkerQueue(ch, stopSignal)
	case "durable":
		consumer.Duration(ch, stopSignal)
		publisher.Duration(ch, stopSignal)
	default:
		log.Println("Unknown Exchange Type")
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-stopSignal
}
