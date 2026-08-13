// Package unit_test groups unit tests that do not require a real
// database (unlike the integrationtest package).
package unit_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"gojo/internal/cryptography/vault"
)

// TestEncryptDecryptRoundTrip verifies that Decrypt recovers the
// exact plaintext Encrypt produced, and that two encryptions of the
// same plaintext under the same key never produce the same
// ciphertext (each call must use a fresh random nonce).
func TestEncryptDecryptRoundTrip(t *testing.T) {
	var key [vault.KeySize]byte
	copy(key[:], bytes.Repeat([]byte{0x42}, vault.KeySize))

	plaintext := []byte("super secret private key material")

	ciphertext1, err := vault.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	ciphertext2, err := vault.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	if bytes.Equal(ciphertext1, ciphertext2) {
		t.Error("two encryptions of the same plaintext produced identical ciphertext (nonce reuse?)")
	}

	got, err := vault.Decrypt(key, ciphertext1)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}
	if !bytes.Equal(got, plaintext) {
		t.Errorf("Decrypt() = %q, want %q", got, plaintext)
	}
}

// TestDecryptWrongKeyFails verifies that decrypting with a different
// key than the one used to encrypt fails (AES-GCM authentication
// tag mismatch), instead of silently returning garbage.
func TestDecryptWrongKeyFails(t *testing.T) {
	var key [vault.KeySize]byte
	copy(key[:], bytes.Repeat([]byte{0x01}, vault.KeySize))

	var wrongKey [vault.KeySize]byte
	copy(wrongKey[:], bytes.Repeat([]byte{0x02}, vault.KeySize))

	ciphertext, err := vault.Encrypt(key, []byte("payload"))
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	if _, err := vault.Decrypt(wrongKey, ciphertext); err == nil {
		t.Error("Decrypt() with the wrong key succeeded, want an error")
	}
}

// TestDeriveKeyIsDeterministic verifies that DeriveKey is a pure
// function of (password, salt): the same inputs always yield the
// same key, and either a different password or a different salt
// yields a different key.
func TestDeriveKeyIsDeterministic(t *testing.T) {
	salt := bytes.Repeat([]byte{0xAA}, vault.SaltSize)
	otherSalt := bytes.Repeat([]byte{0xBB}, vault.SaltSize)

	key1 := vault.DeriveKey("correct-horse-battery-staple", salt)
	key2 := vault.DeriveKey("correct-horse-battery-staple", salt)
	if key1 != key2 {
		t.Error("DeriveKey() with the same password and salt produced different keys")
	}

	if key3 := vault.DeriveKey("a-different-password", salt); key3 == key1 {
		t.Error("DeriveKey() with a different password produced the same key")
	}

	if key4 := vault.DeriveKey("correct-horse-battery-staple", otherSalt); key4 == key1 {
		t.Error("DeriveKey() with a different salt produced the same key")
	}
}

// TestKeyFromHex verifies that KeyFromHex round-trips a valid
// 32-byte hex key and rejects malformed or wrong-length input.
func TestKeyFromHex(t *testing.T) {
	valid := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"

	key, err := vault.KeyFromHex(valid)
	if err != nil {
		t.Fatalf("KeyFromHex() error = %v", err)
	}
	if got := hex.EncodeToString(key[:]); got != valid {
		t.Errorf("KeyFromHex() round-trip = %q, want %q", got, valid)
	}

	if _, err := vault.KeyFromHex("not-hex"); err == nil {
		t.Error("KeyFromHex() with invalid hex succeeded, want an error")
	}

	if _, err := vault.KeyFromHex("aabb"); err == nil {
		t.Error("KeyFromHex() with a too-short key succeeded, want an error")
	}
}
