package main

import (
	"asset/config"
	"asset/models"
	"asset/controllers"
	"asset/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectToPostgresDB()
	models.MigrateModels()

	r := gin.Default()
	routers.RegisterRouters(r)

	config.ConnectToMongoDB()

	// RabbitMQ
	conn := config.ConnectToRabbitMQ()
	channel := config.CreateRabbitMQChannel(conn)
	queue := config.CreateRabbitMQQueue(channel, "asset-measurements")
	go config.RegisterRabbitMQConsumer(channel, queue.Name, controllers.CreateMeasurement)
	defer conn.Close()
	defer channel.Close()

	// TODO read from .env.
	r.Run("asset:8080")
}
