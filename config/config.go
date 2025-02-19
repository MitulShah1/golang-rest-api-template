package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	//Logger
	Port   string
	DBHost string
	DBPort int
	DBUser string
	DBPass string
	DBName string
}

// LoadConfig loads environment variables from .env file
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		log.Fatalf("Invalid DB_PORT: %v", err)
	}

	return &Config{
		Port:   getEnv("PORT", "8080"),
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: dbPort,
		DBUser: getEnv("DB_USER", "user"),
		DBPass: getEnv("DB_PASSWORD", "password"),
		DBName: getEnv("DB_NAME", "mydatabase"),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
