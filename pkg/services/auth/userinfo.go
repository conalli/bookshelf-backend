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

// HashPassword takes a password and returns the hashed version of it.
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

// CheckHashedPassword takes a hashed password and compares it to an unhashed password.
func CheckHashedPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
