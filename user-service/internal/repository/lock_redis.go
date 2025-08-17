package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"

	"github.com/compendium-tech/compendium/common/pkg/log"
	myerror "github.com/compendium-tech/compendium/user-service/internal/error"
)

const authLockTtl = 60 * time.Second

type authLock struct {
	lock  *redislock.Lock
	email string
}

func (e *authLock) Release(ctx context.Context) {
	err := e.lock.Release(ctx)
	if err != nil {
		panic(err)
	}

	log.L(ctx).
		WithField("email", e.email).
		Infof("Successfully released auth lock for %s", e.email)
}

type redisAuthLockRepository struct {
	client *redislock.Client
}

func NewRedisAuthLockRepository(rdb *redis.Client) AuthLockRepository {
	return &redisAuthLockRepository{
		client: redislock.New(rdb),
	}
}

func (r *redisAuthLockRepository) ObtainLock(ctx context.Context, email string) AuthLock {
	logger := log.L(ctx).WithField("email", email)
	logger.Infof("Obtaining auth lock for %s", email)

	lock, err := r.client.Obtain(ctx, fmt.Sprintf("auth_locks:%s", email), authLockTtl, nil)
	if err != nil {
		if errors.Is(err, redislock.ErrNotObtained) {
			logger.Error("Failed to obtain email lock")

			myerror.New(myerror.TooManyRequestsError).Throw()
		}

		panic(err)
	}

	logger.Infof("Successfully obtained auth lock for %s", email)

	return &authLock{lock: lock, email: email}
}
