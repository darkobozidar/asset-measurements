package config

import (
	"asset/utils"

	"log"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

// TODO collect data from .env
func ConnectToRabbitMQ() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"asset-measurements", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	utils.FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}
		go func() {
			for m := range msgs {
				// log.Printf("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA %v", m.Body)
				var assetMeasurement Measurement
				err := json.Unmarshal(m.Body, &assetMeasurement)
				utils.FailOnError(err, "Failed to decode JSON")

				// TODO check if assetID exists

				log.Printf("BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB %+v", assetMeasurement)

				if err := InsertMeasurement(assetMeasurement); err != nil {
					log.Fatalf("Error inserting measurement: %v", err)
				}
			}
		}()

	//   log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

