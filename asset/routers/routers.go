package routers

import (
	"asset/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRouters(r *gin.Engine) {
	r.GET("/asset/:id", controllers.GetAsset)
	r.GET("/asset/", controllers.GetAssets)
	r.POST("/asset/", controllers.CreateAsset)
	r.PUT("/asset/:id", controllers.UpdateAsset)
	r.DELETE("asset/:id", controllers.DeleteAsset)
}
