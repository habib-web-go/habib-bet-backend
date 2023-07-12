package db

import (
	"fmt"
	"github.com/habib-web-go/habib-bet-backend/config"
	"github.com/lib/pq"
	"gorm.io/gorm"
)
import "gorm.io/driver/postgres"

var db *gorm.DB

func Init() error {
	conf := config.GetConfig()
	dsnFormat := "host=%s port=%d user=%s password=%s dbname=%s"
	dsn := fmt.Sprintf(
		dsnFormat,
		conf.GetString("db.host"),
		conf.GetInt("db.port"),
		conf.GetString("db.user"),
		conf.GetString("db.password"),
		conf.GetString("db.name"),
	)
	dbTmp, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	db = dbTmp
	return nil
}

func GetDB() *gorm.DB {
	return db
}

func IsDuplicateKeyError(err error) bool {
	pgErr, ok := err.(*pq.Error)
	if ok {
		return pgErr.Code.Name() == "unique_violation"

	}
	return false
}
