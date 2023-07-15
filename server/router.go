package server

import (
	"github.com/gin-contrib/cors"
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
	initV1Router(router.Group("v1"))
	return router
}

func initV1Router(router *gin.RouterGroup) {
	conf := config.GetConfig()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(
		sessions.Sessions(
			conf.GetString("session.name"),
			cookie.NewStore([]byte(conf.GetString("session.secret"))),
		),
	)
	router.Use(cors.Default())
	controllers.InitHealthController(router)
	controllers.InitUserController(router.Group("user"))
	controllers.InitPublicContestController(router.Group("public"))
	controllers.InitContestController(router.Group("contest"))
	controllers.InitQuestionController(router.Group("question"))
}
