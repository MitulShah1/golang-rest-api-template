// Package validation provides data validation utilities for the application.
// It includes struct validation, custom error messages, and validation rules.
package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// CustomErrorMessages contains custom error messages for validation
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

func ValidateStruct(s any) []ValidationError {
	var validationErrors []ValidationError
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			msg := "Validation failed for field: " + err.Field()

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
