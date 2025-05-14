package main

import (
	"simulator/config"
	"simulator/models"
	"simulator/simulation"

	"github.com/gin-gonic/gin"
)

// TODO read from .env where possible.
func main() {
	// PostgreSQL
	config.ConnectToPostgresDB()
	models.MigrateModels()
	defer config.SQLDB.Close()

	// RabbitMQ
	conn := config.ConnectToRabbitMQ()
    channel := config.CreateRabbitMQChannel(conn)
    queue := config.CreateRabbitMQQueue(channel, "asset-measurements")
	defer conn.Close()
    defer channel.Close()

	// Asset measurement simulation
	simulation.StartSimulation(func(obj any) {
		config.PublishToQueue(channel, queue.Name, obj)
	})

	// Start server
	r := gin.Default()
	r.Run("simulator:8080")
}
