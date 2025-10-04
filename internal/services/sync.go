package services

import (
	"fmt"

	"chat-assistant-backend/internal/models"
	"chat-assistant-backend/internal/repositories"
)

// SyncService 处理数据同步业务逻辑
type SyncService struct {
	conversationRepo *repositories.ConversationRepository
	indexer          repositories.ElasticsearchIndexer
}

// NewSyncService 创建同步服务
func NewSyncService(conversationRepo *repositories.ConversationRepository, indexer repositories.ElasticsearchIndexer) *SyncService {
	return &SyncService{
		conversationRepo: conversationRepo,
		indexer:          indexer,
	}
}

// SyncAll 同步所有数据到 Elasticsearch
func (s *SyncService) SyncAll() error {
	// 1. 从数据库读取所有 conversations 和 messages
	conversations, err := s.conversationRepo.FindAll()
	if err != nil {
		return fmt.Errorf("failed to get conversations: %w", err)
	}

	// 2. 转换为 ES 文档
	docs := s.convertToESDocuments(conversations)

	// 3. 批量索引到 ES
	if err := s.indexer.BulkIndexConversations(docs); err != nil {
		return fmt.Errorf("failed to bulk index conversations: %w", err)
	}

	return nil
}

// convertToESDocuments 转换 conversations 为 ES 文档
func (s *SyncService) convertToESDocuments(conversations []*models.Conversation) []*models.ConversationDocument {
	docs := make([]*models.ConversationDocument, len(conversations))
	for i, conv := range conversations {
		docs[i] = conv.ToESDocument()
	}
	return docs
}
