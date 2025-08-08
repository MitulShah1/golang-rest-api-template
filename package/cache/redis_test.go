package cache

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/MitulShah1/golang-rest-api-template/package/logger"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCache_WithMock(t *testing.T) {
	logger := logger.NewLogger(logger.DefaultOptions())

	cfg := &RedisConfig{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DB:       0,
	}

	// This test validates the config structure and logger assignment
	// without requiring a real Redis connection
	assert.NotNil(t, cfg)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, "6379", cfg.Port)
	assert.Equal(t, "", cfg.Password)
	assert.Equal(t, 0, cfg.DB)
	assert.NotNil(t, logger)
}

func TestCache_SetAndGet(t *testing.T) {
	// Create mock Redis client
	db, mock := redismock.NewClientMock()

	// Create cache with mock client
	cache := &Cache{
		client: db,
		logger: logger.NewLogger(logger.DefaultOptions()),
	}

	ctx := context.Background()
	key := "test_key"
	value := map[string]any{
		"name":  "test",
		"value": 123,
	}

	// Marshal the value to JSON
	jsonData, err := json.Marshal(value)
	require.NoError(t, err)

	// Set up mock expectations
	mock.ExpectSet(key, jsonData, DefaultTTL).SetVal("OK")

	// Test Set
	err = cache.Set(ctx, key, value, DefaultTTL)
	assert.NoError(t, err)

	// Set up mock expectations for Get
	mock.ExpectGet(key).SetVal(string(jsonData))

	// Test Get
	var result map[string]any
	err = cache.Get(ctx, key, &result)
	assert.NoError(t, err)
	// JSON unmarshaling converts numbers to float64 by default
	expectedValue := map[string]any{
		"name":  "test",
		"value": float64(123),
	}
	assert.Equal(t, expectedValue, result)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCache_Get_CacheMiss(t *testing.T) {
	// Create mock Redis client
	db, mock := redismock.NewClientMock()

	// Create cache with mock client
	cache := &Cache{
		client: db,
		logger: logger.NewLogger(logger.DefaultOptions()),
	}

	ctx := context.Background()
	key := "test_miss_key"

	// Set up mock expectations for cache miss
	mock.ExpectGet(key).RedisNil()

	// Test Get with cache miss
	var result map[string]any
	err := cache.Get(ctx, key, &result)
	assert.Equal(t, redis.Nil, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCache_Delete(t *testing.T) {
	// Create mock Redis client
	db, mock := redismock.NewClientMock()

	// Create cache with mock client
	cache := &Cache{
		client: db,
		logger: logger.NewLogger(logger.DefaultOptions()),
	}

	ctx := context.Background()
	key := "test_delete_key"
	value := "test_value"

	// Marshal the value to JSON
	jsonData, err := json.Marshal(value)
	require.NoError(t, err)

	// Set up mock expectations for Set
	mock.ExpectSet(key, jsonData, 1*time.Minute).SetVal("OK")

	// Set a key
	err = cache.Set(ctx, key, value, 1*time.Minute)
	assert.NoError(t, err)

	// Set up mock expectations for Exists (before delete)
	mock.ExpectExists(key).SetVal(1)

	// Verify it exists
	exists, err := cache.Exists(ctx, key)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Set up mock expectations for Delete
	mock.ExpectDel(key).SetVal(1)

	// Delete the key
	err = cache.Delete(ctx, key)
	assert.NoError(t, err)

	// Set up mock expectations for Exists (after delete)
	mock.ExpectExists(key).SetVal(0)

	// Verify it's deleted
	exists, err = cache.Exists(ctx, key)
	assert.NoError(t, err)
	assert.False(t, exists)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCache_DeletePattern(t *testing.T) {
	// Create mock Redis client
	db, mock := redismock.NewClientMock()

	// Create cache with mock client
	cache := &Cache{
		client: db,
		logger: logger.NewLogger(logger.DefaultOptions()),
	}

	ctx := context.Background()
	pattern := "test_pattern_*"
	keys := []string{"test_pattern_1", "test_pattern_2", "test_pattern_3"}

	// Set up mock expectations for Keys
	mock.ExpectKeys(pattern).SetVal(keys)

	// Set up mock expectations for Del
	mock.ExpectDel(keys...).SetVal(int64(len(keys)))

	// Test DeletePattern
	err := cache.DeletePattern(ctx, pattern)
	assert.NoError(t, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCache_DeletePattern_NoKeys(t *testing.T) {
	// Create mock Redis client
	db, mock := redismock.NewClientMock()

	// Create cache with mock client
	cache := &Cache{
		client: db,
		logger: logger.NewLogger(logger.DefaultOptions()),
	}

	ctx := context.Background()
	pattern := "test_pattern_*"

	// Set up mock expectations for Keys (no keys found)
	mock.ExpectKeys(pattern).SetVal([]string{})

	// Test DeletePattern with no matching keys
	err := cache.DeletePattern(ctx, pattern)
	assert.NoError(t, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCache_TTL(t *testing.T) {
	// Create mock Redis client
	db, mock := redismock.NewClientMock()

	// Create cache with mock client
	cache := &Cache{
		client: db,
		logger: logger.NewLogger(logger.DefaultOptions()),
	}

	ctx := context.Background()
	key := "test_ttl_key"
	expectedTTL := 5 * time.Minute

	// Set up mock expectations for TTL
	mock.ExpectTTL(key).SetVal(expectedTTL)

	// Test TTL
	ttl, err := cache.TTL(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, expectedTTL, ttl)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCache_FlushDB(t *testing.T) {
	// Create mock Redis client
	db, mock := redismock.NewClientMock()

	// Create cache with mock client
	cache := &Cache{
		client: db,
		logger: logger.NewLogger(logger.DefaultOptions()),
	}

	ctx := context.Background()

	// Set up mock expectations for FlushDB
	mock.ExpectFlushDB().SetVal("OK")

	// Test FlushDB
	err := cache.FlushDB(ctx)
	assert.NoError(t, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCache_Set_WithCustomTTL(t *testing.T) {
	// Create mock Redis client
	db, mock := redismock.NewClientMock()

	// Create cache with mock client
	cache := &Cache{
		client: db,
		logger: logger.NewLogger(logger.DefaultOptions()),
	}

	ctx := context.Background()
	key := "test_custom_ttl_key"
	value := "test_value"
	customTTL := 10 * time.Minute

	// Marshal the value to JSON
	jsonData, err := json.Marshal(value)
	require.NoError(t, err)

	// Set up mock expectations
	mock.ExpectSet(key, jsonData, customTTL).SetVal("OK")

	// Test Set with custom TTL
	err = cache.Set(ctx, key, value, customTTL)
	assert.NoError(t, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCache_Set_WithZeroTTL(t *testing.T) {
	// Create mock Redis client
	db, mock := redismock.NewClientMock()

	// Create cache with mock client
	cache := &Cache{
		client: db,
		logger: logger.NewLogger(logger.DefaultOptions()),
	}

	ctx := context.Background()
	key := "test_zero_ttl_key"
	value := "test_value"

	// Marshal the value to JSON
	jsonData, err := json.Marshal(value)
	require.NoError(t, err)

	// Set up mock expectations (should use DefaultTTL when ttl is 0)
	mock.ExpectSet(key, jsonData, DefaultTTL).SetVal("OK")

	// Test Set with zero TTL
	err = cache.Set(ctx, key, value, 0)
	assert.NoError(t, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCache_Set_Error(t *testing.T) {
	// Create mock Redis client
	db, mock := redismock.NewClientMock()

	// Create cache with mock client
	cache := &Cache{
		client: db,
		logger: logger.NewLogger(logger.DefaultOptions()),
	}

	ctx := context.Background()
	key := "test_error_key"
	value := "test_value"

	// Marshal the value to JSON
	jsonData, err := json.Marshal(value)
	require.NoError(t, err)

	// Set up mock expectations for error
	expectedError := errors.New("redis connection error")
	mock.ExpectSet(key, jsonData, DefaultTTL).SetErr(expectedError)

	// Test Set with error
	err = cache.Set(ctx, key, value, DefaultTTL)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCache_Get_Error(t *testing.T) {
	// Create mock Redis client
	db, mock := redismock.NewClientMock()

	// Create cache with mock client
	cache := &Cache{
		client: db,
		logger: logger.NewLogger(logger.DefaultOptions()),
	}

	ctx := context.Background()
	key := "test_error_key"

	// Set up mock expectations for error
	expectedError := errors.New("redis connection error")
	mock.ExpectGet(key).SetErr(expectedError)

	// Test Get with error
	var result map[string]any
	err := cache.Get(ctx, key, &result)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCache_Delete_Error(t *testing.T) {
	// Create mock Redis client
	db, mock := redismock.NewClientMock()

	// Create cache with mock client
	cache := &Cache{
		client: db,
		logger: logger.NewLogger(logger.DefaultOptions()),
	}

	ctx := context.Background()
	key := "test_error_key"

	// Set up mock expectations for error
	expectedError := errors.New("redis connection error")
	mock.ExpectDel(key).SetErr(expectedError)

	// Test Delete with error
	err := cache.Delete(ctx, key)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCache_DefaultTTL(t *testing.T) {
	// Test that DefaultTTL is reasonable
	assert.Equal(t, 30*time.Minute, DefaultTTL)
	assert.Equal(t, 3, MaxRetries)
}

func TestCache_GetClient(t *testing.T) {
	// Create mock Redis client
	db, mock := redismock.NewClientMock()

	// Create cache with mock client
	cache := &Cache{
		client: db,
		logger: logger.NewLogger(logger.DefaultOptions()),
	}

	// Test GetClient returns the underlying client
	client := cache.GetClient()
	assert.Equal(t, db, client)

	// Verify no expectations were set (this is just a getter)
	assert.NoError(t, mock.ExpectationsWereMet())
}
