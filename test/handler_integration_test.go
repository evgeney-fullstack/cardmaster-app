package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/evgeney-fullstack/cardmaster-app/internal/app/handler"
	"github.com/evgeney-fullstack/cardmaster-app/internal/app/repository/mongodb"
	"github.com/evgeney-fullstack/cardmaster-app/internal/app/repository/redis"
	"github.com/evgeney-fullstack/cardmaster-app/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestMain installing the test environment
func TestMain(m *testing.M) {
	// Set Gin mode to test mode
	gin.SetMode(gin.TestMode)

	// Run the tests
	m.Run()
}

// setupTestContainers configures MongoDB and Redis containers for testing
func setupTestContainers(ctx context.Context) (mongodb.Config, redis.Config, func(), error) {
	// Launching a MongoDB container using GenericContainer
	mongoReq := testcontainers.ContainerRequest{
		Image:        "mongo:6.0",
		ExposedPorts: []string{"27017/tcp"},
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": "testuser",
			"MONGO_INITDB_ROOT_PASSWORD": "testpass",
		},
		WaitingFor: wait.ForLog("Waiting for connections").WithStartupTimeout(30 * time.Second),
	}

	mongoContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: mongoReq,
		Started:          true,
	})
	if err != nil {
		return mongodb.Config{}, redis.Config{}, nil, fmt.Errorf("failed to start MongoDB container: %w", err)
	}

	// Launching a Redis container using GenericContainer
	redisReq := testcontainers.ContainerRequest{
		Image:        "redis:7.0-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections").WithStartupTimeout(30 * time.Second),
	}

	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: redisReq,
		Started:          true,
	})
	if err != nil {
		return mongodb.Config{}, redis.Config{}, nil, fmt.Errorf("failed to start Redis container: %w", err)
	}

	// Getting MongoDB connection settings
	mongoHost, err := mongoContainer.Host(ctx)
	if err != nil {
		return mongodb.Config{}, redis.Config{}, nil, err
	}
	mongoPort, err := mongoContainer.MappedPort(ctx, "27017")
	if err != nil {
		return mongodb.Config{}, redis.Config{}, nil, err
	}

	// Getting the connection parameters to Redis
	redisHost, err := redisContainer.Host(ctx)
	if err != nil {
		return mongodb.Config{}, redis.Config{}, nil, err
	}
	redisPort, err := redisContainer.MappedPort(ctx, "6379")
	if err != nil {
		return mongodb.Config{}, redis.Config{}, nil, err
	}

	// Configuration for connecting to test containers
	mongoCfg := mongodb.Config{
		User:     "testuser",
		Password: "testpass",
		Host:     mongoHost,
		Port:     mongoPort.Port(),
	}

	redisCfg := redis.Config{
		Addr:     redisHost,
		Port:     redisPort.Port(),
		Password: "",
		DB:       "0",
	}

	// Cleaning function for stopping containers
	cleanup := func() {
		if terminateErr := mongoContainer.Terminate(ctx); terminateErr != nil {
			log.Printf("failed to terminate MongoDB container: %v", terminateErr)
		}
		if terminateErr := redisContainer.Terminate(ctx); terminateErr != nil {
			log.Printf("failed to terminate Redis container: %v", terminateErr)
		}
	}

	return mongoCfg, redisCfg, cleanup, nil
}

// setupTestServer creates and configures a test server
func setupTestServer(mongoCfg mongodb.Config, redisCfg redis.Config) (*gin.Engine, error) {
	// MongoDB Initialization
	mongoClient, err := mongodb.NewMongoDB(mongoCfg)
	if err != nil {
		return nil, err
	}

	// Initializing Redis
	redisClient, err := redis.NewRedisDB(redisCfg)
	if err != nil {
		return nil, err
	}

	// Initializing repositories
	repos := mongodb.NewRepository(mongoClient)
	redisRepo := redis.NewCacheRepository(redisClient)

	// Initialization of services
	services := service.NewService(repos, redisRepo)

	// Initializing handlers
	handler := handler.NewHandler(services)

	// Setting up routes
	router := handler.InitRoutes()

	return router, nil
}

// TestSignUpIntegration is testing the endpoint of user registration
func TestSignUpIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	ctx := context.Background()

	// Setting up test containers
	mongoCfg, redisCfg, cleanup, err := setupTestContainers(ctx)
	defer cleanup()

	if err != nil {
		t.Fatalf("Failed to set up test containers: %v", err)
	}

	// Test Server Setup
	router, err := setupTestServer(mongoCfg, redisCfg)
	if err != nil {
		t.Fatalf("The test server could not be configured: %v", err)
	}

	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Successful registration",
			payload: map[string]interface{}{
				"email":         "test@new.com",
				"password_hash": "password123",
				"username":      "test_user",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "successful")
			},
		},
		{
			name: "Error - the email already exists",
			payload: map[string]interface{}{
				"email":         "test@new.com",
				"password_hash": "password123",
				"username":      "anotheruser",
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "error")
			},
		},
		{
			name: "Error - incorrect data",
			payload: map[string]interface{}{
				"email":         "invalid-email", // Invalid email format
				"password_hash": "short",         // The password is too short
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "Error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Request preparation
			body, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/auth/sign-up", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Request execution
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			// Checking the response status
			assert.Equal(t, tt.expectedStatus, recorder.Code)

			// Checking the response body
			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder)
			}
		})
	}
}
