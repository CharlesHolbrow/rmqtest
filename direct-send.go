package main

import (
	"fmt"
	"log"

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

	err = ch.Publish("mapOrigin", "hunk0|0", false, false, amqp.Publishing{
		Body:        []byte("this is it!!"),
		ContentType: "text/plain",
	})
	failOnError(err, "failed to pushlish")
}
