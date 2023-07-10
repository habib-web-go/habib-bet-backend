package server

import (
	"github.com/gin-gonic/gin"
	"github.com/habib-web-go/habib-bet-backend/config"
	"github.com/habib-web-go/habib-bet-backend/controllers"
)

func NewRouter() *gin.Engine {
	conf := config.GetConfig()

	gin.SetMode(conf.GetString("server.mode"))
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	health := new(controllers.HealthController)

	router.GET("/health", health.Status)
	return router
}
