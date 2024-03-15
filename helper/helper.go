package helper

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func PasswordHasing(password string) (string, error) {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return "", errors.New("generate hash password issue")
	}

	hash := string(hashPassword)

	return hash, nil
}
