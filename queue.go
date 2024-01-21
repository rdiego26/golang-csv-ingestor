package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func acquireChannel() *amqp.Channel {
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

	return ch
}

func launchConsumer() {
	logger := &Logger{}
	ch := acquireChannel()

	q, err := ch.QueueDeclare(
		os.Getenv("QUEUE_NAME"), // name
		false,                   // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		logger.Log(Fatal, fmt.Sprintf("Failed to declare a queue: %v", err))
	}

	messages, err := ch.Consume(
		q.Name, // queue
		"",     // launchConsumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		logger.Log(Fatal, "Failed to register a launchConsumer")
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var user User
			// Unmarshal the JSON data into the struct
			err := json.Unmarshal(d.Body, &user)
			if err != nil {
				logger.Log(Error, fmt.Sprintf("Error while parsing message data: %v", err))
				return
			}
			db := getConnection()
			insertUser(db, user)
			logger.Log(Info, fmt.Sprintf("Imported user: %s", user.Email))
		}
	}()
	<-forever
}

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
		logger.Log(Fatal, fmt.Sprintf("Failed to declare a queue: %v", err))
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

	defer ch.Close()
	logger.Log(Info, fmt.Sprintf("Message was sent: %s", messageBody))
}
