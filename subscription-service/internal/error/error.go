package myerror

import (
	"fmt"
	"net/http"
)

const (
	RequestValidationError            = 1
	TooManyRequestsError              = 5
	InvalidWebhookSignatureError      = 100
	SubscriptionIsRequiredError       = 102
	AlreadySubscribedError            = 103
	PayerPermissionRequired           = 104
	InvalidSubscriptionInvitationCode = 105
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
	case SubscriptionIsRequiredError:
		return http.StatusPaymentRequired
	case PayerPermissionRequired:
		return http.StatusForbidden
	case InvalidSubscriptionInvitationCode:
		return http.StatusForbidden
	default:
		return http.StatusBadRequest
	}
}
