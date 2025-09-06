package service

import (
	"crypto/sha1"
	"fmt"
	"os"
	"time"

	"github.com/evgeney-fullstack/cardmaster-app/internal/app/models"
	"github.com/evgeney-fullstack/cardmaster-app/internal/app/repository/mongodb"
	"github.com/evgeney-fullstack/cardmaster-app/internal/app/repository/redis"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Authentication-related constants
const (
	// Salt for password hashing to prevent rainbow table attacks
	// In production, consider using per-user salts and more secure storage
	salt = "ljknsdfkgiovmsdlk&984kjsdlfj"

	// Access token time-to-live duration (2 hours)
	// Short-lived token for API authorization
	accessTokenTTL = 2 * time.Hour

	// Refresh token time-to-live duration (7 days)
	// Long-lived token for obtaining new access tokens without re-authentication
	refreshTokenTTL = 7 * 24 * time.Hour
)

// accessTokenClaims defines the JWT claims structure for access tokens
// Embeds jwt.StandardClaims for standard JWT fields (exp, iat, etc.)
type accessTokenClaims struct {
	jwt.StandardClaims
	// UserID contains the MongoDB ObjectID of the authenticated user
	// JSON and BSON tags exclude this field from serialization for security
	UserID primitive.ObjectID `json:"-"  bson:"_id"`
}

// refreshTokenClaims defines the JWT claims structure for refresh tokens
// Includes additional TokenID field for token revocation capabilities
type refreshTokenClaims struct {
	jwt.StandardClaims
	// UserID associates the token with a specific user account
	UserID primitive.ObjectID `json:"-"  bson:"_id"`
	// TokenID provides a unique identifier for individual refresh tokens
	// This allows for token-specific revocation and management
	TokenID string `json:"token_id"`
}

// AuthService implements business logic for authentication
type AuthService struct {
	repos     mongodb.Authorization
	cacheRepo redis.Authorization
}

// NewAuthService creates a new authentication service instance
func NewAuthService(repos mongodb.Authorization, cacheRepo redis.Authorization) *AuthService {
	return &AuthService{
		repos:     repos,
		cacheRepo: cacheRepo,
	}
}

// CreateUser creates a new user with hashed password
func (s *AuthService) CreateUser(user models.User) error {
	// Generate unique ID for the user
	user.ID = primitive.NewObjectID()

	// Hash password before storing
	user.Password = generatePasswordHash(user.Password)

	// Set creation and update timestamps
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Delegate to repository for database operation
	return s.repos.CreateUser(user)
}

// generatePasswordHash creates a SHA1 hash of the password with salt
func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	// Return hexadecimal representation of the hash
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

// generateAccessToken creates a JWT access token for the specified user
// Uses HS256 signing method and includes user ID in custom claims
// Token expiration is set based on accessTokenTTL duration
func generateAccessToken(user models.User) (string, error) {
	// Create new token with custom claims containing user ID and standard JWT claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &accessTokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenTTL).Unix(), // Token expiration time
			IssuedAt:  time.Now().Unix(),                     // Token creation time
		},
		user.ID, // User identifier included in token claims
	})

	// Sign the token using secret key from environment variables
	return token.SignedString([]byte(os.Getenv("JWT_ACCESS_SECRET_KEY")))
}

// generateRefreshToken creates both a JWT refresh token and a corresponding database record
// Generates a unique token ID and sets expiration based on refreshTokenTTL
// Returns complete refresh token object with signed token and metadata
func generateRefreshToken(user models.User) (models.RefreshToken, error) {
	var refreshToken models.RefreshToken

	// Generate unique identifier for the refresh token
	tokenID := uuid.New().String()
	// Calculate token expiration time
	expiresAt := time.Now().Add(refreshTokenTTL)

	// Create JWT token with custom claims including user ID and token ID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &refreshTokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),  // Token expiration timestamp
			IssuedAt:  time.Now().Unix(), // Token creation timestamp
		},
		UserID:  user.ID, // Associate token with user
		TokenID: tokenID, // Unique token identifier
	})

	// Sign the token using refresh token secret from environment
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET_KEY")))
	if err != nil {
		return refreshToken, err
	}

	// Populate refresh token model with generated data
	refreshToken.TokenID = tokenID
	refreshToken.UserID = user.ID
	refreshToken.ExpiresAt = expiresAt
	refreshToken.SignedToken = signedToken

	return refreshToken, nil
}

// AuthenticateUser verifies user credentials and generates authentication tokens
// Returns access token, refresh token, and user information upon successful authentication
func (s *AuthService) AuthenticateUser(input models.SignInRequest) (map[string]interface{}, error) {
	// Retrieve user from database using hashed password for security
	user, err := s.repos.GetUser(input.Email, input.Username, generatePasswordHash(input.Password))
	if err != nil {
		return map[string]interface{}{}, err
	}

	// Generate short-lived access token
	accessToken, err := generateAccessToken(user)
	if err != nil {
		return map[string]interface{}{}, err
	}

	// Generate long-lived refresh token with database persistence
	refreshToken, err := generateRefreshToken(user)
	if err != nil {
		return map[string]interface{}{}, err
	}

	// Store refresh token in database for future validation
	err = s.repos.SaveRefreshTokenToDB(refreshToken)
	if err != nil {
		return map[string]interface{}{}, err
	}

	// Return authentication response with tokens and user information
	return map[string]interface{}{
		"access_token":  accessToken,              // Short-lived token for API access
		"refresh_token": refreshToken.SignedToken, // Long-lived token for obtaining new access tokens
		"username":      user.Username,            // Authenticated user's username
	}, nil
}
