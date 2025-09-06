package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user entity in the system
// Contains authentication and profile information
type User struct {
	ID         primitive.ObjectID `json:"-"  bson:"_id"` // ID is excluded from JSON for security
	Email      string             `json:"email" bson:"email" binding:"required"`
	Password   string             `json:"password" bson:"password" binding:"required"` // Hashed password, never store plain text
	Username   string             `json:"username" bson:"username" binding:"required"`
	IsVerified bool               `json:"is_verified" bson:"is_verified"` // Email verification status
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`   // Timestamp when user was created
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`   // Timestamp when user was last updated
}

// SignInRequest represents the authentication request payload
// Requires either email or username along with password
type SignInRequest struct {
	Email    string `json:"email"`                       // User's email address (optional if username is provided)
	Password string `json:"password" binding:"required"` // User's password (will be hashed on server)
	Username string `json:"username"`                    // User's username (optional if email is provided)
}

// RefreshToken represents a refresh token entity for authentication
// Used to maintain user sessions and obtain new access tokens
type RefreshToken struct {
	TokenID     string             `json:"token_id" bson:"token_id"`       // Unique token identifier
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`         // Reference to user ID
	ExpiresAt   time.Time          `json:"expires_at" bson:"expires_at"`   // Token expiration timestamp
	SignedToken string             `json:"signed_token" bson:"sign_token"` // Signed JWT string
}
