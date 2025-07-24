package repository

import "context"

// AuthLockRepository provides an interface for obtaining and managing authentication locks.
// These locks are used to prevent race conditions in authentication-related endpoints,
// ensuring that sensitive operations like sign-up, sign-in, and password resets
// are processed sequentially for a given user.
type AuthLockRepository interface {
	// ObtainLock attempts to acquire an exclusive lock for a given email address.
	// It returns an AuthLock instance if the lock is successfully obtained,
	// or an error if the lock cannot be acquired (e.g., if it's already held).
	ObtainLock(ctx context.Context, email string) (AuthLock, error)
}

type AuthLock interface {
	// Release manually releases the lock. Returns error if lock is no longer held.
	Release(ctx context.Context) error

	// ReleaseAndHandleErr releases the acquired lock and also handles a potential
	// error from the deferred function. This is typically used in defer statements
	// to ensure the lock is always released, even if an error occurs during the
	// function's execution. Example:
	//
	//  func foo() (finalErr error) {
	//    lock, err := s.authLockRepository.ObtainLock(ctx, email)
	//      if err != nil {
	//      return err
	//    }
	//
	//    defer lock.ReleaseAndHandleErr(ctx, &finalErr)
	//    ...
	//  }
	//
	// The [err] parameter should be a pointer to the
	// named return error variable of the function.
	ReleaseAndHandleErr(ctx context.Context, err *error)
}
