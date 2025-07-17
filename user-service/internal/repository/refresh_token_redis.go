package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/seacite-tech/compendium/user-service/internal/model"
	"github.com/ztrue/tracerr"
)

const (
	refreshTokenKeyPrefix = "refresh_token:"
	userTokensKeyPrefix   = "user:refresh_tokens_idx:"
	maxTokensPerUser      = 5

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

	expiryDuration := time.Until(token.Expiry)
	if expiryDuration <= 0 {
		return fmt.Errorf("refresh token has an expiry in the past or present")
	}

	// The score for the ZSET will be the expiry time in Unix seconds.
	// This allows us to sort by expiry and easily identify the oldest.
	score := float64(token.Expiry.Unix())

	// Step 1: Add the new token and its entry to the user's ZSET
	pipe := r.client.Pipeline()

	// Store RefreshToken fields as a Redis Hash
	pipe.HSet(ctx, tokenKey,
		tokenHashField, token.Token,
		userIdHashField, token.UserId.String(),
		expiryHashField, token.Expiry.Unix(),
	)
	pipe.Expire(ctx, tokenKey, expiryDuration) // Set expiry for the hash key
	pipe.ZAdd(ctx, userTokensKey, redis.Z{Score: score, Member: token.Token})
	addCmds, err := pipe.Exec(ctx)
	if err != nil {
		return tracerr.Errorf("failed to add new token to Redis (HSet, Expire, and ZAdd): %w", err)
	}
	for _, cmd := range addCmds {
		if cmd.Err() != nil {
			return tracerr.Errorf("error in pipeline command for adding token: %w", cmd.Err())
		}
	}

	// Step 2: Check current size and trim if necessary
	currentSize, err := r.client.ZCard(ctx, userTokensKey).Result()
	if err != nil {
		return tracerr.Errorf("failed to get current size of user tokens zset: %w", err)
	}

	if currentSize > maxTokensPerUser {
		numToRemove := currentSize - maxTokensPerUser

		// Get the members (token strings) that are oldest and need to be removed.
		// ZRANGE key start stop - returns members by rank (score order).
		// 0 to numToRemove-1 gives the oldest `numToRemove` elements.
		membersToRemove, err := r.client.ZRange(ctx, userTokensKey, 0, numToRemove-1).Result()
		if err != nil {
			return tracerr.Errorf("failed to get members to remove from user tokens zset: %w", err)
		}

		if len(membersToRemove) > 0 {
			removalPipe := r.client.Pipeline()

			// Remove the individual token data for each token being removed
			for _, memberTokenStr := range membersToRemove {
				removedTokenKey := r.createRefreshTokenKey(token.UserId, memberTokenStr)
				removalPipe.Del(ctx, removedTokenKey)
			}

			// Remove these members from the sorted set itself
			removalPipe.ZRemRangeByRank(ctx, userTokensKey, 0, numToRemove-1)

			_, err = removalPipe.Exec(ctx)
			if err != nil {
				return tracerr.Errorf("failed to execute pipeline for removing old tokens: %w", err)
			}
		}
	}

	return nil
}

func (r *RedisRefreshTokenRepository) TryRemoveRefreshTokenByUserIdAndToken(ctx context.Context, userId uuid.UUID, token string) (bool, error) {
	tokenKey := r.createRefreshTokenKey(userId, token)
	userTokensZSetKey := r.createUserTokensKey(userId)

	pipe := r.client.Pipeline()

	// Delete the individual token hash data
	delCmd := pipe.Del(ctx, tokenKey)

	// Remove the token from the user's sorted set
	// Note: We use the `token` string itself as the member in the ZSET.
	zRemCmd := pipe.ZRem(ctx, userTokensZSetKey, token)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, tracerr.Errorf("failed to remove refresh token from Redis: %w", err)
	}

	// Check the number of affected keys from the Del command.
	// If the token hash was deleted, it means the token was found.
	// We also check ZRem for consistency, though Del is usually sufficient for existence.
	if delCmd.Val() > 0 || zRemCmd.Val() > 0 {
		return true, nil // Token was found and removed
	}

	return false, nil // Token was not found
}

func (r *RedisRefreshTokenRepository) createRefreshTokenKey(userId uuid.UUID, token string) string {
	return fmt.Sprintf("%s%s:%s", refreshTokenKeyPrefix, userId.String(), token)
}

func (r *RedisRefreshTokenRepository) createUserTokensKey(userId uuid.UUID) string {
	return fmt.Sprintf("%s%s", userTokensKeyPrefix, userId.String())
}
