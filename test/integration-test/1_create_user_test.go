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

// tojPassword est le mot de passe du compte de test "toj" créé ci-dessous.
// TestLogin et TestDeleteUser réutilisent ce même compte "toj" au lieu
// d'en créer un chacun, ils ont donc aussi besoin du mot de passe en clair.
// Les tests de ce package s'exécutent dans l'ordre des fichiers (les
// préfixes numériques l'imposent : 1_create, 2_login, 3_delete), donc
// "toj" existe déjà quand ces tests s'exécutent, et TestDeleteUser ne le
// supprime qu'une fois que TestLogin a fini de l'utiliser.
const tojPassword = "toj-integration-test-password"

// TestCreateUser est un test d'intégration qui démarre une vraie App
// (config + connexion base + serveur) et exerce user.CreateUser contre la
// base de données réelle : création d'un utilisateur normal, création d'un
// utilisateur avec un mot de passe vide (bcrypt hache sans problème une
// chaîne vide, donc ce cas réussit -- il n'y a pas de validation côté
// application qui le rejette), et rejet d'un nom d'utilisateur en double.
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

	// rules décrit le scénario testé, wantErr indique si CreateUser est
	// censé échouer pour ce scénario.
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
