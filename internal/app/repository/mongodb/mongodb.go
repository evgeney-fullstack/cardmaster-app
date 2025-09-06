package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

const (
	dbName                 = "cardmasterDB" // Database name
	userCollection         = "users"        // User collection name
	refreshTokenCollection = "refresh_tokens"
)

var userCol *mongo.Collection         // Global user collection reference
var refreshTokenCol *mongo.Collection // Global refreshTokenCol collection reference

// NewMongoDB establishes connection to MongoDB database
// Uses context with timeout for connection attempt
// Returns MongoDB client instance or error if connection fails
func NewMongoDB(cfg Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Construct MongoDB connection URI from config parameters
	mdb, err := mongo.Connect(ctx, options.Client().ApplyURI(
		fmt.Sprintf("mongodb://%s:%s@%s:%s", cfg.User, cfg.Password, cfg.Host, cfg.Port)))
	if err != nil {
		return nil, err
	}

	// Verify connection is established by pinging the database
	err = mdb.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Initialize collections
	initCollections(mdb)

	// Create database indexes
	err = initIndexModels()
	if err != nil {
		log.Fatal("Failed to create indexes:", err)
	}

	fmt.Println("Connected to MongoDB!")

	return mdb, nil
}

// initCollections initializes all database collections
func initCollections(mdb *mongo.Client) {
	userCol = mdb.Database(dbName).Collection(userCollection)
	refreshTokenCol = mdb.Database(dbName).Collection(refreshTokenCollection)
}

// initIndexModels creates database indexes for optimal query performance
func initIndexModels() error {
	// Define indexes for user collection
	indexUserModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}}, // Index on email field
			Options: options.Index().SetUnique(true),  // Ensure email uniqueness
		},
		{
			Keys:    bson.D{{Key: "username", Value: 1}}, // Index on username field
			Options: options.Index().SetUnique(true),     // Ensure username uniqueness
		},
	}

	// Create indexes in database
	_, err := userCol.Indexes().CreateMany(context.Background(), indexUserModels)
	if err != nil {
		return err
	}

	return err
}
