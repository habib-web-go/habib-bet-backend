package controllers

import (
	"github.com/gin-gonic/gin"
	_db "github.com/habib-web-go/habib-bet-backend/db"
	"github.com/habib-web-go/habib-bet-backend/forms"
	"github.com/habib-web-go/habib-bet-backend/middlewares"
	"github.com/habib-web-go/habib-bet-backend/models"
	"net/http"
	"strconv"
	"time"
)

type questionController struct{}

func (q *questionController) submitAnswer(c *gin.Context) {
	db := _db.GetDB()
	user := middlewares.GetUser(c)
	now := time.Now()
	var requestBody forms.QuestionAnswer
	if err := c.BindJSON(&requestBody); err != nil {
		handleBadRequest(c, err)
		return
	}
	if requestBody.Option != "A" && requestBody.Option != "B" {
		handleBadRequestWithMessage(c, nil, "option should A or B")
		return
	}
	questionId, err := strconv.Atoi(c.Param("questionId"))
	if err != nil {
		handleBadRequest(c, err)
		return
	}
	question, err := models.GetQuestionById(uint(questionId))
	if err != nil {
		panic(err)
	}
	if question == nil {
		handleError(c, nil, "question not found", http.StatusNotFound)
		return
	}
	contest, err := models.GetContestById(question.ContestId)
	var userContest = &models.UserContest{ContestID: contest.ID, UserID: user.ID}
	db.Where(userContest).Find(userContest)
	if userContest.ID == 0 {
		handleBadRequestWithMessage(c, nil, "you are not part of this contest")
		return
	}
	userContest.Contest = contest
	question.Contest = contest
	questions, err := contest.Questions()
	if err != nil {
		panic(err)
	}
	start, end := question.Interval(len(*questions))
	if now.Before(start) {
		handleBadRequestWithMessage(c, nil, "question is not started")
		return
	}
	if now.After(end) {
		handleBadRequestWithMessage(c, nil, "question time is ended")
		return
	}

	var userAnswer = &models.UserAnswer{QuestionID: question.ID, UserContestID: userContest.ID}
	db.Where(userAnswer).Find(userAnswer)
	userAnswer.Question = question
	userAnswer.UserContest = userContest

	answers, err := userContest.Answers()
	if err != nil {
		panic(err)
	}
	if checkLastQuestions(answers, questions, question) {
		handleBadRequestWithMessage(c, nil, "you are out of contest")
		return
	}
	answer := &(*answers)[question.Order]
	answer.Answer = requestBody.Option
	if err := db.Save(answer).Error; err != nil {
		panic(err)
	}
	publicContest, err := createPublicContest(contest, userContest)
	if err != nil {
		panic(err)
	}
	contestResponse := &forms.Contest{
		PublicContest: *publicContest,
		Registered:    true,
		RewardPaid:    false,
	}
	c.JSON(http.StatusOK, contestResponse)
}

func checkLastQuestions(answers *[]models.UserAnswer, questions *[]models.Question, question *models.Question) bool {
	for i, q := range *questions {
		answer := (*answers)[i]
		if q.Order >= question.Order {
			break
		}
		var a string
		if q.Answer {
			a = "A"
		} else {
			a = "B"
		}
		if answer.Answer != a {
			return true
		}
	}
	return false
}

func InitQuestionController(router *gin.RouterGroup) {
	q := questionController{}
	router.Use(middlewares.AuthMiddleware)
	router.POST(":questionId", q.submitAnswer)
}
