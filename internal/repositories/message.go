package repositories

import (
	"chat-assistant-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MessageRepository defines the interface for message repository
type MessageRepository interface {
	GetByID(id uuid.UUID) (*models.Message, error)
	GetByConversationID(conversationID uuid.UUID, page, limit int) ([]*models.Message, int64, error)
	GetAll(page, limit int) ([]*models.Message, int64, error)
	Delete(id uuid.UUID) error
}

// MessageRepositoryImpl handles message data access
type MessageRepositoryImpl struct {
	db *gorm.DB
}

// NewMessageRepository creates a new message repository
func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &MessageRepositoryImpl{
		db: db,
	}
}

// GetByID retrieves a message by ID
func (r *MessageRepositoryImpl) GetByID(id uuid.UUID) (*models.Message, error) {
	var message models.Message
	err := r.db.Where("id = ?", id).First(&message).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil message and nil error for not found
		}
		return nil, err
	}
	return &message, nil
}

// GetByConversationID retrieves messages by conversation ID with pagination
func (r *MessageRepositoryImpl) GetByConversationID(conversationID uuid.UUID, page, limit int) ([]*models.Message, int64, error) {
	var messages []*models.Message
	var total int64

	// Count total messages for this conversation
	err := r.db.Model(&models.Message{}).Where("conversation_id = ?", conversationID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated messages
	offset := (page - 1) * limit
	err = r.db.Where("conversation_id = ?", conversationID).
		Order("created_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&messages).Error
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// GetAll retrieves all messages with pagination
func (r *MessageRepositoryImpl) GetAll(page, limit int) ([]*models.Message, int64, error) {
	var messages []*models.Message
	var total int64

	// Count total messages
	err := r.db.Model(&models.Message{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated messages
	offset := (page - 1) * limit
	err = r.db.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&messages).Error
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// Delete soft deletes a message by ID
func (r *MessageRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Message{}, id).Error
}
