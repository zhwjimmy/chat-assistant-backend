package services

import (
	"chat-assistant-backend/internal/errors"
	"chat-assistant-backend/internal/models"
	"chat-assistant-backend/internal/repositories"

	"github.com/google/uuid"
)

// UserService handles user business logic
type UserService struct {
	userRepo *repositories.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}
