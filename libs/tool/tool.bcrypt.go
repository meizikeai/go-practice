package tool

import (
	"golang.org/x/crypto/bcrypt"
)

func Encrypt(encrypt string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(encrypt), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), err
}

func CompareHashAndPassword(hashPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))

	return err == nil
}
