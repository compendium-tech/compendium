package error

import (
	"net/http"

	"github.com/compendium-tech/compendium/common/pkg/log"
	"github.com/compendium-tech/compendium/common/pkg/validate"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ztrue/tracerr"
)

func Handle(f func(c *gin.Context) error) func(c *gin.Context) {
	return func(c *gin.Context) {
		err := f(c)

		if err != nil {
			status, kind, details := func() (int, AppErrorKind, any) {
				if appErr, ok := err.(AppError); ok {
					return appErr.kind.httpStatus(), appErr.kind, appErr.details
				} else if errs, ok := err.(validator.ValidationErrors); ok {
					var validationErrors []string
					for _, err := range errs {
						validationErrors = append(validationErrors, validate.BuildValidationErrorMessage(err))
					}

					return http.StatusBadRequest, RequestValidationError, map[string][]string{
						"validationErrors": validationErrors,
					}
				} else {
					if errs, ok := err.(tracerr.Error); ok {
						log.L(c.Request.Context()).Printf("Cause of internal server error: %s\nStacktrace: %s", errs, errs.StackTrace())
					} else {
						log.L(c.Request.Context()).Printf("Cause of internal server error: %s", err)
					}

					return http.StatusInternalServerError, InternalServerError, nil
				}
			}()

			c.AbortWithStatusJSON(status, map[string]any{
				"errorDetails": details,
				"errorKind":    kind,
			})
		}
	}
}

func (k AppErrorKind) httpStatus() int {
	switch k {
	case ApplicationNotFoundError:
		return http.StatusNotFound
	default:
		return http.StatusBadRequest
	}
}
