package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// This is an experiment with rebinding a topic exchange.
//
// It shows that it is trivial to rebind a quene to a new topic pattern without
// dropping any messages or getting double deliveries.
func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	defer func() {
		fmt.Println("Closing the channel and connection")
		conn.Close()
		ch.Close()
	}()

	err = ch.ExchangeDeclare(
		"hunks", // name
		"topic", // type
		false,   // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")
	fmt.Println("Created Queue with name:", q.Name)

	err = ch.QueueBind(
		q.Name,       // queue name
		"objs.alpha", // routing key
		"hunks",      // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for {
			select {
			case d := <-msgs:
				log.Printf(" [x] %s", d.Body)
			}
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	time.Sleep(time.Second * 3)
	fmt.Println("Re-binding channel to objs.alpha")

	err = ch.QueueBind(q.Name, "objs.beta", "hunks", false, nil)
	failOnError(err, "failed to rebind object re-bind1")

	err = ch.QueueBind(q.Name, "objs.alpha", "hunks", false, nil)
	failOnError(err, "failed to rebind object re-bind2")

	err = ch.QueueBind(q.Name, "objs.alpha", "hunks", false, nil)
	failOnError(err, "failed to rebind object re-bind3")

	err = ch.QueueUnbind(q.Name, "objs.alpha", "hunks", nil)
	failOnError(err, "failed to rebind object re-unbind4")

	wait := sync.WaitGroup{}
	wait.Add(1)
	wait.Wait()
}
