package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register function to get json tag as field name
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return fld.Name
		}
		return name
	})
}

func ValidateStruct(s interface{}) []ValidationError {
	var errors []ValidationError

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			element := ValidationError{
				Field: err.Field(),
				Error: generateValidationMessage(err),
			}
			errors = append(errors, element)
		}
	}
	return errors
}

func generateValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Should be at least %s characters", err.Param())
	case "max":
		return fmt.Sprintf("Should be at most %s characters", err.Param())
	case "oneof":
		return fmt.Sprintf("Should be one of: %s", err.Param())
	case "ltefield":
		return fmt.Sprintf("Should be less than %s", err.Param())
	case "gtefield":
		return fmt.Sprintf("Should be greater than %s", err.Param())
	case "eqfield":
		return fmt.Sprintf("Should be equal to %s", err.Param())
	default:
		return fmt.Sprintf("Failed validation on %s", err.Tag())
	}
}
