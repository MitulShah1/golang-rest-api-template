// Package cache provides Redis caching functionality for the application.
// It includes connection management, basic operations, and caching utilities.
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MitulShah1/golang-rest-api-template/package/logger"
	"github.com/redis/go-redis/v9"
)

const (
	DefaultTTL = 30 * time.Minute
	MaxRetries = 3
)

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// Cache wraps the Redis client and provides caching operations
type Cache struct {
	client *redis.Client
	logger *logger.Logger
}

// NewCache initializes a new Redis cache connection
func NewCache(cfg *RedisConfig, logger *logger.Logger) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: 10,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Redis connection established successfully")

	return &Cache{
		client: client,
		logger: logger,
	}, nil
}

// Close gracefully closes the Redis connection
func (c *Cache) Close() error {
	return c.client.Close()
}

// Set stores a key-value pair in Redis with optional TTL
func (c *Cache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if ttl == 0 {
		ttl = DefaultTTL
	}

	err = c.client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		c.logger.Error("failed to set cache key", "key", key, "error", err)
		return err
	}

	c.logger.Debug("cache set successful", "key", key, "ttl", ttl)
	return nil
}

// Get retrieves a value from Redis by key
func (c *Cache) Get(ctx context.Context, key string, dest any) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			c.logger.Debug("cache miss", "key", key)
			return redis.Nil
		}
		c.logger.Error("failed to get cache key", "key", key, "error", err)
		return err
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal cached value: %w", err)
	}

	c.logger.Debug("cache hit", "key", key)
	return nil
}

// Delete removes a key from Redis
func (c *Cache) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		c.logger.Error("failed to delete cache key", "key", key, "error", err)
		return err
	}

	c.logger.Debug("cache delete successful", "key", key)
	return nil
}

// DeletePattern removes all keys matching a pattern
func (c *Cache) DeletePattern(ctx context.Context, pattern string) error {
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		c.logger.Error("failed to get keys for pattern", "pattern", pattern, "error", err)
		return err
	}

	if len(keys) > 0 {
		err = c.client.Del(ctx, keys...).Err()
		if err != nil {
			c.logger.Error("failed to delete keys by pattern", "pattern", pattern, "error", err)
			return err
		}
		c.logger.Debug("cache delete pattern successful", "pattern", pattern, "count", len(keys))
	}

	return nil
}

// Exists checks if a key exists in Redis
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		c.logger.Error("failed to check key existence", "key", key, "error", err)
		return false, err
	}

	return exists > 0, nil
}

// TTL gets the remaining time to live for a key
func (c *Cache) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := c.client.TTL(ctx, key).Result()
	if err != nil {
		c.logger.Error("failed to get TTL", "key", key, "error", err)
		return 0, err
	}

	return ttl, nil
}

// FlushDB clears all keys in the current database
func (c *Cache) FlushDB(ctx context.Context) error {
	err := c.client.FlushDB(ctx).Err()
	if err != nil {
		c.logger.Error("failed to flush database", "error", err)
		return err
	}

	c.logger.Info("cache database flushed successfully")
	return nil
}

// GetClient returns the underlying Redis client for advanced operations
func (c *Cache) GetClient() *redis.Client {
	return c.client
}
