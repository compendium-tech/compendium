package hash

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	HashPassword(password string) []byte
	IsPasswordHashValid(passwordHash []byte, password string) bool
}

type bcryptPasswordHasher struct {
	cost int
}

func NewBcryptPasswordHasher(cost int) PasswordHasher {
	return &bcryptPasswordHasher{
		cost: cost,
	}
}

func (b *bcryptPasswordHasher) HashPassword(password string) []byte {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return passwordHash
}

func (b *bcryptPasswordHasher) IsPasswordHashValid(passwordHash []byte, password string) bool {
	err := bcrypt.CompareHashAndPassword(passwordHash, []byte(password))

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false
		}

		panic(err)
	}

	return true
}
