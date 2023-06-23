package publisher

import (
	"bufio"
	"fmt"
	
	"strings"

)

func PublishMessages(ch *amqp.Channel, reader *bufio.Reader) error {
	fmt.Print("Enter the queue name to publish messages to: ")
	queueName, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading input:", err)
		return err
	}

	// Trim any leading/trailing whitespaces and newlines from the queue name
	queueName = strings.TrimSpace(queueName)

	// Declare the RabbitMQ queue for publishing
	publishQueue, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a publish queue", err)
		return err
	}

	for {
		fmt.Print("Enter a message to send (or type 'exit' to quit): ")
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading input:", err)
			continue
		}

		// Trim any leading/trailing whitespaces and newlines from the message
		message = strings.TrimSpace(message)

		if message == "exit" {
			break
		}

		// Publish the message to the publish queue
		err = ch.Publish(
			"",               // exchange
			publishQueue.Name, // routing key (queue name)
			false,            // mandatory
			false,            // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(message),
			},
		)
		if err != nil {
			log.Fatalf("%s: %s", "Failed to publish a message", err)
			return err
		}

		fmt.Println("Publish success!")
	}

	return nil
}