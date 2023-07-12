package server

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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
	router.Use(
		sessions.Sessions(
			conf.GetString("session.name"),
			cookie.NewStore([]byte(conf.GetString("session.secret"))),
		),
	)
	rootRouter := router.Group("/")
	controllers.InitHealthController(rootRouter)
	controllers.InitUserController(router.Group("user"))
	return router
}
