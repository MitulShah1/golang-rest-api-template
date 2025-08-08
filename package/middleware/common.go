// Package middleware provides HTTP middleware components for the application.
// It includes authentication, CORS, logging, and telemetry middleware.
package middleware

import "net/http"

// responseWriterDelegator is a wrapper around http.ResponseWriter that captures
// the status code and other response information for middleware processing.
type responseWriterDelegator struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

func (r *responseWriterDelegator) WriteHeader(code int) {
	r.status = code
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseWriterDelegator) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	n, err := r.ResponseWriter.Write(b)
	r.written += int64(n)
	return n, err
}
