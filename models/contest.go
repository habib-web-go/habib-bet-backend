package models

import (
	_db "github.com/habib-web-go/habib-bet-backend/db"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type Contest struct {
	gorm.Model
	Name      string    `gorm:"notNull"`
	Start     time.Time `gorm:"notNull"`
	End       time.Time `gorm:"notNull"`
	EntryFee  uint      `gorm:"notNull"`
	questions *[]Question
}

func CreateContest() error {
	db := _db.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		contest := &Contest{
			Name:     "c" + string(rune(rand.Int())),
			Start:    time.Now().Add(time.Minute * 1),
			End:      time.Now().Add(time.Minute * 6),
			EntryFee: 10,
		}
		if err := db.Create(contest).Error; err != nil {
			return err
		}
		for i := 0; i < 5; i++ {
			question := Question{
				OptionA: "option A :)))",
				OptionB: "option B :)))",
				//Answer:  true,
				Answer:  rand.Intn(2) == 0,
				Order:   uint(i),
				Contest: contest,
			}
			if err := db.Create(&question).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (c *Contest) Questions() (*[]Question, error) {
	db := _db.GetDB()
	if c.questions == nil {
		result := db.Where("contest_id = ?", c.ID).Order("\"order\"").Find(&c.questions)
		if result.Error != nil {
			return nil, result.Error
		}
		for i := range *c.questions {
			(*c.questions)[i].Contest = c
		}
	}
	return c.questions, nil
}

func (c *Contest) UserCount() (int64, error) {
	db := _db.GetDB()
	var count int64
	result := db.Model(&UserContest{}).Where(&UserContest{ContestID: c.ID}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func GetContestById(id uint) (*Contest, error) {
	db := _db.GetDB()
	contest := Contest{}
	result := db.Find(&contest, id)
	if result.Error != nil {
		return nil, result.Error
	}
	if contest.ID == 0 {
		return nil, nil
	}
	return &contest, nil
}
