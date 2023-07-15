package controllers

import (
	"github.com/gin-gonic/gin"
	_db "github.com/habib-web-go/habib-bet-backend/db"
	"github.com/habib-web-go/habib-bet-backend/forms"
	"github.com/habib-web-go/habib-bet-backend/models"
	"net/http"
	"time"
)

type publicContestController struct{}

func (h publicContestController) ongoing(c *gin.Context) {
	db := _db.GetDB()
	var contests []models.Contest
	now := time.Now()
	result := db.Where("\"start\" <= ? AND ?  <= \"end\"", now, now).Find(&contests)
	if result.Error != nil {
		panic(result.Error)
	}
	publicContests := make([]*forms.PublicContest, len(contests))
	for i, contest := range contests {
		publicContest, err := createPublicContest(&contest, nil)
		if err != nil {
			panic(err)
		}
		publicContests[i] = publicContest
	}
	c.JSON(http.StatusOK, gin.H{"data": publicContests})
}

func (h publicContestController) coming(c *gin.Context) {
	db := _db.GetDB()
	var contests []models.Contest
	now := time.Now()
	paginateQuery, metadata, err := paginate(
		c,
		db.Model(&models.Contest{}).Where("? <= \"start\"", now).Order("start, \"end\""),
	)
	if err != nil {
		panic(err)
	}
	result := paginateQuery.Find(&contests)
	if result.Error != nil {
		panic(result.Error)
	}
	publicContests := make([]*forms.PublicContest, len(contests))
	for i, contest := range contests {
		publicContest, err := createPublicContest(&contest, nil)
		if err != nil {
			panic(err)
		}
		publicContests[i] = publicContest
	}
	c.JSON(http.StatusOK, gin.H{"data": publicContests, "metadata": metadata})
}

func createPublicContest(contest *models.Contest, userContest *models.UserContest) (*forms.PublicContest, error) {
	userCount, err := contest.UserCount()
	if err != nil {
		return nil, err
	}
	questions, err := contest.Questions()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	var questionStates []forms.QuestionState
	x := 0
	if contest.Start.Before(now) {
		questionStates = make([]forms.QuestionState, len(*questions))
		input := userCount
		for i, q := range *questions {
			start, end := q.Interval(len(*questions))
			if end.Before(now) {
				output, err := q.TrueAnswerCount()
				if err != nil {
					return nil, err
				}
				var answer string
				if q.Answer {
					answer = "A"
				} else {
					answer = "B"
				}
				questionState := forms.QuestionState{
					ID:      q.ID,
					OptionA: q.OptionA,
					OptionB: q.OptionB,
					Start:   start,
					End:     end,
					Order:   q.Order,
					Input:   input,
					Output:  output,
					Answer:  answer,
				}
				if userContest != nil {
					answers, err := userContest.Answers()
					if err != nil {
						return nil, err
					}
					questionState.UserAnswer = (*answers)[i].Answer
					if len(questionState.UserAnswer) == 0 {
						questionState.UserAnswer = "_"
					}
				}
				questionStates[i] = questionState
				input = output
			} else if start.Before(now) && end.After(now) {
				questionState := forms.QuestionState{
					ID:      q.ID,
					OptionA: q.OptionA,
					OptionB: q.OptionB,
					Start:   start,
					End:     end,
					Order:   q.Order,
					Input:   input,
				}
				if userContest != nil {
					answers, err := userContest.Answers()
					if err != nil {
						return nil, err
					}
					questionState.UserAnswer = (*answers)[i].Answer
					if len(questionState.UserAnswer) == 0 {
						questionState.UserAnswer = "_"
					}
				}
				questionStates[i] = questionState
			} else {
				x++
			}
		}
		questionStates = questionStates[0 : len(*questions)-x]
	}
	return &forms.PublicContest{
		ID:            contest.ID,
		Name:          contest.Name,
		Start:         contest.Start,
		End:           contest.End,
		UserCount:     userCount,
		Questions:     &questionStates,
		QuestionCount: len(*questions),
	}, nil
}

func InitPublicContestController(router *gin.RouterGroup) {
	c := publicContestController{}
	router.GET("ongoing", c.ongoing)
	router.GET("coming", c.coming)
}
