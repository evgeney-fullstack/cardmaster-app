package service

import (
	"github.com/evgeney-fullstack/cardmaster-app/internal/app/repository/mongodb"
	"github.com/evgeney-fullstack/cardmaster-app/internal/app/repository/redis"
)

// Service represents the business logic layer of the application
// It orchestrates data flow between repositories and handlers
type Service struct {
	// TODO: Add fields for repositories when implemented
	// repos *mongodb.Repository
	// cache *redis.CacheRepository
}

// NewService creates a new Service instance with dependencies injection
// Currently accepts repositories but doesn't utilize them (to be implemented)
func NewService(repos *mongodb.Repository, cacheRepo *redis.CacheRepository) *Service {
	return &Service{}
	// TODO: Store dependencies in Service struct for future use
	// return &Service{
	//     repos: repos,
	//     cache: cacheRepo,
	// }
}
