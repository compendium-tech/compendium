package httputils

import (
	"github.com/compendium-tech/compendium/common/pkg/validate"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func MustBindWith[T any](c *gin.Context, binding binding.Binding, validateStruct bool) T {
	var t T

	if err := c.MustBindWith(&t, binding); err != nil {
		panic(err)
	}

	if !validateStruct {
		return t
	}

	if err := validate.Validate.Struct(t); err != nil {
		panic(err)
	}

	return t
}
