package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func sendMessage(user User) {
	logger := &Logger{}
	conn, err := amqp.Dial(os.Getenv("QUEUE_URL"))

	if err != nil {
		log.Println(err)
		logger.Log(Fatal, fmt.Sprintf("Failed to connect to broker err: %v", err))
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.Log(Fatal, fmt.Sprintf("Failed to open a channel err: %v", err))
	}
	defer ch.Close()

	// We create a Queue to send the message to.
	queue, err := ch.QueueDeclare(
		os.Getenv("QUEUE_NAME"), // name
		false,                   // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		logger.Log(Fatal, "Failed to declare a queue")
	}

	// Convert struct to JSON
	messageBody, err := json.Marshal(user)
	if err != nil {
		logger.Log(Fatal, fmt.Sprintf("Failed to marshal struct to JSON: %v", err))
	}

	err = ch.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBody,
		})

	logger.Log(Info, fmt.Sprintf("Message was sent: %s", messageBody))
}
