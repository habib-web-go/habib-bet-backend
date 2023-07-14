package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/habib-web-go/habib-bet-backend/forms"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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

func paginate(c *gin.Context, db *gorm.DB) (*gorm.DB, *forms.PaginationMetadata, error) {
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var count int64
	if result := db.Count(&count); result.Error != nil {
		return nil, nil, result.Error
	}
	metadata := &forms.PaginationMetadata{
		Page:      page,
		PageSize:  pageSize,
		PageCount: (int(count) / pageSize) + 1,
	}
	return db.Offset(offset).Limit(pageSize), metadata, nil
}
