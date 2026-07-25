package middleware

import "testing"

func TestGenerateSecret(t *testing.T) {
	secret, err := GeneratePassword()
	if err != nil {
		panic(err)
	}

	want := make([]byte, 32)

	if want != btoa(secret) {
		t.Errorf("%d something different with %v", want, secret)
	}
}

func btoa(sentence *string) []byte {
}

// we put a string and then we have an array of byte
// I need to create a function that take the string, remove "", put one letter by one letter in the test
