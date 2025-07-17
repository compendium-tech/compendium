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
	revokedTokensKeyPrefix = "revoked_token:"
	maxTokensPerUser       = 5

	revokedTokenTTL = 3 * 24 * time.Hour

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

func (r *RedisRefreshTokenRepository) TryRemoveRefreshTokenByToken(ctx context.Context, token string) (uuid.UUID, bool, error) {
	tokenToUserIDKey := r.createTokenToUserIDKey(token)
	revokedTokenKey := r.createRevokedTokenKey(token)

	// First, check if the token is in the active mapping
	userIDStr, err := r.client.Get(ctx, tokenToUserIDKey).Result()

	if err == redis.Nil {
		// Token is NOT active. Now check if it's in the revoked list.
		revokedUserIDStr, err := r.client.Get(ctx, revokedTokenKey).Result()
		if err == nil {
			// Token found in revoked list! This is the reuse detection.
			parsedUserID, parseErr := uuid.Parse(revokedUserIDStr)
			if parseErr != nil {
				return uuid.Nil, false, tracerr.Errorf("failed to parse revoked userId from string '%s': %w", revokedUserIDStr, parseErr)
			}
			return parsedUserID, true, nil // Returns userID and `true` for `isReused`
		}
		if err != redis.Nil {
			// Some other error when checking revoked tokens
			return uuid.Nil, false, tracerr.Errorf("failed to check revoked token %s: %w", token, err)
		}
		// Token is neither active nor in revoked list. Truly not found.
		return uuid.Nil, false, nil // Returns uuid.Nil and `false` for `isReused`
	}

	if err != nil {
		return uuid.Nil, false, tracerr.Errorf("failed to get userId for token %s: %w", token, err)
	}

	// Token IS active, proceed with removal
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, false, tracerr.Errorf("failed to parse userId from string '%s': %w", userIDStr, err)
	}

	tokenKey := r.createRefreshTokenKey(userID, token)
	userTokensZSetKey := r.createUserTokensKey(userID)

	pipe := r.client.Pipeline()

	pipe.Del(ctx, tokenKey)
	pipe.ZRem(ctx, userTokensZSetKey, token)
	pipe.Del(ctx, tokenToUserIDKey)

	// NEW: Add the token to a revoked list with its original expiry TTL
	pipe.Set(ctx, revokedTokenKey, userID.String(), revokedTokenTTL)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return uuid.Nil, false, tracerr.Errorf("failed to remove refresh token from Redis: %w", err)
	}

	// If the active token was successfully removed, it was NOT a reuse.
	return userID, false, nil // Returns userID and `false` for `isReused`
}

func (r *RedisRefreshTokenRepository) RemoveAllRefreshTokensForUser(ctx context.Context, userId uuid.UUID) error {
	userTokensKey := r.createUserTokensKey(userId)

	// Get all tokens for the user from the ZSET
	tokens, err := r.client.ZRange(ctx, userTokensKey, 0, -1).Result()
	if err != nil {
		return tracerr.Errorf("failed to get all tokens for user %s: %w", userId, err)
	}

	pipe := r.client.Pipeline()

	// Delete each individual refresh token and its token-to-userID mapping
	for _, token := range tokens {
		pipe.Del(ctx, r.createRefreshTokenKey(userId, token))
		pipe.Del(ctx, r.createTokenToUserIDKey(token))
	}

	// Delete the user's ZSET itself
	pipe.Del(ctx, userTokensKey)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return tracerr.Errorf("failed to remove all refresh tokens for user %s: %w", userId, err)
	}

	return nil
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

func (r *RedisRefreshTokenRepository) createRevokedTokenKey(token string) string {
	return fmt.Sprintf("%s%s", revokedTokensKeyPrefix, token)
}
