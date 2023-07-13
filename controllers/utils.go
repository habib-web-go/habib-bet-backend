package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/habib-web-go/habib-bet-backend/forms"
	"net/http"
)

func getErrorMessage(message string, err error) string {
	if gin.Mode() == gin.DebugMode && err != nil {
		message += " " + err.Error()
	}
	return message
}

func handleBadRequest(c *gin.Context, err error) {
	handleBadRequestWithMessage(c, err, "bad request.")
}

func handleBadRequestWithMessage(c *gin.Context, err error, message string) {
	handleError(c, err, message, http.StatusBadRequest)
}

func handleError(c *gin.Context, err error, message string, statusCode int) {
	c.JSON(statusCode, forms.ErrorResponse{Error: getErrorMessage(message, err)})
}
