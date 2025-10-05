package models

import (
	"time"

	"github.com/google/uuid"
)

// ConversationDocument 是 ES 中的统一文档模型
// 包含 conversation 信息和嵌套的 messages
type ConversationDocument struct {
	// Conversation 字段
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Provider    string    `json:"provider"`
	Model       string    `json:"model"`
	SourceID    string    `json:"source_id"`
	SourceTitle string    `json:"source_title"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 嵌套的 Messages 和 Tags
	Messages []MessageDocument `json:"messages,omitempty"`
	Tags     []TagDocument     `json:"tags,omitempty"`
}

// MessageDocument 是 ES 中的消息文档
type MessageDocument struct {
	ID             uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	SourceID       string    `json:"source_id"`
	SourceContent  string    `json:"source_content"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TagDocument 是 ES 中的标签文档
type TagDocument struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 转换方法：从 ES 文档提取 Conversation 模型
func (d *ConversationDocument) ToConversation() *Conversation {
	return &Conversation{
		Base: Base{
			ID:        d.ID,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		},
		UserID:      d.UserID,
		Title:       d.Title,
		Provider:    d.Provider,
		Model:       d.Model,
		SourceID:    d.SourceID,
		SourceTitle: d.SourceTitle,
	}
}

// 转换方法：从 ES 文档提取 Messages 模型
func (d *ConversationDocument) ToMessages() []*Message {
	if len(d.Messages) == 0 {
		return nil
	}

	messages := make([]*Message, len(d.Messages))
	for i, msgDoc := range d.Messages {
		messages[i] = &Message{
			Base: Base{
				ID:        msgDoc.ID,
				CreatedAt: msgDoc.CreatedAt,
				UpdatedAt: msgDoc.UpdatedAt,
			},
			ConversationID: msgDoc.ConversationID,
			Role:           msgDoc.Role,
			Content:        msgDoc.Content,
			SourceID:       msgDoc.SourceID,
			SourceContent:  msgDoc.SourceContent,
		}
	}

	return messages
}
