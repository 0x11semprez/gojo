// Package unit_test groups unit tests that do not require a real
// database.
package unit_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"gojo/internal/cryptography/secret"
)

// TestGenerateSecret verifies that GenerateSecret returns a
// hexadecimal string that decodes to exactly 32 bytes (the expected
// size of the random secret).
func TestGenerateSecret(t *testing.T) {
	// Verify that the generated secret decodes to the expected size.
	// The generator returns a hexadecimal representation of 32 random bytes.
	generate, err := secret.GenerateSecret()
	if err != nil {
		t.Fatalf("failed to generate secret: %v", err)
	}

	fmt.Printf("generated secret sample: %v\n", generate)

	got, err := hex.DecodeString(generate)
	if err != nil {
		t.Fatalf("failed to decode generated secret as hexadecimal: %v", err)
	}

	want := 32

	if len(got) != want {
		t.Errorf("invalid secret length: got %d bytes, want %d\n", len(got), want)
	}

	fmt.Printf("secret generation test completed successfully\n")
}
