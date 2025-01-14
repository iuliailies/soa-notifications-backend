package main

import (
	"log"

	"github.com/streadway/amqp"
)

// StartConsumer connects to RabbitMQ, listens to the specified queue, and broadcasts messages to the WebSocket server.
func StartConsumer(rabbitMQURL, queueName string, broadcast chan<- string) error {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declare a queue
	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	// Start consuming messages
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	log.Println("RabbitMQ consumer started, waiting for messages...")

	// Process messages
	for msg := range msgs {
		log.Printf("Received message: %s", msg.Body)
		// Send the message to the WebSocket broadcast channel
		broadcast <- string(msg.Body)
	}

	return nil
}
