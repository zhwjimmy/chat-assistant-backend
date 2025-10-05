package services

import (
	"strings"
	"time"

	"chat-assistant-backend/internal/repositories"
	"chat-assistant-backend/internal/response"

	"github.com/google/uuid"
)

// SearchService defines the interface for search service
type SearchService interface {
	SearchWithMatchedMessages(query string, userID *uuid.UUID, providerID *string, tagID *uuid.UUID, startDate, endDate *time.Time, page, limit int) (*response.SearchResponse, int64, error)
}

// SearchServiceImpl handles search business logic
type SearchServiceImpl struct {
	searchRepo repositories.SearchRepository
}

// NewSearchService creates a new search service
func NewSearchService(searchRepo repositories.SearchRepository) SearchService {
	return &SearchServiceImpl{
		searchRepo: searchRepo,
	}
}

// SearchWithMatchedMessages performs a search and returns conversations with matched messages
func (s *SearchServiceImpl) SearchWithMatchedMessages(query string, userID *uuid.UUID, providerID *string, tagID *uuid.UUID, startDate, endDate *time.Time, page, limit int) (*response.SearchResponse, int64, error) {
	// Validate and clean query
	query = strings.TrimSpace(query)

	// Search conversations with matched messages and field information
	conversationDocs, matchedMessagesMap, matchedFieldsMap, total, err := s.searchRepo.SearchConversationsWithMatchedMessages(query, userID, providerID, tagID, startDate, endDate, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Convert to new search response format
	return response.NewSearchResponse(query, conversationDocs, matchedMessagesMap, matchedFieldsMap), total, nil
}
