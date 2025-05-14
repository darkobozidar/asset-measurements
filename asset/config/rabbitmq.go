package config

import (
    "asset/utils"

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

