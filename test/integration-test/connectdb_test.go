package integrationtest

import (
	"context"
	"testing"
	"time"

	"gojo/internal/config"
	"gojo/internal/database"
)

// TestConnectDb is an integration test that verifies database.ConnectDB
// opens a working connection to the real database (config loaded
// from .env) and that the connection responds to a ping within
// 5 seconds.
func TestConnectDb(t *testing.T) {
	cfg, err := config.NewConfigTest()
	if err != nil {
		t.Fatal(err)
	}
	db, err := database.ConnectDB(cfg)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}
}
