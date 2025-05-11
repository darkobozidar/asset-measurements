package main

import (
	"simulator/config"
	"simulator/models"
	"simulator/routers"
	"simulator/simulation"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectToDB()
	models.MigrateModels()

	config.ConnectToRabbitMQ()

	r := gin.Default()
	routers.RegisterRouters(r)

	simulation.StartSimulation()

	// TODO read from .env
	r.Run("asset:8080")
}
