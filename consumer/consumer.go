package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"os"
	"log"
)

//Get the enviroment varaibles passed into the container, using the "os" library
var rabbit_host = os.Getenv("RABBIT_HOST")
var rabbit_port = os.Getenv("RABBIT_PORT") 
var rabbit_user = os.Getenv("RABBIT_USERNAME")
var rabbit_password = os.Getenv("RABBIT_PASSWORD")

func main(){
	consume()

}

func consume(){

	//Stablish rabbitmq connection
	connection, err := amqp.Dial("amqp://" + rabbit_user + ":" +rabbit_password + "@" + rabbit_host + ":" + rabbit_port +"/")

	if err != nil{
		log.Fatalf("%s:%s","Failed to connect to rabbitMQ", err)
	}

	channel, err := connection.Channel()

	if err != nil {
		log.Fatalf("%s:%s","Failed to open a channel",err)
	}

	queue, err := channel.QueueDeclare(
		"publisher", // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	fmt.Println("Channel and Queue established")

	defer connection.Close()
	defer channel.Close()

	messages, err := channel.Consume(
		queue.Name, // queue
		"",     // consumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	  )

	  if err != nil {
		log.Fatalf("%s: %s", "Failed to register consumer", err)
	}

	forever := make(chan bool)

	//routine that runs concurrently
	go func() {
		//loop to check the messages incoming in the queue
		for d := range messages {
			log.Printf("Received a message: %s", d.Body)
			
			d.Ack(false)
		}
	  }()
	  
	  fmt.Println("Running...")
	  <-forever
}