package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL          string
	DatabaseURLMigration string
	Port                 string
	Jwt                  string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		return nil, errors.New("cannot load	.env file")
	}
	return &Config{
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		DatabaseURLMigration: os.Getenv("DATABASE_URL_MIGRATION"),
		Port:                 os.Getenv("PORT"),
		Jwt:                  os.Getenv("JWT_TOKEN"),
	}, nil
}
