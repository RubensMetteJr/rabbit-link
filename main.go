package main

import ("bufio"
"fmt"
"log"
"os"
"strings"
"github.com/RubensMetteJr/rabbit-link/rabbitmq"
"github.com/RubensMetteJr/rabbit-link/publisher"
"github.com/RubensMetteJr/rabbit-link/consumer"
)

func main() {
reader := bufio.NewReader(os.Stdin)

conn, err := rabbitmq.Connect()
if err != nil {
	log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
}
defer rabbitmq.CloseConnection(conn)

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
	err := publisher.PublishMessages(ch, reader)
	if err != nil {
		log.Println("Error publishing messages:", err)
	}
}()

// Goroutine to consume messages
go func() {
	consumer.ConsumeMessages(ch, messages)
}()

fmt.Println("Running...")
<-forever
}