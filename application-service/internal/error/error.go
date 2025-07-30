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
	case ApplicationNotFoundError:
		return http.StatusNotFound
	default:
		return http.StatusBadRequest
	}
}
