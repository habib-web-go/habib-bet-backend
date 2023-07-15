package controllers

import (
	"github.com/gin-gonic/gin"
	_db "github.com/habib-web-go/habib-bet-backend/db"
	"github.com/habib-web-go/habib-bet-backend/forms"
	"github.com/habib-web-go/habib-bet-backend/middlewares"
	"github.com/habib-web-go/habib-bet-backend/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type contestController struct{}

func (cc *contestController) register(c *gin.Context) {
	contestId, err := strconv.Atoi(c.Param("contestId"))
	if err != nil {
		handleBadRequest(c, err)
		return
	}
	contest, err := models.GetContestById(uint(contestId))
	if err != nil {
		panic(err)
	}
	if contest == nil {
		handleError(c, nil, "contest not found", http.StatusNotFound)
		return
	}
	now := time.Now()
	if now.After(contest.Start) {
		handleBadRequestWithMessage(c, nil, "contest is already began.")
		return
	}
	user := middlewares.GetUser(c)
	if user.Coins < contest.EntryFee {
		handleBadRequestWithMessage(c, nil, "not enough coins.")
		return
	}
	_, err = models.AddUserToContest(user, contest)
	if err != nil {
		if _db.IsDuplicateKeyError(err) {
			handleBadRequestWithMessage(c, err, "you already joined this contest")
			return
		}
		panic(err)
	}
	publicContest, err := createPublicContest(contest, nil)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, publicContest)
}

func (cc *contestController) coming(c *gin.Context) {
	db := _db.GetDB()
	now := time.Now()
	type result struct {
		models.Contest
		Registered bool
	}
	var contests []result
	user := middlewares.GetUser(c)
	filter := c.Query("registered") == "true"
	query := db.Model(&models.Contest{}).Where("? <= \"start\"", now).Order("start, \"end\"").Joins(
		"Left Join user_contests on contests.id = user_contests.contest_id AND user_contests.user_id = ?",
		user.ID,
	).Select("contests.*, user_contests.id is not null as registered")
	if filter {
		query = query.Where("user_contests.id is not null")
	}
	paginateQuery, metadata, err := paginate(c, query)
	if err != nil {
		panic(err)
	}
	r := paginateQuery.Scan(&contests)
	if r.Error != nil {
		panic(r.Error)
	}
	contestResponses := make([]*forms.Contest, len(contests))
	for i, contest := range contests {
		publicContest, err := createPublicContest(&contest.Contest, nil)
		if err != nil {
			panic(err)
		}

		contestResponses[i] = &forms.Contest{
			PublicContest: *publicContest,
			Registered:    contest.Registered,
			RewardPaid:    false,
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": contestResponses, "metadata": metadata})
}

func (cc *contestController) ongoing(c *gin.Context) {
	db := _db.GetDB()
	now := time.Now()
	var userContests []models.UserContest
	user := middlewares.GetUser(c)
	r := db.Model(&models.UserContest{UserID: user.ID}).Joins("Contest").Where(
		"\"start\" <= ? AND ?  <= \"end\"", now, now,
	).Order("start, \"end\"").Find(&userContests)
	if r.Error != nil {
		panic(r.Error)
	}
	contestResponses := make([]*forms.Contest, len(userContests))
	for i, userContest := range userContests {
		contest, err := models.GetContestById(userContest.ContestID)
		if err != nil {
			panic(err)
		}
		userContest.Contest = contest
		publicContest, err := createPublicContest(userContest.Contest, &userContest)
		if err != nil {
			panic(err)
		}
		contestResponses[i] = &forms.Contest{
			PublicContest: *publicContest,
			Registered:    true,
			RewardPaid:    false,
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": contestResponses})
}

func (cc *contestController) archive(c *gin.Context) {
	db := _db.GetDB()
	now := time.Now()
	var userContests []models.UserContest
	user := middlewares.GetUser(c)
	query := db.Model(&models.UserContest{UserID: user.ID}).Where(
		"\"end\" <= ?", now,
	).Order("\"end\" desc, start desc").Joins("Contest")
	paginateQuery, metadata, err := paginate(c, query)
	if err != nil {
		panic(err)
	}
	r := paginateQuery.Find(&userContests)
	if r.Error != nil {
		panic(r.Error)
	}
	contestResponses := make([]*forms.Contest, len(userContests))
	for i, userContest := range userContests {
		userContest.Contest, err = models.GetContestById(userContest.ContestID)
		if err != nil {
			panic(err)
		}
		publicContest, err := createPublicContest(userContest.Contest, &userContest)
		if err != nil {
			panic(err)
		}
		contestResponses[i] = &forms.Contest{
			PublicContest: *publicContest,
			Registered:    true,
			RewardPaid:    userContest.Claimed,
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": contestResponses, "metadata": metadata})
}

func (cc *contestController) contest(c *gin.Context) {
	db := _db.GetDB()
	user := middlewares.GetUser(c)
	contestId, err := strconv.Atoi(c.Param("contestId"))
	if err != nil {
		handleBadRequest(c, err)
		return
	}
	contest, err := models.GetContestById(uint(contestId))
	if err != nil {
		panic(err)
	}
	if contest == nil {
		handleError(c, nil, "contest not found", http.StatusNotFound)
		return
	}
	var userContest = &models.UserContest{ContestID: contest.ID, UserID: user.ID}
	db.Where(userContest).Find(userContest)
	if userContest.ID == 0 {
		handleBadRequestWithMessage(c, nil, "you are not part of this contest")
		return
	}
	userContest.Contest = contest
	publicContest, err := createPublicContest(contest, userContest)
	if err != nil {
		panic(err)
	}
	contestResponse := &forms.Contest{
		PublicContest: *publicContest,
		Registered:    true,
		RewardPaid:    userContest.Claimed,
	}
	c.JSON(http.StatusOK, contestResponse)
}

func (cc *contestController) reward(c *gin.Context) {
	db := _db.GetDB()
	user := middlewares.GetUser(c)
	contestId, err := strconv.Atoi(c.Param("contestId"))
	if err != nil {
		handleBadRequest(c, err)
		return
	}
	now := time.Now()
	contest, err := models.GetContestById(uint(contestId))
	if err != nil {
		panic(err)
	}
	if contest == nil {
		handleError(c, nil, "contest not found", http.StatusNotFound)
		return
	}
	if now.Before(contest.End) {
		handleBadRequestWithMessage(c, nil, "contest not ended")
		return
	}
	var userContest = &models.UserContest{ContestID: contest.ID, UserID: user.ID}
	db.Where(userContest).Find(userContest)
	if userContest.ID == 0 {
		handleBadRequestWithMessage(c, nil, "you are not part of this contest")
		return
	}
	userContest.Contest = contest
	if userContest.Claimed {
		handleBadRequestWithMessage(c, nil, "reward already paid")
		return
	}
	questions, err := contest.Questions()
	if err != nil {
		panic(err)
	}
	lastQuestion := (*questions)[len(*questions)-1]
	userAnswer := &models.UserAnswer{UserContestID: userContest.ID, QuestionID: lastQuestion.ID}
	err = db.Where(userAnswer).Find(userAnswer).Error
	if err != nil {
		panic(err)
	}
	var answer string
	if lastQuestion.Answer {
		answer = "A"
	} else {
		answer = "B"
	}
	if userAnswer.Answer != answer {
		handleBadRequestWithMessage(c, nil, "you are not win this contest")
		return
	}
	userCount, err := contest.UserCount()
	if err != nil {
		panic(err)
	}
	winnersCount, err := lastQuestion.TrueAnswerCount()
	if err != nil {
		panic(err)
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		userContest.Claimed = true
		if err := tx.Save(userContest).Error; err != nil {
			return err
		}
		amount := int64(contest.EntryFee) * userCount / winnersCount
		result := db.Model(user).Update("coins", gorm.Expr("coins + ?", amount))
		if result.Error != nil {
			return err
		}
		user.Coins += uint(amount)
		return nil
	})
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, createUserResponse(user))
}

func InitContestController(router *gin.RouterGroup) {
	cc := contestController{}
	router.Use(middlewares.AuthMiddleware)
	router.POST(":contestId/register", cc.register)
	router.POST(":contestId/reward", cc.reward)
	router.GET(":contestId", cc.contest)
	router.GET("coming", cc.coming)
	router.GET("ongoing", cc.ongoing)
	router.GET("archive", cc.archive)
}
