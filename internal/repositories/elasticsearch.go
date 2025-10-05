package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"chat-assistant-backend/internal/models"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
)

// SearchRepository defines the interface for search repository
type SearchRepository interface {
	SearchConversationsWithMatchedMessages(query string, userID *uuid.UUID, providerID *string, startDate, endDate *time.Time, page, limit int) ([]*models.ConversationDocument, map[uuid.UUID][]*models.MessageDocument, map[uuid.UUID][]string, int64, error)
}

// ElasticsearchRepositoryImpl handles Elasticsearch search operations
type ElasticsearchRepositoryImpl struct {
	esClient  *es.Client
	indexName string
}

// NewElasticsearchRepository creates a new Elasticsearch repository
func NewElasticsearchRepository(esClient *es.Client, indexName string) SearchRepository {
	return &ElasticsearchRepositoryImpl{
		esClient:  esClient,
		indexName: indexName,
	}
}

// SearchConversationsWithMatchedMessages searches conversations and returns matched messages
func (r *ElasticsearchRepositoryImpl) SearchConversationsWithMatchedMessages(query string, userID *uuid.UUID, providerID *string, startDate, endDate *time.Time, page, limit int) ([]*models.ConversationDocument, map[uuid.UUID][]*models.MessageDocument, map[uuid.UUID][]string, int64, error) {
	// 1. 在 ES 中搜索
	esDocs, highlights, total, err := r.searchConversationDocumentsWithHighlights(query, userID, providerID, startDate, endDate, page, limit)
	if err != nil {
		return nil, nil, nil, 0, err
	}

	// 2. 提取匹配的消息和字段信息
	matchedMessagesMap := make(map[uuid.UUID][]*models.MessageDocument)
	matchedFieldsMap := make(map[uuid.UUID][]string)

	for i, doc := range esDocs {
		conversationID := doc.ID
		var matchedFields []string

		// 检查哪些字段有高亮（即匹配）
		if _, exists := highlights[i]["title"]; exists {
			matchedFields = append(matchedFields, "title")
		}
		if _, exists := highlights[i]["source_title"]; exists {
			matchedFields = append(matchedFields, "source_title")
		}
		if _, exists := highlights[i]["messages.content"]; exists {
			matchedFields = append(matchedFields, "messages.content")
		}
		if _, exists := highlights[i]["messages.source_content"]; exists {
			matchedFields = append(matchedFields, "messages.source_content")
		}
		if _, exists := highlights[i]["tags.name"]; exists {
			matchedFields = append(matchedFields, "tags.name")
		}

		// 提取匹配的消息
		_, hasContent := highlights[i]["messages.content"]
		_, hasSourceContent := highlights[i]["messages.source_content"]
		if hasContent || hasSourceContent {
			// 如果 ES 返回了消息字段的高亮，说明有消息匹配
			// 最多返回 3 条消息，优先选择包含匹配关键词的消息
			const maxMessages = 3
			matchedMessages := make([]*models.MessageDocument, 0, maxMessages)

			// 首先尝试找到真正包含匹配关键词的消息
			for _, msgDoc := range doc.Messages {
				if len(matchedMessages) >= maxMessages {
					break
				}

				// 检查消息是否包含匹配的关键词（通过高亮信息判断）
				hasMatch := false
				if contentHighlights, exists := highlights[i]["messages.content"]; exists {
					if highlights, ok := contentHighlights.([]interface{}); ok {
						for _, highlight := range highlights {
							if highlightStr, ok := highlight.(string); ok {
								cleanHighlight := removeHighlightTags(highlightStr)
								if contains(msgDoc.Content, cleanHighlight) || contains(msgDoc.SourceContent, cleanHighlight) {
									hasMatch = true
									break
								}
							}
						}
					}
				}

				if hasMatch {
					matchedMessages = append(matchedMessages, &msgDoc)
				}
			}

			// 如果没有找到匹配的消息，或者匹配的消息少于3条，则补充前几条消息
			if len(matchedMessages) < maxMessages {
				for _, msgDoc := range doc.Messages {
					if len(matchedMessages) >= maxMessages {
						break
					}

					// 检查是否已经包含这条消息
					alreadyIncluded := false
					for _, existing := range matchedMessages {
						if existing.ID == msgDoc.ID {
							alreadyIncluded = true
							break
						}
					}

					if !alreadyIncluded {
						matchedMessages = append(matchedMessages, &msgDoc)
					}
				}
			}

			matchedMessagesMap[conversationID] = matchedMessages
			fmt.Printf("DEBUG: Added %d messages (max 3) for conversation %s\n", len(matchedMessages), conversationID)
		}

		// 总是设置 matched_fields，即使为空
		matchedFieldsMap[conversationID] = matchedFields
		fmt.Printf("DEBUG: Matched fields for conversation %s: %v\n", conversationID, matchedFields)
	}

	return esDocs, matchedMessagesMap, matchedFieldsMap, total, nil
}

// buildSearchQuery 构建 ES 搜索查询
func (r *ElasticsearchRepositoryImpl) buildSearchQuery(query string, userID *uuid.UUID, providerID *string, startDate, endDate *time.Time, page, limit int) []byte {
	// 计算偏移量
	offset := (page - 1) * limit

	// 构建查询条件
	var mustQueries []map[string]interface{}

	// 用户过滤
	if userID != nil {
		mustQueries = append(mustQueries, map[string]interface{}{
			"term": map[string]interface{}{
				"user_id": userID.String(),
			},
		})
	}

	// Provider过滤
	if providerID != nil {
		mustQueries = append(mustQueries, map[string]interface{}{
			"term": map[string]interface{}{
				"provider": *providerID,
			},
		})
	}

	// 日期范围过滤
	if startDate != nil || endDate != nil {
		dateRange := map[string]interface{}{}
		if startDate != nil {
			dateRange["gte"] = startDate.Format("2006-01-02T15:04:05Z07:00")
		}
		if endDate != nil {
			dateRange["lte"] = endDate.Format("2006-01-02T15:04:05Z07:00")
		}
		mustQueries = append(mustQueries, map[string]interface{}{
			"range": map[string]interface{}{
				"created_at": dateRange,
			},
		})
	}

	// 搜索查询 - 使用精确匹配
	searchQueries := []map[string]interface{}{
		// 1. 完全精确匹配 - 最高优先级
		{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"title.exact^5", "source_title.exact^4"},
				"type":   "phrase",
				"slop":   0,
			},
		},
		{
			"nested": map[string]interface{}{
				"path": "messages",
				"query": map[string]interface{}{
					"multi_match": map[string]interface{}{
						"query":  query,
						"fields": []string{"messages.content.exact^5", "messages.source_content.exact^4"},
						"type":   "phrase",
						"slop":   0,
					},
				},
			},
		},
		{
			"nested": map[string]interface{}{
				"path": "tags",
				"query": map[string]interface{}{
					"multi_match": map[string]interface{}{
						"query":  query,
						"fields": []string{"tags.name.exact^3"},
						"type":   "phrase",
						"slop":   0,
					},
				},
			},
		},
		// 添加词级别的精确匹配作为备选
		{
			"multi_match": map[string]interface{}{
				"query":    query,
				"fields":   []string{"title^2", "source_title"},
				"type":     "cross_fields",
				"operator": "and", // 所有词都必须匹配
			},
		},
		{
			"nested": map[string]interface{}{
				"path": "messages",
				"query": map[string]interface{}{
					"multi_match": map[string]interface{}{
						"query":    query,
						"fields":   []string{"messages.content^2", "messages.source_content"},
						"type":     "cross_fields",
						"operator": "and", // 所有词都必须匹配
					},
				},
			},
		},
		{
			"nested": map[string]interface{}{
				"path": "tags",
				"query": map[string]interface{}{
					"multi_match": map[string]interface{}{
						"query":    query,
						"fields":   []string{"tags.name^2"},
						"type":     "cross_fields",
						"operator": "and", // 所有词都必须匹配
					},
				},
			},
		},
	}

	// 构建完整的查询
	searchBody := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":                 mustQueries,
				"should":               searchQueries,
				"minimum_should_match": 1,
			},
		},
		"from": offset,
		"size": limit,
		"sort": []map[string]interface{}{
			{
				"_score": map[string]interface{}{
					"order": "desc",
				},
			},
			{
				"created_at": map[string]interface{}{
					"order": "desc",
				},
			},
		},
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"title":                   map[string]interface{}{},
				"source_title":            map[string]interface{}{},
				"messages.content":        map[string]interface{}{},
				"messages.source_content": map[string]interface{}{},
				"tags.name":               map[string]interface{}{},
			},
			"pre_tags":  []string{"<mark>"},
			"post_tags": []string{"</mark>"},
		},
	}

	// 序列化查询
	queryBytes, _ := json.Marshal(searchBody)
	return queryBytes
}

// parseSearchResponse 解析 ES 搜索响应
func (r *ElasticsearchRepositoryImpl) parseSearchResponse(response map[string]interface{}) ([]*models.ConversationDocument, int64, error) {
	// 提取总数
	hits, ok := response["hits"].(map[string]interface{})
	if !ok {
		return nil, 0, fmt.Errorf("invalid search response format")
	}

	totalValue, ok := hits["total"]
	if !ok {
		return nil, 0, fmt.Errorf("missing total in search response")
	}

	var total int64
	if totalMap, ok := totalValue.(map[string]interface{}); ok {
		if value, ok := totalMap["value"].(float64); ok {
			total = int64(value)
		}
	} else if value, ok := totalValue.(float64); ok {
		total = int64(value)
	}

	// 提取文档
	hitsList, ok := hits["hits"].([]interface{})
	if !ok {
		return nil, 0, fmt.Errorf("invalid hits format in search response")
	}

	documents := make([]*models.ConversationDocument, 0, len(hitsList))
	for _, hit := range hitsList {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}

		source, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			continue
		}

		// 解析文档
		doc := &models.ConversationDocument{}
		if err := r.parseDocument(source, doc); err != nil {
			continue // 跳过解析失败的文档
		}

		documents = append(documents, doc)
	}

	return documents, total, nil
}

// searchConversationDocumentsWithHighlights 在 ES 中搜索 conversation 文档并返回高亮信息
func (r *ElasticsearchRepositoryImpl) searchConversationDocumentsWithHighlights(query string, userID *uuid.UUID, providerID *string, startDate, endDate *time.Time, page, limit int) ([]*models.ConversationDocument, []map[string]interface{}, int64, error) {
	ctx := context.Background()

	// 构建 ES 查询
	searchQuery := r.buildSearchQuery(query, userID, providerID, startDate, endDate, page, limit)

	// 执行搜索
	req := esapi.SearchRequest{
		Index: []string{r.indexName},
		Body:  bytes.NewReader(searchQuery),
	}

	res, err := req.Do(ctx, r.esClient)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to execute search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, nil, 0, fmt.Errorf("search request failed with status: %s", res.Status())
	}

	// 解析响应
	var searchResponse map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&searchResponse); err != nil {
		return nil, nil, 0, fmt.Errorf("failed to decode search response: %w", err)
	}

	// 提取结果和高亮信息
	return r.parseSearchResponseWithHighlights(searchResponse)
}

// parseDocument 解析单个文档
func (r *ElasticsearchRepositoryImpl) parseDocument(source map[string]interface{}, doc *models.ConversationDocument) error {
	// 解析基础字段
	if id, ok := source["id"].(string); ok {
		if parsed, err := uuid.Parse(id); err == nil {
			doc.ID = parsed
		}
	}

	if userID, ok := source["user_id"].(string); ok {
		if parsed, err := uuid.Parse(userID); err == nil {
			doc.UserID = parsed
		}
	}

	if title, ok := source["title"].(string); ok {
		doc.Title = title
	}

	if provider, ok := source["provider"].(string); ok {
		doc.Provider = provider
	}

	if model, ok := source["model"].(string); ok {
		doc.Model = model
	}

	if sourceID, ok := source["source_id"].(string); ok {
		doc.SourceID = sourceID
	}

	if sourceTitle, ok := source["source_title"].(string); ok {
		doc.SourceTitle = sourceTitle
	}

	// 解析时间字段
	if createdAt, ok := source["created_at"].(string); ok {
		if parsed, err := json.Marshal(createdAt); err == nil {
			json.Unmarshal(parsed, &doc.CreatedAt)
		}
	}

	if updatedAt, ok := source["updated_at"].(string); ok {
		if parsed, err := json.Marshal(updatedAt); err == nil {
			json.Unmarshal(parsed, &doc.UpdatedAt)
		}
	}

	// 解析嵌套的 messages
	if messages, ok := source["messages"].([]interface{}); ok {
		doc.Messages = make([]models.MessageDocument, 0, len(messages))
		for _, msg := range messages {
			if msgMap, ok := msg.(map[string]interface{}); ok {
				messageDoc := models.MessageDocument{}
				if err := r.parseMessageDocument(msgMap, &messageDoc); err == nil {
					doc.Messages = append(doc.Messages, messageDoc)
				}
			}
		}
	}

	// 解析嵌套的 tags
	if tags, ok := source["tags"].([]interface{}); ok {
		doc.Tags = make([]models.TagDocument, 0, len(tags))
		for _, tag := range tags {
			if tagMap, ok := tag.(map[string]interface{}); ok {
				tagDoc := models.TagDocument{}
				if err := r.parseTagDocument(tagMap, &tagDoc); err == nil {
					doc.Tags = append(doc.Tags, tagDoc)
				}
			}
		}
	}

	return nil
}

// parseMessageDocument 解析消息文档
func (r *ElasticsearchRepositoryImpl) parseMessageDocument(source map[string]interface{}, doc *models.MessageDocument) error {
	// 解析消息字段
	if id, ok := source["id"].(string); ok {
		if parsed, err := uuid.Parse(id); err == nil {
			doc.ID = parsed
		}
	}

	if conversationID, ok := source["conversation_id"].(string); ok {
		if parsed, err := uuid.Parse(conversationID); err == nil {
			doc.ConversationID = parsed
		}
	}

	if role, ok := source["role"].(string); ok {
		doc.Role = role
	}

	if content, ok := source["content"].(string); ok {
		doc.Content = content
	}

	if sourceID, ok := source["source_id"].(string); ok {
		doc.SourceID = sourceID
	}

	if sourceContent, ok := source["source_content"].(string); ok {
		doc.SourceContent = sourceContent
	}

	// 解析时间字段
	if createdAt, ok := source["created_at"].(string); ok {
		if parsed, err := json.Marshal(createdAt); err == nil {
			json.Unmarshal(parsed, &doc.CreatedAt)
		}
	}

	if updatedAt, ok := source["updated_at"].(string); ok {
		if parsed, err := json.Marshal(updatedAt); err == nil {
			json.Unmarshal(parsed, &doc.UpdatedAt)
		}
	}

	return nil
}

// parseTagDocument 解析标签文档
func (r *ElasticsearchRepositoryImpl) parseTagDocument(source map[string]interface{}, doc *models.TagDocument) error {
	// 解析标签字段
	if id, ok := source["id"].(string); ok {
		if parsed, err := uuid.Parse(id); err == nil {
			doc.ID = parsed
		}
	}

	if name, ok := source["name"].(string); ok {
		doc.Name = name
	}

	// 解析时间字段
	if createdAt, ok := source["created_at"].(string); ok {
		if parsed, err := json.Marshal(createdAt); err == nil {
			json.Unmarshal(parsed, &doc.CreatedAt)
		}
	}

	if updatedAt, ok := source["updated_at"].(string); ok {
		if parsed, err := json.Marshal(updatedAt); err == nil {
			json.Unmarshal(parsed, &doc.UpdatedAt)
		}
	}

	return nil
}

// parseSearchResponseWithHighlights 解析 ES 搜索响应并提取高亮信息
func (r *ElasticsearchRepositoryImpl) parseSearchResponseWithHighlights(response map[string]interface{}) ([]*models.ConversationDocument, []map[string]interface{}, int64, error) {
	// 提取总数
	hits, ok := response["hits"].(map[string]interface{})
	if !ok {
		return nil, nil, 0, fmt.Errorf("invalid search response format")
	}

	totalValue, ok := hits["total"]
	if !ok {
		return nil, nil, 0, fmt.Errorf("missing total in search response")
	}

	var total int64
	if totalMap, ok := totalValue.(map[string]interface{}); ok {
		if value, ok := totalMap["value"].(float64); ok {
			total = int64(value)
		}
	} else if value, ok := totalValue.(float64); ok {
		total = int64(value)
	}

	// 提取文档和高亮信息
	hitsList, ok := hits["hits"].([]interface{})
	if !ok {
		return nil, nil, 0, fmt.Errorf("invalid hits format in search response")
	}

	documents := make([]*models.ConversationDocument, 0, len(hitsList))
	highlights := make([]map[string]interface{}, 0, len(hitsList))

	for _, hit := range hitsList {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}

		source, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			continue
		}

		// 解析文档
		doc := &models.ConversationDocument{}
		if err := r.parseDocument(source, doc); err != nil {
			continue // 跳过解析失败的文档
		}

		// 提取高亮信息
		highlight := make(map[string]interface{})
		if highlightData, exists := hitMap["highlight"]; exists {
			if highlightMap, ok := highlightData.(map[string]interface{}); ok {
				highlight = highlightMap
			}
		}

		documents = append(documents, doc)
		highlights = append(highlights, highlight)
	}

	return documents, highlights, total, nil
}

// contains 检查字符串是否包含子字符串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && indexOf(s, substr) >= 0))
}

// indexOf 查找子字符串在字符串中的位置
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// removeHighlightTags 移除高亮标签，提取纯文本内容
func removeHighlightTags(text string) string {
	// 移除 <mark> 和 </mark> 标签
	text = strings.ReplaceAll(text, "<mark>", "")
	text = strings.ReplaceAll(text, "</mark>", "")
	return text
}
