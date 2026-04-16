package shared

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateStruct validates a struct based on `validate` tags.
func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	// Return the first validation error found
	validationErrors := err.(validator.ValidationErrors)
	if len(validationErrors) > 0 {
		fieldErr := validationErrors[0]
		return errors.New(fieldErr.Field())
	}
	return errors.New("validation failed")
}

func init() {
	// Register a custom tag name function to handle JSON names in errors if needed
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}
