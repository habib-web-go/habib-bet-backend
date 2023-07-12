package models

import _db "github.com/habib-web-go/habib-bet-backend/db"

func AutoMigrate() error {
	db := _db.GetDB()
	err := db.AutoMigrate(&User{})
	if err != nil {
		return err
	}
	return nil
}
