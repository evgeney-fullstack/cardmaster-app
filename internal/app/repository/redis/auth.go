package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisRepository implements Redis-specific cache operations
type RedisRepository struct {
	client *redis.Client
}

// NewRedisRepository creates a new Redis repository instance
func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{client: client}
}

// Get retrieves a value from Redis cache and unmarshals it into the destination
func (r *RedisRepository) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("key %s: %w", key, errors.New("not found"))
		}
		return err
	}

	// Unmarshal JSON data into destination object
	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("json unmarshal error: %w", err)
	}
	return nil
}

// Set marshals a value to JSON and stores it in Redis cache with expiration
func (r *RedisRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	// Store value in Redis with TTL
	if err := r.client.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}
	return nil
}

// Delete removes one or more keys from Redis cache
func (r *RedisRepository) Delete(ctx context.Context, keys ...string) error {
	if err := r.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("redis delete error: %w", err)
	}
	return nil
}
