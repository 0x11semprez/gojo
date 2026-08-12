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

// TestLogin is an integration test that starts a real App against the
// test database and exercises middleware.Login on the "toj" account
// created by TestCreateUser: the correct username/password pair
// succeeds, an unknown username fails, a misspelled "toj" fails (like
// an unknown user), and the correct username with a wrong password
// fails.
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
