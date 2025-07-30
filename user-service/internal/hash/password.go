package hash

import (
	"errors"
	"github.com/ztrue/tracerr"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	HashPassword(password string) ([]byte, error)
	IsPasswordHashValid(passwordHash []byte, password string) (bool, error)
}

type bcryptPasswordHasher struct {
	cost int
}

func NewBcryptPasswordHasher(cost int) PasswordHasher {
	return &bcryptPasswordHasher{
		cost: cost,
	}
}

func (b *bcryptPasswordHasher) HashPassword(password string) ([]byte, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return passwordHash, nil
}

func (b *bcryptPasswordHasher) IsPasswordHashValid(passwordHash []byte, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(passwordHash, []byte(password))

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}

		return false, tracerr.Wrap(err)
	}

	return true, nil
}
