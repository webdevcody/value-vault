package messaging

import (
	"log"

	"github.com/streadway/amqp"
)

func PublishEvent(event string) error {

	// Publish the event message to the topic exchange
	err := rabbitCh.Publish(
		topicName, // exchange
		"",        // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Event '%s' published to topic '%s'\n", event, topicName)
	return nil
}
