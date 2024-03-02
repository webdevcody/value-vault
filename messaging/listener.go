package messaging

import (
	"log"
)

func InitializeRedistributedEventListener() {

	if rabbitConn == nil {
		log.Fatal("RabbitMQ connection is not initialized")
	}

	// Consume messages from the queue
	msgs, err := rabbitCh.Consume(
		getRedistributeQueueName(), // queue name
		"",                         // consumer
		true,                       // auto-ack
		false,                      // exclusive
		false,                      // no-local
		false,                      // no-wait
		nil,                        // arguments
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Read messages from the channel
	for msg := range msgs {
		event := string(msg.Body)
		handleRedistributeEvent(event)
	}
}

func InitializeReplicationEventListener() {

	if rabbitConn == nil {
		log.Fatal("RabbitMQ connection is not initialized")
	}

	// Consume messages from the queue
	msgs, err := rabbitCh.Consume(
		rabbitQueue.Name, // queue name
		"",               // consumer
		true,             // auto-ack
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Read messages from the channel
	for msg := range msgs {
		event := string(msg.Body)
		handleReplicationEvent(event)
	}
}
