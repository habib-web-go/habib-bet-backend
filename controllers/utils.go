package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func getErrorMessage(message string, err error) string {
	if gin.Mode() == gin.DebugMode {
		message += " " + err.Error()
	}
	return message
}

func handleBadRequest(c *gin.Context, err error) {
	handleBadRequestWithMessage(c, err, "bad request.")
}

func handleBadRequestWithMessage(c *gin.Context, err error, message string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": getErrorMessage(message, err)})
}
