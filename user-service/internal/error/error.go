package myerror

import (
	"fmt"
	"net/http"
)

const (
	RequestValidationError            = 1
	InvalidCredentialsError           = 2
	EmailTakenError                   = 3
	UserNotFoundError                 = 4
	TooManyRequestsError              = 5
	MfaNotRequestedError              = 6
	InvalidMfaOtpError                = 7
	InvalidSessionError               = 8
	SessionNotFoundError              = 9
	FailedToRemoveCurrentSessionError = 10
)

type MyError struct {
	ty      int
	details any
}

func New(ty int) MyError {
	return NewWithDetails(ty, nil)
}

func NewWithDetails(ty int, details any) MyError {
	return MyError{
		ty:      ty,
		details: details,
	}
}

func NewWithReason(ty int, reason string) MyError {
	return NewWithDetails(ty, map[string]any{"reason": reason})
}

func (e MyError) Error() string {
	return fmt.Sprintf("application error, type: %d, details: %v", e.ty, e.details)
}

func (e MyError) ErrorType() int {
	return e.ty
}

func (e MyError) ErrorDetails() any {
	return e.details
}

func (e MyError) HttpStatus() int {
	switch e.ty {
	case InvalidCredentialsError:
		return http.StatusUnauthorized
	case InvalidMfaOtpError:
		return http.StatusUnauthorized
	case InvalidSessionError:
		return http.StatusUnauthorized
	case TooManyRequestsError:
		return http.StatusTooManyRequests
	case UserNotFoundError:
		return http.StatusNotFound
	case EmailTakenError:
		return http.StatusConflict
	case SessionNotFoundError:
		return http.StatusNotFound
	case FailedToRemoveCurrentSessionError:
		return http.StatusForbidden
	default:
		return http.StatusBadRequest
	}
}
