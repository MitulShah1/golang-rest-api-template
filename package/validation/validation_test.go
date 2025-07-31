package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateStruct(t *testing.T) {
	tests := []struct {
		name          string
		input         any
		expectedErrs  int
		expectedField string
	}{
		{
			name: "Valid struct",
			input: struct {
				Name string `validate:"required"`
			}{
				Name: "test",
			},
			expectedErrs: 0,
		},
		{
			name: "Invalid struct - required field missing",
			input: struct {
				Name string `validate:"required"`
			}{
				Name: "",
			},
			expectedErrs:  1,
			expectedField: "Name",
		},
		{
			name: "Multiple validation errors",
			input: struct {
				Name  string `validate:"required"`
				Email string `validate:"required,email"`
			}{
				Name:  "",
				Email: "invalid-email",
			},
			expectedErrs: 2,
		},
		{
			name: "Custom error message",
			input: struct {
				Age int `validate:"min=18"`
			}{
				Age: 16,
			},
			expectedErrs:  1,
			expectedField: "Age",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateStruct(tt.input)

			assert.Len(t, errors, tt.expectedErrs)

			if tt.expectedErrs > 0 && tt.expectedField != "" {
				assert.Equal(t, tt.expectedField, errors[0].Field)
				assert.NotEmpty(t, errors[0].Message)
			}
		})
	}
}
