package ratelimit

import (
	"context"
	"time"
)

type RateLimiter interface {
	IsRateLimited(ctx context.Context, key string, window time.Duration, maxRequests uint) bool
}
