package models

import "github.com/google/uuid"

// Conversation represents a chat conversation
type Conversation struct {
	Base
	UserID   uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Title    string    `gorm:"type:varchar(500);not null" json:"title"`
	Provider string    `gorm:"type:varchar(50);not null" json:"provider"`         // openai, gemini, local 等
	Model    string    `gorm:"type:varchar(50)" json:"model"`                     // gpt-4, gemini-pro, llama-3 等
	SourceID string    `gorm:"type:varchar(255);not null;index" json:"source_id"` // 原始数据中的ID，用于关联导入内容
}

// TableName returns the table name for the Conversation model
func (Conversation) TableName() string {
	return "conversations"
}
