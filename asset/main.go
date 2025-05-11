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

	// TODO hack - Run is called before RabbitMQ because of forever
	// Try with defer
	defer config.ConnectToRabbitMQ()
	r.Run("asset:8080")

	// TODO read from .env
}
