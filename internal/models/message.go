package models

import "github.com/google/uuid"

// Message represents a message in a conversation
type Message struct {
	Base
	ConversationID uuid.UUID `gorm:"type:uuid;not null;index" json:"conversation_id"`
	Role           string    `gorm:"type:varchar(20);not null" json:"role"` // user, assistant, system
	Content        string    `gorm:"type:text;not null" json:"content"`
}

// TableName returns the table name for the Message model
func (Message) TableName() string {
	return "messages"
}
