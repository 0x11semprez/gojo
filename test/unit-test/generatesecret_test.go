// Package unit_tests regroupe des tests unitaires ne nécessitant pas de
// base de données réelle.
//
// ATTENTION : ce fichier déclare "package unit_tests" (avec un "s") alors
// que config_test.go du même dossier déclare "package unit_test" (sans
// "s"). Deux fichiers .go d'un même dossier doivent appartenir au même
// package : ceci empêche actuellement le dossier de compiler
// ("go build ./..." échoue avec "found packages unit and unit_tests").
// De plus, GenerateSecret() est appelée ici sans import du package
// gojo/internal/cryptography/secret qui la définit. Ces deux points sont
// des anomalies préexistantes, non corrigées ici (seuls des commentaires
// ont été ajoutés à la demande) : à corriger avant de pouvoir exécuter ce test.
package unit_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"gojo/internal/cryptography/secret"
)

// TestGenerateSecret vérifie que GenerateSecret retourne une chaîne
// hexadécimale qui se décode en exactement 32 octets (la taille attendue
// du secret aléatoire).
func TestGenerateSecret(t *testing.T) {
	// Vérifie que le secret généré se décode à la taille attendue.
	// Le générateur retourne une représentation hexadécimale de 32 octets aléatoires.
	generate, err := secret.GenerateSecret()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("example of secret generated for the user: %v\n", generate)

	got, err := hex.DecodeString(generate)
	if err != nil {
		t.Fatal(err)
	}

	want := 32

	if len(got) != want {
		t.Errorf("invalid secret length: got %d bytes, want %d\n", len(got), want)
	}

	fmt.Printf("gojo loves you, thanks for testing GenerateSecret() \n")
}
