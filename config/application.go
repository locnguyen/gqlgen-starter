package config

import (
	"github.com/spf13/viper"
)

type applicationConfig struct {
	ServerPort        string
	LogLevel          string
	StructuredLogging bool
	GoEnv             string
	DatabaseURL       string
}

var Application applicationConfig

func init() {
	// Load envars into Viper
	viper.AutomaticEnv()
	Application = applicationConfig{}

	viper.SetDefault("SERVER_PORT", 8080)
	Application.ServerPort = viper.GetString("SERVER_PORT")

	viper.SetDefault("LOG_LEVEL", "info")
	Application.LogLevel = viper.GetString("LOG_LEVEL")

	viper.SetDefault("STRUCTURED_LOGGING", true)
	Application.StructuredLogging = viper.GetBool("STRUCTURED_LOGGING")

	viper.SetDefault("GO_ENV", "development")
	Application.GoEnv = viper.GetString("GO_ENV")

	Application.DatabaseURL = viper.GetString("DATABASE_URL")
}

func (c *applicationConfig) IsDevelopment() bool {
	return c.GoEnv == "development"
}

func (c *applicationConfig) IsProduction() bool {
	return c.GoEnv == "production"
}
