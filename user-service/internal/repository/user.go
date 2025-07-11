package repository

import (
	"context"
	"time"

	"github.com/seacite-tech/compendium/user-service/internal/model"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateIsEmailVerifiedByEmail(ctx context.Context, email string, isEmailVerified bool) error
	UpdatePasswordHashAndCreatedAtByEmail(ctx context.Context, email string, passwordHash []byte, createdAt time.Time) error
	CreateUser(ctx context.Context, user model.User) error
}
