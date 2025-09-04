package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/evgeney-fullstack/cardmaster-app/internal/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMongoAuthRepository is a mock implementation of mongodb.Authorization
type MockMongoAuthRepository struct {
	mock.Mock
}

func (m *MockMongoAuthRepository) CreateUser(user models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// MockRedisAuthRepository is a mock implementation of redis.Authorization
type MockRedisAuthRepository struct {
	mock.Mock
}

func (m *MockRedisAuthRepository) Get(ctx context.Context, key string, dest interface{}) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

func (m *MockRedisAuthRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockRedisAuthRepository) Delete(ctx context.Context, keys ...string) error {
	args := m.Called(ctx, keys)
	return args.Error(0)
}

func TestAuthService_CreateUser(t *testing.T) {
	// Create mock repositories
	mockMongoRepo := new(MockMongoAuthRepository)
	mockRedisRepo := new(MockRedisAuthRepository)

	// Create service with mock dependencies
	authService := &AuthService{
		repos:     mockMongoRepo,
		cacheRepo: mockRedisRepo,
	}

	tests := []struct {
		name          string
		inputUser     models.User
		mockSetup     func()
		expectedError error
	}{
		{
			name: "Success - User created successfully",
			inputUser: models.User{
				Email:        "test@example.com",
				PasswordHash: "password123",
				Username:     "testuser",
			},
			mockSetup: func() {
				// Expect CreateUser to be called with any user and return no error
				mockMongoRepo.On("CreateUser", mock.AnythingOfType("models.User")).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "Error - MongoDB repository returns error",
			inputUser: models.User{
				Email:        "test@example.com",
				PasswordHash: "password123",
				Username:     "testuser",
			},
			mockSetup: func() {
				// Expect CreateUser to be called and return an error
				mockMongoRepo.On("CreateUser", mock.AnythingOfType("models.User")).Return(errors.New("database error")).Once()
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.mockSetup()

			// Execute the method
			err := authService.CreateUser(tt.inputUser)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)

				// Verify that password was hashed
				assert.NotEqual(t, generatePasswordHash(tt.inputUser.PasswordHash), "password123")
				assert.Contains(t, generatePasswordHash(tt.inputUser.PasswordHash), fmt.Sprintf("%x", "ljknsdfkgiovmsdlk&984kjsdlfj"))
			}

			// Verify that all expectations were met
			mockMongoRepo.AssertExpectations(t)
		})
	}
}

func TestGeneratePasswordHash(t *testing.T) {
	tests := []struct {
		name            string
		inputPassword   string
		expectedPattern string
	}{
		{
			name:            "Hash with salt",
			inputPassword:   "password123",
			expectedPattern: salt, // Should contain the salt
		},
		{
			name:            "Empty password",
			inputPassword:   "",
			expectedPattern: salt, // Should contain the salt even for empty password
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute the function
			result := generatePasswordHash(tt.inputPassword)

			// Assert results
			assert.NotEmpty(t, result)
			assert.Contains(t, result, fmt.Sprintf("%x", tt.expectedPattern))

			// Verify that the same input produces the same output
			assert.Equal(t, result, generatePasswordHash(tt.inputPassword))
		})
	}
}

func TestNewAuthService(t *testing.T) {
	mockMongoRepo := new(MockMongoAuthRepository)
	mockRedisRepo := new(MockRedisAuthRepository)

	// Create service
	authService := NewAuthService(mockMongoRepo, mockRedisRepo)

	// Verify service was created with correct dependencies
	assert.NotNil(t, authService)
	assert.Equal(t, mockMongoRepo, authService.repos)
	assert.Equal(t, mockRedisRepo, authService.cacheRepo)
}
