package services

import (
	"chat-assistant-backend/internal/errors"
	"chat-assistant-backend/internal/models"
	"chat-assistant-backend/internal/repositories"

	"github.com/google/uuid"
)

// ConversationService defines the interface for conversation service
type ConversationService interface {
	GetConversationByID(id uuid.UUID) (*models.Conversation, error)
	GetConversationsByUserID(userID uuid.UUID, page, limit int) ([]*models.Conversation, int64, error)
	DeleteConversation(id uuid.UUID) error
	CreateConversationWithTags(conversation *models.Conversation, tagNames []string) (*models.Conversation, error)
	UpdateConversationTags(conversationID uuid.UUID, tagNames []string) error
}

// ConversationServiceImpl handles conversation business logic
type ConversationServiceImpl struct {
	conversationRepo repositories.ConversationRepository
	tagRepo          repositories.TagRepository
}

// NewConversationService creates a new conversation service
func NewConversationService(conversationRepo repositories.ConversationRepository, tagRepo repositories.TagRepository) ConversationService {
	return &ConversationServiceImpl{
		conversationRepo: conversationRepo,
		tagRepo:          tagRepo,
	}
}

// GetConversationByID retrieves a conversation by ID
func (s *ConversationServiceImpl) GetConversationByID(id uuid.UUID) (*models.Conversation, error) {
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
func (s *ConversationServiceImpl) GetConversationsByUserID(userID uuid.UUID, page, limit int) ([]*models.Conversation, int64, error) {
	conversations, total, err := s.conversationRepo.GetByUserID(userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return conversations, total, nil
}

// DeleteConversation deletes a conversation by ID
func (s *ConversationServiceImpl) DeleteConversation(id uuid.UUID) error {
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

// CreateConversationWithTags creates a new conversation with tags
func (s *ConversationServiceImpl) CreateConversationWithTags(conversation *models.Conversation, tagNames []string) (*models.Conversation, error) {
	// 创建对话
	err := s.conversationRepo.Create(conversation)
	if err != nil {
		return nil, err
	}

	// 处理标签
	if len(tagNames) > 0 {
		tags, err := s.tagRepo.CreateOrGetTags(tagNames)
		if err != nil {
			return nil, err
		}

		// 建立标签关系
		tagIDs := make([]string, len(tags))
		for i, tag := range tags {
			tagIDs[i] = tag.ID.String()
		}

		err = s.conversationRepo.ReplaceTags(conversation.ID, tagIDs)
		if err != nil {
			return nil, err
		}
	}

	// 重新获取对话以包含标签
	return s.conversationRepo.GetByID(conversation.ID)
}

// UpdateConversationTags updates tags for a conversation
func (s *ConversationServiceImpl) UpdateConversationTags(conversationID uuid.UUID, tagNames []string) error {
	// 检查对话是否存在
	conversation, err := s.conversationRepo.GetByID(conversationID)
	if err != nil {
		return err
	}

	if conversation == nil {
		return errors.ErrConversationNotFound
	}

	// 处理标签
	var tagIDs []string
	if len(tagNames) > 0 {
		tags, err := s.tagRepo.CreateOrGetTags(tagNames)
		if err != nil {
			return err
		}

		tagIDs = make([]string, len(tags))
		for i, tag := range tags {
			tagIDs[i] = tag.ID.String()
		}
	}

	// 更新标签关系
	return s.conversationRepo.ReplaceTags(conversationID, tagIDs)
}
