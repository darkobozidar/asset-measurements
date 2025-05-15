package main

import (
	"simulator/config"
	"simulator/simulation"

	"os"
    "fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// PostgreSQL
	config.ConnectToPostgresDB()
	defer config.SQLDB.Close()

	// RabbitMQ
	conn := config.ConnectToRabbitMQ()
    channel := config.CreateRabbitMQChannel(conn)
    queue := config.CreateRabbitMQQueue(channel, os.Getenv("RABBITMQ_MEASUREMENTS_QUEUE_NAME"))
	defer conn.Close()
    defer channel.Close()

	// Asset measurement simulation
	simManager := &(simulation.SimulationManager{})
	simManager.StartSimulation(func(obj any) {
		config.PublishToQueue(channel, queue.Name, obj)
	})

	// Start server
	r := gin.Default()
	r.Run(fmt.Sprintf("%s:%s", os.Getenv("SIMULATOR_SERVICE_HOST"), os.Getenv("SIMULATOR_SERVICE_CONTAINER_PORT")))
}
