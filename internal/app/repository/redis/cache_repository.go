package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Authorization interface defines cache operations for authentication
type Authorization interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, keys ...string) error
}

// CacheRepository provides caching operations using Redis
// Implements the Authorization interface
type CacheRepository struct {
	Authorization
}

// NewCacheRepository creates a new Redis cache repository instance
// Initializes the Authorization interface implementation
func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{
		Authorization: NewRedisRepository(client),
	}
}
