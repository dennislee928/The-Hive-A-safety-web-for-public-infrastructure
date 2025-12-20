package redis

import (
	"context"
	"fmt"

	"github.com/erh-safety-system/poc/internal/config"
	"github.com/go-redis/redis/v8"
)

// Client is the global Redis client instance
var Client *redis.Client

// Init initializes the Redis connection
func Init(cfg *config.RedisConfig) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	
	// Test connection
	ctx := context.Background()
	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	
	return nil
}

// Close closes the Redis connection
func Close() error {
	if Client == nil {
		return nil
	}
	return Client.Close()
}

