package models

import "github.com/google/uuid"

// Conversation represents a chat conversation
type Conversation struct {
	Base
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Title       string    `gorm:"type:varchar(500)" json:"title"`
	Provider    string    `gorm:"type:varchar(50);not null" json:"provider"`         // openai, gemini, local 等
	Model       string    `gorm:"type:varchar(50)" json:"model"`                     // gpt-4, gemini-pro, llama-3 等
	SourceID    string    `gorm:"type:varchar(255);not null;index" json:"source_id"` // 原始数据中的ID，用于关联导入内容
	SourceTitle string    `gorm:"type:varchar(500);not null" json:"source_title"`
	Metadata    string    `gorm:"type:text" json:"metadata"` // 可选元信息
	Messages    []Message `gorm:"foreignKey:ConversationID" json:"messages,omitempty"`
	Tags        []Tag     `gorm:"many2many:conversation_tags;" json:"tags,omitempty"`
}

// TableName returns the table name for the Conversation model
func (Conversation) TableName() string {
	return "conversations"
}

// ToESDocument converts Conversation to ConversationDocument for Elasticsearch
func (c *Conversation) ToESDocument() *ConversationDocument {
	doc := &ConversationDocument{
		ID:          c.ID,
		UserID:      c.UserID,
		Title:       c.Title,
		Provider:    c.Provider,
		Model:       c.Model,
		SourceID:    c.SourceID,
		SourceTitle: c.SourceTitle,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		Messages:    []MessageDocument{},
		Tags:        []TagDocument{},
	}

	// 如果有预加载的 Messages，转换它们
	if c.Messages != nil {
		doc.Messages = make([]MessageDocument, len(c.Messages))
		for i, msg := range c.Messages {
			doc.Messages[i] = msg.ToESDocument()
		}
	}

	// 如果有预加载的 Tags，转换它们
	if c.Tags != nil {
		doc.Tags = make([]TagDocument, len(c.Tags))
		for i, tag := range c.Tags {
			doc.Tags[i] = tag.ToESDocument()
		}
	}

	return doc
}
