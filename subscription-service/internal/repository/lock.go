package repository

import (
	"context"

	"github.com/google/uuid"
)

// Prevents race conditions in billing webhook processing.
type BillingLockRepository interface {
	ObtainLock(ctx context.Context, userID uuid.UUID) (BillingLock, error)
}

type BillingLock interface {
	Release(ctx context.Context) error
	ReleaseAndHandleErr(ctx context.Context, err *error)
}
