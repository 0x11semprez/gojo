// Package middleware fournit les fonctions liées à l'authentification :
// hachage des mots de passe et vérification des identifiants (login).
package middleware

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

// ErrInvalidCredentials est retournée quand le nom d'utilisateur n'existe pas
// ou que le mot de passe ne correspond pas. Les deux cas retournent
// volontairement la même erreur, pour qu'un appelant (ou un attaquant) ne
// puisse pas s'en servir pour deviner quels noms d'utilisateur existent.
var ErrInvalidCredentials = errors.New("invalid username or password")

// HashPassword transforme un mot de passe en clair en un hash bcrypt,
// qui peut être stocké en base de données sans risque (colonne "secret",
// de type BYTEA, dans la table users). bcrypt intègre un sel aléatoire
// dans son résultat, donc hacher deux fois le même mot de passe ne
// produit jamais les mêmes octets.
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// Login vérifie un couple nom d'utilisateur/mot de passe par rapport à la
// base de données et retourne l'id de l'utilisateur correspondant en cas
// de succès.
func Login(ctx context.Context, db *bun.DB, username, password string) (string, error) {
	var (
		id   string
		hash []byte
	)

	// Un seul aller-retour en base : on récupère l'id et le hash bcrypt
	// stocké pour ce nom d'utilisateur. On interroge directement la table
	// plutôt que le modèle user.User pour éviter un cycle d'import
	// (le package user importe déjà middleware).
	err := db.NewSelect().
		Table("users").
		Column("id", "secret").
		Where("username = ?", username).
		Where("deleted_at IS NULL").
		Scan(ctx, &id, &hash)

	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrInvalidCredentials
	}
	if err != nil {
		return "", err
	}

	// CompareHashAndPassword re-hache le mot de passe fourni avec le sel
	// extrait du hash stocké et compare le résultat en temps constant,
	// ce qui protège contre les attaques par mesure de temps (timing attacks).
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	return id, nil
}
