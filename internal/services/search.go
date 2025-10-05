package services

import (
	"strings"
	"time"

	"chat-assistant-backend/internal/models"
	"chat-assistant-backend/internal/response"

	"github.com/google/uuid"
)

// SearchRepository interface abstracts search functionality
type SearchRepository interface {
	SearchConversationsWithMatchedMessages(query string, userID *uuid.UUID, providerID *string, startDate, endDate *time.Time, page, limit int) ([]*models.ConversationDocument, map[uuid.UUID][]*models.MessageDocument, map[uuid.UUID][]string, int64, error)
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

// SearchWithMatchedMessages performs a search and returns conversations with matched messages
func (s *SearchService) SearchWithMatchedMessages(query string, userID *uuid.UUID, providerID *string, startDate, endDate *time.Time, page, limit int) (*response.SearchResponse, int64, error) {
	// Validate and clean query
	query = strings.TrimSpace(query)
	if query == "" {
		return &response.SearchResponse{Query: query, Conversations: []response.SearchConversationResponse{}}, 0, nil
	}

	// Search conversations with matched messages and field information
	conversationDocs, matchedMessagesMap, matchedFieldsMap, total, err := s.searchRepo.SearchConversationsWithMatchedMessages(query, userID, providerID, startDate, endDate, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Convert to new search response format
	return response.NewSearchResponse(query, conversationDocs, matchedMessagesMap, matchedFieldsMap), total, nil
}
