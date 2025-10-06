package models

import "github.com/google/uuid"

// Message represents a message in a conversation
type Message struct {
	Base
	ConversationID uuid.UUID `gorm:"type:uuid;not null;index" json:"conversation_id"`
	Role           string    `gorm:"type:varchar(20);not null" json:"role"` // user, assistant, system
	Content        string    `gorm:"type:text;not null" json:"content"`
	SourceID       string    `gorm:"type:varchar(255);not null;index" json:"source_id"` // 原始数据中的ID，用于关联导入内容
	SourceContent  string    `gorm:"type:text;not null" json:"source_content"`          // 原始数据中的内容，用于对比和调试
	Metadata       string    `gorm:"type:text" json:"metadata"`                         // 可选元信息
}

// TableName returns the table name for the Message model
func (Message) TableName() string {
	return "messages"
}

// ToESDocument converts Message to MessageDocument for Elasticsearch
func (m *Message) ToESDocument() MessageDocument {
	return MessageDocument{
		ID:             m.ID,
		ConversationID: m.ConversationID,
		Role:           m.Role,
		Content:        m.Content,
		SourceID:       m.SourceID,
		SourceContent:  m.SourceContent,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}
