package sec

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password []byte) ([]byte, error) {
	// zero cost use default
	bytes, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	return bytes, err
}

func CheckHashPassword(hash, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)

	return err == nil
}
