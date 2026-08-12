package integrationtest

import (
	"context"
	"testing"
	"time"

	"gojo/internal/config"
	"gojo/internal/database"
)

// TestConnectDb est un test d'intégration qui vérifie que database.ConnectDB
// ouvre une connexion fonctionnelle à la base de données réelle (config
// chargée depuis .env) et que la connexion répond à un ping en moins de
// 5 secondes.
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
