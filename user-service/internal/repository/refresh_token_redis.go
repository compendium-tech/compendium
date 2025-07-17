package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/seacite-tech/compendium/user-service/internal/model" // Assuming this path is correct
	"github.com/ztrue/tracerr"
)

const (
	refreshTokenKeyPrefix  = "refresh_token:"
	userTokensKeyPrefix    = "user:refresh_tokens_idx:"
	tokenToUserIDKeyPrefix = "token:user_id:"
	maxTokensPerUser       = 5

	// Hash field names
	tokenHashField   = "token"
	userIdHashField  = "userId"
	expiryHashField  = "expiry"
	createdHashField = "createdAt"
)

type RedisRefreshTokenRepository struct {
	client *redis.Client
}

func NewRedisRefreshTokenRepository(client *redis.Client) *RedisRefreshTokenRepository {
	return &RedisRefreshTokenRepository{
		client: client,
	}
}

func (r *RedisRefreshTokenRepository) AddRefreshToken(ctx context.Context, token model.RefreshToken) error {
	tokenKey := r.createRefreshTokenKey(token.UserId, token.Token)
	userTokensKey := r.createUserTokensKey(token.UserId)
	tokenToUserIDKey := r.createTokenToUserIDKey(token.Token)

	expiryDuration := time.Until(token.Expiry)
	if expiryDuration <= 0 {
		return fmt.Errorf("refresh token has an expiry in the past or present")
	}

	score := float64(token.Expiry.Unix())

	pipe := r.client.Pipeline()

	pipe.HSet(ctx, tokenKey,
		tokenHashField, token.Token,
		userIdHashField, token.UserId.String(),
		expiryHashField, token.Expiry.Unix(),
	)
	pipe.Expire(ctx, tokenKey, expiryDuration)
	pipe.ZAdd(ctx, userTokensKey, redis.Z{Score: score, Member: token.Token})

	pipe.Set(ctx, tokenToUserIDKey, token.UserId.String(), expiryDuration)

	addCmds, err := pipe.Exec(ctx)
	if err != nil {
		return tracerr.Errorf("failed to add new token to Redis (HSet, Expire, ZAdd, and Set tokenToUserID): %w", err)
	}
	for _, cmd := range addCmds {
		if cmd.Err() != nil {
			return tracerr.Errorf("error in pipeline command for adding token: %w", cmd.Err())
		}
	}

	currentSize, err := r.client.ZCard(ctx, userTokensKey).Result()
	if err != nil {
		return tracerr.Errorf("failed to get current size of user tokens zset: %w", err)
	}

	if currentSize > maxTokensPerUser {
		numToRemove := currentSize - maxTokensPerUser

		membersToRemove, err := r.client.ZRange(ctx, userTokensKey, 0, numToRemove-1).Result()
		if err != nil {
			return tracerr.Errorf("failed to get members to remove from user tokens zset: %w", err)
		}

		if len(membersToRemove) > 0 {
			removalPipe := r.client.Pipeline()

			for _, memberTokenStr := range membersToRemove {
				removedTokenKey := r.createRefreshTokenKey(token.UserId, memberTokenStr)
				removalPipe.Del(ctx, removedTokenKey)
				removalPipe.Del(ctx, r.createTokenToUserIDKey(memberTokenStr))
			}

			removalPipe.ZRemRangeByRank(ctx, userTokensKey, 0, numToRemove-1)

			_, err = removalPipe.Exec(ctx)
			if err != nil {
				return tracerr.Errorf("failed to execute pipeline for removing old tokens: %w", err)
			}
		}
	}

	return nil
}

func (r *RedisRefreshTokenRepository) TryRemoveRefreshTokenByToken(ctx context.Context, token string) (uuid.UUID, error) {
	tokenToUserIDKey := r.createTokenToUserIDKey(token)
	userIDStr, err := r.client.Get(ctx, tokenToUserIDKey).Result()

	if err == redis.Nil {
		return uuid.Nil, nil
	}

	if err != nil {
		return uuid.Nil, tracerr.Errorf("failed to get userId for token %s: %w", token, err)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, tracerr.Errorf("failed to parse userId from string '%s': %w", userIDStr, err)
	}

	tokenKey := r.createRefreshTokenKey(userID, token)
	userTokensZSetKey := r.createUserTokensKey(userID)

	pipe := r.client.Pipeline()

	delCmd := pipe.Del(ctx, tokenKey)

	zRemCmd := pipe.ZRem(ctx, userTokensZSetKey, token)

	delTokenToUserCmd := pipe.Del(ctx, tokenToUserIDKey)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return uuid.Nil, tracerr.Errorf("failed to remove refresh token from Redis: %w", err)
	}

	// Check if any of the deletion commands actually removed something.
	// If the global mapping was deleted (delTokenToUserCmd.Val() > 0), it implies the token existed.
	if delCmd.Val() > 0 || zRemCmd.Val() > 0 || delTokenToUserCmd.Val() > 0 {
		return userID, nil // Token was found and removed, return its userId
	}

	// This case should ideally be covered by the initial redis.Nil check,
	// but as a fallback, if for some reason the token was not found despite
	// the initial Get command succeeding, we'd return false.
	return uuid.Nil, nil
}

func (r *RedisRefreshTokenRepository) createRefreshTokenKey(userId uuid.UUID, token string) string {
	return fmt.Sprintf("%s%s:%s", refreshTokenKeyPrefix, userId.String(), token)
}

func (r *RedisRefreshTokenRepository) createUserTokensKey(userId uuid.UUID) string {
	return fmt.Sprintf("%s%s", userTokensKeyPrefix, userId.String())
}

func (r *RedisRefreshTokenRepository) createTokenToUserIDKey(token string) string {
	return fmt.Sprintf("%s%s", tokenToUserIDKeyPrefix, token)
}
