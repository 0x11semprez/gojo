// Package config loads the application configuration from .env files
// and exposes it as a typed struct (Config).
package config

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

// repoRootPath resolves a path relative to the repository root,
// anchored to this source file's own location (like
// database.StartMigrations does for the migrations directory)
// instead of the process's working directory, which depends on
// where `go run`/`go test`/the binary is launched from and does not
// reliably point anywhere useful.
func repoRootPath(name string) string {
	_, thisFile, _, _ := runtime.Caller(0)
	// config.go lives at api/internal/config: three levels below the
	// repository root.
	repoRoot := filepath.Join(filepath.Dir(thisFile), "..", "..", "..")
	return filepath.Join(repoRoot, name)
}

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
	// WalletEncryptionKey is the hex-encoded 32-byte AES-256-GCM key
	// used as the outer layer of a wallet's envelope encryption (see
	// package wallet). It never leaves the server, unlike the
	// password-derived inner key.
	WalletEncryptionKey string
	// GeneratorBinPath optionally overrides the path to the compiled
	// rust key generator binary (see package
	// cryptography/keygen). When empty, the binary is located
	// relative to the repository layout instead.
	GeneratorBinPath string
}

// NewConfig loads the .env file at the project root and builds the
// production configuration from environment variables.
func NewConfig() (*Config, error) {
	err := godotenv.Load(repoRootPath(".env"))
	if err != nil {
		return nil, errors.New("cannot load .env file")
	}
	return &Config{
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		DatabaseURLMigration: os.Getenv("DATABASE_URL_MIGRATION"),
		Port:                 os.Getenv("PORT"),
		WalletEncryptionKey:  os.Getenv("WALLET_ENCRYPTION_KEY"),
		GeneratorBinPath:     os.Getenv("GENERATOR_BIN_PATH"),
	}, nil
}

// NewConfigTest loads the .env.test file and builds a configuration
// dedicated to tests, isolated from the production configuration
// (test database, test migration DSN).
func NewConfigTest() (*Config, error) {
	err := godotenv.Load(repoRootPath(".env.test"))
	if err != nil {
		return nil, errors.New("cannot load .env.test file")
	}
	return &Config{
		DatabaseURL:          os.Getenv("DATABASE_URL_TEST"),
		DatabaseURLMigration: os.Getenv("DATABASE_URL_MIGRATION_TEST"),
		Port:                 os.Getenv("PORT"),
		WalletEncryptionKey:  os.Getenv("WALLET_ENCRYPTION_KEY"),
		GeneratorBinPath:     os.Getenv("GENERATOR_BIN_PATH"),
	}, nil
}
