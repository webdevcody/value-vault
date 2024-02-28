package messaging

import (
	"log"
)

func InitializeEventListener(handleEvent func(event string)) {

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
		handleEvent(event)
	}
}
