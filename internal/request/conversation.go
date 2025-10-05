package request

import "github.com/google/uuid"

// CreateConversationRequest represents a request to create a conversation
type CreateConversationRequest struct {
	UserID      uuid.UUID    `json:"user_id" binding:"required"`
	Title       string       `json:"title"`
	Provider    string       `json:"provider" binding:"required"`
	Model       string       `json:"model"`
	SourceID    string       `json:"source_id" binding:"required"`
	SourceTitle string       `json:"source_title" binding:"required"`
	Tags        []TagRequest `json:"tags,omitempty"`
}

// UpdateConversationRequest represents a request to update a conversation
type UpdateConversationRequest struct {
	Title       string       `json:"title"`
	Provider    string       `json:"provider"`
	Model       string       `json:"model"`
	SourceTitle string       `json:"source_title"`
	Tags        []TagRequest `json:"tags,omitempty"`
}
