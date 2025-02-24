package middleware

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		shouldProceed  bool
	}{
		{
			name:           "Valid credentials",
			authHeader:     "Basic " + base64.StdEncoding.EncodeToString([]byte("admin:password")),
			expectedStatus: http.StatusOK,
			shouldProceed:  true,
		},
		{
			name:           "Missing auth header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			shouldProceed:  false,
		},
		{
			name:           "Invalid auth format",
			authHeader:     "Bearer token123",
			expectedStatus: http.StatusUnauthorized,
			shouldProceed:  false,
		},
		{
			name:           "Invalid base64 encoding",
			authHeader:     "Basic invalid-base64",
			expectedStatus: http.StatusUnauthorized,
			shouldProceed:  false,
		},
		{
			name:           "Invalid credentials format",
			authHeader:     "Basic " + base64.StdEncoding.EncodeToString([]byte("invalid")),
			expectedStatus: http.StatusUnauthorized,
			shouldProceed:  false,
		},
		{
			name:           "Wrong username",
			authHeader:     "Basic " + base64.StdEncoding.EncodeToString([]byte("wronguser:password")),
			expectedStatus: http.StatusUnauthorized,
			shouldProceed:  false,
		},
		{
			name:           "Wrong password",
			authHeader:     "Basic " + base64.StdEncoding.EncodeToString([]byte("admin:wrongpassword")),
			expectedStatus: http.StatusUnauthorized,
			shouldProceed:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextCalled := false
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})

			handler := AuthMiddleware(nextHandler)
			req := httptest.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.shouldProceed, nextCalled)
		})
	}
}
