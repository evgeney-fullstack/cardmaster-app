package mongodb

import (
	"github.com/evgeney-fullstack/cardmaster-app/internal/app/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// Authorization interface defines user authentication methods
type Authorization interface {
	CreateUser(user models.User) error
}

// Repository provides data access methods for MongoDB
// Implements the Authorization interface
type Repository struct {
	Authorization
}

// NewRepository creates a new MongoDB repository instance
// Initializes the Authorization interface implementation
func NewRepository(mdb *mongo.Client) *Repository {
	return &Repository{
		Authorization: NewAuthRepository(mdb),
	}
}
