package importer

import (
	"context"
	"fmt"

	"chat-assistant-backend/internal/config"
	"chat-assistant-backend/internal/models"
	"chat-assistant-backend/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Loader 数据加载器
type Loader struct {
	config           *config.Config
	db               *gorm.DB
	conversationRepo *repositories.ConversationRepository
	messageRepo      *repositories.MessageRepository
}

// NewLoader 创建加载器
func NewLoader(cfg *config.Config) *Loader {
	return &Loader{
		config: cfg,
	}
}

// SetDependencies 设置依赖（用于依赖注入）
func (l *Loader) SetDependencies(db *gorm.DB, conversationRepo *repositories.ConversationRepository, messageRepo *repositories.MessageRepository) {
	l.db = db
	l.conversationRepo = conversationRepo
	l.messageRepo = messageRepo
}

// Load 逐个处理数据到数据库，使用upsert确保幂等性
func (l *Loader) Load(ctx context.Context, conversations []*models.Conversation, messagesWithSource []*MessageWithConversationSource) error {
	if l.db == nil {
		return fmt.Errorf("database connection not initialized")
	}

	// 开始事务
	tx := l.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 逐个处理对话，先查询再更新/创建
	conversationIDMap := make(map[string]uuid.UUID) // 用于映射source_id到实际的conversation_id
	for _, conv := range conversations {
		var existingConv models.Conversation
		// 根据业务唯一键查询：user_id + source_id
		err := tx.Where("user_id = ? AND source_id = ?", conv.UserID, conv.SourceID).First(&existingConv).Error

		if err == gorm.ErrRecordNotFound {
			// 记录不存在，创建新记录
			if err := tx.Create(conv).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create conversation %s: %w", conv.SourceID, err)
			}
			conversationIDMap[conv.SourceID] = conv.ID
		} else if err != nil {
			// 查询出错
			tx.Rollback()
			return fmt.Errorf("failed to query conversation %s: %w", conv.SourceID, err)
		} else {
			// 记录存在，更新现有记录
			conv.ID = existingConv.ID               // 保持原有ID
			conv.CreatedAt = existingConv.CreatedAt // 保持原有创建时间
			if err := tx.Save(conv).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update conversation %s: %w", conv.SourceID, err)
			}
			conversationIDMap[conv.SourceID] = conv.ID
		}
	}

	// 逐个处理消息，先查询再更新/创建
	for _, msgWithSource := range messagesWithSource {
		msg := msgWithSource.Message
		// 使用正确的conversation_id（从conversationIDMap获取）
		actualConversationID, exists := conversationIDMap[msgWithSource.ConversationSourceID]
		if !exists {
			tx.Rollback()
			return fmt.Errorf("conversation source_id %s not found in mapping", msgWithSource.ConversationSourceID)
		}
		msg.ConversationID = actualConversationID
		var existingMsg models.Message
		// 根据业务唯一键查询：conversation_id + source_id
		err := tx.Where("conversation_id = ? AND source_id = ?", msg.ConversationID, msg.SourceID).First(&existingMsg).Error

		if err == gorm.ErrRecordNotFound {
			// 记录不存在，创建新记录
			if err := tx.Create(msg).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create message %s: %w", msg.SourceID, err)
			}
		} else if err != nil {
			// 查询出错
			tx.Rollback()
			return fmt.Errorf("failed to query message %s: %w", msg.SourceID, err)
		} else {
			// 记录存在，更新现有记录
			msg.ID = existingMsg.ID               // 保持原有ID
			msg.CreatedAt = existingMsg.CreatedAt // 保持原有创建时间
			if err := tx.Save(msg).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update message %s: %w", msg.SourceID, err)
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
