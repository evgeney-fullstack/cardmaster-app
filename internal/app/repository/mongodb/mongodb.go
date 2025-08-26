package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config holds MongoDB connection parameters
type Config struct {
	User     string
	Password string
	Host     string
	Port     string
}

// NewMongoDB establishes connection to MongoDB database
// Uses context with timeout for connection attempt
// Returns MongoDB client instance or error if connection fails
func NewMongoDB(cfg Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Construct MongoDB connection URI from config parameters
	mdb, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", cfg.User, cfg.Password, cfg.Host, cfg.Port)))
	if err != nil {
		return nil, err
	}

	// Verify connection is established by pinging the database
	err = mdb.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")

	return mdb, nil
}
