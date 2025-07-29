package error

import (
	"fmt"
)

const (
	InternalServerError               AppErrorKind = 0
	RequestValidationError            AppErrorKind = 1
	InvalidCredentialsError           AppErrorKind = 2
	EmailTakenError                   AppErrorKind = 3
	UserNotFoundError                 AppErrorKind = 4
	TooManyRequestsError              AppErrorKind = 5
	MfaNotRequestedError              AppErrorKind = 6
	InvalidMfaOtpError                AppErrorKind = 7
	InvalidSessionError               AppErrorKind = 8
	SessionNotFoundError              AppErrorKind = 9
	FailedToRemoveCurrentSessionError AppErrorKind = 10
)

type AppErrorKind int

type AppError struct {
	kind    AppErrorKind
	details any
}

func New(kind AppErrorKind) AppError {
	return NewWithDetails(kind, nil)
}

func NewWithDetails(kind AppErrorKind, details any) AppError {
	return AppError{
		kind:    kind,
		details: details,
	}
}

func NewWithReason(kind AppErrorKind, reason string) AppError {
	return NewWithDetails(kind, map[string]any{"reason": reason})
}

func (e AppError) Error() string {
	return fmt.Sprintf("application error, kind: %d, details: %v", e.kind, e.details)
}

func (e AppError) Kind() AppErrorKind {
	return e.kind
}
