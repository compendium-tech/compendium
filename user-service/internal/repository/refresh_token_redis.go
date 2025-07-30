package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/compendium-tech/compendium/user-service/internal/model"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	refreshTokenKeyPrefix        = "refresh_token:"
	userTokensKeyPrefix          = "user:refresh_tokens_idx:"
	revokedTokensKeyPrefix       = "revoked_token:"
	sessionRefreshTokenKeyPrefix = "session_refresh_token:"
	maxTokensPerUser             = 5
	revokedTokenTTL              = 3 * 24 * time.Hour
	tokenHashField               = "token"
	userIDHashField              = "userID"
	expiresAtHashField           = "expiresAt"
	sessionIDHashField           = "sessionID"
	userAgentHashField           = "userAgent"
	ipAddressHashField           = "ipAddress"
	sessionCreatedAtHashField    = "sessionCreatedAt"
)

type redisRefreshTokenRepository struct {
	client *redis.Client
}

func NewRedisRefreshTokenRepository(client *redis.Client) RefreshTokenRepository {
	return &redisRefreshTokenRepository{
		client: client,
	}
}

func (r *redisRefreshTokenRepository) CreateRefreshToken(ctx context.Context, token model.RefreshToken) error {
	tokenKey := refreshTokenKeyPrefix + token.Token
	sessionRefreshTokenKey := sessionRefreshTokenKeyPrefix + token.Session.ID.String()
	tokenDetails := map[string]any{
		tokenHashField:            token.Token,
		userIDHashField:           token.UserID.String(),
		sessionIDHashField:        token.Session.ID.String(),
		expiresAtHashField:        token.ExpiresAt.Unix(),
		userAgentHashField:        token.Session.UserAgent,
		ipAddressHashField:        token.Session.IPAddress,
		sessionCreatedAtHashField: token.Session.CreatedAt.Unix(),
	}

	pipe := r.client.TxPipeline()

	pipe.SetArgs(ctx, sessionRefreshTokenKey, token.Token, redis.SetArgs{
		ExpireAt: token.ExpiresAt,
	})
	pipe.HSet(ctx, tokenKey, tokenDetails)
	pipe.ExpireAt(ctx, tokenKey, token.ExpiresAt)

	userTokensKey := userTokensKeyPrefix + token.UserID.String()
	pipe.ZAdd(ctx, userTokensKey, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: token.Token,
	})

	pipe.ZRemRangeByRank(ctx, userTokensKey, 0, -int64(maxTokensPerUser+1))

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to add refresh token: %w", err)
	}

	return nil
}

func (r *redisRefreshTokenRepository) GetRefreshToken(ctx context.Context, tokenString string) (*model.RefreshToken, bool, error) {
	tokenKey := refreshTokenKeyPrefix + tokenString
	details, err := r.client.HGetAll(ctx, tokenKey).Result()
	if err != nil {
		return nil, false, fmt.Errorf("failed to get refresh token details: %w", err)
	}

	if len(details) == 0 {
		return nil, false, nil
	}

	refreshToken, parseErr := r.parseRefreshTokenDetails(details)
	if parseErr != nil {
		return nil, false, parseErr
	}

	revokedKey := revokedTokensKeyPrefix + tokenString
	revokedStatus, err := r.client.Exists(ctx, revokedKey).Result()
	if err != nil {
		fmt.Printf("Warning: Failed to check revoked status for token %s: %v\n", tokenString, err)
	}

	isRevoked := revokedStatus > 0

	return refreshToken, isRevoked, nil
}

func (r *redisRefreshTokenRepository) GetRefreshTokenBySessionID(ctx context.Context, sessionID uuid.UUID) (*model.RefreshToken, bool, error) {
	tokenString, err := r.client.Get(ctx, sessionRefreshTokenKeyPrefix+sessionID.String()).Result()
	if errors.Is(err, redis.Nil) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("failed to get token string by session ID: %w", err)
	}

	return r.GetRefreshToken(ctx, tokenString)
}

func (r *redisRefreshTokenRepository) RemoveRefreshToken(ctx context.Context, token string, userID uuid.UUID) error {
	tokenKey := refreshTokenKeyPrefix + token
	userTokensKey := userTokensKeyPrefix + userID.String()
	revokedKey := revokedTokensKeyPrefix + token

	details, err := r.client.HGetAll(ctx, tokenKey).Result()
	if err != nil {
		return fmt.Errorf("failed to get refresh token details before removal: %w", err)
	}
	sessionID := details[sessionIDHashField]

	pipe := r.client.TxPipeline()

	pipe.Del(ctx, tokenKey)
	pipe.ZRem(ctx, userTokensKey, token)
	pipe.Set(ctx, revokedKey, "true", revokedTokenTTL)

	if sessionID != "" {
		pipe.Del(ctx, sessionRefreshTokenKeyPrefix+sessionID)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to remove refresh token: %w", err)
	}
	return nil
}

func (r *redisRefreshTokenRepository) RemoveAllRefreshTokensForUser(ctx context.Context, userID uuid.UUID) error {
	userTokensKey := userTokensKeyPrefix + userID.String()

	tokens, err := r.client.ZRange(ctx, userTokensKey, 0, -1).Result()
	if err != nil {
		return fmt.Errorf("failed to get all tokens for user %s: %w", userID.String(), err)
	}

	if len(tokens) == 0 {
		return nil
	}

	pipe := r.client.TxPipeline()

	for _, tokenString := range tokens {
		tokenKey := refreshTokenKeyPrefix + tokenString
		details, err := r.client.HGetAll(ctx, tokenKey).Result()
		if err == nil && len(details) > 0 {
			sessionID := details[sessionIDHashField]
			pipe.Del(ctx, tokenKey)
			pipe.Set(ctx, revokedTokensKeyPrefix+tokenString, "true", revokedTokenTTL)

			if sessionID != "" {
				pipe.Del(ctx, sessionRefreshTokenKeyPrefix+sessionID)
			}
		}
	}

	pipe.Del(ctx, userTokensKey)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to remove all refresh tokens for user %s: %w", userID.String(), err)
	}
	return nil
}

func (r *redisRefreshTokenRepository) GetAllRefreshTokensForUser(ctx context.Context, userID uuid.UUID) ([]model.RefreshToken, error) {
	userTokensKey := userTokensKeyPrefix + userID.String()
	tokenStrings, err := r.client.ZRange(ctx, userTokensKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token strings for user %s: %w", userID.String(), err)
	}

	var refreshTokens []model.RefreshToken
	for _, tokenString := range tokenStrings {
		tokenKey := refreshTokenKeyPrefix + tokenString
		details, err := r.client.HGetAll(ctx, tokenKey).Result()
		if err != nil {
			fmt.Printf("Warning: Failed to get details for token %s: %v\n", tokenString, err)
			continue
		}
		if len(details) == 0 {
			continue
		}

		refreshToken, parseErr := r.parseRefreshTokenDetails(details)
		if parseErr != nil {
			fmt.Printf("Warning: Failed to parse details for token %s: %v\n", tokenString, parseErr)
			continue
		}

		if time.Now().Before(refreshToken.ExpiresAt) {
			refreshTokens = append(refreshTokens, *refreshToken)
		}
	}

	return refreshTokens, nil
}

func (r *redisRefreshTokenRepository) parseRefreshTokenDetails(details map[string]string) (*model.RefreshToken, error) {
	userID, err := uuid.Parse(details[userIDHashField])
	if err != nil {
		return nil, fmt.Errorf("failed to parse userID: %w", err)
	}

	sessionID, err := uuid.Parse(details[sessionIDHashField])
	if err != nil {
		return nil, fmt.Errorf("failed to parse sessionID: %w", err)
	}

	expiresAtUnix, err := strconv.ParseInt(details[expiresAtHashField], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expiresAt: %w", err)
	}

	createdAtUnix, err := strconv.ParseInt(details[sessionCreatedAtHashField], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse createdAt: %w", err)
	}

	return &model.RefreshToken{
		UserID:    userID,
		Token:     details[tokenHashField],
		ExpiresAt: time.Unix(expiresAtUnix, 0),
		Session: model.Session{
			ID:        sessionID,
			UserAgent: details[userAgentHashField],
			IPAddress: details[ipAddressHashField],
			CreatedAt: time.Unix(createdAtUnix, 0),
		},
	}, nil
}
