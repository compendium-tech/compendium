package repository

import (
	"context"
	"time"

	"github.com/compendium-tech/compendium/user-service/internal/model"
	"github.com/google/uuid"
)

// UserRepository defines the interface for interacting with user data.
// This repository handles user verification status very specifically:
//
// GetUser and FindUserByVerifiedEmail will only return a user if their email is verified,
// treating unverified users as non-existent.
//
// However, FindUserByEmail is an exception, designed for the authentication service to work
// directly with unverified users during sign-up or password reset flows.
type UserRepository interface {
	GetUser(ctx context.Context, id uuid.UUID) *model.User
	FindUserByEmail(ctx context.Context, email string) *model.User
	FindUserByVerifiedEmail(ctx context.Context, email string) *model.User
	UpdateUserName(ctx context.Context, id uuid.UUID, name string) model.User
	UpdateIsEmailVerifiedByEmail(ctx context.Context, email string, isEmailVerified bool)
	UpdatePasswordHash(ctx context.Context, id uuid.UUID, passwordHash []byte)
	UpdatePasswordHashAndCreatedAt(ctx context.Context, id uuid.UUID, passwordHash []byte, createdAt time.Time)
	CreateUser(ctx context.Context, user model.User, isEmailTaken *bool)
}
