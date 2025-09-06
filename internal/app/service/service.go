package service

import (
	"github.com/evgeney-fullstack/cardmaster-app/internal/app/models"
	"github.com/evgeney-fullstack/cardmaster-app/internal/app/repository/mongodb"
	"github.com/evgeney-fullstack/cardmaster-app/internal/app/repository/redis"
)

// Authorization interface defines authentication service methods
type Authorization interface {
	CreateUser(user models.User) error
	AuthenticateUser(input models.SignInRequest) (map[string]interface{}, error)
}

// Service represents the business logic layer of the application
// It orchestrates data flow between repositories and handlers
type Service struct {
	Authorization
}

// NewService creates a new Service instance with dependencies injection
// Initializes the Authorization service with MongoDB and Redis repositories
func NewService(repos *mongodb.Repository, cacheRepo *redis.CacheRepository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, cacheRepo.Authorization),
	}
}
