package errorutils

import (
	"context"

	"github.com/ztrue/tracerr"
)

type DeferFunc func() error
type DeferFuncWithContext func(ctx context.Context) error

// DeferTry executes a deferred function and if it returns an error and the caller
// doesn't, the caller's error is set to error of deferred function.
//
//	func locked() (finalErr error) {
//	  lock, err := Acquire()
//	  if err != nil { return err }
//
//	  defer errorutils.DeferTry(&finalErr, lock.Release)
func DeferTry(e *error, f DeferFunc) {
	err := f()

	if err != nil {
		if *e == nil {
			*e = tracerr.Wrap(err)
		}
	}
}

// DeferTryWithContext is a convenience wrapper for DeferTry, allowing deferred functions
// that require a context.Context. It passes the provided context to the deferred function before
// handling any returned error.
//
//	func locked(ctx context.Context) (finalErr error) {
//	  lock, err := Acquire(ctx)
//	  if err != nil { return err }
//
//	  defer errorutils.DeferTryWithContext(ctx, &finalErr, lock.Release)
func DeferTryWithContext(ctx context.Context, e *error, f DeferFuncWithContext) {
	DeferTry(e, func() error { return f(ctx) })
}
