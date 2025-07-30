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
	kind    int
	details any
}

func New(kind int) MyError {
	return NewWithDetails(kind, nil)
}

func NewWithDetails(kind int, details any) MyError {
	return MyError{
		kind:    kind,
		details: details,
	}
}

func NewWithReason(kind int, reason string) MyError {
	return NewWithDetails(kind, map[string]any{"reason": reason})
}

func (e MyError) Error() string {
	return fmt.Sprintf("application error, kind: %d, details: %v", e.kind, e.details)
}

func (e MyError) Kind() int {
	return e.kind
}

func (e MyError) Details() any {
	return e.details
}

func (e MyError) HttpStatus() int {
	switch e.kind {
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
