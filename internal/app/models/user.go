package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user entity in the system
// Contains authentication and profile information
type User struct {
	Id           primitive.ObjectID `json:"-"  bson:"_id"` // ID is excluded from JSON for security
	Email        string             `json:"email" bson:"email" binding:"required"`
	PasswordHash string             `json:"password_hash" bson:"password_hash" binding:"required"` // Hashed password, never store plain text
	Username     string             `json:"username" bson:"username" binding:"required"`
	IsVerified   bool               `json:"is_verified" bson:"is_verified"` // Email verification status
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`   // Timestamp when user was created
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`   // Timestamp when user was last updated
}

// RefreshToken represents a refresh token entity for authentication
// Used to maintain user sessions and obtain new access tokens
type RefreshToken struct {
	TokenId   string             `json:"token_id" bson:"token_id"`     // Unique token identifier
	UserId    primitive.ObjectID `json:"user_id" bson:"user_id"`       // Reference to user ID
	ExpiresAt time.Time          `json:"expires_at" bson:"expires_at"` // Token expiration timestamp
}
