// Package config loads the application configuration from .env files
// and exposes it as a typed struct (Config).
package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// Config groups the configuration parameters needed to start the
// application (database access, listen port).
type Config struct {
	// DatabaseURL is the DSN used for normal application queries.
	DatabaseURL string
	// DatabaseURLMigration is the DSN used specifically to run
	// migrations (may require different privileges than DatabaseURL).
	DatabaseURLMigration string
	// Port is the port the HTTP server listens on.
	Port string
}

// NewConfig loads the .env file at the project root and builds the
// production configuration from environment variables.
func NewConfig() (*Config, error) {
	// "../../.env" is resolved relative to the process's working
	// directory at runtime, not relative to this file's directory.
	err := godotenv.Load("../../.env")
	if err != nil {
		return nil, errors.New("cannot load .env file")
	}
	return &Config{
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		DatabaseURLMigration: os.Getenv("DATABASE_URL_MIGRATION"),
		Port:                 os.Getenv("PORT"),
	}, nil
}

// NewConfigTest loads the .env.test file and builds a configuration
// dedicated to tests, isolated from the production configuration
// (test database, test migration DSN).
func NewConfigTest() (*Config, error) {
	err := godotenv.Load("../../.env.test")
	if err != nil {
		return nil, errors.New("cannot load .env.test file")
	}
	return &Config{
		DatabaseURL:          os.Getenv("DATABASE_URL_TEST"),
		DatabaseURLMigration: os.Getenv("DATABASE_URL_MIGRATION_TEST"),
		Port:                 os.Getenv("PORT"),
	}, nil
}
