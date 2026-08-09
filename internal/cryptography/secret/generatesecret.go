package secret

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateSecret() (string, error) {
	b := make([]byte, 32)
	rand.Read(b)

	secret := hex.EncodeToString(b)

	return secret, nil
}
