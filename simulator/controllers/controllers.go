package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartSimulation(c *gin.Context) {
	c.JSON(http.StatusOK, "")
}

func StopSimulation(c *gin.Context) {
	c.JSON(http.StatusOK, "")
}
