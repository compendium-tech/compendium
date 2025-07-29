package error

import (
	"net/http"

	"github.com/compendium-tech/compendium/common/pkg/validate"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"github.com/ztrue/tracerr"
)

func Handle(f func(c *gin.Context) error) func(c *gin.Context) {
	return func(c *gin.Context) {
		err := f(c)

		if err != nil {
			status, kind, message := func() (int, AppErrorKind, string) {
				if appErr, ok := err.(AppError); ok {
					return appErr.Kind().httpStatus(), appErr.Kind(), appErr.Message()
				} else if errs, ok := err.(validator.ValidationErrors); ok {
					return http.StatusBadRequest, RequestValidationError, validate.BuildErrorMessage(errs)
				} else {
					if errs, ok := err.(tracerr.Error); ok {
						log.Printf("Cause of internal server error: %s\nStacktrace: %s", errs, errs.StackTrace())
					} else {
						log.Printf("Cause of internal server error: %s", err)
					}

					return http.StatusInternalServerError, InternalServerError, "Internal server error"
				}
			}()

			c.AbortWithStatusJSON(status, map[string]any{
				"errorMessage": message,
				"errorKind":    kind,
			})
		}
	}
}

func (k AppErrorKind) httpStatus() int {
	switch k {
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
