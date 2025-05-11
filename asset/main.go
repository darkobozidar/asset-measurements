package main

import (
	"asset/config"
	"asset/models"
	"asset/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectToDB()
	models.MigrateModels()

	r := gin.Default()
	routers.RegisterRouters(r)

	config.ConnectToMongoDB()

	// TODO how to clean this up?
	go config.ConnectToRabbitMQ()
	// TODO read from .env.
	r.Run("asset:8080")
}
