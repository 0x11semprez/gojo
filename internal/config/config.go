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
		return nil, errors.New("error loading .env file")
	}
	return &Config{
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		DatabaseURLMigration: os.Getenv("DatabaseURLMigration"),
		Port:                 os.Getenv("PORT"),
		Jwt:                  os.Getenv("JWT_TOKEN"),
	}, nil
}
