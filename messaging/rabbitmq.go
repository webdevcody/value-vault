package messaging

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

var (
	rabbitConn  *amqp.Connection
	rabbitCh    *amqp.Channel
	rabbitQueue amqp.Queue
	connErr     error
	queueErr    error
)

var topicName = "events"

func getRabbitMqUrl() string {
	var rabbitMqHost = os.Getenv("RABBIT_MQ_HOST")
	var rabbitMqPass = os.Getenv("RABBIT_MQ_PASSWORD")

	if rabbitMqPass == "" {
		log.Fatal("RABBIT_MQ_PASSWORD env variable is not set")
	}

	return "amqp://user:" + rabbitMqPass + "@" + rabbitMqHost + ":5672/"
}

func Initialize() {
	if rabbitConn != nil {
		return
	}

	rabbitConn, connErr = amqp.Dial(getRabbitMqUrl())
	if connErr != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", connErr)
	}

	rabbitCh, connErr = rabbitConn.Channel()
	if connErr != nil {
		log.Fatalf("Failed to open a channel: %v", connErr)
	}

	// Declare the topic exchange if it's not already declared
	err := rabbitCh.ExchangeDeclare(
		topicName, // exchange name
		"topic",   // exchange type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %v", err)
	}

	rabbitQueue, queueErr = rabbitCh.QueueDeclare(
		"",    // queue name (empty to let RabbitMQ generate a unique name)
		false, // durable
		true,  // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if queueErr != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Bind the queue to the exchange
	err = rabbitCh.QueueBind(
		rabbitQueue.Name, // queue name
		"#",              // routing key (listen to all topics)
		"events",         // exchange name
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Fatalf("Failed to bind the queue to the exchange: %v", err)
	}
}