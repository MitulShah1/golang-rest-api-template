package config

import (
	"context"
	"errors"
	"fmt"
	"golang-rest-api-template/internal/handlers"
	"golang-rest-api-template/package/logger"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Logger   *logger.Logger
	Server   *handlers.Server
	dbEnv    DBConfig
	srvConfg ServerConf
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type ServerConf struct {
	Address string
	Port    string
}

// Init Configs
func (cnf *Config) Init() (err error) {

	//Load Env variables
	if err := cnf.LoadConfig(); err != nil {
		return err
	}

	//initiale logger
	cnf.Logger = logger.NewLogger(logger.DefaultOptions())

	//initiale server
	serverAddr := cnf.srvConfg.Address + ":" + cnf.srvConfg.Port
	if cnf.Server, err = handlers.NewServer(serverAddr, cnf.Logger); err != nil {
		return err
	}

	return nil
}

func (cnf *Config) Run() error {

	// Channel to listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Run server in a goroutine
	go func() {
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

	if err := cnf.Server.ServerDown(ctx); err != nil {
		cnf.Logger.Error("Server shutdown failed", "error", err)
	} else {
		cnf.Logger.Info("Server exited gracefully")
	}
	return nil
}

func (cnf *Config) LoadConfig() error {
	//loads environment variables from .env file
	if err := godotenv.Load(); err != nil {
		return errors.New("no .env file found, using system environment variables")
	}

	//DB config
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "3306"))
	if err != nil {
		return fmt.Errorf("invalid DB_PORT: %v", err)
	}

	cnf.dbEnv = DBConfig{
		Port:     dbPort,
		Host:     getEnv("DB_HOST", "localhost"),
		User:     getEnv("DB_USER", "user"),
		Password: getEnv("DB_PASSWORD", "password"),
		Name:     getEnv("DB_NAME", "mydatabase"),
	}

	//Server config
	cnf.srvConfg = ServerConf{
		Address: getEnv("SERVER_ADDR", "localhost"),
		Port:    getEnv("SERVER_PORT", "8080"),
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
