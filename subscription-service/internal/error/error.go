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
	return fmt.Sprintf("application error, ty: %d, details: %v", e.ty, e.details)
}

func (e MyError) ErrorType() int {
	return e.ty
}

func (e MyError) ErrorDetails() any {
	return e.details
}

func (e MyError) HttpStatus() int {
	switch e.ty {
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
