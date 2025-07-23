package repository

import (
	"context"
	"time"

	"github.com/compendium-tech/compendium/user-service/internal/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateName(ctx context.Context, id uuid.UUID, name string) (*model.User, error)
	UpdateIsEmailVerifiedByEmail(ctx context.Context, email string, isEmailVerified bool) error
	UpdatePasswordHash(ctx context.Context, id uuid.UUID, passwordHash []byte) error
	UpdatePasswordHashAndCreatedAt(ctx context.Context, id uuid.UUID, passwordHash []byte, createdAt time.Time) error
	CreateUser(ctx context.Context, user model.User) error
}
