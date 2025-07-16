package repository

import "context"

// Prevents race conditions in authentication endpoints.
type AuthEmailLockRepository interface {
	ObtainLock(ctx context.Context, email string) (AuthEmailLock, error)
}

type AuthEmailLock interface {
	Release(ctx context.Context) error
	ReleaseAndHandleErr(ctx context.Context, err *error)
}
