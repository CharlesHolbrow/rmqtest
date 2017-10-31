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

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
