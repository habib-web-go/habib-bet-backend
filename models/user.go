package models

import (
	_db "github.com/habib-web-go/habib-bet-backend/db"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string `gorm:"unique;index"`
	HashedPassword string `gorm:"notNull"`
	Coins          uint   `gorm:"default:0"`
}

func CreateUser(username, password *string) (*User, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*password), 14)
	if err != nil {
		return nil, err
	}
	db := _db.GetDB()
	user := User{Username: *username}
	user.HashedPassword = string(bytes)
	result := db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, result.Error
}

func (u *User) CheckPasswordHash(password *string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(*password))
	return err == nil
}

func GetUserByUsername(username *string) (*User, error) {
	db := _db.GetDB()
	user := User{Username: *username}
	result := db.Where(&user).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if user.ID == 0 {
		return nil, nil
	}
	return &user, nil
}

func GetUserById(id uint) (*User, error) {
	db := _db.GetDB()
	user := User{}
	result := db.Find(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	if user.ID == 0 {
		return nil, nil
	}
	return &user, nil
}
