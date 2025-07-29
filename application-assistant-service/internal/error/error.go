package error

import (
	"fmt"
)

const (
	InternalServerError      AppErrorKind = 0
	RequestValidationError   AppErrorKind = 1
	ApplicationNotFoundError AppErrorKind = 300
	ActivityNotFoundError    AppErrorKind = 301
	HonorNotFoundError       AppErrorKind = 302
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
