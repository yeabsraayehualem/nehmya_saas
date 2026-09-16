package utils

import "golang.org/x/crypto/bcrypt"

// HashText hashes a plain-text value (e.g. a password) with bcrypt.
func HashText(text string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// CompareHash reports whether the plain-text value matches the bcrypt hash.
func CompareHash(hash, text string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(text))
}
