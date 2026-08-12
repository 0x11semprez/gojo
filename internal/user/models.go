package user

import "github.com/uptrace/bun"

// User représente une ligne de la table "users" et sert de modèle à l'ORM
// Bun pour générer les requêtes SQL correspondantes.
type User struct {
	// bun.BaseModel associe cette struct à la table "users" (alias SQL "u").
	bun.BaseModel `bun:"table:users,alias:u"`

	// Id est un UUID généré par Postgres (DEFAULT gen_random_uuid() dans la
	// migration). "nullzero" indique à bun d'omettre la colonne lors d'un
	// insert quand la valeur Go est encore "", pour laisser la base
	// appliquer sa valeur par défaut.
	Id string `bun:",pk,nullzero"`
	// Username doit être unique et non nul (contraintes appliquées côté base).
	Username string `bun:",notnull,unique"`
	// Secret stocke le hash bcrypt du mot de passe, jamais le mot de passe en clair.
	Secret []byte `bun:",notnull"`
}
