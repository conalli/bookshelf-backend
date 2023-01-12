package auth

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// GenerateAPIKey generates a random URL-safe string of random length for use as an API key.
func GenerateAPIKey() (string, error) {
	key, err := uuid.NewRandom()
	return key.String(), err
}

// Hash takes a password and returns the hashed version of it.
func Hash(input string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	return string(hashed), err
}

// CheckHash takes a hashed password and compares it to an unhashed password.
func CheckHash(hash, input string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input))
	return err == nil
}
