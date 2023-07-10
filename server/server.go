package server

import "github.com/habib-web-go/habib-bet-backend/config"

func Init() error {
	conf := config.GetConfig()
	r := NewRouter()
	err := r.Run(conf.GetString("server.addr"))
	if err != nil {
		return err
	}
	return nil
}
