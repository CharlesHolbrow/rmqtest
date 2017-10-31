# RabbitMQ Tests

These small one-file programs are intended to be run like this:

`$ go run receive.go`

They are experiments as I am learning RabbitMQ and amqp

## Best Practices

### When a client only cares about real time messages

The presence of a queue that is bound to an exchange will cause unhandled messages to to accumulate in that queue.

To create a client for processing real time messages, that client should have a dedicated queue with:

- `durable = false`
- `autoDelete = true`

That way the the queue will be deleted when all channels listening to it are closed.

## Durable and Auto-Delete

From *Getting Started with RabbitMQ and CloudAMQP*:

> Exchanges, connections and queues can be configured with parameters such as durable, temporary, and auto delete upon creation.

- Durable exchanges will survive server restarts and will last until they are explicitly deleted.
- Temporary exchanges exist until RabbitMQ is shutdown.
- Auto deleted exchanges are removed once the last bound object unbound from the exchange.

### Queues vs Exchanges

Clients:

- **can only** receive (consume) messages from queues
- **cannot** receive from an exchange

Exchanges:

- **can** can receive messages from clients
- **can** can receive messages from other exchanges

Queues:

- **can only** receive messages from exchanges

### Routing

- Messages have a **routing key** that helps determine where the message goes
- `QueueBind(queueName, bindingKey, exchangeName string, noWait bool, args Table)`
  - QueueBind string arguments:
    - `queue name`
    - `binding key`
    - `exchange name`

## Exchanges

There is a global namespace for exchanges on a RabbitMQ server (or a RabbitMQ virtual server). Each exchange type has a default exchange.

### Direct Exchange

- Identified by a string (an empty string indicated the default direct exchange)
- Delivers messages to queues based on a message's routing key
- If the message routing key does not match any binding key, the message will be discarded.

- Queues have 0 or more **bindings**
- Each Binding has a **binding key** and a **cue name**
- Messages (aka Publishings) have a **routing key**

### Topic Exchange

- Routes messages to queues (or exchanges) based on **routing key** and something called the **routing pattern**
- The **routing pattern** is specified by the queue binding.
- We can create multiple bindings between a topic exchange and a queue

### Fanout Exchange

- A fanout exchange routes messages to all of the queues that are bound to it.
- Keys provided will simply be ignored.

### Headers Exchange

Headers exchanges use the message header attributes for routing.

## Memory & Performance

Some notes from this [Article](https://spring.io/blog/2011/04/01/routing-topologies-for-performance-and-scalability-with-rabbitmq/):

- Exchanges and bindings are cheap (in terms of memory cost)
- There are no hard limits to the number of exchanges and queues, one can create, and running 100,000 queues on one broker is fine.
- An exchange is not an erlang process for scalability reasons, it is simply a row in RabbitMQ’s built-in Mnesia database.
- With the right tuning and enough RAM you can run well over a million queues.
- Fanout exchanges are very fast because they have no routing to process yet if bound to a large number of queues that changes
- After a queue is idle for 10 seconds or more it will “hibernate” (written to disk) which causes GC on that queue.
