package repositories

import (
	"fmt"
	"strings"

	"chat-assistant-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SearchRepository handles search-related database operations
type SearchRepository struct {
	db *gorm.DB
}

// NewSearchRepository creates a new search repository
func NewSearchRepository(db *gorm.DB) *SearchRepository {
	return &SearchRepository{
		db: db,
	}
}

// SearchConversations searches conversations by title
func (r *SearchRepository) SearchConversations(query string, userID *uuid.UUID, page, limit int) ([]*models.Conversation, int64, error) {
	var conversations []*models.Conversation
	var total int64

	// Build the search query
	db := r.db.Model(&models.Conversation{})

	// Add user filter if provided
	if userID != nil {
		db = db.Where("user_id = ?", *userID)
	}

	// Add search condition - search in both title and source_title
	searchPattern := "%" + strings.ToLower(query) + "%"
	db = db.Where("LOWER(title) LIKE ? OR LOWER(source_title) LIKE ?", searchPattern, searchPattern)

	// Get total count
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and get results
	offset := (page - 1) * limit
	if err := db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&conversations).Error; err != nil {
		return nil, 0, err
	}

	return conversations, total, nil
}

// SearchMessages searches messages by content
func (r *SearchRepository) SearchMessages(query string, userID *uuid.UUID, page, limit int) ([]*models.Message, []*models.Conversation, int64, error) {
	var messages []*models.Message
	var conversations []*models.Conversation
	var total int64

	// Build the search query with JOIN to conversations table
	db := r.db.Model(&models.Message{}).
		Joins("JOIN conversations ON messages.conversation_id = conversations.id")

	// Add user filter if provided
	if userID != nil {
		db = db.Where("conversations.user_id = ?", *userID)
	}

	// Add search condition - search in both content and source_content
	searchPattern := "%" + strings.ToLower(query) + "%"
	db = db.Where("LOWER(messages.content) LIKE ? OR LOWER(messages.source_content) LIKE ?", searchPattern, searchPattern)

	// Get total count
	if err := db.Count(&total).Error; err != nil {
		return nil, nil, 0, err
	}

	// Apply pagination and get results
	offset := (page - 1) * limit
	if err := db.Order("messages.created_at DESC").Offset(offset).Limit(limit).Find(&messages).Error; err != nil {
		return nil, nil, 0, err
	}

	// Get conversation details for the found messages
	if len(messages) > 0 {
		conversationIDs := make([]uuid.UUID, len(messages))
		for i, msg := range messages {
			conversationIDs[i] = msg.ConversationID
		}

		if err := r.db.Where("id IN ?", conversationIDs).Find(&conversations).Error; err != nil {
			return nil, nil, 0, err
		}
	}

	return messages, conversations, total, nil
}

// SearchConversationsWithMessages searches conversations that match either title or have messages with matching content
func (r *SearchRepository) SearchConversationsWithMessages(query string, userID *uuid.UUID, page, limit int) ([]*models.Conversation, int64, error) {
	var conversations []*models.Conversation
	var total int64

	// Build the search query using UNION to combine:
	// 1. Conversations with matching titles
	// 2. Conversations that have messages with matching content
	searchPattern := "%" + strings.ToLower(query) + "%"

	// Use a subquery to find conversations that match either condition
	subQuery := r.db.Model(&models.Conversation{}).
		Select("DISTINCT conversations.id").
		Joins("LEFT JOIN messages ON conversations.id = messages.conversation_id").
		Where("LOWER(conversations.title) LIKE ? OR LOWER(conversations.source_title) LIKE ? OR LOWER(messages.content) LIKE ? OR LOWER(messages.source_content) LIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern)

	// Add user filter if provided
	if userID != nil {
		subQuery = subQuery.Where("conversations.user_id = ?", *userID)
	}

	// Main query to get full conversation objects
	db := r.db.Model(&models.Conversation{}).
		Where("conversations.id IN (?)", subQuery)

	// Get total count
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and get results
	offset := (page - 1) * limit
	if err := db.Order("conversations.created_at DESC").Offset(offset).Limit(limit).Find(&conversations).Error; err != nil {
		return nil, 0, err
	}

	return conversations, total, nil
}

// HighlightText highlights search terms in text
func (r *SearchRepository) HighlightText(text, query string) string {
	if query == "" {
		return text
	}

	// Simple highlighting - wrap matching terms with <mark> tags
	queryLower := strings.ToLower(query)
	textLower := strings.ToLower(text)

	// Find the first occurrence
	index := strings.Index(textLower, queryLower)
	if index == -1 {
		return text
	}

	// Get the original case version of the matched text
	matchedText := text[index : index+len(query)]
	highlighted := fmt.Sprintf("<mark>%s</mark>", matchedText)

	// Replace the first occurrence
	result := text[:index] + highlighted + text[index+len(query):]

	return result
}
