package repositories

import (
	"chat-assistant-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepositoryInterface defines the interface for user repository
type UserRepositoryInterface interface {
	GetByID(id uuid.UUID) (*models.User, error)
}

// UserRepository handles user data access
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil user and nil error for not found
		}
		return nil, err
	}
	return &user, nil
}
