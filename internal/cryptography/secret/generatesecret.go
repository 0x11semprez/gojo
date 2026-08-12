// Package secret fournit des utilitaires de génération de valeurs
// aléatoires cryptographiquement sûres (ex: secrets d'application).
package secret

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateSecret génère un secret aléatoire de 32 octets (256 bits) à
// partir du générateur cryptographique sécurisé crypto/rand, puis
// l'encode en chaîne hexadécimale (64 caractères) pour un usage simple
// (ex: stockage dans une variable d'environnement).
func GenerateSecret() (string, error) {
	b := make([]byte, 32)
	rand.Read(b)

	secret := hex.EncodeToString(b)

	return secret, nil
}
