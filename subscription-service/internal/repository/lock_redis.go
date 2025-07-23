package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/compendium-tech/compendium/common/pkg/log"
	appErr "github.com/compendium-tech/compendium/subscription-service/internal/error"
	"github.com/redis/go-redis/v9"
	"github.com/ztrue/tracerr"
)

const billingLockTtl = 60 * time.Second

type authLock struct {
	lock       *redislock.Lock
	customerId string
}

func (e *authLock) Release(ctx context.Context) error {
	err := e.lock.Release(ctx)

	if err != nil {
		return tracerr.Wrap(err)
	}

	log.L(ctx).Infof("Successfully released auth lock for %s", e.customerId)

	return nil
}

func (e *authLock) ReleaseAndHandleErr(ctx context.Context, err *error) {
	// If an error already exists, don't overwrite it with a potential lock release error.
	// The original error is usually more important.
	if *err != nil {
		// Log the lock release error if it occurs, but don't change the primary error.
		if lockErr := e.Release(ctx); lockErr != nil {
			// Log this, as it's a problem, but *err already holds a more primary error.
			log.L(ctx).Warnf("Warning: Failed to release billing lock, but original error already present: %v (release error: %v)\n", (*err), lockErr)
		}

		return
	}

	// If no error existed, attempt to release the lock and set *err if release fails.
	lockErr := e.Release(ctx)
	if lockErr != nil {
		*err = lockErr // Only set *err if it was nil and release failed.
	}
}

type redisBillingLockRepository struct {
	client *redislock.Client
}

func NewRedisBillingLockRepository(rdb *redis.Client) BillingLockRepository {
	return &redisBillingLockRepository{
		client: redislock.New(rdb),
	}
}

func (r *redisBillingLockRepository) ObtainLock(ctx context.Context, customerId string) (BillingLock, error) {
	log.L(ctx).Infof("Obtaining billing lock for %s", customerId)

	lock, err := r.client.Obtain(ctx, fmt.Sprintf("billing_locks:%s", customerId), billingLockTtl, nil)
	if err != nil {
		if errors.Is(err, redislock.ErrNotObtained) {
			log.L(ctx).Error("Failed to obtain billing lock")

			return nil, appErr.Errorf(appErr.TooManyRequestsError, "Too many requests")
		}

		return nil, tracerr.Wrap(err)
	}

	log.L(ctx).Infof("Successfully obtained billing lock for %s", customerId)

	return &authLock{lock: lock, customerId: customerId}, nil
}
