package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort     string
	DatabaseURL string
	JWTSecret   string
	Environment string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)
		return nil, err
	}

	return &Config{
		AppPort:     viper.GetString("APP_PORT"),
		DatabaseURL: viper.GetString("DATABASE_URL"),
		JWTSecret:   viper.GetString("JWT_SECRET"),
		Environment: viper.GetString("ENVIRONMENT"),
	}, nil
}
