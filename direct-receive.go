package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {

	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

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

	err = ch.ExchangeDeclare("mapOrigin", "direct", false, true, false, false, nil)
	failOnError(err, "Failed creating simple exchange")
	queue, err := ch.QueueDeclare("", false, true, false, false, nil)
	failOnError(err, "Failed to create queue")

	err = ch.QueueBind(queue.Name, "hunk0|0", "mapOrigin", false, nil)
	failOnError(err, "failed to create queue")

	messages, err := ch.Consume(queue.Name, "", true, false, false, false, nil)

	for m := range messages {
		log.Printf("%s - %s", m.RoutingKey, m.Body)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
