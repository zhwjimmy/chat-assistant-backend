package services

import (
	"strings"

	"chat-assistant-backend/internal/models"
	"chat-assistant-backend/internal/response"

	"github.com/google/uuid"
)

// SearchRepository interface abstracts search functionality
type SearchRepository interface {
	SearchConversationsWithMessages(query string, userID *uuid.UUID, page, limit int) ([]*models.Conversation, int64, error)
	SearchConversationsWithMatchedMessages(query string, userID *uuid.UUID, page, limit int) ([]*models.ConversationDocument, map[uuid.UUID][]*models.MessageDocument, map[uuid.UUID][]string, int64, error)
}

// SearchService handles search business logic
type SearchService struct {
	searchRepo SearchRepository
}

// NewSearchService creates a new search service
func NewSearchService(searchRepo SearchRepository) *SearchService {
	return &SearchService{
		searchRepo: searchRepo,
	}
}

// Search performs a search across conversations and messages, returns conversation list
func (s *SearchService) Search(query string, userID *uuid.UUID, page, limit int) (*response.ConversationListResponse, int64, error) {
	// Validate and clean query
	query = strings.TrimSpace(query)
	if query == "" {
		return &response.ConversationListResponse{Conversations: []response.ConversationResponse{}}, 0, nil
	}

	// Search conversations that match either title or have messages with matching content
	conversations, total, err := s.searchRepo.SearchConversationsWithMessages(query, userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response format using existing ConversationResponse
	return response.NewConversationListResponse(conversations), total, nil
}

// SearchWithMatchedMessages performs a search and returns conversations with matched messages
func (s *SearchService) SearchWithMatchedMessages(query string, userID *uuid.UUID, page, limit int) (*response.SearchResponse, int64, error) {
	// Validate and clean query
	query = strings.TrimSpace(query)
	if query == "" {
		return &response.SearchResponse{Query: query, Conversations: []response.SearchConversationResponse{}}, 0, nil
	}

	// Search conversations with matched messages and field information
	conversationDocs, matchedMessagesMap, matchedFieldsMap, total, err := s.searchRepo.SearchConversationsWithMatchedMessages(query, userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Convert to new search response format
	return response.NewSearchResponse(query, conversationDocs, matchedMessagesMap, matchedFieldsMap), total, nil
}
