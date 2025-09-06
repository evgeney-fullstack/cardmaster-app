package handler

import (
	"github.com/evgeney-fullstack/cardmaster-app/internal/app/service"
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests and manages routing.
// Contains dependencies required for request handlers (future fields).
type Handler struct {
	services *service.Service
}

// NewHandler creates and returns a new Handler instance.
// Constructor function for initializing a handler with possible dependencies.
func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

// InitRoutes configures and returns the Gin router with defined endpoints.
// Adds middleware and registers handlers for all API paths.
func (h *Handler) InitRoutes() *gin.Engine {

	router := gin.New()
	// Auth group handles all authentication-related endpoints
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)              // User registration
		auth.POST("/sign-in", h.signIn)              // User login
		auth.POST("/refresh-token", h.refreshTokens) // Token refresh
		auth.POST("/logout", h.logout)               // Single device logout
		auth.POST("/logout-all", h.logoutAll)        // All devices logout
	}

	return router
}
