package response

import (
	"net/http"

	"golang-rest-api-template/package/logger"
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
