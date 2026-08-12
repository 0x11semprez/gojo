// Package integrationtest groups the integration tests that exercise
// the application against a real test database (.env.test).
// Files are prefixed with a number (0_, 1_, 2_, 3_) because Go tests
// within a package run in declaration order within a file, and in
// alphabetical file order: some tests depend on data created by
// earlier ones (see 1_create_user_test.go).
package integrationtest

import (
	"testing"

	"gojo/internal/config"
	"gojo/internal/database"
)

// TestStartMigrations is an integration test that runs
// database.StartMigrations against the test database (config loaded
// from .env.test) and checks that there is no configuration error.
// It does not verify the resulting schema.
func TestStartMigrations(t *testing.T) {
	cfg, err := config.NewConfigTest()
	if err != nil {
		t.Errorf("failed to load test configuration: %v", err)
	}

	database.StartMigrations(cfg)
}
