// Package health provides HTTP handlers for health check operations.
// It includes endpoints for monitoring application health and status.
package health

import (
	"net/http"

	"github.com/MitulShah1/golang-rest-api-template/internal/response"
	"github.com/MitulShah1/golang-rest-api-template/package/logger"
	"github.com/gorilla/mux"
)

// HealthCheckPath is the path for health check endpoint
const HealthCheckPath = "/health-check"

type HealthAPI struct {
	logger *logger.Logger
}

func NewHealthAPI(logger *logger.Logger) *HealthAPI {
	return &HealthAPI{
		logger: logger,
	}
}

// RegisterHandlers registers the health check API handler on the provided router.
// The health check API handler responds with a 200 OK status when the /health-check
// endpoint is accessed, which can be used for health checks or liveness probes.

func (api *HealthAPI) RegisterHandlers(router *mux.Router) {
	router.HandleFunc(HealthCheckPath, api.HealthCheckAPIHandler).Methods("GET")
}

// HealthCheckAPIHandler godoc
// @Summary ping example
// @Description do ping
// @Tags HealthCheck
// @Accept json
// @Produce json
// @Success 200 {string} string "OK"
// @Router /api/health-check [get]
// HealthCheckAPIHandler handles HTTP requests for health check endpoints.
// It returns a simple OK response to indicate the service is running.
func (api *HealthAPI) HealthCheckAPIHandler(w http.ResponseWriter, r *http.Request) {
	response.SendResponseRaw(w, http.StatusOK, []byte("OK"))
}
