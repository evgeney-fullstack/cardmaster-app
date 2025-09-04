package service

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/evgeney-fullstack/cardmaster-app/internal/app/models"
	"github.com/evgeney-fullstack/cardmaster-app/internal/app/repository/mongodb"
	"github.com/evgeney-fullstack/cardmaster-app/internal/app/repository/redis"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// Salt for password hashing to prevent rainbow table attacks
	salt = "ljknsdfkgiovmsdlk&984kjsdlfj"
)

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
	user.Id = primitive.NewObjectID()

	// Hash password before storing
	user.PasswordHash = generatePasswordHash(user.PasswordHash)

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
