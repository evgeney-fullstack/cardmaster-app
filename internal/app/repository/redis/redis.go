package redis

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// Config holds Redis connection parameters
type Config struct {
	Addr     string
	Port     string
	Password string
	DB       string // DB number as string (will be converted to int)
}

// NewRedisDB establishes connection to Redis database
// Converts DB string to integer and validates connection with ping
// Returns Redis client instance or error if connection fails
func NewRedisDB(cfg Config) (*redis.Client, error) {

	// Convert DB string to integer
	dbEnv := cfg.DB
	db, err := strconv.Atoi(dbEnv)
	if err != nil {
		return nil, err
	}

	// Create Redis client with configuration
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       db,
	})

	ctx := context.Background()

	// Verify connection is established by pinging Redis
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
