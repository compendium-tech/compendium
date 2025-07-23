package ratelimit

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ztrue/tracerr"
)

type redisRateLimiter struct {
	client *redis.Client
}

func NewRedisRateLimiter(client *redis.Client) RateLimiter {
	return &redisRateLimiter{
		client: client,
	}
}

func (r *redisRateLimiter) IsRateLimited(
	ctx context.Context, key string,
	window time.Duration, maxRequests uint) (bool, error) {
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return false, tracerr.Errorf("failed to increment rate limit counter for key %s: %w", key, err)
	}

	if count == 1 {
		err = r.client.Expire(ctx, key, window).Err()
		if err != nil {
			return false, tracerr.Errorf("failed to set expiry for rate limit key for key %s: %w", key, err)
		}
	}

	if count > int64(maxRequests) {
		return true, nil
	}

	return false, nil
}
