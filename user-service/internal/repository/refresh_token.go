package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/seacite-tech/compendium/user-service/internal/model"
)

type RefreshTokenRepository interface {
	AddRefreshToken(ctx context.Context, token model.RefreshToken) error
	GetRefreshTokenByUserIdAndToken(ctx context.Context, userId uuid.UUID, token string) (*model.RefreshToken, error)
	RemoveRefreshTokenByUserIdAndToken(ctx context.Context, userId uuid.UUID, token string) error
}
