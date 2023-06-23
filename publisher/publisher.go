package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/streadway/amqp"
)

// Get the environment variables passed into the container using the "os" library
var rabbit_host = os.Getenv("RABBIT_HOST")
var rabbit_port = os.Getenv("RABBIT_PORT")
var rabbit_user = os.Getenv("RABBIT_USERNAME")
var rabbit_password = os.Getenv("RABBIT_PASSWORD")

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Establish RabbitMQ connection
	conn, err := amqp.Dial("amqp://" + rabbit_user + ":" + rabbit_password + "@" + rabbit_host + ":" + rabbit_port + "/")
	if err != nil {
		log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a channel", err)
	}
	defer ch.Close()

	fmt.Print("Enter your nickname: ")
	nickname, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading nickname:", err)
		return
	}
	nickname = strings.TrimSpace(nickname)

	// Declare the RabbitMQ queue for consuming
	consumeQueue, err := ch.QueueDeclare(
		nickname, // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a consume queue", err)
	}

	// Create a channel to receive consumed messages
	messages, err := ch.Consume(
		consumeQueue.Name, // queue
		"",                // consumer
		false,             // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to register consumer", err)
	}

	forever := make(chan bool)

	// Goroutine to publish messages
	go func() {
		fmt.Print("Enter the queue name to publish messages to: ")
		queueName, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading input:", err)
			return
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
			}

			fmt.Println("Publish success!")
		}
	}()

	// Goroutine to consume messages
	go func() {
		// Loop to check the messages incoming in the consume queue
		for d := range messages {
			log.Printf("Received a message: %s", d.Body)

			// Acknowledge the message
			d.Ack(false)
		}
	}()

	fmt.Println("Running...")
	<-forever
}