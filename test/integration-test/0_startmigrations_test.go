// Package integrationtest regroupe les tests d'intégration qui exercent
// l'application contre une vraie base de données de test (.env.test).
// Les fichiers sont préfixés par un numéro (0_, 1_, 2_, 3_) car les tests
// Go d'un même package s'exécutent dans l'ordre de déclaration au sein
// d'un fichier, et par ordre alphabétique de fichier : certains tests
// dépendent des données créées par les précédents (voir 1_create_user_test.go).
package integrationtest

import (
	"testing"

	"gojo/internal/config"
	"gojo/internal/database"
)

// TestStartMigrations est un test d'intégration qui exécute
// database.StartMigrations contre la base de données de test (config
// chargée depuis .env.test) et vérifie qu'il n'y a pas d'erreur de
// configuration. Il ne vérifie pas le schéma résultant.
func TestStartMigrations(t *testing.T) {
	cfg, err := config.NewConfigTest()
	if err != nil {
		t.Errorf("error with cfg")
	}

	database.StartMigrations(cfg)
}
