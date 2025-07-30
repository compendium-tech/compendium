package myerror

import (
	"fmt"
	"net/http"
)

const (
	RequestValidationError   = 1
	ApplicationNotFoundError = 300
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
	case ApplicationNotFoundError:
		return http.StatusNotFound
	default:
		return http.StatusBadRequest
	}
}
