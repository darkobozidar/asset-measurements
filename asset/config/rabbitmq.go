package config

import (
	"log"
	"time"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type AssetMeasurement struct {
    AssetID   uint      `json:"asset_id"`
    Timestamp time.Time `json:"timestamp"`
    Power     float64   `json:"power"`
    SOE       float64   `json:"soe"`
}

func ConnectToRabbitMQ() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"asset-measurements", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}
		go func() {
			for m := range msgs {
				var assetMeasurement AssetMeasurement
				err := json.Unmarshal(m.Body, &assetMeasurement)
				failOnError(err, "Failed to decode JSON")

				persistAssetMeasurement(assetMeasurement)
			}
		}()

	  log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func persistAssetMeasurement(assetMeasurement AssetMeasurement) {
	log.Printf("Received a message: %v", assetMeasurement)
}
