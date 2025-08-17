package repository

import "context"

// AuthLockRepository provides an interface for obtaining and managing authentication locks.
// These locks are used to prevent race conditions in authentication-related endpoints,
// ensuring that sensitive operations like sign-up, sign-in, and password resets
// are processed sequentially for a given user.
//
// ObtainLock attempts to acquire an exclusive lock for a given email address.
// It returns an AuthLock instance if the lock is successfully obtained,
// or an error if the lock cannot be acquired (e.g., if it's already held).
type AuthLockRepository interface {
	ObtainLock(ctx context.Context, email string) AuthLock
}

// AuthLock represents an acquired authentication lock.
// See [AuthLockRepository] for more details.
//
// Release manually releases the lock. Returns error if lock is no longer held.
type AuthLock interface {
	Release(ctx context.Context)
}
