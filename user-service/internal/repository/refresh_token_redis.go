package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/compendium-tech/compendium/user-service/internal/model"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	refreshTokenKeyPrefix  = "refresh_token:"
	userTokensKeyPrefix    = "user:refresh_tokens_idx:"
	revokedTokensKeyPrefix = "revoked_token:"
	maxTokensPerUser       = 5
	revokedTokenTTL        = 3 * 24 * time.Hour
	tokenHashField         = "token"
	userIdHashField        = "userId"
	expiryHashField        = "expiry"
	sessionIdHashField     = "sessionId"
)

type redisRefreshTokenRepository struct {
	client *redis.Client
}

func NewRedisRefreshTokenRepository(client *redis.Client) RefreshTokenRepository {
	return &redisRefreshTokenRepository{
		client: client,
	}
}

func (r *redisRefreshTokenRepository) AddRefreshToken(ctx context.Context, token model.RefreshToken) error {
	tokenKey := refreshTokenKeyPrefix + token.Token
	tokenDetails := map[string]interface{}{
		tokenHashField:     token.Token,
		userIdHashField:    token.UserId.String(),
		sessionIdHashField: token.SessionID.String(),
		expiryHashField:    token.Expiry.Unix(),
	}

	pipe := r.client.TxPipeline()

	pipe.HSet(ctx, tokenKey, tokenDetails)
	pipe.ExpireAt(ctx, tokenKey, token.Expiry)

	userTokensKey := userTokensKeyPrefix + token.UserId.String()
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
		return nil, false, redis.Nil
	}

	userId, err := uuid.Parse(details[userIdHashField])
	if err != nil {
		return nil, false, fmt.Errorf("failed to parse userId: %w", err)
	}

	sessionId, err := uuid.Parse(details[sessionIdHashField])
	if err != nil {
		return nil, false, fmt.Errorf("failed to parse sessionId: %w", err)
	}

	expiryUnix, err := strconv.ParseInt(details[expiryHashField], 10, 64)
	if err != nil {
		return nil, false, fmt.Errorf("failed to parse expiry: %w", err)
	}

	// Check if the token has expired
	if time.Now().After(time.Unix(expiryUnix, 0)) {
		return nil, false, nil
	}

	refreshToken := model.RefreshToken{
		UserId:    userId,
		Token:     details[tokenHashField],
		SessionID: sessionId,
		Expiry:    time.Unix(expiryUnix, 0),
	}

	revokedKey := revokedTokensKeyPrefix + tokenString
	revokedStatus, err := r.client.Exists(ctx, revokedKey).Result()
	if err != nil {
		fmt.Printf("Warning: Failed to check revoked status for token %s: %v\n", tokenString, err)
	}

	isRevoked := revokedStatus > 0

	return &refreshToken, isRevoked, nil
}

func (r *redisRefreshTokenRepository) RemoveRefreshToken(ctx context.Context, tokenString string, userId uuid.UUID) error {
	tokenKey := refreshTokenKeyPrefix + tokenString
	userTokensKey := userTokensKeyPrefix + userId.String()
	revokedKey := revokedTokensKeyPrefix + tokenString

	pipe := r.client.TxPipeline()

	pipe.Del(ctx, tokenKey)

	pipe.ZRem(ctx, userTokensKey, tokenString)

	pipe.Set(ctx, revokedKey, "true", revokedTokenTTL)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to remove refresh token: %w", err)
	}
	return nil
}

func (r *redisRefreshTokenRepository) RemoveAllRefreshTokensForUser(ctx context.Context, userId uuid.UUID) error {
	userTokensKey := userTokensKeyPrefix + userId.String()

	tokens, err := r.client.ZRange(ctx, userTokensKey, 0, -1).Result()
	if err != nil {
		return fmt.Errorf("failed to get all tokens for user %s: %w", userId.String(), err)
	}

	if len(tokens) == 0 {
		return nil
	}

	pipe := r.client.TxPipeline()

	for _, token := range tokens {
		pipe.Del(ctx, refreshTokenKeyPrefix+token)
		pipe.Set(ctx, revokedTokensKeyPrefix+token, "true", revokedTokenTTL)
	}

	pipe.Del(ctx, userTokensKey)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to remove all refresh tokens for user %s: %w", userId.String(), err)
	}
	return nil
}
