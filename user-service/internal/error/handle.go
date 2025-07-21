package error

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/seacite-tech/compendium/common/pkg/validate"
	log "github.com/sirupsen/logrus"
	"github.com/ztrue/tracerr"
)

func HandleAppErr(f func(c *gin.Context) error) func(c *gin.Context) {
	return func(c *gin.Context) {
		err := f(c)

		if err != nil {
			status, kind, message := func() (int, int, string) {
				if appErr, ok := err.(AppError); ok {
					return appErr.Kind().httpStatus(), int(appErr.Kind()), appErr.Message()
				} else if errs, ok := err.(validator.ValidationErrors); ok {
					return http.StatusBadRequest, 1, validate.BuildErrorMessage(errs)
				} else {
					if errs, ok := err.(tracerr.Error); ok {
						log.Printf("Cause of internal server error: %s\nStacktrace: %s", errs, errs.StackTrace())
					} else {
						log.Printf("Cause of internal server error: %s", err)
					}

					return http.StatusInternalServerError, 0, "Internal server error"
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
	case InvalidCredentialsError:
		return http.StatusUnauthorized
	case InvalidMfaOtpError:
		return http.StatusUnauthorized
	case InvalidSessionError:
		return http.StatusUnauthorized
	case TooManyRequestsError:
		return http.StatusTooManyRequests
	case UserNotFoundError:
		return http.StatusNotFound
	case EmailTakenError:
		return http.StatusConflict
	default:
		return http.StatusBadRequest
	}
}
