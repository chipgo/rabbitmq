# RabbitMQ is a Message Queue System

Common components

- Publisher
- Exchange
- Queue
- Consumer

### By default, RabbitMQ declare a nameless exchange ("") type called Default, routing my keys match with queue name.

```
  err = ch.Publish(
	"",     // exchange | nameless exchange
	q.Name, // routing key | queue's name
	false,  // mandatory
	false,  // immediate
	amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(MESSAGE),
	})
```

## Declare a queue

- You can redefine an existing queue, but RabbitMQ doesn't allow you to redefine with different parameters (just same all attribute)

```
ourQueue, err := ch.QueueDeclare(
  "supper_queue",   // name
  false,             // durable
  false,            // delete when unused
  false,            // exclusive
  false,            // no-wait
  nil,              // arguments
)
```

### Message durability

About how to make sure that even if the consumer dies, the message isn't lost. Our messages will still be lost if RabbitMQ server (RabbitMQ Broker) stops. When RabbitMQ quits or crashes it will forget the queues and messages unless you tell it not to.

Two things are required to make sure that messages aren't lost:

- we need to mark the queue (declare step)

```
q, err := ch.QueueDeclare(
  "my_queue", // name
  true,       // durable
  false, false, false, nil,
)
```

- And mark messages as durable (publishing step)

```
err = ch.Publish(
  "",           // default nameless exchange
  q.Name,       // routing key
  false, false,
  amqp.Publishing {
    DeliveryMode: amqp.Persistent, // using the amqp.Persistent option amqp.Publishing
    ContentType:  "text/plain",
    Body:         []byte(body),
})
```

Although it tells RabbitMQ to save the message to disk, there is still a short time window when RabbitMQ has accepted a message and hasn't saved it yet. Also, RabbitMQ doesn't do fsync(2) for every message -- it may be just saved to cache and not really written to the disk. The persistence guarantees aren't strong, but it's more than enough for our simple task queue. If you need a stronger guarantee then you can use [publisher confirms](https://www.rabbitmq.com/confirms.html).
