package repository

import (
	"context"

	"github.com/compendium-tech/compendium/user-service/internal/model"
	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	GetRefreshToken(ctx context.Context, token string) (*model.RefreshToken, bool)
	CreateRefreshToken(ctx context.Context, token model.RefreshToken)

	RemoveRefreshToken(ctx context.Context, token string, userID uuid.UUID)
	RemoveAllRefreshTokensForUser(ctx context.Context, userID uuid.UUID)

	GetAllRefreshTokensForUser(ctx context.Context, userID uuid.UUID) []model.RefreshToken
	GetRefreshTokenBySessionID(ctx context.Context, sessionID uuid.UUID) (*model.RefreshToken, bool)
}
