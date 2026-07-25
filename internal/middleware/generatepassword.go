package middleware

import (
	"crypto/rand"
	"encoding/hex"
)

func GeneratePassword() string {
	b := make([]byte, 32)
	rand.Read(b)

	secret := hex.EncodeToString(b)

	return secret
}
