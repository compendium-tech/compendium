package repository

import (
	"context"

	"github.com/compendium-tech/compendium/user-service/internal/model"
	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	AddRefreshToken(ctx context.Context, token model.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (model *model.RefreshToken, isReused bool, err error)
	RemoveRefreshToken(ctx context.Context, token string, userID uuid.UUID) (err error)
	RemoveAllRefreshTokensForUser(ctx context.Context, userID uuid.UUID) error
}
