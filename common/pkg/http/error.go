package httputils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ztrue/tracerr"

	"github.com/compendium-tech/compendium/common/pkg/log"
	"github.com/compendium-tech/compendium/common/pkg/validate"
)

// ErrorHandler is a default solution to centralizing error handling for Gin handlers.
// See Handle for more details.
type ErrorHandler struct{}

type HandlerFuncWithError func(c *gin.Context) error

type CustomError interface {
	ErrorType() int
	ErrorDetails() any
	HttpStatus() int
}

// Handle intercepts errors, categorizes them using errors.As, and sends appropriate HTTP
// status codes and JSON responses to the client, while logging internal server errors.
//
// # Usage with Custom Errors
//
// To be handled gracefully by this middleware, your custom error types should implement CustomError.
//
// # Example
//
//	// Define application-specific error types as integers.
//	const (
//	  UserNotFoundError = 1001
//	  InvalidInputError = 1002
//	  ServiceUnavailableError = 2001
//	)
//
//	type APIError struct {
//	  Code       int    `json:"code"`        // Corresponds to Kind()
//	  Status     int    `json:"-"`           // Corresponds to HttpStatus(), typically hidden from JSON
//	  Message    string `json:"message"`     // A human-readable error message
//	  ExtraInfo  any    `json:"details,omitempty"` // Corresponds to Details()
//	}
//
//	// Error implements the standard error interface.
//	func (e APIError) Error() string {
//	  if e.wrappedErr != nil {
//	    return e.Message + ": " + e.wrappedErr.Error()
//	  }
//
//	  return e.Message
//	}
//
//	// Kind returns the application-specific error code.
//	func (e APIError) ErrorType() int { return e.Code }
//
//	// HttpStatus returns the HTTP status code for the response.
//	func (e APIError) HttpStatus() int { return e.Status }
//
//	// Details returns any extra information about the error.
//	func (e APIError) ErrorDetails() any { return e.ExtraInfo }
//
//	// Helper function to create a new APIError (optional, but good practice).
//	func NewAPIError(ty, status int, msg string, details any, err error) APIError {
//	  return APIError{
//	    Code:       ty,
//	    Status:     status,
//	    Message:    msg,
//	    ExtraInfo:  details,
//	    wrappedErr: err,
//	  }
//	}
//
//	func GetUserByID(c *gin.Context) error {
//	  userID := c.Param("id")
//	    if userID == "invalid" {
//	      return NewAPIError(
//	        InvalidInputError,
//	        http.StatusBadRequest,
//	        "Invalid user ID format.",
//	        map[string]string{"input": userID, "reason": "non-numeric"}
//	      )
//	  } else if userID == "nonexistent" {
//	      return NewAPIError(
//	         UserNotFoundError,
//	         http.StatusNotFound,
//	         "User with the provided ID does not exist.",
//	         map[string]string{"requested_id": userID}
//	    )
//	  }
//
//	  c.JSON(http.StatusOK, gin.H{"message": "User found", "id": userID})
//	  return nil
//	}
//
//	func main() {
//	  var e httputils.ErrorHandler
//	  router := gin.Default()
//	  router.GET("/users/:id", e.Handle(GetUserByID))
//	  router.Run(":8080")
//	}
func (h ErrorHandler) Handle(f HandlerFuncWithError) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := f(c)

		if err != nil {
			var customErr CustomError
			var validationErrs validator.ValidationErrors

			if errors.As(err, &customErr) {
				h.handleCustomError(c, customErr)
			} else if errors.As(err, &validationErrs) {
				h.handleValidationErrors(c, validationErrs)
			} else {
				h.handleISE(c, err)
			}
		}
	}
}

func (h ErrorHandler) handleCustomError(c *gin.Context, err CustomError) {
	h.abort(c, err.HttpStatus(), err.ErrorType(), err.ErrorDetails())
}

func (h ErrorHandler) handleValidationErrors(c *gin.Context, errs validator.ValidationErrors) {
	var validationMessages []string
	for _, err := range errs {
		validationMessages = append(validationMessages, validate.BuildErrorMessage(err))
	}

	h.abort(c, http.StatusBadRequest, 1, map[string][]string{
		"reasons": validationMessages,
	})
}

func (h ErrorHandler) handleISE(c *gin.Context, err error) {
	var trcErr tracerr.Error

	if errors.As(err, &trcErr) {
		log.L(c.Request.Context()).Printf("Cause of internal server error: %s\nStacktrace: %s", trcErr, trcErr.StackTrace())
	} else {
		log.L(c.Request.Context()).Printf("Cause of internal server error: %s", err)
	}

	h.abort(c, http.StatusInternalServerError, 0, nil)
}

func (h ErrorHandler) abort(c *gin.Context, httpStatus int, ty int, details any) {
	c.AbortWithStatusJSON(httpStatus, map[string]any{
		"errorDetails": details,
		"errorType":    ty,
	})
}
