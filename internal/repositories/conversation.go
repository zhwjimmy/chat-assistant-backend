package repositories

import (
	"strings"

	"chat-assistant-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ConversationRepository defines the interface for conversation repository
type ConversationRepository interface {
	GetByID(id uuid.UUID) (*models.Conversation, error)
	GetByUserID(userID uuid.UUID, page, limit int) ([]*models.Conversation, int64, error)
	Create(conversation *models.Conversation) error
	Update(conversation *models.Conversation) error
	Delete(id uuid.UUID) error
	FindAll() ([]*models.Conversation, error)
	ReplaceTags(conversationID uuid.UUID, tagIDs []string) error
}

// ConversationRepositoryImpl handles conversation data access
type ConversationRepositoryImpl struct {
	db *gorm.DB
}

// NewConversationRepository creates a new conversation repository
func NewConversationRepository(db *gorm.DB) ConversationRepository {
	return &ConversationRepositoryImpl{
		db: db,
	}
}

// GetByID retrieves a conversation by ID
func (r *ConversationRepositoryImpl) GetByID(id uuid.UUID) (*models.Conversation, error) {
	var conversation models.Conversation
	err := r.db.Preload("Tags").Where("id = ?", id).First(&conversation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil conversation and nil error for not found
		}
		return nil, err
	}
	return &conversation, nil
}

// GetByUserID retrieves conversations by user ID with pagination
func (r *ConversationRepositoryImpl) GetByUserID(userID uuid.UUID, page, limit int) ([]*models.Conversation, int64, error) {
	var conversations []*models.Conversation
	var total int64

	// Count total conversations for this user
	err := r.db.Model(&models.Conversation{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated conversations
	offset := (page - 1) * limit
	err = r.db.Preload("Tags").Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&conversations).Error
	if err != nil {
		return nil, 0, err
	}

	return conversations, total, nil
}

// Create creates a new conversation
func (r *ConversationRepositoryImpl) Create(conversation *models.Conversation) error {
	return r.db.Create(conversation).Error
}

// Update updates an existing conversation
func (r *ConversationRepositoryImpl) Update(conversation *models.Conversation) error {
	return r.db.Save(conversation).Error
}

// Delete soft deletes a conversation by ID
func (r *ConversationRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Conversation{}, id).Error
}

// ReplaceTags replaces all tags for a conversation
func (r *ConversationRepositoryImpl) ReplaceTags(conversationID uuid.UUID, tagIDs []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除所有现有的标签关系
		err := tx.Exec("DELETE FROM conversation_tags WHERE conversation_id = ?", conversationID).Error
		if err != nil {
			return err
		}

		// 如果有新的标签，插入新的关系
		if len(tagIDs) > 0 {
			// 构建批量插入的 SQL
			values := make([]string, len(tagIDs))
			args := make([]interface{}, len(tagIDs)*2)

			for i, tagID := range tagIDs {
				values[i] = "(?, ?)"
				args[i*2] = conversationID
				args[i*2+1] = tagID
			}

			query := "INSERT INTO conversation_tags (conversation_id, tag_id) VALUES " +
				strings.Join(values, ", ") + " ON CONFLICT DO NOTHING"

			err = tx.Exec(query, args...).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *ConversationRepositoryImpl) FindAll() ([]*models.Conversation, error) {
	var conversations []*models.Conversation

	// 预加载 messages 和 tags，按创建时间排序
	err := r.db.Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).Preload("Tags").Order("created_at ASC").Find(&conversations).Error
	if err != nil {
		return nil, err
	}

	return conversations, nil
}
