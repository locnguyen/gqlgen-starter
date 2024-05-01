package config

import (
	"github.com/spf13/viper"
)

type ApplicationConfig struct {
	DatabaseURL        string
	GoEnv              string
	LogLevel           string
	NatsURL            string
	RedisAuth          string
	RedisURL           string
	ServerPort         string
	StructuredLogging  bool
	DefaultReferrerURL string
}

var Application ApplicationConfig

func (c *ApplicationConfig) IsDevelopment() bool {
	return c.GoEnv == "development"
}

func init() {
	viper.AutomaticEnv()
	Application = ApplicationConfig{}

	Application.DatabaseURL = viper.GetString("DATABASE_URL")

	viper.SetDefault("GO_ENV", "development")
	Application.GoEnv = viper.GetString("GO_ENV")

	viper.SetDefault("LOG_LEVEL", "debug")
	Application.LogLevel = viper.GetString("LOG_LEVEL")

	Application.NatsURL = viper.GetString("NATS_URL")

	Application.RedisAuth = viper.GetString("REDIS_AUTH")

	Application.RedisURL = viper.GetString("REDIS_URL")

	Application.ServerPort = viper.GetString("SERVER_PORT")

	viper.SetDefault("SERVER_PORT", 8080)
	Application.ServerPort = viper.GetString("SERVER_PORT")

	viper.SetDefault("STRUCTURED_LOGGING", true)
	Application.StructuredLogging = viper.GetBool("STRUCTURED_LOGGING")

	viper.SetDefault("WEBAPP_URL", "http://localhost:9000")
	Application.DefaultReferrerURL = viper.GetString("WEBAPP_URL")
}
