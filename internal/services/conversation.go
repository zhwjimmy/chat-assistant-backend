package services

import (
	"chat-assistant-backend/internal/errors"
	"chat-assistant-backend/internal/logger"
	"chat-assistant-backend/internal/models"
	"chat-assistant-backend/internal/repositories"

	"github.com/google/uuid"
	"go.uber.org/zap"
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
	indexer          repositories.ElasticsearchIndexer
}

// NewConversationService creates a new conversation service
func NewConversationService(conversationRepo repositories.ConversationRepository, tagRepo repositories.TagRepository, indexer repositories.ElasticsearchIndexer) ConversationService {
	return &ConversationServiceImpl{
		conversationRepo: conversationRepo,
		tagRepo:          tagRepo,
		indexer:          indexer,
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

	// Delete the conversation from PostgreSQL
	if err := s.conversationRepo.Delete(id); err != nil {
		return err
	}

	// Delete the conversation from Elasticsearch
	if err := s.indexer.DeleteConversation(id); err != nil {
		// Log the error but don't fail the operation
		// ES is used for search, so we can tolerate temporary inconsistency
		logger.GetLogger().Error("Failed to delete conversation from Elasticsearch",
			zap.String("conversation_id", id.String()),
			zap.Error(err),
		)
	}

	return nil
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
	createdConversation, err := s.conversationRepo.GetByID(conversation.ID)
	if err != nil {
		return nil, err
	}

	// 索引到 Elasticsearch
	if err := s.indexer.IndexConversation(createdConversation.ToESDocument()); err != nil {
		// Log the error but don't fail the operation
		// ES is used for search, so we can tolerate temporary inconsistency
		logger.GetLogger().Error("Failed to index conversation to Elasticsearch",
			zap.String("conversation_id", conversation.ID.String()),
			zap.Error(err),
		)
	}

	return createdConversation, nil
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
	if err := s.conversationRepo.ReplaceTags(conversationID, tagIDs); err != nil {
		return err
	}

	// 重新获取对话以包含更新后的标签
	updatedConversation, err := s.conversationRepo.GetByID(conversationID)
	if err != nil {
		return err
	}

	// 更新 Elasticsearch 中的对话文档
	if err := s.indexer.UpdateConversation(updatedConversation.ToESDocument()); err != nil {
		// Log the error but don't fail the operation
		// ES is used for search, so we can tolerate temporary inconsistency
		logger.GetLogger().Error("Failed to update conversation in Elasticsearch",
			zap.String("conversation_id", conversationID.String()),
			zap.Error(err),
		)
	}

	return nil
}
