// Package handlers provides HTTP server setup and routing functionality.
// It includes server initialization, middleware configuration, and route registration.
package handlers

import (
	"context"
	"net"
	"net/http"

	_ "github.com/MitulShah1/golang-rest-api-template/docs"
	catApi "github.com/MitulShah1/golang-rest-api-template/internal/handlers/category"
	"github.com/MitulShah1/golang-rest-api-template/internal/handlers/health"
	prodApi "github.com/MitulShah1/golang-rest-api-template/internal/handlers/product"
	"github.com/MitulShah1/golang-rest-api-template/internal/repository"
	"github.com/MitulShah1/golang-rest-api-template/internal/services/category"
	"github.com/MitulShah1/golang-rest-api-template/internal/services/product"
	"github.com/MitulShah1/golang-rest-api-template/package/cache"
	"github.com/MitulShah1/golang-rest-api-template/package/database"
	"github.com/MitulShah1/golang-rest-api-template/package/logger"
	"github.com/MitulShah1/golang-rest-api-template/package/middleware"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Server struct {
	httpList net.Listener
	httpSrvr *http.Server
	logger   *logger.Logger
}

func NewServer(address string, logger *logger.Logger, db *database.Database, cache *cache.Cache, tm *middleware.TelemetryConfig) (*Server, error) {
	// Create a new router
	router := mux.NewRouter()

	// swagger docs
	// Serve Swagger UI
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Prometheus metrics
	prometheusMiddleware := middleware.NewPrometheusMiddleware(middleware.Config{
		DoNotUseRequestPathFor404: true,
	})

	mw := func(handler http.Handler) http.Handler {
		return prometheusMiddleware.Middleware(
			tm.OpenTelemetryMiddleware(handler),
		)
	}

	router.Use(mw)

	router.Handle("/metrics", promhttp.Handler())

	r := router.PathPrefix("/api").Subrouter()

	// health check API
	healthAPI := health.NewHealthAPI(logger)
	healthAPI.RegisterHandlers(r)

	// cache health check API
	cacheHealthAPI := health.NewCacheHealthAPI(logger, cache)
	cacheHealthAPI.RegisterHandlers(r)

	// Create versioned subrouter (e.g., /v1)
	apiV1 := r.PathPrefix("/v1").Subrouter()

	// Register all middlewares
	middlewares := func(handler http.Handler) http.Handler {
		return middleware.CorsMiddleware(
			middleware.AuthMiddleware(handler),
		)
	}

	// Protected routes uses authentication middleware
	apiV1.Use(middlewares)

	// initialize repository
	repo := repository.NewDBRepository(db)

	// initialize product service with cache
	productService := product.NewProductService(repo, logger, cache)

	// initialize product handler
	productHandler := prodApi.NewProductAPI(logger, productService)

	// Register product handlers
	productHandler.RegisterHandlers(apiV1)

	// initialize category service with cache
	categoryService := category.NewCategoryService(repo, logger, cache)

	// initialize category handler
	categoryHandler := catApi.NewCategoryAPI(logger, categoryService)

	/// Register category handlers
	categoryHandler.RegisterHandlers(apiV1)

	httpLis, err := net.Listen(`tcp`, address)
	if err != nil {
		return nil, err
	}

	return &Server{
		httpList: httpLis,
		httpSrvr: &http.Server{
			Addr:    address,
			Handler: router,
		},
		logger: logger,
	}, nil
}

// ListenAndServe starts the HTTP server and listens for incoming requests.
// It returns an error if the server fails to start.
func (srv *Server) ListenAndServe() error {
	return srv.httpSrvr.Serve(srv.httpList)
}

// ServerDown gracefully shuts down the HTTP server. It takes a context.Context
// argument to allow for cancellation of the shutdown process.
func (srv *Server) ServerDown(ctx context.Context) error {
	return srv.httpSrvr.Shutdown(ctx)
}
