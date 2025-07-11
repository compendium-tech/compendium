package repository

import (
	"context"
	"errors"
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

type RedisEmailLockRepository struct {
	client *redislock.Client
}

func NewRedisEmailLockRepository(rdb *redis.Client) *RedisEmailLockRepository {
	return &RedisEmailLockRepository{
		client: redislock.New(rdb),
	}
}

func (r *RedisEmailLockRepository) ObtainEmailLock(ctx context.Context, email string) (EmailLock, error) {
	lock, err := r.client.Obtain(ctx, email, emailLockTtl, nil)
	if err != nil {
		if errors.Is(err, redislock.ErrNotObtained) {
			log.L(ctx).Error("Failed to obtain email lock")
			return nil, apperr.Errorf(apperr.TooManyRequestsError, "Too many requests")
		}

		return nil, tracerr.Wrap(err)
	}

	return &emailLock{lock: lock}, nil
}
