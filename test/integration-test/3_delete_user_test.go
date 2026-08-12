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

// idForUsername recherche l'id d'un utilisateur non supprimé à partir de
// son nom d'utilisateur ; utilisé pour retrouver les comptes "toj" et
// "coeurco" créés par TestCreateUser.
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

// TestDeleteUser est un test d'intégration qui démarre une vraie App
// contre la base de données de test et exerce user.DeleteUser sur : le
// compte "toj" (créé par TestCreateUser, et déjà utilisé par TestLogin à
// ce stade), le compte "coeurco" (créé par TestCreateUser avec un mot de
// passe vide), un compte déjà supprimé, et un id de compte qui n'a jamais
// existé.
//
// DeleteUser est une suppression logique (soft delete) qui ne met à jour
// que les lignes où deleted_at IS NULL ; avec l'implémentation actuelle,
// elle ne retourne donc pas d'erreur quand 0 ligne correspond -- supprimer
// un compte déjà supprimé ou inconnu est une opération silencieuse
// (no-op), pas un échec.
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
