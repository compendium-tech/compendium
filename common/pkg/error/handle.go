package errorutils

import (
	"errors"
	"github.com/compendium-tech/compendium/common/pkg/log"
	"github.com/compendium-tech/compendium/common/pkg/validate"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ztrue/tracerr"
	"net/http"
)

type HandlerFuncWithError func(c *gin.Context) error

// Handle is a solution to centralizing error handling for handler functions.
//
// It intercepts errors, categorizes them using `errors.As`, and sends appropriate HTTP status codes
// and JSON responses to the client, while logging internal server errors.
//
// # Usage with Custom Errors
//
// To be handled gracefully by this middleware, your custom error types should implement the following interface (or a struct with methods matching it):
//
//	type HandledError interface {
//	  Kind() int        // Returns an application-specific integer code for the error category.
//	  HttpStatus() int  // Returns the appropriate HTTP status code (e.g., http.StatusBadRequest, http.StatusNotFound).
//	  Details() any     // Returns any additional, structured data relevant to the error (e.g., validation messages, specific IDs).
//	}
//
// # Example
//
//	// Define application-specific error kinds as integers.
//	const (
//	  UserNotFoundErrKind = 1001
//	  InvalidInputErrKind = 1002
//	  ServiceUnavailableKind = 2001
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
//	func (e APIError) Kind() int { return e.Code }
//
//	// HttpStatus returns the HTTP status code for the response.
//	func (e APIError) HttpStatus() int { return e.Status }
//
//	// Details returns any extra information about the error.
//	func (e APIError) Details() any { return e.ExtraInfo }
//
//	// Helper function to create a new APIError (optional, but good practice).
//	func NewAPIError(kind, status int, msg string, details any, err error) APIError {
//	  return APIError{
//	    Code:       kind,
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
//	        InvalidInputErrKind,
//	        http.StatusBadRequest,
//	        "Invalid user ID format.",
//	        map[string]string{"input": userID, "reason": "non-numeric"}
//	      )
//	  } else if userID == "nonexistent" {
//	      return NewAPIError(
//	         UserNotFoundErrKind,
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
//	router := gin.Default()
//	router.GET("/users/:id", errorutils.Handle(GetUserByID))
//	router.Run(":8080")
func Handle(f HandlerFuncWithError) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := f(c)

		if err != nil {
			status, kind, details := func() (int, int, any) {
				var myErr interface {
					Kind() int
					HttpStatus() int
					Details() any
				}

				var validationErrs validator.ValidationErrors
				var trcErr tracerr.Error

				if errors.As(err, &myErr) {
					return myErr.HttpStatus(), myErr.Kind(), myErr.Details()
				} else if errors.As(err, &validationErrs) {
					var validationMessages []string
					for _, ve := range validationErrs {
						validationMessages = append(validationMessages, validate.BuildErrorMessage(ve))
					}

					return http.StatusBadRequest, 1, map[string][]string{
						"reasons": validationMessages,
					}
				} else {
					if errors.As(err, &trcErr) {
						log.L(c.Request.Context()).Printf("Cause of internal server error: %s\nStacktrace: %s", trcErr, trcErr.StackTrace())
					} else {
						log.L(c.Request.Context()).Printf("Cause of internal server error: %s", err)
					}

					return http.StatusInternalServerError, 0, nil
				}
			}()

			c.AbortWithStatusJSON(status, map[string]any{
				"errorDetails": details,
				"errorKind":    kind,
			})
		}
	}
}
