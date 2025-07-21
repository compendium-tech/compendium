package repository

import (
	"context"

	"github.com/compendium-tech/compendium/user-service/internal/model"
	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	AddRefreshToken(ctx context.Context, token model.RefreshToken) error
	TryRemoveRefreshTokenByToken(ctx context.Context, token string) (userId uuid.UUID, isReused bool, err error)
	RemoveAllRefreshTokensForUser(ctx context.Context, userId uuid.UUID) error
}
