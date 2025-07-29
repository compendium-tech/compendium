package validate

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func InitValidator() {
	Validate = validator.New()

	Validate.RegisterValidation("password", ValidatePassword)
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" || name == "" {
			return fld.Name
		}

		return name
	})
}

func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return false
	}

	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return false
	}

	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return false
	}

	if !regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]`).MatchString(password) {
		return false
	}

	return true
}

func BuildErrorMessage(errs validator.ValidationErrors) string {
	buff := bytes.NewBufferString("")

	for i := 0; i < len(errs); i++ {
		buff.WriteString(fmt.Sprintf("field <%s> doesn't follow rule <%s>", errs[i].Field(), errs[i].Tag()))
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}

func BuildValidationErrorMessage(err validator.FieldError) string {
	buff := bytes.NewBufferString("")

	buff.WriteString(fmt.Sprintf("field <%s> doesn't follow rule <%s>", err.Field(), err.Tag()))

	return buff.String()
}
