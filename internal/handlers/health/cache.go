// Package health provides health check functionality for the application.
// This file includes cache health check endpoints.
package health

import (
	"context"
	"net/http"
	"time"

	"github.com/MitulShah1/golang-rest-api-template/internal/response"
	"github.com/MitulShah1/golang-rest-api-template/package/cache"
	"github.com/MitulShah1/golang-rest-api-template/package/logger"
	"github.com/gorilla/mux"
)

// CacheHealthAPI provides cache health check endpoints
type CacheHealthAPI struct {
	logger *logger.Logger
	cache  *cache.Cache
}

// NewCacheHealthAPI creates a new cache health API instance
func NewCacheHealthAPI(logger *logger.Logger, cache *cache.Cache) *CacheHealthAPI {
	return &CacheHealthAPI{
		logger: logger,
		cache:  cache,
	}
}

// RegisterHandlers registers the cache health check routes
func (h *CacheHealthAPI) RegisterHandlers(router *mux.Router) {
	router.HandleFunc(CacheHealthPath, h.CacheHealth).Methods(http.MethodGet)
	router.HandleFunc(CacheStatsPath, h.CacheStats).Methods(http.MethodGet)
	router.HandleFunc(FlushCachePath, h.FlushCache).Methods(http.MethodPost)
}

// CacheHealth checks if Redis cache is healthy
func (h *CacheHealthAPI) CacheHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Test cache connection with a simple ping
	testKey := "health_check"
	testValue := "ok"

	// Try to set a test value
	if err := h.cache.Set(ctx, testKey, testValue, 1*time.Minute); err != nil {
		h.logger.Error("cache health check failed - set operation", "error", err)
		response.Error(w, http.StatusServiceUnavailable, "Cache is unhealthy - set operation failed")
		return
	}

	// Try to get the test value
	var retrievedValue string
	if err := h.cache.Get(ctx, testKey, &retrievedValue); err != nil {
		h.logger.Error("cache health check failed - get operation", "error", err)
		response.Error(w, http.StatusServiceUnavailable, "Cache is unhealthy - get operation failed")
		return
	}

	// Clean up test key
	if err := h.cache.Delete(ctx, testKey); err != nil {
		h.logger.Error("failed to delete test key", "error", err)
	}

	if retrievedValue != testValue {
		h.logger.Error("cache health check failed - value mismatch")
		response.Error(w, http.StatusServiceUnavailable, "Cache is unhealthy - value mismatch")
		return
	}

	response.Success(w, http.StatusOK, "Cache is healthy", map[string]any{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"service":   "redis",
	})
}

// CacheStats returns cache statistics
func (h *CacheHealthAPI) CacheStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get Redis client for advanced stats
	client := h.cache.GetClient()

	// Get Redis info
	info, err := client.Info(ctx).Result()
	if err != nil {
		h.logger.Error("failed to get cache stats", "error", err)
		response.Error(w, http.StatusInternalServerError, "Failed to get cache statistics")
		return
	}

	// Get database size
	dbSize, err := client.DBSize(ctx).Result()
	if err != nil {
		h.logger.Error("failed to get database size", "error", err)
		response.Error(w, http.StatusInternalServerError, "Failed to get database size")
		return
	}

	response.Success(w, http.StatusOK, "Cache statistics", map[string]any{
		"db_size":   dbSize,
		"timestamp": time.Now().Unix(),
		"info":      info,
	})
}

// FlushCache clears all cache data
func (h *CacheHealthAPI) FlushCache(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := h.cache.FlushDB(ctx); err != nil {
		h.logger.Error("failed to flush cache", "error", err)
		response.Error(w, http.StatusInternalServerError, "Failed to flush cache")
		return
	}

	response.Success(w, http.StatusOK, "Cache flushed successfully", map[string]any{
		"message":   "All cache data cleared",
		"timestamp": time.Now().Unix(),
	})
}
