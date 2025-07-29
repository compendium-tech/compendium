package error

import "fmt"

const (
	InternalServerError            AppErrorKind = 0
	RequestValidationError         AppErrorKind = 1
	ApplicationNotFoundError       AppErrorKind = 300
	ActivityNotFoundError          AppErrorKind = 301
	HonorNotFoundError             AppErrorKind = 302
	EssayNotFoundError             AppErrorKind = 303
	SupplementalEssayNotFoundError AppErrorKind = 304
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
