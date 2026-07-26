package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		return nil, errors.New("error loading .env file")
	}
	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
	}, nil
}
