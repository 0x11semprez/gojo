package integrationtest

import (
	"context"
	"net/http"
	"testing"

	"gojo/internal/app"
	"gojo/internal/config"
	"gojo/internal/database"
	"gojo/internal/middleware"
	"gojo/internal/server"
)

// TestLogin est un test d'intégration qui démarre une vraie App contre la
// base de données de test et exerce middleware.Login sur le compte "toj"
// créé par TestCreateUser : le bon couple utilisateur/mot de passe
// réussit, un nom d'utilisateur inconnu échoue, "toj" mal orthographié
// échoue (comme un utilisateur inconnu), et le bon nom d'utilisateur avec
// un mauvais mot de passe échoue.
func TestLogin(t *testing.T) {
	cfg, err := config.NewConfigTest()
	if err != nil {
		t.Fatal(err)
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	srv, err := server.NewServe(cfg, mux)
	if err != nil {
		t.Fatal(err)
	}

	testApp := app.App{}.NewApp(srv, db, cfg)

	tests := []struct {
		rules    string
		username string
		password string
		wantErr  bool
	}{
		{"login with the toj account", "toj", tojPassword, false},
		{"login with an account that does not exist", "does-not-exist", tojPassword, true},
		{"login with toj misspelled", "tojj", tojPassword, true},
		{"login with toj but the wrong password", "toj", "wrong-password", true},
	}

	for _, tt := range tests {
		_, err := middleware.Login(context.Background(), testApp.Database, tt.username, tt.password)
		if (err != nil) != tt.wantErr {
			t.Errorf("%s: Login() error = %v, wantErr %v", tt.rules, err, tt.wantErr)
		}
	}
}
