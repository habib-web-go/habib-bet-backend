package models

import (
	_db "github.com/habib-web-go/habib-bet-backend/db"
	"gorm.io/gorm"
	"time"
)

type Question struct {
	gorm.Model
	OptionA   string   `gorm:"notNull"`
	OptionB   string   `gorm:"notNull"`
	Answer    bool     `gorm:"notNull"` // true stand for option a
	Order     uint     `gorm:"notNull"`
	ContestId uint     `gorm:"notNull"`
	Contest   *Contest `gorm:"notNull"`
}

func (q *Question) Interval(questionCount int) (time.Time, time.Time) {
	step := q.Contest.End.Sub(q.Contest.Start).Nanoseconds() / int64(questionCount)
	start := q.Contest.Start.Add(time.Duration(int64(q.Order) * step))
	end := q.Contest.Start.Add(time.Duration(int64(q.Order+1) * step))
	return start, end
}

func (q *Question) TrueAnswerCount() (int64, error) {
	var count int64
	db := _db.GetDB()
	var answer string
	if q.Answer {
		answer = "A"
	} else {
		answer = "B"
	}
	result := db.Model(&UserAnswer{}).Where(
		&UserAnswer{QuestionID: q.ID, Answer: answer},
	).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func GetQuestionById(id uint) (*Question, error) {
	db := _db.GetDB()
	question := Question{}
	result := db.Find(&question, id)
	if result.Error != nil {
		return nil, result.Error
	}
	if question.ID == 0 {
		return nil, nil
	}
	return &question, nil
}
