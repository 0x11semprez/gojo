// Package unit_test regroupe les tests unitaires ne nécessitant pas de
// base de données réelle (contrairement au package integrationtest).
package unit_test

import (
	"fmt"
	"testing"

	"gojo/internal/config"
)

// TestNewConfig vérifie que config.NewConfig charge le fichier .env sans
// erreur et affiche les valeurs Port/DatabaseURL/DatabaseURLMigration
// obtenues pour une inspection visuelle ; il ne vérifie pas leur contenu.
func TestNewConfig(t *testing.T) {
	cfg, err := config.NewConfig()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("thank you for testing config, this is crucial to launch our application\n")
	fmt.Printf("your listening on port -> %q\n", cfg.Port)
	fmt.Printf("your database url is -> %q\n", cfg.DatabaseURL)
	fmt.Printf("your database url for migrations is -> %q\n", cfg.DatabaseURLMigration)
}
