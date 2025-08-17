package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/compendium-tech/compendium/common/pkg/log"
	myerror "github.com/compendium-tech/compendium/subscription-service/internal/error"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const billingLockTtl = 60 * time.Second

type billingLock struct {
	lock   *redislock.Lock
	userID uuid.UUID
}

func (e *billingLock) Release(ctx context.Context) {
	err := e.lock.Release(ctx)
	if err != nil {
		panic(err)
	}

	log.L(ctx).Infof("Successfully released billing lock for %s", e.userID)
}

type redisBillingLockRepository struct {
	client *redislock.Client
}

func NewRedisBillingLockRepository(rdb *redis.Client) BillingLockRepository {
	return &redisBillingLockRepository{
		client: redislock.New(rdb),
	}
}

func (r *redisBillingLockRepository) ObtainLock(ctx context.Context, userID uuid.UUID) BillingLock {
	log.L(ctx).Infof("Obtaining billing lock for %s", userID)

	lock, err := r.client.Obtain(ctx, fmt.Sprintf("billing_locks:%s", userID), billingLockTtl, nil)
	if err != nil {
		if errors.Is(err, redislock.ErrNotObtained) {
			log.L(ctx).Error("Failed to obtain billing lock")

			myerror.New(myerror.TooManyRequestsError).Throw()
		}

		panic(err)
	}

	log.L(ctx).Infof("Successfully obtained billing lock for %s", userID)

	return &billingLock{lock: lock, userID: userID}
}
