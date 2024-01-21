package server

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func launchConsumer() {
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

	queueName := os.Getenv("QUEUE_NAME")

	messages, err := ch.Consume(
		queueName, // queue
		"",        // launchConsumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		logger.Log(Fatal, "Failed to register a launchConsumer")
	}

	logger.Log(Debug, fmt.Sprintf("Fetched %d messages from queue %s", len(messages), queueName))

	forever := make(chan bool)
	go func() {
		var users []User
		for message := range messages {
			var user User
			// Unmarshal the JSON data into the struct
			err := json.Unmarshal(message.Body, &user)
			if err != nil {
				logger.Log(Error, fmt.Sprintf("Error while parsing message data: %v", err))
				return
			}
			users = append(users, user)
		}
		log.Println(fmt.Sprintf("Parsed %d users", len(users)))
		db := getConnection()
		bulkInsertUsers(db, users)
		logger.Log(Info, fmt.Sprintf("Imported %d users", len(users)))
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

	channel, err := conn.Channel()
	if err != nil {
		logger.Log(Fatal, fmt.Sprintf("Failed to open a channel err: %v", err))
	}

	defer channel.Close()

	// We create a Queue to send the message to.
	queue, err := channel.QueueDeclare(
		os.Getenv("QUEUE_NAME"), // name
		true,                    // durable
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

	err = channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         messageBody,
		})

	logger.Log(Info, fmt.Sprintf("Message was sent: %s", messageBody))
}
