package middlewares

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/habib-web-go/habib-bet-backend/config"
	"github.com/habib-web-go/habib-bet-backend/models"
	"net/http"
)

func AuthMiddleware(c *gin.Context) {
	conf := config.GetConfig()
	session := sessions.Default(c)
	userId := session.Get(conf.GetString("session.userKey"))
	if userId == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user, err := models.GetUserById(userId.(uint))
	if err != nil {
		panic(err)
	}
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.Set("user", user)
	c.Next()
}

func GetUser(c *gin.Context) *models.User {
	user, exists := c.Get("user")
	if !exists {
		panic(errors.New("user key not exists. pleas use AuthMiddleware"))
	}
	return user.(*models.User)
}
