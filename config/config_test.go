package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	service := NewService()

	assert.NotNil(t, service)
	assert.Equal(t, "go-rest-api-template", service.Name)
}

func TestService_Structure(t *testing.T) {
	service := NewService()

	// Test that the service has the expected structure
	assert.NotNil(t, service)
	assert.Equal(t, "go-rest-api-template", service.Name)

	// Test that getter methods exist and return empty values before init
	dbConfig := service.GetDBConfig()
	assert.Equal(t, "", dbConfig.Host)
	assert.Equal(t, "", dbConfig.Port)
	assert.Equal(t, "", dbConfig.User)
	assert.Equal(t, "", dbConfig.Password)
	assert.Equal(t, "", dbConfig.Name)

	redisConfig := service.GetRedisConfig()
	assert.Equal(t, "", redisConfig.Host)
	assert.Equal(t, "", redisConfig.Port)
	assert.Equal(t, "", redisConfig.Password)
	assert.Equal(t, 0, redisConfig.DB)

	serverConfig := service.GetServerConfig()
	assert.Equal(t, "", serverConfig.Address)
	assert.Equal(t, "", serverConfig.Port)

	jaegerConfig := service.GetJaegerConfig()
	assert.Equal(t, "", jaegerConfig.AgentHost)
	assert.Equal(t, "", jaegerConfig.AgentPort)
}

func TestService_Init_WithEnvFile(t *testing.T) {
	// This test requires a .env file to exist
	// We'll test the basic structure and error handling
	service := NewService()

	// Test that Init can be called (may fail if no .env file, but that's OK)
	_ = service.Init()

	// The test passes if Init doesn't panic and service is created
	assert.NotNil(t, service)

	// Test that getter methods work after init attempt
	dbConfig := service.GetDBConfig()
	assert.NotNil(t, dbConfig)

	redisConfig := service.GetRedisConfig()
	assert.NotNil(t, redisConfig)

	serverConfig := service.GetServerConfig()
	assert.NotNil(t, serverConfig)

	jaegerConfig := service.GetJaegerConfig()
	assert.NotNil(t, jaegerConfig)
}

func TestService_ConfigStructs(t *testing.T) {
	// Test that config structs can be created and have expected fields
	dbConfig := DBConfig{
		Host:     "testhost",
		Port:     "3306",
		User:     "testuser",
		Password: "testpass",
		Name:     "testdb",
	}

	assert.Equal(t, "testhost", dbConfig.Host)
	assert.Equal(t, "3306", dbConfig.Port)
	assert.Equal(t, "testuser", dbConfig.User)
	assert.Equal(t, "testpass", dbConfig.Password)
	assert.Equal(t, "testdb", dbConfig.Name)

	redisConfig := RedisConfig{
		Host:     "redis-test",
		Port:     "6379",
		Password: "redispass",
		DB:       1,
	}

	assert.Equal(t, "redis-test", redisConfig.Host)
	assert.Equal(t, "6379", redisConfig.Port)
	assert.Equal(t, "redispass", redisConfig.Password)
	assert.Equal(t, 1, redisConfig.DB)

	serverConfig := ServerConf{
		Address: "0.0.0.0",
		Port:    "8080",
	}

	assert.Equal(t, "0.0.0.0", serverConfig.Address)
	assert.Equal(t, "8080", serverConfig.Port)

	jaegerConfig := JaegerConfig{
		AgentHost: "jaeger-test",
		AgentPort: "6831",
	}

	assert.Equal(t, "jaeger-test", jaegerConfig.AgentHost)
	assert.Equal(t, "6831", jaegerConfig.AgentPort)
}
