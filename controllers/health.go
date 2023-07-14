package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type healthController struct{}

func (h healthController) status(c *gin.Context) {
	c.JSONP(http.StatusOK, gin.H{"status": "ok"})
}

func InitHealthController(router *gin.RouterGroup) {
	h := healthController{}
	router.GET("health", h.status)
}
