package importer

import (
	"context"
	"fmt"

	"chat-assistant-backend/internal/config"
	"chat-assistant-backend/internal/models"
	"chat-assistant-backend/internal/repositories"

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

// Load 批量加载数据到数据库
func (l *Loader) Load(ctx context.Context, conversations []*models.Conversation, messages []*models.Message) error {
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

	// 批量插入对话
	if len(conversations) > 0 {
		if err := tx.CreateInBatches(conversations, l.config.Import.BatchSize).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create conversations: %w", err)
		}
	}

	// 批量插入消息
	if len(messages) > 0 {
		if err := tx.CreateInBatches(messages, l.config.Import.BatchSize).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create messages: %w", err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
