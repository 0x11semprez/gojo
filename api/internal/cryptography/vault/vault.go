// Package vault provides the primitives used to envelope-encrypt
// wallet private keys before they reach the database (see package
// wallet): a password-based key derivation function (Argon2id) for
// the inner, user-specific layer, and AES-256-GCM for both the inner
// (password-derived key) and outer (server WalletEncryptionKey)
// layers.
package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// KeySize is the size in bytes of an AES-256 key.
const KeySize = 32

// SaltSize is the size in bytes of the salt fed to DeriveKey.
const SaltSize = 16

// argon2id tuning: 64 MiB memory, 1 pass, 4 lanes. Chosen as a
// reasonable default for deriving a key on every wallet
// creation/decryption request without becoming a noticeable latency
// source; revisit if profiling shows otherwise.
const (
	argon2Time      = 1
	argon2MemoryKiB = 64 * 1024
	argon2Threads   = 4
)

// DeriveKey derives a 32-byte AES-256 key from a password and salt
// using Argon2id. The same password and salt always yield the same
// key; a different salt yields an unrelated key even for the same
// password.
func DeriveKey(password string, salt []byte) [KeySize]byte {
	var key [KeySize]byte
	derived := argon2.IDKey([]byte(password), salt, argon2Time, argon2MemoryKiB, argon2Threads, KeySize)
	copy(key[:], derived)
	return key
}

// NewSalt generates a random salt suitable for DeriveKey.
func NewSalt() ([]byte, error) {
	salt := make([]byte, SaltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("cannot generate salt: %w", err)
	}
	return salt, nil
}

// Encrypt encrypts plaintext with AES-256-GCM under key. The returned
// blob is prefixed with the random nonce Decrypt needs, so it is
// self-contained and safe to store as-is.
func Encrypt(key [KeySize]byte, plaintext []byte) ([]byte, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("cannot generate nonce: %w", err)
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt reverses Encrypt: ciphertext must be a blob Encrypt
// produced under the same key (nonce-prefixed).
func Decrypt(key [KeySize]byte, ciphertext []byte) ([]byte, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, sealed := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, sealed, nil)
}

func newGCM(key [KeySize]byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("cannot build AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cannot build AES-GCM: %w", err)
	}

	return gcm, nil
}

// KeyFromHex decodes a hex-encoded 32-byte key, such as
// config.Config.WalletEncryptionKey.
func KeyFromHex(s string) ([KeySize]byte, error) {
	var key [KeySize]byte

	b, err := hex.DecodeString(s)
	if err != nil {
		return key, fmt.Errorf("invalid hex key: %w", err)
	}
	if len(b) != KeySize {
		return key, fmt.Errorf("invalid key size: got %d bytes, want %d", len(b), KeySize)
	}

	copy(key[:], b)
	return key, nil
}
