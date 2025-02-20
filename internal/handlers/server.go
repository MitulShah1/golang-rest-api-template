package handlers

import (
	"context"
	"golang-rest-api-template/internal/handlers/health"
	"golang-rest-api-template/package/logger"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	httpList net.Listener
	httpSrvr *http.Server
	logger   *logger.Logger
}

func NewServer(address string, logger *logger.Logger) (*Server, error) {

	router := mux.NewRouter()

	healthAPI := health.NewHealthAPI(logger)
	healthAPI.RegisterHandlers(router)

	httpLis, err := net.Listen(`tcp`, address)
	if err != nil {
		return nil, err
	}

	return &Server{
		httpList: httpLis,
		httpSrvr: &http.Server{
			Addr: address,
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
