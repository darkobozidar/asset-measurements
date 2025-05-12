package config

import (
	"asset/utils"

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

func RegisterRabbitMQConsumer(ch *amqp.Channel, queueName string, messageHandler func (queueMessage []byte)) {
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	utils.FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}
		go func() {
			for m := range msgs {
				messageHandler(m.Body)
			}
		}()
	<-forever
}

