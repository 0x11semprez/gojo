package database

import (
	"testing"

	"gojo/internal/config"
)

func TestStartMigrations(t *testing.T) {
	cfg, err := config.NewConfig()
	if err != nil {
		t.Errorf("error with cfg")
	}

	StartMigrations(cfg)
}
