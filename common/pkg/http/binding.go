package httputils

import (
	"github.com/compendium-tech/compendium/common/pkg/validate"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type mustBindWith[T any] struct {
	t T
}

// MustBindWith binds the body of the request to the given type T with the given gin binding.
//
// If the binding fails, it will panic.
// It will return a value of type mustBindWith[T].
//
// You can then call .Validated() to get the validated value of type T, or .NotValidated() to get the non-validated value of type T.
func MustBindWith[T any](c *gin.Context, binding binding.Binding) mustBindWith[T] {
	var t T

	err := c.MustBindWith(&t, binding)
	if err != nil {
		panic(err)
	}

	return mustBindWith[T]{t: t}
}

// Validated returns the validated value of type T.
//
// It will panic if the validation fails.
func (m mustBindWith[T]) Validated() T {
	err := validate.Validate.Struct(m.t)
	if err != nil {
		panic(err)
	}

	return m.t
}

// NotValidated returns the non-validated value of type T.
//
// It does not do any validation on the value, and will not panic.
func (m mustBindWith[T]) NotValidated() T {
	return m.t
}
