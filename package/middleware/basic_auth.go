package middleware

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/MitulShah1/golang-rest-api-template/internal/response"
)

// AuthMiddleware is a simple Basic Authentication middleware
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			sendResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Basic Auth Format: "Basic base64(username:password)"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Basic" {
			sendResponse(w, http.StatusUnauthorized, "Invalid authentication format")
			return
		}

		// Decode Base64 (username:password)
		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			sendResponse(w, http.StatusUnauthorized, "Invalid base64 encoding")
			return
		}

		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 {
			sendResponse(w, http.StatusUnauthorized, "Invalid credentials format")
			return
		}

		username, password := credentials[0], credentials[1]

		// Validate credentials (Replace with DB check or Config-based auth)
		if username != "admin" || password != "password" { // Replace with actual auth logic
			sendResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Proceed to next handler
		next.ServeHTTP(w, r)
	})
}

func sendResponse(w http.ResponseWriter, code int, message string) {
	type StandardResponse struct {
		IsSuccess bool   `json:"issuccess"`
		Message   string `json:"message"`
	}

	res := StandardResponse{
		IsSuccess: false,
		Message:   message,
	}

	// Send the response
	resp, err := json.Marshal(res)
	if err != nil {
		response.SendResponseRaw(w, http.StatusInternalServerError, nil)
		return
	}

	response.SendResponseRaw(w, code, resp)
}
