package sec

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// zero cost use default
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 0)

	return string(bytes), err
}

func CheckHashPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
