package middleware

import (
	"encoding/hex"
	"testing"
)

func TestGenerateSecret(t *testing.T) {
	// Verify that the generated secret decodes to the expected size.
	// The generator returns a hexadecimal representation of 32 random bytes.
	generate, err := GenerateSecret()
	if err != nil {
		t.Fatal(err)
	}

	got, err := hex.DecodeString(generate)
	if err != nil {
		t.Fatal(err)
	}

	want := 32

	if len(got) != want {
		t.Errorf("invalid secret length: got %d bytes, want %d", len(got), want)
	}
}
