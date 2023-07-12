package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/habib-web-go/habib-bet-backend/db"
	"github.com/habib-web-go/habib-bet-backend/forms"
	"github.com/habib-web-go/habib-bet-backend/models"
	"net/http"
)

const userKey = "user"

type UserController struct{}

func (u UserController) signup(c *gin.Context) {
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
	addUserToSession(c, user)
}

func (u UserController) login(c *gin.Context) {
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
	addUserToSession(c, user)
}

func addUserToSession(c *gin.Context, u *models.User) {
	session := sessions.Default(c)
	session.Clear()
	session.Set(userKey, u.ID)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": u.ID, "username": u.Username})
}

func InitUserController(router *gin.RouterGroup) {
	u := UserController{}
	router.POST("signup", u.signup)
	router.POST("login", u.login)
}
