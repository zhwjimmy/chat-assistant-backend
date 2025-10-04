package repositories

import (
	"chat-assistant-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ConversationRepositoryInterface defines the interface for conversation repository
type ConversationRepositoryInterface interface {
	GetByID(id uuid.UUID) (*models.Conversation, error)
	GetByUserID(userID uuid.UUID, page, limit int) ([]*models.Conversation, int64, error)
	Delete(id uuid.UUID) error
	FindAll() ([]*models.Conversation, error)
}

// ConversationRepository handles conversation data access
type ConversationRepository struct {
	db *gorm.DB
}

// NewConversationRepository creates a new conversation repository
func NewConversationRepository(db *gorm.DB) *ConversationRepository {
	return &ConversationRepository{
		db: db,
	}
}

// GetByID retrieves a conversation by ID
func (r *ConversationRepository) GetByID(id uuid.UUID) (*models.Conversation, error) {
	var conversation models.Conversation
	err := r.db.Where("id = ?", id).First(&conversation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil conversation and nil error for not found
		}
		return nil, err
	}
	return &conversation, nil
}

// GetByUserID retrieves conversations by user ID with pagination
func (r *ConversationRepository) GetByUserID(userID uuid.UUID, page, limit int) ([]*models.Conversation, int64, error) {
	var conversations []*models.Conversation
	var total int64

	// Count total conversations for this user
	err := r.db.Model(&models.Conversation{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated conversations
	offset := (page - 1) * limit
	err = r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&conversations).Error
	if err != nil {
		return nil, 0, err
	}

	return conversations, total, nil
}

// Delete soft deletes a conversation by ID
func (r *ConversationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Conversation{}, id).Error
}

func (r *ConversationRepository) FindAll() ([]*models.Conversation, error) {
	var conversations []*models.Conversation

	// 预加载 messages，按创建时间排序
	err := r.db.Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).Order("created_at ASC").Find(&conversations).Error
	if err != nil {
		return nil, err
	}

	return conversations, nil
}
