package error

import (
	"fmt"
)

const (
	InternalServerError               AppErrorKind = 0
	RequestValidationError            AppErrorKind = 1
	UserNotFoundError                 AppErrorKind = 4
	TooManyRequestsError              AppErrorKind = 5
	InvalidWebhookSignatureError      AppErrorKind = 100
	LowPrioritySubscriptionLevelError AppErrorKind = 101
	SubscriptionIsRequiredError       AppErrorKind = 102
	AlreadySubscribedError            AppErrorKind = 103
	PayerPermissionRequired           AppErrorKind = 104
	InvalidSubscriptionInvitationCode AppErrorKind = 105
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
