package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/compendium-tech/compendium/common/pkg/log"
	appErr "github.com/compendium-tech/compendium/user-service/internal/error"
	"github.com/redis/go-redis/v9"
	"github.com/ztrue/tracerr"
)

const authLockTtl = 60 * time.Second

type authLock struct {
	lock  *redislock.Lock
	email string
}

func (e *authLock) Release(ctx context.Context) error {
	err := e.lock.Release(ctx)

	if err != nil {
		return tracerr.Wrap(err)
	}

	log.L(ctx).
		WithField("email", e.email).
		Infof("Successfully released auth lock for %s", e.email)

	return nil
}

func (e *authLock) ReleaseAndHandleErr(ctx context.Context, err *error) {
	// If an error already exists, don't overwrite it with a potential lock release error.
	// The original error is usually more important.
	if *err != nil {
		if lockErr := e.Release(ctx); lockErr != nil {
			log.L(ctx).
				WithField("email", e.email).
				Warnf("Warning: Failed to release email lock, but original error already present: %v (release error: %v)\n", (*err), lockErr)
		}

		return
	}

	lockErr := e.Release(ctx)
	if lockErr != nil {
		*err = lockErr
	}
}

type redisAuthLockRepository struct {
	client *redislock.Client
}

func NewRedisAuthLockRepository(rdb *redis.Client) AuthLockRepository {
	return &redisAuthLockRepository{
		client: redislock.New(rdb),
	}
}

func (r *redisAuthLockRepository) ObtainLock(ctx context.Context, email string) (AuthLock, error) {
	logger := log.L(ctx).WithField("email", email)
	logger.Infof("Obtaining auth lock for %s", email)

	lock, err := r.client.Obtain(ctx, fmt.Sprintf("auth_locks:%s", email), authLockTtl, nil)
	if err != nil {
		if errors.Is(err, redislock.ErrNotObtained) {
			logger.Error("Failed to obtain email lock")

			return nil, appErr.New(appErr.TooManyRequestsError, "Too many requests")
		}

		return nil, tracerr.Wrap(err)
	}

	logger.Infof("Successfully obtained auth lock for %s", email)

	return &authLock{lock: lock, email: email}, nil
}
