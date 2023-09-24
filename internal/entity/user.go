package entity

import (
	"errors"

	gorm_provider "example.com/the-boring-to-do-list-1/pkg/gorm_provider"
	"golang.org/x/crypto/bcrypt"
)

const userPasswordBcryptCost = 14

type User struct {
	gorm_provider.EntityWithUUID

	Username          string `json:"username"`
	EncryptedPassword []byte `json:"encryptedPassword"`
}

func (u *User) SetEncryptedPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), userPasswordBcryptCost)
	if err != nil {
		return err
	}
	u.EncryptedPassword = bytes
	return nil
}

func (u *User) ComparePassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(u.EncryptedPassword, []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
