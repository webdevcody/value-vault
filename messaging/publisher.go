package messaging

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func PublishEvent(event string, value string) error {

	mode := os.Getenv("MODE")

	// Publish the event message to the topic exchange
	err := rabbitCh.Publish(
		"events", // exchange
		"",       // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("%s|YOLO|%s|YOLO|%s", mode, event, value)),
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Event '%s' published to topic '%s'\n", event, getTopicName())
	return nil
}
