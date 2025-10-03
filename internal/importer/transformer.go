package importer

import (
	"fmt"
	"time"

	"chat-assistant-backend/internal/importer/types"
	"chat-assistant-backend/internal/models"

	"github.com/google/uuid"
)

// Transformer 数据转换器
type Transformer struct{}

// NewTransformer 创建转换器
func NewTransformer() *Transformer {
	return &Transformer{}
}

// Transform 将标准化格式转换为数据库模型
func (t *Transformer) Transform(data *types.StandardFormat, userID uuid.UUID, platform string) ([]*models.Conversation, []*models.Message, error) {
	var conversations []*models.Conversation
	var messages []*models.Message

	for _, stdConv := range data.Conversations {
		// 转换对话
		conv, err := t.transformConversation(stdConv, userID, platform)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to transform conversation %s: %w", stdConv.ID, err)
		}
		conversations = append(conversations, conv)

		// 转换消息
		for _, stdMsg := range stdConv.Messages {
			msg, err := t.transformMessage(stdMsg, conv.ID)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to transform message: %w", err)
			}
			messages = append(messages, msg)
		}
	}

	return conversations, messages, nil
}

// transformConversation 转换对话
func (t *Transformer) transformConversation(stdConv *types.StandardConversation, userID uuid.UUID, platform string) (*models.Conversation, error) {
	conv := &models.Conversation{
		UserID:   userID,
		Title:    stdConv.Title,
		Provider: platform,
		Model:    stdConv.Model,
		SourceID: stdConv.ID, // 使用原始数据中的ID作为SourceID
	}

	// 设置时间
	if !stdConv.CreatedAt.IsZero() {
		conv.CreatedAt = stdConv.CreatedAt
	} else {
		conv.CreatedAt = time.Now()
	}

	if !stdConv.UpdatedAt.IsZero() {
		conv.UpdatedAt = stdConv.UpdatedAt
	} else {
		conv.UpdatedAt = time.Now()
	}

	return conv, nil
}

// transformMessage 转换消息
func (t *Transformer) transformMessage(stdMsg *types.StandardMessage, conversationID uuid.UUID) (*models.Message, error) {
	msg := &models.Message{
		ConversationID: conversationID,
		Role:           stdMsg.Role,
		Content:        stdMsg.Content,
	}

	// 设置时间
	if !stdMsg.CreatedAt.IsZero() {
		msg.CreatedAt = stdMsg.CreatedAt
	} else {
		msg.CreatedAt = time.Now()
	}

	msg.UpdatedAt = time.Now()

	return msg, nil
}
