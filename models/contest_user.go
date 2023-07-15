package models

import (
	_db "github.com/habib-web-go/habib-bet-backend/db"
	"gorm.io/gorm"
)

type UserContest struct {
	gorm.Model
	UserID    uint `gorm:"uniqueIndex:idx_user_contest"`
	User      *User
	ContestID uint `gorm:"uniqueIndex:idx_user_contest"`
	Contest   *Contest
	answers   *[]UserAnswer
	Claimed   bool
}

func (uc *UserContest) Answers() (*[]UserAnswer, error) {
	db := _db.GetDB()
	if uc.answers == nil {
		result := db.Where("user_contest_id = ?", uc.ID).Order("\"order\"").Find(&uc.answers)
		if result.Error != nil {
			return nil, result.Error
		}
		for i := range *uc.answers {
			(*uc.answers)[i].UserContest = uc

		}
	}
	return uc.answers, nil
}

func AddUserToContest(user *User, contest *Contest) (*UserContest, error) {
	db := _db.GetDB()
	userContest := &UserContest{
		User:    user,
		Contest: contest,
	}
	questions, err := contest.Questions()
	if err != nil {
		return nil, err
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := db.Model(user).Update("coins", gorm.Expr("coins - ?", contest.EntryFee)).Error; err != nil {
			return err
		}
		user.Coins -= contest.EntryFee // make user object sync with db
		if err := tx.Create(userContest).Error; err != nil {
			return err
		}
		for _, question := range *questions {
			userAnswer := &UserAnswer{
				UserContest: userContest,
				Question:    &question,
				Order:       question.Order,
			}
			if err := tx.Create(userAnswer).Error; err != nil {
				return err
			}
		}
		return nil
	})

	return userContest, err
}

type UserAnswer struct {
	gorm.Model
	UserContestID uint `gorm:"uniqueIndex:idx_user_contest_question"`
	UserContest   *UserContest
	QuestionID    uint `gorm:"uniqueIndex:idx_user_contest_question"`
	Question      *Question
	Answer        string // "A" for optionA and "B" fot optionB
	Order         uint   `gorm:"notNull"`
}
