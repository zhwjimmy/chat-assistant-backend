package response

import (
	"chat-assistant-backend/internal/models"

	"github.com/google/uuid"
)

// MessageResponse represents a message in API response
type MessageResponse struct {
	ID             uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
}

// MessageListResponse represents a list of messages in API response
type MessageListResponse struct {
	Messages []MessageResponse `json:"messages"`
}

// NewMessageResponse creates a MessageResponse from models.Message
func NewMessageResponse(message *models.Message) *MessageResponse {
	return &MessageResponse{
		ID:             message.Base.ID,
		ConversationID: message.ConversationID,
		Role:           message.Role,
		Content:        message.Content,
		CreatedAt:      message.Base.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      message.Base.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// NewMessageListResponse creates a MessageListResponse from a slice of models.Message
func NewMessageListResponse(messages []*models.Message) *MessageListResponse {
	messageResponses := make([]MessageResponse, len(messages))
	for i, message := range messages {
		messageResponses[i] = *NewMessageResponse(message)
	}

	return &MessageListResponse{
		Messages: messageResponses,
	}
}
