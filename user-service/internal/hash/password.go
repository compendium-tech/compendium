package hash

import (
	"github.com/ztrue/tracerr"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	HashPassword(password string) ([]byte, error)
	IsPasswordHashValid(passwordHash []byte, password string) (bool, error)
}

type BcryptPasswordHasher struct {
	cost int
}

func NewBcryptPasswordHasher(cost int) *BcryptPasswordHasher {
	return &BcryptPasswordHasher{
		cost: cost,
	}
}

func (b *BcryptPasswordHasher) HashPassword(password string) ([]byte, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return passwordHash, nil
}

func (b *BcryptPasswordHasher) IsPasswordHashValid(passwordHash []byte, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(passwordHash, []byte(password))

	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}

		return false, tracerr.Wrap(err)
	}

	return true, nil
}
