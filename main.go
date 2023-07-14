package main

import (
	"flag"
	"github.com/habib-web-go/habib-bet-backend/config"
	"github.com/habib-web-go/habib-bet-backend/db"
	"github.com/habib-web-go/habib-bet-backend/models"
	"github.com/habib-web-go/habib-bet-backend/server"
)

func main() {
	profile := flag.String("profile", "development", "config profile")
	command := flag.String(
		"command",
		"run_server",
		"command to execute. possible values: run_server,migration,create_contest",
	)
	flag.Parse()
	err := config.Init(*profile)
	if err != nil {
		panic(err)
	}
	err = db.Init()
	if err != nil {
		panic(err)
	}
	if *command == "run_server" {
		err = server.Init()
		if err != nil {
			panic(err)
		}
	}
	if *command == "migrate" {
		err = models.AutoMigrate()
		if err != nil {
			panic(err)
		}
	}
	if *command == "create_contest" {
		err = models.CreateContest()
		if err != nil {
			panic(err)
		}
	}
}
