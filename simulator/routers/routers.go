package routers

import (
	"simulator/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRouters(r *gin.Engine) {
	r.GET("/asset/:id", controllers.StartSimulation)
}
