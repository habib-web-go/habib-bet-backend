package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthController struct{}

func (h HealthController) status(c *gin.Context) {
	c.JSONP(http.StatusOK, gin.H{"status": "ok"})
}

func InitHealthController(router *gin.RouterGroup) {
	h := HealthController{}
	router.GET("health", h.status)
}
