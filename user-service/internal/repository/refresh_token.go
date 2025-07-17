package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/seacite-tech/compendium/user-service/internal/model"
)

type RefreshTokenRepository interface {
	AddRefreshToken(ctx context.Context, token model.RefreshToken) error
	TryRemoveRefreshTokenByToken(ctx context.Context, token string) (userId uuid.UUID, isReused bool, err error)
	RemoveAllRefreshTokensForUser(ctx context.Context, userId uuid.UUID) error
}
