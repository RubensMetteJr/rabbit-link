package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
	"os"
)

// Get the environment variables passed into the container using the "os" library
var rabbit_host = os.Getenv("RABBIT_HOST")
var rabbit_port = os.Getenv("RABBIT_PORT")
var rabbit_user = os.Getenv("RABBIT_USERNAME")
var rabbit_password = os.Getenv("RABBIT_PASSWORD")

func Connect() (*amqp.Connection, error) {
	// Establish RabbitMQ connection
	conn, err := amqp.Dial("amqp://" + rabbit_user + ":" + rabbit_password + "@" + rabbit_host + ":" + rabbit_port + "/")
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func CloseConnection(conn *amqp.Connection) {
	conn.Close()
}