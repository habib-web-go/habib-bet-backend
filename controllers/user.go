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

type userController struct {
	userKey string
}

func (u *userController) signup(c *gin.Context) {
	var requestBody forms.UserRequest
	if err := c.BindJSON(&requestBody); err != nil {
		handleBadRequest(c, err)
		return
	}
	user, err := models.CreateUser(&requestBody.Username, &requestBody.Password)
	if err != nil {
		if db.IsDuplicateKeyError(err) {
			handleBadRequestWithMessage(c, err, "username already taken")
			return
		}
		panic(err)
	}
	u.addUserToSession(c, user)
}

func (u *userController) login(c *gin.Context) {
	var requestBody forms.UserRequest
	if err := c.BindJSON(&requestBody); err != nil {
		handleBadRequest(c, err)
		return
	}
	user, err := models.GetUserByUsername(&requestBody.Username)
	if err != nil {
		panic(err)
	}
	if user == nil {
		handleError(c, nil, "User not found.", http.StatusNotFound)
		return
	}
	if !user.CheckPasswordHash(&requestBody.Password) {
		handleError(c, nil, "wrong password", http.StatusForbidden)
		return
	}
	u.addUserToSession(c, user)
}

func (u *userController) logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (u *userController) addUserToSession(c *gin.Context, user *models.User) {
	session := sessions.Default(c)
	session.Clear()
	session.Set(u.userKey, user.ID)
	if err := session.Save(); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, createUserResponse(user))
}

func (u *userController) me(c *gin.Context) {
	user := middlewares.GetUser(c)
	c.JSON(http.StatusOK, createUserResponse(user))
}

func (u *userController) increaseCoin(c *gin.Context) {
	var requestBody forms.IncreaseCoinsRequest
	if err := c.BindJSON(&requestBody); err != nil {
		handleBadRequest(c, err)
		return
	}
	user := middlewares.GetUser(c)
	err := user.IncreaseCoins(requestBody.Amount)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, createUserResponse(user))
}

func createUserResponse(u *models.User) *forms.UserResponse {
	return &forms.UserResponse{
		Id:       u.ID,
		Username: u.Username,
		Coins:    u.Coins,
	}
}

func InitUserController(router *gin.RouterGroup) {
	conf := config.GetConfig()
	u := userController{userKey: conf.GetString("session.userKey")}
	router.POST("signup", u.signup)
	router.POST("login", u.login)
	router.POST("logout", u.logout)
	withAuthRouter := router.Group("")
	withAuthRouter.Use(middlewares.AuthMiddleware)
	withAuthRouter.GET("me", u.me)
	withAuthRouter.POST("increase-coins", u.increaseCoin)
}
