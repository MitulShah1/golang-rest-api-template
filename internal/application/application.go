// Package application provides application lifecycle management.
// It handles initialization, configuration, and graceful shutdown of all application components.
package application

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/MitulShah1/golang-rest-api-template/config"
	"github.com/MitulShah1/golang-rest-api-template/internal/handlers"
	"github.com/MitulShah1/golang-rest-api-template/package/cache"
	"github.com/MitulShah1/golang-rest-api-template/package/database"
	"github.com/MitulShah1/golang-rest-api-template/package/logger"
	"github.com/MitulShah1/golang-rest-api-template/package/middleware"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

// Application represents the main application instance
type Application struct {
	Name         string
	Logger       *logger.Logger
	Server       *handlers.Server
	Config       *config.Service
	Database     *database.Database
	Cache        *cache.Cache
	Tracer       *tracesdk.TracerProvider
	ShutdownChan chan os.Signal
}

// NewApplication creates a new application instance
func NewApplication() *Application {
	return &Application{
		Name:         "go-rest-api-template",
		ShutdownChan: make(chan os.Signal, 1),
	}
}

// Initialize sets up all application components
func (app *Application) Initialize() error {
	// Initialize logger first
	app.Logger = logger.NewLogger(logger.DefaultOptions())
	app.Logger.Info("Starting application initialization")

	// Initialize components in order
	initializers := []struct {
		name string
		fn   func() error
	}{
		{"configuration", app.initializeConfiguration},
		{"database", app.initializeDatabase},
		{"cache", app.initializeCache},
		{"telemetry", app.initializeTelemetry},
		{"server", app.initializeServer},
	}

	for _, init := range initializers {
		if err := init.fn(); err != nil {
			app.Logger.Fatal("failed to initialize "+init.name, "error", err.Error())
			return err
		}
	}

	app.Logger.Info("Application initialization completed successfully")
	return nil
}

// Run starts the application and handles graceful shutdown
func (app *Application) Run() error {
	app.Logger.Info("Starting application", "port", app.Config.GetServerConfig().Port)

	// Set up signal handling for graceful shutdown
	signal.Notify(app.ShutdownChan, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.Logger.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	<-app.ShutdownChan
	app.Logger.Info("Shutting down application...")

	// Perform graceful shutdown
	return app.Shutdown()
}

// Shutdown gracefully shuts down all application components
func (app *Application) Shutdown() error {
	// Check if logger is available, if not create a basic one
	if app.Logger == nil {
		app.Logger = logger.NewLogger(logger.DefaultOptions())
	}

	app.Logger.Info("Starting graceful shutdown")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Define shutdown components in reverse order
	shutdownComponents := []struct {
		name string
		fn   func(context.Context) error
	}{
		{"HTTP server", app.shutdownServer},
		{"tracer", app.shutdownTracer},
		{"cache", app.shutdownCache},
		{"database", app.shutdownDatabase},
	}

	var shutdownErrors []error
	for _, component := range shutdownComponents {
		if err := component.fn(ctx); err != nil {
			app.Logger.Error("Failed to shutdown "+component.name, "error", err)
			shutdownErrors = append(shutdownErrors, err)
		} else {
			app.Logger.Info(component.name + " shutdown completed")
		}
	}

	if len(shutdownErrors) > 0 {
		app.Logger.Error("Some components failed to shutdown gracefully", "errors", shutdownErrors)
		return shutdownErrors[0]
	}

	app.Logger.Info("Application shutdown completed successfully")
	return nil
}

// GetLogger returns the application logger
func (app *Application) GetLogger() *logger.Logger {
	return app.Logger
}

// GetConfig returns the application configuration
func (app *Application) GetConfig() *config.Service {
	return app.Config
}

// GetDatabase returns the database instance
func (app *Application) GetDatabase() *database.Database {
	return app.Database
}

// GetCache returns the cache instance
func (app *Application) GetCache() *cache.Cache {
	return app.Cache
}

// GetServer returns the HTTP server instance
func (app *Application) GetServer() *handlers.Server {
	return app.Server
}

// createTelemetryConfig creates a TelemetryConfig with the application's settings
func (app *Application) createTelemetryConfig() middleware.TelemetryConfig {
	agentPort, _ := strconv.Atoi(app.Config.GetJaegerConfig().AgentPort)
	return middleware.TelemetryConfig{
		Host:        app.Config.GetJaegerConfig().AgentHost,
		Port:        agentPort,
		ServiceName: app.Name,
	}
}

// initializeConfiguration sets up the application configuration
func (app *Application) initializeConfiguration() error {
	app.Config = config.NewService()
	return app.Config.Init()
}

// initializeDatabase sets up the database connection
func (app *Application) initializeDatabase() error {
	app.Logger.Info("Initializing database connection")

	db, err := database.NewDatabase(&database.DBConfig{
		Host:     app.Config.GetDBConfig().Host,
		Port:     app.Config.GetDBConfig().Port,
		User:     app.Config.GetDBConfig().User,
		Password: app.Config.GetDBConfig().Password,
		DBName:   app.Config.GetDBConfig().Name,
	})
	if err != nil {
		return err
	}

	app.Database = db
	app.Logger.Info("Database connection established successfully")
	return nil
}

// initializeCache sets up the Redis cache connection
func (app *Application) initializeCache() error {
	app.Logger.Info("Initializing Redis cache connection")

	cache, err := cache.NewCache(&cache.RedisConfig{
		Host:     app.Config.GetRedisConfig().Host,
		Port:     app.Config.GetRedisConfig().Port,
		Password: app.Config.GetRedisConfig().Password,
		DB:       app.Config.GetRedisConfig().DB,
	}, app.Logger)
	if err != nil {
		return err
	}

	app.Cache = cache
	app.Logger.Info("Redis cache connection established successfully")
	return nil
}

// initializeTelemetry sets up OpenTelemetry tracing
func (app *Application) initializeTelemetry() error {
	app.Logger.Info("Initializing telemetry")

	tmConfig := app.createTelemetryConfig()
	tracer, err := tmConfig.InitTracer()
	if err != nil {
		return err
	}

	app.Tracer = tracer
	app.Logger.Info("Telemetry initialized successfully")
	return nil
}

// initializeServer sets up the HTTP server
func (app *Application) initializeServer() error {
	app.Logger.Info("Initializing HTTP server")

	serverAddr := app.Config.GetServerConfig().Address + ":" + app.Config.GetServerConfig().Port

	// Create telemetry config for server
	tmConfig := app.createTelemetryConfig()

	// Initialize the tracer for this config
	tracer, err := tmConfig.InitTracer()
	if err != nil {
		return err
	}

	// Set the tracer provider
	tmConfig.TraceProvider = tracer

	server, err := handlers.NewServer(serverAddr, app.Logger, app.Database, app.Cache, &tmConfig)
	if err != nil {
		return err
	}

	app.Server = server
	app.Logger.Info("HTTP server initialized successfully")
	return nil
}

// shutdownServer shuts down the HTTP server
func (app *Application) shutdownServer(ctx context.Context) error {
	if app.Server == nil {
		return nil
	}
	return app.Server.ServerDown(ctx)
}

// shutdownTracer shuts down the tracer
func (app *Application) shutdownTracer(ctx context.Context) error {
	if app.Tracer == nil {
		return nil
	}
	return app.Tracer.Shutdown(ctx)
}

// shutdownCache shuts down the cache
func (app *Application) shutdownCache(ctx context.Context) error {
	if app.Cache == nil {
		return nil
	}
	return app.Cache.Close()
}

// shutdownDatabase shuts down the database
func (app *Application) shutdownDatabase(ctx context.Context) error {
	if app.Database == nil {
		return nil
	}
	app.Database.Close()
	return nil
}
