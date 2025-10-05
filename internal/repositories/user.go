package repositories

import (
	"chat-assistant-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user repository
type UserRepository interface {
	GetByID(id uuid.UUID) (*models.User, error)
}

// UserRepositoryImpl handles user data access
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

// GetByID retrieves a user by ID
func (r *UserRepositoryImpl) GetByID(id uuid.UUID) (*models.User, error) {
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
