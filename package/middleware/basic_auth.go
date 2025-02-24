package middleware

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// AuthMiddleware is a simple Basic Authentication middleware
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Basic Auth Format: "Basic base64(username:password)"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Basic" {
			http.Error(w, "Invalid authentication format", http.StatusUnauthorized)
			return
		}

		// Decode Base64 (username:password)
		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			http.Error(w, "Invalid base64 encoding", http.StatusUnauthorized)
			return
		}

		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 {
			http.Error(w, "Invalid credentials format", http.StatusUnauthorized)
			return
		}

		username, password := credentials[0], credentials[1]

		// Validate credentials (Replace with DB check or Config-based auth)
		if username != "admin" || password != "password" { // Replace with actual auth logic
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Proceed to next handler
		next.ServeHTTP(w, r)
	})
}
