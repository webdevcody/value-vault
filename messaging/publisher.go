package messaging

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func PublishEvent(event string, value string) error {

	// Publish the event message to the topic exchange
	err := rabbitCh.Publish(
		"events", // exchange
		"",       // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("%s|YOLO|%s", event, value)),
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Event '%s' published to topic '%s'\n", event, getTopicName())
	return nil
}
