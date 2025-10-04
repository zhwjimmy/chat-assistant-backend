package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"chat-assistant-backend/internal/models"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
)

// ElasticsearchRepository handles Elasticsearch search operations
type ElasticsearchRepository struct {
	esClient  *es.Client
	indexName string
}

// NewElasticsearchRepository creates a new Elasticsearch repository
func NewElasticsearchRepository(esClient *es.Client, indexName string) *ElasticsearchRepository {
	return &ElasticsearchRepository{
		esClient:  esClient,
		indexName: indexName,
	}
}

// SearchConversationsWithMessages searches conversations using Elasticsearch
func (r *ElasticsearchRepository) SearchConversationsWithMessages(query string, userID *uuid.UUID, page, limit int) ([]*models.Conversation, int64, error) {
	// 1. 在 ES 中搜索
	esDocs, total, err := r.searchConversationDocuments(query, userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// 2. 转换为 PostgreSQL 模型
	conversations := make([]*models.Conversation, len(esDocs))
	for i, doc := range esDocs {
		conversations[i] = doc.ToConversation()
	}

	return conversations, total, nil
}

// searchConversationDocuments 在 ES 中搜索 conversation 文档
func (r *ElasticsearchRepository) searchConversationDocuments(query string, userID *uuid.UUID, page, limit int) ([]*models.ConversationDocument, int64, error) {
	ctx := context.Background()

	// 构建 ES 查询
	searchQuery := r.buildSearchQuery(query, userID, page, limit)

	// 执行搜索
	req := esapi.SearchRequest{
		Index: []string{r.indexName},
		Body:  bytes.NewReader(searchQuery),
	}

	res, err := req.Do(ctx, r.esClient)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, fmt.Errorf("search request failed with status: %s", res.Status())
	}

	// 解析响应
	var searchResponse map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&searchResponse); err != nil {
		return nil, 0, fmt.Errorf("failed to decode search response: %w", err)
	}

	// 提取结果
	return r.parseSearchResponse(searchResponse)
}

// buildSearchQuery 构建 ES 搜索查询
func (r *ElasticsearchRepository) buildSearchQuery(query string, userID *uuid.UUID, page, limit int) []byte {
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
func (r *ElasticsearchRepository) parseSearchResponse(response map[string]interface{}) ([]*models.ConversationDocument, int64, error) {
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

// parseDocument 解析单个文档
func (r *ElasticsearchRepository) parseDocument(source map[string]interface{}, doc *models.ConversationDocument) error {
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

	return nil
}

// parseMessageDocument 解析消息文档
func (r *ElasticsearchRepository) parseMessageDocument(source map[string]interface{}, doc *models.MessageDocument) error {
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
