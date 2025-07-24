package repository

import (
	"context"

	"github.com/compendium-tech/compendium/user-service/internal/model"
	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	GetRefreshToken(ctx context.Context, token string) (*model.RefreshToken, bool, error)
	CreateRefreshToken(ctx context.Context, token model.RefreshToken) error

	RemoveRefreshToken(ctx context.Context, token string, userID uuid.UUID) error
	RemoveAllRefreshTokensForUser(ctx context.Context, userID uuid.UUID) error

	GetAllRefreshTokensForUser(ctx context.Context, userID uuid.UUID) ([]model.RefreshToken, error)
	GetRefreshTokenBySessionID(ctx context.Context, sessionID uuid.UUID) (*model.RefreshToken, bool, error) // New method
}
