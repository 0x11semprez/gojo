package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL              string
	DatabaseURLMigration     string
	DatabaseURLTest          string
	DatabaseURLMigrationTest string
	Port                     string
	Jwt                      string
}

func (c Config) NewConfig() (*Config, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		return nil, errors.New("error loading .env file")
	}
	return &Config{
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		DatabaseURLMigration: os.Getenv("DATABASE_URL_MIGRATION"),
		Port:                 os.Getenv("PORT"),
		Jwt:                  os.Getenv("JWT_TOKEN"),
	}, nil
}

func (c Config) NewConfigTest() (*Config, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		return nil, errors.New("error loading .env.test file")
	}
	return &Config{
		DatabaseURL:          os.Getenv("DATABASE_URL_TEST"),
		DatabaseURLMigration: os.Getenv("DATABASE_URL_MIGRATION_TEST"),
		Port:                 os.Getenv("PORT"),
		Jwt:                  os.Getenv("JWT_TOKEN"),
	}, nil
}
