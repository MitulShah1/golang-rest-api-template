package config

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/MitulShah1/golang-rest-api-template/internal/handlers"
	"github.com/MitulShah1/golang-rest-api-template/package/database"
	"github.com/MitulShah1/golang-rest-api-template/package/logger"
	"github.com/MitulShah1/golang-rest-api-template/package/middleware"
	"github.com/joho/godotenv"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

// Config holds application configuration
type Service struct {
	Name         string
	Logger       *logger.Logger
	Server       *handlers.Server
	dbEnv        DBConfig
	srvConfg     ServerConf
	jaegerConfig JaegerConfig
	db           *database.Database
	tp           *tracesdk.TracerProvider
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type ServerConf struct {
	Address string
	Port    string
}

type JaegerConfig struct {
	AgentHost string
	AgentPort string
}

func NewService() *Service {
	return &Service{
		Name: "go-rest-api-template",
	}
}

// Init initializes the application configuration, including loading environment variables,
// initializing the logger, database connection, and server. It returns an error if any
// of the initialization steps fail.
func (cnf *Service) Init() (err error) {

	//initiale logger
	cnf.Logger = logger.NewLogger(logger.DefaultOptions())

	//Load Env variables
	if err := cnf.LoadConfig(); err != nil {
		return err
	}

	//initiale database
	cnf.db, err = database.NewDatabase(database.DBConfig{
		Host:     cnf.dbEnv.Host,
		Port:     cnf.dbEnv.Port,
		User:     cnf.dbEnv.User,
		Password: cnf.dbEnv.Password,
		DBName:   cnf.dbEnv.Name,
	})
	if err != nil {
		return err
	}

	cnf.Logger.Info("Database connection successful")

	agentPort, _ := strconv.Atoi(cnf.jaegerConfig.AgentPort)
	tmCnfg := middleware.TelemetryConfig{
		Host:        cnf.jaegerConfig.AgentHost,
		Port:        agentPort,
		ServiceName: "go-rest-api-template",
	}

	// initialize tracer
	cnf.tp, err = tmCnfg.InitTracer()
	if err != nil {
		return err
	}

	cnf.Logger.Info("Tracer initialized")

	//initiale server
	serverAddr := cnf.srvConfg.Address + ":" + cnf.srvConfg.Port
	if cnf.Server, err = handlers.NewServer(serverAddr, cnf.Logger, cnf.db, tmCnfg); err != nil {
		return err
	}

	return nil
}

// Run starts the server and listens for termination signals.
// It runs the server in a goroutine and waits for a termination signal (SIGINT or SIGTERM).
// When a termination signal is received, it gracefully shuts down the server.
func (cnf *Service) Run() error {

	// Channel to listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Run server in a goroutine
	go func() {
		cnf.Logger.Info("Starting server port: " + cnf.srvConfg.Port)
		if err := cnf.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			cnf.Logger.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for termination signal
	<-stop
	cnf.Logger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the tracer provider
	if err := cnf.tp.Shutdown(ctx); err != nil {
		cnf.Logger.Error("Failed to shutdown tracer provider", "error", err)
	}

	if err := cnf.Server.ServerDown(ctx); err != nil {
		cnf.Logger.Error("Server shutdown failed", "error", err)
	} else {
		cnf.Logger.Info("Server exited gracefully")
	}
	return nil
}

func (cnf *Service) LoadConfig() error {
	//loads environment variables from .env file
	if err := godotenv.Load(); err != nil {
		cnf.Logger.Warn("no .env file found, using system environment variables")
	}

	cnf.dbEnv = DBConfig{
		Port:     getEnv("DB_PORT", "3306"),
		Host:     getEnv("DB_HOST", "localhost"),
		User:     getEnv("DB_USER", "user"),
		Password: getEnv("DB_PASSWORD", "password"),
		Name:     getEnv("DB_NAME", "mydatabase"),
	}

	//Server config
	cnf.srvConfg = ServerConf{
		Address: getEnv("SERVER_ADDR", ""),
		Port:    getEnv("SERVER_PORT", "8080"),
	}

	//Jaeger config
	cnf.jaegerConfig = JaegerConfig{
		AgentHost: getEnv("JAEGER_AGENT_HOST", "localhost"),
		AgentPort: getEnv("JAEGER_AGENT_PORT", "6831"),
	}

	return nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
