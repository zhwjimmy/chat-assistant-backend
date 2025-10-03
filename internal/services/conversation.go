package services

import (
	"chat-assistant-backend/internal/errors"
	"chat-assistant-backend/internal/models"
	"chat-assistant-backend/internal/repositories"

	"github.com/google/uuid"
)

// ConversationService handles conversation business logic
type ConversationService struct {
	conversationRepo *repositories.ConversationRepository
}

// NewConversationService creates a new conversation service
func NewConversationService(conversationRepo *repositories.ConversationRepository) *ConversationService {
	return &ConversationService{
		conversationRepo: conversationRepo,
	}
}

// GetConversationByID retrieves a conversation by ID
func (s *ConversationService) GetConversationByID(id uuid.UUID) (*models.Conversation, error) {
	conversation, err := s.conversationRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if conversation == nil {
		return nil, errors.ErrConversationNotFound
	}

	return conversation, nil
}

// GetConversationsByUserID retrieves conversations by user ID with pagination
func (s *ConversationService) GetConversationsByUserID(userID uuid.UUID, page, limit int) ([]*models.Conversation, int64, error) {
	conversations, total, err := s.conversationRepo.GetByUserID(userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return conversations, total, nil
}

// DeleteConversation deletes a conversation by ID
func (s *ConversationService) DeleteConversation(id uuid.UUID) error {
	// First check if conversation exists
	conversation, err := s.conversationRepo.GetByID(id)
	if err != nil {
		return err
	}

	if conversation == nil {
		return errors.ErrConversationNotFound
	}

	// Delete the conversation
	return s.conversationRepo.Delete(id)
}
