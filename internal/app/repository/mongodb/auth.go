package mongodb

import (
	"context"
	"strings"

	"github.com/evgeney-fullstack/cardmaster-app/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
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

	// Insert user document into collection
	_, err := userCol.InsertOne(context.Background(), &user)
	if err != nil {
		return err
	}

	return err
}

// GetUser returns the user by credentials (email/username and password)
// If the email is specified (not empty after trimming spaces), it searches by email and password
// Otherwise it searches by username and password
func (r *AuthRepository) GetUser(email, username, password string) (models.User, error) {
	var user models.User

	var filter interface{}

	// We define the search criteria depending on the availability of email
	if len(strings.TrimSpace(email)) != 0 {
		filter = bson.M{
			"email":    email,
			"password": password,
		}
	} else {
		filter = bson.M{
			"username": username,
			"password": password,
		}

	}

	// Making a request to the user collection
	err := userCol.FindOne(context.Background(), filter).Decode(&user)

	return user, err
}

// SaveRefreshTokenToDB saves the update token to the database
func (r *AuthRepository) SaveRefreshTokenToDB(token models.RefreshToken) error {
	// Inserting a new document with a token into the refreshTokenCol collection
	_, err := refreshTokenCol.InsertOne(context.Background(), &token)

	return err
}
