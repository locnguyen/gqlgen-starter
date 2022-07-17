package config

import (
	"github.com/spf13/viper"
)

type applicationConfig struct {
	ServerPort string
}

var Application applicationConfig

func init() {
	viper.AutomaticEnv()
	Application = applicationConfig{}
	viper.SetDefault("SERVER_PORT", 8080)
	Application.ServerPort = viper.GetString("SERVER_PORT")
}
