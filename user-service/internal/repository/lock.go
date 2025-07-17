package repository

import "context"

// Prevents race conditions in authentication endpoints.
type AuthLockRepository interface {
	ObtainLock(ctx context.Context, email string) (AuthLock, error)
}

type AuthLock interface {
	Release(ctx context.Context) error
	ReleaseAndHandleErr(ctx context.Context, err *error)
}
