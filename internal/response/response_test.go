package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendResponse(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		resp        []byte
		contentType string
		wantHeaders map[string]string
	}{
		{
			name:        "Success with JSON response",
			status:      http.StatusOK,
			resp:        []byte(`{"message":"success"}`),
			contentType: "application/json",
			wantHeaders: map[string]string{
				"Content-Type":           "application/json",
				"X-Content-Type-Options": "nosniff",
				"Cache-Control":          "no-cache, no-store, must-revalidate",
				"Pragma":                 "no-cache",
				"Expires":                "0",
			},
		},
		{
			name:        "Error response with text content",
			status:      http.StatusBadRequest,
			resp:        []byte("error message"),
			contentType: "text/plain",
			wantHeaders: map[string]string{
				"Content-Type":           "text/plain",
				"X-Content-Type-Options": "nosniff",
				"Cache-Control":          "no-cache, no-store, must-revalidate",
				"Pragma":                 "no-cache",
				"Expires":                "0",
			},
		},
		{
			name:        "Empty response",
			status:      http.StatusNoContent,
			resp:        []byte(nil),
			contentType: "application/json",
			wantHeaders: map[string]string{
				"Content-Type":           "application/json",
				"X-Content-Type-Options": "nosniff",
				"Cache-Control":          "no-cache, no-store, must-revalidate",
				"Pragma":                 "no-cache",
				"Expires":                "0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			sendResponse(w, tt.status, tt.resp, tt.contentType)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.status, result.StatusCode)

			for key, value := range tt.wantHeaders {
				assert.Equal(t, value, result.Header.Get(key))
			}

			body := w.Body.Bytes()
			assert.Equal(t, tt.resp, body)
		})
	}
}
