package services

import (
	"chat-assistant-backend/internal/errors"
	"chat-assistant-backend/internal/models"
	"chat-assistant-backend/internal/repositories"

	"github.com/google/uuid"
)

// MessageService defines the interface for message service
type MessageService interface {
	GetMessageByID(id uuid.UUID) (*models.Message, error)
	GetMessagesByConversationID(conversationID uuid.UUID, page, limit int) ([]*models.Message, int64, error)
	GetAllMessages(page, limit int) ([]*models.Message, int64, error)
	DeleteMessage(id uuid.UUID) error
}

// MessageServiceImpl handles message business logic
type MessageServiceImpl struct {
	messageRepo repositories.MessageRepository
}

// NewMessageService creates a new message service
func NewMessageService(messageRepo repositories.MessageRepository) MessageService {
	return &MessageServiceImpl{
		messageRepo: messageRepo,
	}
}

// GetMessageByID retrieves a message by ID
func (s *MessageServiceImpl) GetMessageByID(id uuid.UUID) (*models.Message, error) {
	message, err := s.messageRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if message == nil {
		return nil, errors.ErrMessageNotFound
	}

	return message, nil
}

// GetMessagesByConversationID retrieves messages by conversation ID with pagination
func (s *MessageServiceImpl) GetMessagesByConversationID(conversationID uuid.UUID, page, limit int) ([]*models.Message, int64, error) {
	messages, total, err := s.messageRepo.GetByConversationID(conversationID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// GetAllMessages retrieves all messages with pagination
func (s *MessageServiceImpl) GetAllMessages(page, limit int) ([]*models.Message, int64, error) {
	messages, total, err := s.messageRepo.GetAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// DeleteMessage deletes a message by ID
func (s *MessageServiceImpl) DeleteMessage(id uuid.UUID) error {
	// First check if message exists
	message, err := s.messageRepo.GetByID(id)
	if err != nil {
		return err
	}

	if message == nil {
		return errors.ErrMessageNotFound
	}

	// Delete the message
	return s.messageRepo.Delete(id)
}
