package mongodb

import (
	"context"
	"time"

	"github.com/evgeney-fullstack/cardmaster-app/internal/app/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AuthRepository implements MongoDB-specific authentication operations
type AuthRepository struct {
	mdb *mongo.Client
}

// NewAuthRepository creates a new authentication repository instance
func NewAuthRepository(mdb *mongo.Client) *AuthRepository {
	return &AuthRepository{mdb: mdb}
}

// CreateUser inserts a new user into the database
// Generates unique ID and timestamps before insertion
func (r *AuthRepository) CreateUser(user models.User) error {
	// Generate unique ID for the user
	user.Id = primitive.NewObjectID()

	// Set creation and update timestamps
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Insert user document into collection
	_, err := userCol.InsertOne(context.Background(), &user)
	if err != nil {
		return err
	}

	return err
}
