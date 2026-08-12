// Package unit_test groups unit tests that do not require a real
// database (unlike the integrationtest package).
package unit_test

import (
	"fmt"
	"testing"

	"gojo/internal/config"
)

// TestNewConfig verifies that config.NewConfig loads the .env file
// without error and prints the resulting Port/DatabaseURL/DatabaseURLMigration
// values for visual inspection; it does not verify their content.
func TestNewConfig(t *testing.T) {
	cfg, err := config.NewConfig()
	if err != nil {
		t.Fatalf("failed to load configuration: %v", err)
	}

	fmt.Printf("configuration loaded successfully, values below are ready for the application to start\n")
	fmt.Printf("listening port is set to -> %q\n", cfg.Port)
	fmt.Printf("database URL is set to -> %q\n", cfg.DatabaseURL)
	fmt.Printf("database migration URL is set to -> %q\n", cfg.DatabaseURLMigration)
}
