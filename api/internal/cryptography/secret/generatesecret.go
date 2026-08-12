// Package secret provides utilities for generating cryptographically
// secure random values (e.g. application secrets).
package secret

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateSecret generates a random 32-byte (256-bit) secret using
// the crypto/rand secure random generator, then encodes it as a
// hexadecimal string (64 characters) for simple usage (e.g. storing
// it in an environment variable).
func GenerateSecret() (string, error) {
	b := make([]byte, 32)
	rand.Read(b)

	secret := hex.EncodeToString(b)

	return secret, nil
}
