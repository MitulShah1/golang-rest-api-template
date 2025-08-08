// Package response provides HTTP response utilities for the application.
// It includes standardized response formatting and error handling.
package response

import (
	"encoding/json"
	"net/http"

	"github.com/MitulShah1/golang-rest-api-template/package/logger"
)

func sendResponse(w http.ResponseWriter, status int, resp []byte, contentType string) {
	w.Header().Set(`Content-Type`, contentType)
	w.Header().Set(`X-Content-Type-Options`, `nosniff`)
	w.Header().Set(`Cache-Control`, `no-cache, no-store, must-revalidate`)
	w.Header().Set(`Pragma`, `no-cache`)
	w.Header().Set(`Expires`, `0`)

	w.WriteHeader(status)

	if _, err := w.Write(resp); err != nil {
		logger.NewLogger(logger.DefaultOptions()).Error(`failed to write response body`, `error`, err.Error())
	}
}

func SendXMLResponseRaw(w http.ResponseWriter, status int, resp []byte) {
	sendResponse(w, status, resp, `application/soap+xml; charset=utf-8`)
}

func SendResponseRaw(w http.ResponseWriter, status int, resp []byte) {
	sendResponse(w, status, resp, `application/json; charset=utf-8`)
}

// Success sends a successful JSON response
func Success(w http.ResponseWriter, status int, message string, data any) {
	response := map[string]any{
		"success": true,
		"message": message,
		"data":    data,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		Error(w, http.StatusInternalServerError, "Failed to marshal response")
		return
	}

	SendResponseRaw(w, status, jsonData)
}

// Error sends an error JSON response
func Error(w http.ResponseWriter, status int, message string) {
	response := map[string]any{
		"success": false,
		"message": message,
		"error":   true,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		// Fallback to plain text if JSON marshaling fails
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(status)
		if _, err := w.Write([]byte(message)); err != nil {
			logger.NewLogger(logger.DefaultOptions()).Error(`failed to write fallback response`, `error`, err.Error())
		}
		return
	}

	SendResponseRaw(w, status, jsonData)
}
