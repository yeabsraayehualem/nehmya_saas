package utils

import (
	"crypto/rand"
	"errors"
)

// GenerateSecret returns n cryptographically random bytes, suitable for use
// as a session store secret.
func GenerateSecret(n int) ([]byte, error) {
	if n <= 0 {
		return nil, errors.New("secret size must be positive")
	}
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}
