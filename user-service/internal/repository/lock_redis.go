package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"github.com/seacite-tech/compendium/common/pkg/log"
	apperr "github.com/seacite-tech/compendium/user-service/internal/error"
	"github.com/ztrue/tracerr"
)

const emailLockTtl = 60 * time.Second

type emailLock struct {
	lock *redislock.Lock
}

func (e *emailLock) Release(ctx context.Context) error {
	return tracerr.Wrap(e.lock.Release(ctx))
}

func (e *emailLock) ReleaseAndHandleErr(ctx context.Context, err *error) {
	// If an error already exists, don't overwrite it with a potential lock release error.
	// The original error is usually more important.
	if *err != nil {
		// Log the lock release error if it occurs, but don't change the primary error.
		if lockErr := e.Release(ctx); lockErr != nil {
			// Log this, as it's a problem, but *err already holds a more primary error.
			log.L(ctx).Warnf("Warning: Failed to release email lock, but original error already present: %v (release error: %v)\n", (*err), lockErr)
		}

		return
	}

	// If no error existed, attempt to release the lock and set *err if release fails.
	lockErr := e.Release(ctx)
	if lockErr != nil {
		*err = lockErr // Only set *err if it was nil and release failed.
	}
}

type RedisEmailLockRepository struct {
	client *redislock.Client
}

func NewRedisEmailLockRepository(rdb *redis.Client) *RedisEmailLockRepository {
	return &RedisEmailLockRepository{
		client: redislock.New(rdb),
	}
}

func (r *RedisEmailLockRepository) ObtainEmailLock(ctx context.Context, actionKey string, email string) (EmailLock, error) {
	lock, err := r.client.Obtain(ctx, fmt.Sprintf("%s:%s", actionKey, email), emailLockTtl, nil)
	if err != nil {
		if errors.Is(err, redislock.ErrNotObtained) {
			log.L(ctx).Error("Failed to obtain email lock")
			return nil, apperr.Errorf(apperr.TooManyRequestsError, "Too many requests")
		}

		return nil, tracerr.Wrap(err)
	}

	return &emailLock{lock: lock}, nil
}
