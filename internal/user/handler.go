// Package user contient la logique métier et les handlers HTTP liés aux
// utilisateurs : création de compte, suppression de compte et vérification
// de bon fonctionnement du service (Health).
package user

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"gojo/internal/app"
	"gojo/internal/middleware"
)

// Health est un handler HTTP simple ("healthcheck") qui répond que le
// service fonctionne. Utile pour les sondes de supervision/monitoring.
func Health(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "All is good\n")
	if err != nil {
		log.Fatal(err)
	}
}

// CreateUser crée un nouveau compte : il hache le mot de passe avec bcrypt
// (via middleware.HashPassword) puis enregistre le nom d'utilisateur et le
// hash. Le mot de passe en clair ne touche jamais la base de données ni
// les logs. Elle retourne l'id du nouvel utilisateur (un UUID généré par
// Postgres).
func CreateUser(db *app.App, username, password string) (string, error) {
	hash, err := middleware.HashPassword(password)
	if err != nil {
		return "", fmt.Errorf("cannot hash password: %w", err)
	}

	ctx := context.Background()

	user := &User{Username: username, Secret: hash}
	if _, err := db.Database.NewInsert().Model(user).Exec(ctx); err != nil {
		return "", err
	}

	return user.Id, nil
}

// DeleteUser supprime logiquement (soft delete) un compte à partir de son
// id : elle renseigne deleted_at plutôt que de supprimer la ligne, de
// façon cohérente avec la colonne deleted_at et avec middleware.Login, qui
// exclut déjà les lignes où deleted_at est renseigné.
func DeleteUser(db *app.App, id string) error {
	ctx := context.Background()

	_, err := db.Database.NewUpdate().
		Model((*User)(nil)).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Exec(ctx)

	return err
}
