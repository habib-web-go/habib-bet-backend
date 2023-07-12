package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/habib-web-go/habib-bet-backend/config"
	"github.com/habib-web-go/habib-bet-backend/db"
	"github.com/habib-web-go/habib-bet-backend/forms"
	"github.com/habib-web-go/habib-bet-backend/middlewares"
	"github.com/habib-web-go/habib-bet-backend/models"
	"net/http"
)

type UserController struct {
	userKey string
}

func (u *UserController) signup(c *gin.Context) {
	var requestBody forms.User
	if err := c.BindJSON(&requestBody); err != nil {
		handleBadRequest(c, err)
		return
	}
	user, err := models.CreateUser(&requestBody.Username, &requestBody.Password)
	if err != nil {
		if db.IsDuplicateKeyError(err) {
			handleBadRequestWithMessage(c, err, "username already taken")
		}
		panic(err)
	}
	u.addUserToSession(c, user)
}

func (u *UserController) login(c *gin.Context) {
	var requestBody forms.User
	if err := c.BindJSON(&requestBody); err != nil {
		handleBadRequest(c, err)
		return
	}
	user, err := models.GetUserByUsername(&requestBody.Username)
	if err != nil {
		panic(err)
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if !user.CheckPasswordHash(&requestBody.Password) {
		c.JSON(http.StatusForbidden, gin.H{"error": "wrong password"})
		return
	}
	u.addUserToSession(c, user)
}

func (u *UserController) logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (u *UserController) addUserToSession(c *gin.Context, user *models.User) {
	session := sessions.Default(c)
	session.Clear()
	session.Set(u.userKey, user.ID)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": user.ID, "username": user.Username})
}

func (u *UserController) me(c *gin.Context) {
	user := middlewares.GetUser(c)
	c.JSON(http.StatusOK, gin.H{"id": user.ID, "username": user.Username})
}

func InitUserController(router *gin.RouterGroup) {
	conf := config.GetConfig()
	u := UserController{userKey: conf.GetString("session.userKey")}
	router.POST("signup", u.signup)
	router.POST("login", u.login)
	router.POST("logout", u.logout)
	withAuthRouter := router.Group("")
	withAuthRouter.Use(middlewares.AuthMiddleware)
	withAuthRouter.GET("me", u.me)
}
