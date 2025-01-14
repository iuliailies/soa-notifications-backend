package main

import (
	"log"
	"os"
	"time"
)

func main() {
	// Load RabbitMQ URL from environment variables
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		log.Fatal("RABBITMQ_URL environment variable is not set")
	}

	// Start WebSocket server
	broadcast := make(chan string)
	go StartWebSocketServer(":8084", broadcast)

	time.Sleep(10 * time.Second)

	// Start RabbitMQ consumer
	err := StartConsumer(rabbitMQURL, "notifications_queue", broadcast)
	if err != nil {
		log.Fatalf("Failed to start RabbitMQ consumer: %v", err)
	}
}
