package health

import (
	"testing"

	"github.com/MitulShah1/golang-rest-api-template/package/cache"
	"github.com/MitulShah1/golang-rest-api-template/package/logger"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestNewCacheHealthAPI(t *testing.T) {
	logger := logger.NewLogger(logger.DefaultOptions())
	cache := &cache.Cache{}

	api := NewCacheHealthAPI(logger, cache)

	assert.NotNil(t, api)
	assert.Equal(t, logger, api.logger)
	assert.Equal(t, cache, api.cache)
}

func TestCacheHealthAPI_RegisterHandlers(t *testing.T) {
	logger := logger.NewLogger(logger.DefaultOptions())
	cache := &cache.Cache{}
	api := NewCacheHealthAPI(logger, cache)

	router := mux.NewRouter()
	api.RegisterHandlers(router)

	// Test that the API was created successfully
	assert.NotNil(t, api)
	assert.NotNil(t, api.logger)
	assert.NotNil(t, api.cache)
}

func TestCacheHealthAPI_Structure(t *testing.T) {
	logger := logger.NewLogger(logger.DefaultOptions())
	cache := &cache.Cache{}
	api := NewCacheHealthAPI(logger, cache)

	// Test that the API has the expected structure
	assert.NotNil(t, api.logger)
	assert.NotNil(t, api.cache)
}

func TestCacheHealthAPI_HandlerMethods(t *testing.T) {
	logger := logger.NewLogger(logger.DefaultOptions())
	cache := &cache.Cache{}
	api := NewCacheHealthAPI(logger, cache)

	// Test that the API structure is correct
	assert.NotNil(t, api)
	assert.NotNil(t, api.logger)
	assert.NotNil(t, api.cache)

	// Test that the methods exist (we won't call them due to nil client)
	assert.NotNil(t, api.CacheHealth)
	assert.NotNil(t, api.CacheStats)
	assert.NotNil(t, api.FlushCache)
}
