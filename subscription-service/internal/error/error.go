package error

import (
	"fmt"
)

const (
	// 0 is for internal server error.
	RequestValidationError AppErrorKind = 1

	UserNotFoundError    = 4
	TooManyRequestsError = 5

	// 100
	InvalidWebhookSignatureError      = 100
	LowPrioritySubscriptionLevelError = 101
	SubscriptionIsRequiredError       = 102
	AlreadySubscribedError            = 103
	PayerPermissionRequired           = 104
	InvalidSubscriptionInvitationCode = 105
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

func Errorf(kind AppErrorKind, format string, args ...any) errorWrap {
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
