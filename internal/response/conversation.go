package response

import (
	"chat-assistant-backend/internal/models"

	"github.com/google/uuid"
)

// ConversationResponse represents a conversation in API response
type ConversationResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Provider  string    `json:"provider"`
	Model     string    `json:"model"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

// ConversationListResponse represents a list of conversations in API response
type ConversationListResponse struct {
	Conversations []ConversationResponse `json:"conversations"`
}

// NewConversationResponse creates a ConversationResponse from models.Conversation
func NewConversationResponse(conversation *models.Conversation) *ConversationResponse {
	title := conversation.Title
	if title == "" {
		title = conversation.SourceTitle
	}

	return &ConversationResponse{
		ID:        conversation.Base.ID,
		UserID:    conversation.UserID,
		Title:     title,
		Provider:  conversation.Provider,
		Model:     conversation.Model,
		CreatedAt: conversation.Base.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: conversation.Base.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// NewConversationListResponse creates a ConversationListResponse from a slice of models.Conversation
func NewConversationListResponse(conversations []*models.Conversation) *ConversationListResponse {
	conversationResponses := make([]ConversationResponse, len(conversations))
	for i, conversation := range conversations {
		conversationResponses[i] = *NewConversationResponse(conversation)
	}

	return &ConversationListResponse{
		Conversations: conversationResponses,
	}
}
