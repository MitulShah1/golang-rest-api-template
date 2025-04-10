package health

import (
	"net/http"

	"github.com/MitulShah1/golang-rest-api-template/internal/response"
	"github.com/MitulShah1/golang-rest-api-template/package/logger"

	"github.com/gorilla/mux"
)

const HEALTH_CHECK_PATH = "/health-check"

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
	router.HandleFunc(HEALTH_CHECK_PATH, api.HealthCheckApiHandler).Methods("GET")
}

// HealthCheckApiHandler is an HTTP handler that responds with a 200 OK status
// when the /health-check endpoint is accessed. This can be used for health
// checks or liveness probes.

// Healthcheck godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags HealthCheck
// @Accept json
// @Produce json
// @Success 200 {string} ok
// @Router /health-check [get]
func (api *HealthAPI) HealthCheckApiHandler(w http.ResponseWriter, r *http.Request) {
	response.SendResponseRaw(w, http.StatusOK, []byte("OK"))
}
