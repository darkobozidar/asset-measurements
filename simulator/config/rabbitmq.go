package config

import (
	"simulator/utils"

	"time"
	"context"
	"encoding/json"
	"fmt"
    "os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectToRabbitMQ() *amqp.Connection {
    const BASE_CONN_URL string = "amqp://%s:%s@rabbitmq:%s/"

    connUrl := fmt.Sprintf(
        BASE_CONN_URL,
        os.Getenv("RABBITMQ_DEFAULT_USER"),
        os.Getenv("RABBITMQ_DEFAULT_PASS"),
        os.Getenv("RABBITMQ_CONTAINER_PORT"),
    )

    conn, err := amqp.Dial(connUrl)
    utils.FailOnError(err, "Failed to connect to RabbitMQ")
    return conn
}

func CreateRabbitMQChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	return ch
}

func CreateRabbitMQQueue(ch *amqp.Channel, queueName string) amqp.Queue {
	queue, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")

	return queue
}

func PublishToQueue(ch *amqp.Channel, queueName string, obj any) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	body, err := json.Marshal(obj)
	utils.FailOnError(err, "Failed to marshal obj")

	err = ch.PublishWithContext(
		ctx,
        "",        // exchange
        queueName, // routing key (queue name)
        false,     // mandatory
        false,     // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
		},
	)

	utils.FailOnError(err, "Failed to publish a message")
}
