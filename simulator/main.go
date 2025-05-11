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

	r := gin.Default()
	routers.RegisterRouters(r)

	simulation.StartSimulation()

	// TODO read from .env
	r.Run("simulator:8080")
}
