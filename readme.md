# MESSAGE QUEUE

Common concepts

- Publisher
- Subscriber
- Worker Queue

# RabbitMQ is a Message Queue System

Common components

- Publisher :: A producer is a user application that sends messages.
- Consumer :: A consumer is a user application that receives messages.
- Queue :: A queue is a buffer that stores messages.
- Exchange :: An exchange is a very simple thing. On one side it receives messages from producers and the other side it pushes them to queues.
- Work queue ::

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

## Declare a channel

```
	ch, err := conn.Channel()
	defer ch.Close()
```

### Work Queue, Distribute workload

By default, RabbitMQ Queue support to dispatches every n-th message to the n-th consumer. Sequential queue will dispatch message to all consumer.

### Fair dispatch

RabbitMQ just dispatches a message when the message enters the queue. It doesn't look at the number of unacknowledged messages for a consumer. It just blindly dispatches every n-th message to the n-th consumer (so by default, RabbitMQ support WorkerQueue).

```
err = ch.Qos(
  1,     // prefetch count
  0,     // prefetch size
  false, // global
)
```

This code tells RabbitMQ not to give more than one message to a worker at a time. Or, in other words, don't dispatch a new message to a worker until it has processed and acknowledged the previous one. Instead, it will dispatch it to the next worker that is not still busy.

prefetch size : If all the workers are busy, your queue can fill up. You will want to keep an eye on that, and maybe add more workers, or have some other strategy.

## Declare a queue

RabbitMQ Queue Guide is [here](https://www.rabbitmq.com/queues.html).

You can redefine an existing queue, but RabbitMQ doesn't allow you to redefine with different parameters (just same all attribute)

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

### Temporary queues

## Declare Exchange

Exchange have a few types available: Direct, Topic, Headers, Fanout.

Also we have a default exchange type, Nameless exchange. When we publish message to nameless exchange.

- The default exchange : RabbitMQ was using a default exchange, which is identified by the empty string ("").
- Fanout : As you can probably guess from the name, it just broadcasts all the messages it receives to all the queues it knows
