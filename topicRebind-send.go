package main

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("send!")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"hunks", // name
		"topic", // type
		false,   // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare exchange")

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			fmt.Println("sending to objs.alpha and objs.beta")
			err = ch.Publish("hunks", "objs.alpha", false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("sent to objs.alpha"),
			})
			failOnError(err, "Failed sending to objs.alpha")

			err = ch.Publish("hunks", "objs.beta", false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("sent to objs.beta"),
			})
			failOnError(err, "Failed sending to objs.beta")
		}
	}

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
