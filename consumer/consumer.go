package consumer

import (
	"log"
	"github.com/streadway/amqp"
)

func ConsumeMessages(ch *amqp.Channel, messages <-chan amqp.Delivery) {
	// Loop to check the messages incoming in the consume queue
	for d := range messages {
		log.Printf("Received a message: %s", d.Body)

		// Acknowledge the message
		d.Ack(false)
	}
}