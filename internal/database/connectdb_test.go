package database

import (
	"context"
	"log"
	"testing"
	"time"

	"gojo/internal/config"
)

func TestConnectDb(t *testing.T) {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := ConnectDB(config)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testdb := db.PingContext(ctx)

	if testdb != nil {
		t.Fatal(testdb)
	}
}
