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

	err := Validate.RegisterValidation("password", validatePassword)
	if err != nil {
		return
	}

	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" || name == "" {
			return fld.Name
		}

		return name
	})
}

func validatePassword(fl validator.FieldLevel) bool {
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

func BuildErrorMessage(err validator.FieldError) string {
	buff := bytes.NewBufferString("")
	buff.WriteString(fmt.Sprintf("field <%s> doesn't follow rule <%s>", err.Field(), err.Tag()))

	return buff.String()
}
