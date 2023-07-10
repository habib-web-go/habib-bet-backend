package config

import (
	"github.com/spf13/viper"
)

var config *viper.Viper

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(env string) error {
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(env)
	config.AddConfigPath("config/")
	err := config.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

func GetConfig() *viper.Viper {
	return config
}
