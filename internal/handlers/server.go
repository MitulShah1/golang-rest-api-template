package handlers

import (
	"context"
	"golang-rest-api-template/internal/handlers/health"
	prodApi "golang-rest-api-template/internal/handlers/product"
	"golang-rest-api-template/internal/repository"
	"golang-rest-api-template/internal/services/category"
	"golang-rest-api-template/internal/services/product"
	"golang-rest-api-template/package/database"
	"golang-rest-api-template/package/logger"
	"golang-rest-api-template/package/middleware"
	"net"
	"net/http"

	_ "golang-rest-api-template/docs"
	catApi "golang-rest-api-template/internal/handlers/category"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Server struct {
	httpList net.Listener
	httpSrvr *http.Server
	logger   *logger.Logger
}

func NewServer(address string, logger *logger.Logger, db *database.Database) (*Server, error) {

	// Create a new router
	router := mux.NewRouter()

	// swagger docs
	// Serve Swagger UI
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r := router.PathPrefix("/api").Subrouter()

	// health check API
	healthAPI := health.NewHealthAPI(logger)
	healthAPI.RegisterHandlers(r)

	//Create versioned subrouter (e.g., /v1)
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

	// initialize product service
	productService := product.NewProductService(repo, logger)

	// initialize product handler
	productHandler := prodApi.NewProductAPI(logger, productService)

	// Register product handlers
	productHandler.RegisterHandlers(apiV1)

	// initialize category service
	categoryService := category.NewCategoryService(repo, logger)

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

// ServerUp starts the HTTP server and listens for incoming requests.
// It returns an error if the server fails to start.
func (srv *Server) ListenAndServe() error {
	return srv.httpSrvr.Serve(srv.httpList)
}

// ServerDown gracefully shuts down the HTTP server. It takes a context.Context
// argument to allow for cancellation of the shutdown process.
func (srv *Server) ServerDown(ctx context.Context) error {
	return srv.httpSrvr.Shutdown(ctx)
}
