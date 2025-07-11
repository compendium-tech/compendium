package repository

import "context"

// Prevents race conditions in authentication endpoints.
type EmailLockRepository interface {
	ObtainEmailLock(ctx context.Context, email string) (EmailLock, error)
}

type EmailLock interface {
	Release(ctx context.Context) error
}
