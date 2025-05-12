package config

import (
	"time"
	"context"
	"encoding/json"

	"simulator/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

// TODO collect data from .env
func ConnectToRabbitMQ() *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
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
