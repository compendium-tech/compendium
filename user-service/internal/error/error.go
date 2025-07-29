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

type AppError interface {
	Message() string
	Kind() AppErrorKind
}

type errorWrap struct {
	message string
	kind    AppErrorKind
}

func New(kind AppErrorKind, format string, args ...any) errorWrap {
	return errorWrap{
		message: fmt.Errorf(format, args...).Error(),
		kind:    kind,
	}
}

func (e errorWrap) Error() string {
	return e.message
}

func (e errorWrap) Message() string {
	return e.message
}

func (e errorWrap) Kind() AppErrorKind {
	return e.kind
}
