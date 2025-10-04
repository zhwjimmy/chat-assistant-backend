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
