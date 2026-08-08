package database

import (
	"testing"

	"gojo/internal/config"
)

func TestStartMigrations(t *testing.T) {
	cfg, err := config.NewConfigTest()
	if err != nil {
		t.Errorf("error with cfg")
	}

	StartMigrations(cfg)
}
