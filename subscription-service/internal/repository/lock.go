package repository

import "context"

// Prevents race conditions in billing webhook processing.
type BillingLockRepository interface {
	ObtainLock(ctx context.Context, customerId string) (BillingLock, error)
}

type BillingLock interface {
	Release(ctx context.Context) error
	ReleaseAndHandleErr(ctx context.Context, err *error)
}
