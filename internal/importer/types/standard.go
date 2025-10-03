package types

import "time"

// StandardFormat 标准化格式
type StandardFormat struct {
	Conversations []*StandardConversation `json:"conversations"`
}

// StandardConversation 标准化对话
type StandardConversation struct {
	ID        string                 `json:"id"`
	Title     string                 `json:"title"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	Provider  string                 `json:"provider"`
	Model     string                 `json:"model"`
	Messages  []*StandardMessage     `json:"messages"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// StandardMessage 标准化消息
type StandardMessage struct {
	ID        string                 `json:"id"`
	Role      string                 `json:"role"`
	Content   string                 `json:"content"`
	CreatedAt time.Time              `json:"created_at"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}
