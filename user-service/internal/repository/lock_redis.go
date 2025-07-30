package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/compendium-tech/compendium/common/pkg/log"
	"github.com/compendium-tech/compendium/user-service/internal/error"
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

			return nil, myerror.New(myerror.TooManyRequestsError)
		}

		return nil, tracerr.Wrap(err)
	}

	logger.Infof("Successfully obtained auth lock for %s", email)

	return &authLock{lock: lock, email: email}, nil
}
