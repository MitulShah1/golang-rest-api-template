package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// Custom error messages
var CustomErrorMessages = map[string]string{
	"required": "The field {{.Field}} is required",
	"email":    "The field {{.Field}} must be a valid email",
	"min":      "The field {{.Field}} must be at least {{.Param}} characters long",
	"max":      "The field {{.Field}} must be at most {{.Param}} characters long",
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateStruct(s interface{}) []ValidationError {
	var validationErrors []ValidationError
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Validation failed for field: %s", err.Field())

			if customMsg, exists := CustomErrorMessages[err.Tag()]; exists {
				msg = customMsg
				msg = replacePlaceholders(msg, err)
			}

			validationErrors = append(validationErrors, ValidationError{
				Field:   err.Field(),
				Message: msg,
			})
		}
	}
	return validationErrors
}

func replacePlaceholders(msg string, err validator.FieldError) string {
	msg = strings.ReplaceAll(msg, "{{.Field}}", err.Field())
	msg = strings.ReplaceAll(msg, "{{.Param}}", err.Param())
	return msg
}
