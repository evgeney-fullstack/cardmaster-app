package redis

import (
	"github.com/redis/go-redis/v9"
)

// CacheRepository provides caching operations using Redis
// Will contain methods for getting/setting cached data
type CacheRepository struct {
	// TODO: Add Redis client as a field
	// client *redis.Client
}

// NewCacheRepository creates a new Redis cache repository instance
// Accepts Redis client but doesn't store it yet
func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{}
	// TODO: Store Redis client for future use
	// return &CacheRepository{
	//     client: client,
	// }
}
