package integrationtest

import (
	"context"
	"net/http"
	"testing"

	"gojo/internal/app"
	"gojo/internal/config"
	"gojo/internal/cryptography/secret"
	"gojo/internal/database"
	"gojo/internal/server"
	"gojo/internal/user"
)

// idForUsername looks up the id of a non-deleted user from their
// username; used to find the "toj" and "coeurco" accounts created by
// TestCreateUser.
func idForUsername(t *testing.T, testApp *app.App, username string) string {
	t.Helper()

	var id string
	err := testApp.Database.NewSelect().
		Table("users").
		Column("id").
		Where("username = ?", username).
		Where("deleted_at IS NULL").
		Scan(context.Background(), &id)
	if err != nil {
		t.Fatalf("lookup id for %q: %v", username, err)
	}
	return id
}

// TestDeleteUser is an integration test that starts a real App
// against the test database and exercises user.DeleteUser on: the
// "toj" account (created by TestCreateUser, and already used by
// TestLogin at this point), the "coeurco" account (created by
// TestCreateUser with an empty password), an already deleted account,
// and an account id that never existed.
//
// DeleteUser is a soft delete that only updates rows where
// deleted_at IS NULL; with the current implementation, it therefore
// does not return an error when 0 rows match -- deleting an already
// deleted or unknown account is a silent no-op, not a failure.
func TestDeleteUser(t *testing.T) {
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

	tojID := idForUsername(t, testApp, "toj")
	coeurcoID := idForUsername(t, testApp, "coeurco")

	throwawayPassword, err := secret.GenerateSecret()
	if err != nil {
		t.Fatal(err)
	}
	alreadyDeletedID, err := user.CreateUser(testApp, "toj-already-deleted", throwawayPassword)
	if err != nil {
		t.Fatal(err)
	}
	if err := user.DeleteUser(testApp, alreadyDeletedID); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		rules   string
		id      string
		wantErr bool
	}{
		{"delete the toj account", tojID, false},
		{"delete the coeurco account (no password)", coeurcoID, false},
		{"delete an already deleted account", alreadyDeletedID, false},
		{"delete an account that does not exist", "00000000-0000-0000-0000-000000000000", false},
	}

	for _, tt := range tests {
		err := user.DeleteUser(testApp, tt.id)
		if (err != nil) != tt.wantErr {
			t.Errorf("%s: DeleteUser() error = %v, wantErr %v", tt.rules, err, tt.wantErr)
		}
	}
}
