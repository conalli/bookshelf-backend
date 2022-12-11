package auth

import "golang.org/x/crypto/bcrypt"

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
