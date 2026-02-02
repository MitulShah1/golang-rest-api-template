// Package config provides configuration management for the application.
// It handles loading environment variables, database configuration,
// server settings, and telemetry configuration.
package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Service holds application configuration.
// It includes database, server, and telemetry configuration.
type Service struct {
	Name         string
	dbEnv        DBConfig
	redisEnv     RedisConfig
	srvConfg     ServerConf
	jaegerConfig JaegerConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
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

// Init initializes the application configuration by loading environment variables.
// It returns an error if the configuration loading fails.
func (cnf *Service) Init() error {
	return cnf.LoadConfig()
}

// LoadConfig loads configuration from environment variables
func (cnf *Service) LoadConfig() error {
	// Load .env file if present (optional - env vars may be set by Docker, etc.)
	_ = godotenv.Load()

	cnf.dbEnv = DBConfig{
		Port:     getEnv("DB_PORT", "3306"),
		Host:     getEnv("DB_HOST", "localhost"),
		User:     getEnv("DB_USER", "user"),
		Password: getEnv("DB_PASSWORD", "password"),
		Name:     getEnv("DB_NAME", "mydatabase"),
	}

	// Redis config
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	cnf.redisEnv = RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnv("REDIS_PORT", "6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       redisDB,
	}

	// Server config
	cnf.srvConfg = ServerConf{
		Address: getEnv("SERVER_ADDR", ""),
		Port:    getEnv("SERVER_PORT", "8080"),
	}

	// Jaeger config
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

// GetDBConfig returns the database configuration
func (cnf *Service) GetDBConfig() DBConfig {
	return cnf.dbEnv
}

// GetRedisConfig returns the Redis configuration
func (cnf *Service) GetRedisConfig() RedisConfig {
	return cnf.redisEnv
}

// GetServerConfig returns the server configuration
func (cnf *Service) GetServerConfig() ServerConf {
	return cnf.srvConfg
}

// GetJaegerConfig returns the Jaeger configuration
func (cnf *Service) GetJaegerConfig() JaegerConfig {
	return cnf.jaegerConfig
}
