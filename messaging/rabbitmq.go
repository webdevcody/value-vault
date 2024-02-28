package messaging

import (
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

var (
	rabbitConn  *amqp.Connection
	rabbitCh    *amqp.Channel
	rabbitQueue amqp.Queue
	connErr     error
	queueErr    error
)

func getRabbitMqUrl() string {
	var rabbitMqHost = os.Getenv("RABBIT_MQ_HOST")
	var rabbitMqPass = os.Getenv("RABBIT_MQ_PASSWORD")

	if rabbitMqPass == "" {
		log.Fatal("RABBIT_MQ_PASSWORD env variable is not set")
	}

	return "amqp://user:" + rabbitMqPass + "@" + rabbitMqHost + ":5672/"
}

func Shutdown() {
	err := rabbitConn.Close()
	if err != nil {
		log.Fatalf("Failed to close rabbitmq connection: %v", err)
	}
	err = rabbitCh.Close()
	if err != nil {
		log.Fatalf("Failed to close rabbitmq channel: %v", err)
	}
}

func getTopicName() string {
	hostname := os.Getenv("HOSTNAME")
	return strings.ReplaceAll(strings.ReplaceAll(hostname, "-primary", ""), "-secondary", "")
}

func Initialize() {
	hostname := os.Getenv("HOSTNAME")

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
		"events", // both primary and secondary need to share the same exchange
		"topic",  // exchange type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %v", err)
	}

	rabbitQueue, queueErr = rabbitCh.QueueDeclare(
		hostname, // queue name (empty to let RabbitMQ generate a unique name)
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if queueErr != nil {
		log.Fatalf("Failed to declare a queue called %s: %v", hostname, queueErr)
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
