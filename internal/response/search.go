package response

import (
	"chat-assistant-backend/internal/models"

	"github.com/google/uuid"
)

// SearchMessageResponse represents a message in search results with highlighting
type SearchMessageResponse struct {
	ID             uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	SourceID       string    `json:"source_id,omitempty"`
	SourceContent  string    `json:"source_content,omitempty"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
	// 匹配信息，用于前端高亮
	MatchedFields []string `json:"matched_fields,omitempty"` // 匹配的字段名
}

// SearchTagResponse represents a tag in search results with highlighting
type SearchTagResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	// 匹配信息，用于前端高亮
	MatchedFields []string `json:"matched_fields,omitempty"` // 匹配的字段名
}

// SearchConversationResponse represents a conversation in search results with matched messages
type SearchConversationResponse struct {
	ID          uuid.UUID           `json:"id"`
	UserID      uuid.UUID           `json:"user_id"`
	Title       string              `json:"title"`
	Provider    string              `json:"provider"`
	Model       string              `json:"model"`
	SourceID    string              `json:"source_id,omitempty"`
	SourceTitle string              `json:"source_title,omitempty"`
	Tags        []SearchTagResponse `json:"tags"`
	CreatedAt   string              `json:"created_at"`
	UpdatedAt   string              `json:"updated_at"`
	// 匹配的消息列表（如果 conversation 匹配但消息不匹配，则为空；最多返回3条消息）
	Messages []SearchMessageResponse `json:"messages"`
	// 匹配信息，用于前端高亮
	MatchedFields []string `json:"matched_fields,omitempty"` // 匹配的字段名，如 ["title", "messages.content"]
}

// SearchResponse represents the search results
type SearchResponse struct {
	Query         string                       `json:"query"` // 搜索关键词，用于前端高亮
	Conversations []SearchConversationResponse `json:"conversations"`
}

// NewSearchMessageResponse creates a SearchMessageResponse from models.MessageDocument
func NewSearchMessageResponse(messageDoc *models.MessageDocument, matchedFields []string) *SearchMessageResponse {
	content := messageDoc.Content
	if content == "" {
		content = messageDoc.SourceContent
	}

	return &SearchMessageResponse{
		ID:             messageDoc.ID,
		ConversationID: messageDoc.ConversationID,
		Role:           messageDoc.Role,
		Content:        content,
		SourceID:       messageDoc.SourceID,
		SourceContent:  messageDoc.SourceContent,
		CreatedAt:      messageDoc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      messageDoc.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		MatchedFields:  matchedFields,
	}
}

// NewSearchTagResponse creates a SearchTagResponse from models.TagDocument
func NewSearchTagResponse(tagDoc *models.TagDocument, matchedFields []string) *SearchTagResponse {
	return &SearchTagResponse{
		ID:            tagDoc.ID,
		Name:          tagDoc.Name,
		CreatedAt:     tagDoc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     tagDoc.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		MatchedFields: matchedFields,
	}
}

// NewSearchConversationResponse creates a SearchConversationResponse from models.ConversationDocument
func NewSearchConversationResponse(conversationDoc *models.ConversationDocument, matchedMessages []*models.MessageDocument, matchedFields []string) *SearchConversationResponse {
	title := conversationDoc.Title
	if title == "" {
		title = conversationDoc.SourceTitle
	}

	// 转换匹配的消息
	messageResponses := make([]SearchMessageResponse, len(matchedMessages))
	for i, msgDoc := range matchedMessages {
		// 为消息添加匹配字段信息
		messageMatchedFields := []string{}
		for _, field := range matchedFields {
			if field == "messages.content" || field == "messages.source_content" {
				messageMatchedFields = append(messageMatchedFields, "content")
			}
		}
		messageResponses[i] = *NewSearchMessageResponse(msgDoc, messageMatchedFields)
	}

	// 转换 Tags
	var tags []SearchTagResponse
	if conversationDoc.Tags != nil {
		tags = make([]SearchTagResponse, len(conversationDoc.Tags))
		for i, tagDoc := range conversationDoc.Tags {
			// 为标签添加匹配字段信息
			tagMatchedFields := []string{}
			for _, field := range matchedFields {
				if field == "tags.name" {
					tagMatchedFields = append(tagMatchedFields, "name")
				}
			}
			tags[i] = *NewSearchTagResponse(&tagDoc, tagMatchedFields)
		}
	}

	return &SearchConversationResponse{
		ID:            conversationDoc.ID,
		UserID:        conversationDoc.UserID,
		Title:         title,
		Provider:      conversationDoc.Provider,
		Model:         conversationDoc.Model,
		SourceID:      conversationDoc.SourceID,
		SourceTitle:   conversationDoc.SourceTitle,
		Tags:          tags,
		CreatedAt:     conversationDoc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     conversationDoc.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Messages:      messageResponses,
		MatchedFields: matchedFields,
	}
}

// NewSearchResponse creates a SearchResponse from a slice of conversation documents
func NewSearchResponse(query string, conversationDocs []*models.ConversationDocument, matchedMessagesMap map[uuid.UUID][]*models.MessageDocument, matchedFieldsMap map[uuid.UUID][]string) *SearchResponse {
	conversationResponses := make([]SearchConversationResponse, len(conversationDocs))

	for i, conversationDoc := range conversationDocs {
		var matchedMessages []*models.MessageDocument
		var matchedFields []string

		if messages, exists := matchedMessagesMap[conversationDoc.ID]; exists {
			matchedMessages = messages
		}

		if fields, exists := matchedFieldsMap[conversationDoc.ID]; exists {
			matchedFields = fields
		}

		conversationResponses[i] = *NewSearchConversationResponse(conversationDoc, matchedMessages, matchedFields)
	}

	return &SearchResponse{
		Query:         query,
		Conversations: conversationResponses,
	}
}
