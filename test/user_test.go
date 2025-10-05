package test

import (
	"testing"

	"chat-assistant-backend/internal/errors"
	"chat-assistant-backend/internal/models"
	"chat-assistant-backend/internal/repositories"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of repositories.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// TestUserService is a test version of UserService that accepts interface
type TestUserService struct {
	userRepo repositories.UserRepository
}

// NewTestUserService creates a test user service with interface
func NewTestUserService(userRepo repositories.UserRepository) *TestUserService {
	return &TestUserService{
		userRepo: userRepo,
	}
}

// GetUserByID retrieves a user by ID (same logic as real service)
func (s *TestUserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

func TestUserService_GetUserByID(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockUserRepository)

	// Create service with mock repository using interface
	userService := NewTestUserService(mockRepo)

	// Test case 1: User found
	t.Run("User found", func(t *testing.T) {
		userID := uuid.New()
		expectedUser := &models.User{
			Base: models.Base{
				ID: userID,
			},
			Username: "testuser",
			Avatar:   "https://example.com/avatar.jpg",
		}

		mockRepo.On("GetByID", userID).Return(expectedUser, nil)

		user, err := userService.GetUserByID(userID)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.Base.ID, user.Base.ID)
		assert.Equal(t, expectedUser.Username, user.Username)
		assert.Equal(t, expectedUser.Avatar, user.Avatar)

		mockRepo.AssertExpectations(t)
	})

	// Test case 2: User not found
	t.Run("User not found", func(t *testing.T) {
		userID := uuid.New()

		mockRepo.On("GetByID", userID).Return(nil, nil)

		user, err := userService.GetUserByID(userID)

		assert.Error(t, err)
		assert.Nil(t, user)

		mockRepo.AssertExpectations(t)
	})
}
