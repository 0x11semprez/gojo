package integrationtest

import (
	"net/http"
	"testing"

	"gojo/internal/app"
	"gojo/internal/config"
	"gojo/internal/cryptography/secret"
	"gojo/internal/database"
	"gojo/internal/server"
	"gojo/internal/user"
)

// tojPassword is the password for the "toj" test account created below.
// TestLogin and TestDeleteUser reuse this same "toj" account instead
// of each creating one, so they also need the plaintext password.
// Tests in this package run in file order (enforced by the numeric
// prefixes: 1_create, 2_login, 3_delete), so "toj" already exists by
// the time these tests run, and TestDeleteUser only removes it once
// TestLogin is done using it.
const tojPassword = "toj-integration-test-password"

// TestCreateUser is an integration test that starts a real App
// (config + database connection + server) and exercises user.CreateUser
// against the real database: creating a normal user, creating a user
// with an empty password (bcrypt hashes an empty string without
// issue, so this case succeeds -- there is no application-side
// validation rejecting it), and rejecting a duplicate username.
func TestCreateUser(t *testing.T) {
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

	duplicateUserPassword, err := secret.GenerateSecret()
	if err != nil {
		t.Fatal(err)
	}

	// rules describes the tested scenario, wantErr indicates whether
	// CreateUser is expected to fail for this scenario.
	tests := []struct {
		rules    string
		username string
		password string
		wantErr  bool
	}{
		{"create a simple user", "toj", tojPassword, false},
		{"create a user with an empty password", "coeurco", "", false},
		{"create a user with same username", "toj", duplicateUserPassword, true},
	}

	for _, tt := range tests {
		_, err := user.CreateUser(testApp, tt.username, tt.password)
		if (err != nil) != tt.wantErr {
			t.Errorf("%s: CreateUser() error = %v, wantErr %v", tt.rules, err, tt.wantErr)
		}
	}
}
