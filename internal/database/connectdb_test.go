package database

import (
	"context"
	"testing"
	"time"

	"gojo/internal/config"
)

func TestConnectDb(t *testing.T) {
	cfg, err := config.NewConfig()
	if err != nil {
		t.Fatal(err)
	}
	db, err := ConnectDB(cfg)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}
}
