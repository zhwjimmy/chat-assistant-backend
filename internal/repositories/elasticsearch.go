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
	SearchConversationsWithMatchedMessages(query string, userID *uuid.UUID, providerID *string, tagID *uuid.UUID, startDate, endDate *time.Time, page, limit int) ([]*models.ConversationDocument, map[uuid.UUID][]*models.MessageDocument, map[uuid.UUID][]string, int64, error)
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
func (r *ElasticsearchRepositoryImpl) SearchConversationsWithMatchedMessages(query string, userID *uuid.UUID, providerID *string, tagID *uuid.UUID, startDate, endDate *time.Time, page, limit int) ([]*models.ConversationDocument, map[uuid.UUID][]*models.MessageDocument, map[uuid.UUID][]string, int64, error) {
	// 1. 在 ES 中搜索
	esDocs, highlights, total, err := r.searchConversationDocumentsWithHighlights(query, userID, providerID, tagID, startDate, endDate, page, limit)
	if err != nil {
		return nil, nil, nil, 0, err
	}

	// 2. 使用精确匹配过滤结果，确保关键词精确匹配
	filteredDocs := make([]*models.ConversationDocument, 0, len(esDocs))
	filteredHighlights := make([]map[string]interface{}, 0, len(highlights))

	for i, doc := range esDocs {
		// 检查是否真正包含关键词
		if r.hasExactMatch(doc, query) {
			filteredDocs = append(filteredDocs, doc)
			filteredHighlights = append(filteredHighlights, highlights[i])
		}
	}

	// 3. 按相关性评分排序
	r.sortByRelevance(filteredDocs, query)

	// 4. 提取匹配的消息和字段信息
	matchedMessagesMap := make(map[uuid.UUID][]*models.MessageDocument)
	matchedFieldsMap := make(map[uuid.UUID][]string)

	for i, doc := range filteredDocs {
		conversationID := doc.ID
		var matchedFields []string

		// 检查哪些字段有高亮（即匹配）
		if _, exists := filteredHighlights[i]["title"]; exists {
			matchedFields = append(matchedFields, "title")
		}
		if _, exists := filteredHighlights[i]["source_title"]; exists {
			matchedFields = append(matchedFields, "source_title")
		}
		if _, exists := filteredHighlights[i]["messages.content"]; exists {
			matchedFields = append(matchedFields, "messages.content")
		}
		if _, exists := filteredHighlights[i]["messages.source_content"]; exists {
			matchedFields = append(matchedFields, "messages.source_content")
		}
		if _, exists := filteredHighlights[i]["tags.name"]; exists {
			matchedFields = append(matchedFields, "tags.name")
		}

		// 提取匹配的消息
		_, hasContent := filteredHighlights[i]["messages.content"]
		_, hasSourceContent := filteredHighlights[i]["messages.source_content"]
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

				// 检查消息是否包含匹配的关键词（通过精确匹配判断）
				content := msgDoc.Content
				if content == "" {
					content = msgDoc.SourceContent
				}

				if countKeywordMatches(content, query) > 0 {
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
		}

		// 总是设置 matched_fields，即使为空
		matchedFieldsMap[conversationID] = matchedFields
	}

	return filteredDocs, matchedMessagesMap, matchedFieldsMap, total, nil
}

// buildSearchQuery 构建 ES 搜索查询
func (r *ElasticsearchRepositoryImpl) buildSearchQuery(query string, userID *uuid.UUID, providerID *string, tagID *uuid.UUID, startDate, endDate *time.Time, page, limit int) []byte {
	// 预处理查询词，确保精确匹配
	query = strings.TrimSpace(query)
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

	// Tag ID过滤 - 使用嵌套查询确保完全匹配
	if tagID != nil {
		mustQueries = append(mustQueries, map[string]interface{}{
			"nested": map[string]interface{}{
				"path": "tags",
				"query": map[string]interface{}{
					"term": map[string]interface{}{
						"tags.id": tagID.String(),
					},
				},
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

	// 搜索查询 - 平衡精确匹配和相关性
	searchQueries := []map[string]interface{}{
		// 1. 完全精确匹配 - 最高优先级 (权重: 10)
		{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"title.exact^10", "source_title.exact^8"},
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
						"fields": []string{"messages.content.exact^10", "messages.source_content.exact^8"},
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
						"fields": []string{"tags.name.exact^6"},
						"type":   "phrase",
						"slop":   0,
					},
				},
			},
		},
		// 2. 标准匹配 - 高优先级 (权重: 8)
		{
			"multi_match": map[string]interface{}{
				"query":     query,
				"fields":    []string{"title^8", "source_title^6"},
				"type":      "best_fields",
				"fuzziness": "AUTO",
			},
		},
		{
			"nested": map[string]interface{}{
				"path": "messages",
				"query": map[string]interface{}{
					"multi_match": map[string]interface{}{
						"query":     query,
						"fields":    []string{"messages.content^8", "messages.source_content^6"},
						"type":      "best_fields",
						"fuzziness": "AUTO",
					},
				},
			},
		},
		{
			"nested": map[string]interface{}{
				"path": "tags",
				"query": map[string]interface{}{
					"multi_match": map[string]interface{}{
						"query":     query,
						"fields":    []string{"tags.name^6"},
						"type":      "best_fields",
						"fuzziness": "AUTO",
					},
				},
			},
		},
		// 3. 词级别匹配 - 中等优先级 (权重: 5)
		{
			"multi_match": map[string]interface{}{
				"query":    query,
				"fields":   []string{"title^5", "source_title^4"},
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
						"fields":   []string{"messages.content^5", "messages.source_content^4"},
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
						"fields":   []string{"tags.name^4"},
						"type":     "cross_fields",
						"operator": "and", // 所有词都必须匹配
					},
				},
			},
		},
		// 4. 部分匹配 - 低优先级 (权重: 2)
		{
			"multi_match": map[string]interface{}{
				"query":    query,
				"fields":   []string{"title^2", "source_title^1"},
				"type":     "best_fields",
				"operator": "or", // 任意词匹配即可
			},
		},
		{
			"nested": map[string]interface{}{
				"path": "messages",
				"query": map[string]interface{}{
					"multi_match": map[string]interface{}{
						"query":    query,
						"fields":   []string{"messages.content^2", "messages.source_content^1"},
						"type":     "best_fields",
						"operator": "or", // 任意词匹配即可
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
						"type":     "best_fields",
						"operator": "or", // 任意词匹配即可
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
			"pre_tags":            []string{"<mark>"},
			"post_tags":           []string{"</mark>"},
			"fragment_size":       150, // 限制高亮片段长度
			"number_of_fragments": 3,   // 最多返回3个高亮片段
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
func (r *ElasticsearchRepositoryImpl) searchConversationDocumentsWithHighlights(query string, userID *uuid.UUID, providerID *string, tagID *uuid.UUID, startDate, endDate *time.Time, page, limit int) ([]*models.ConversationDocument, []map[string]interface{}, int64, error) {
	ctx := context.Background()

	// 构建 ES 查询
	searchQuery := r.buildSearchQuery(query, userID, providerID, tagID, startDate, endDate, page, limit)

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
		// 读取错误响应体以获取更详细的错误信息
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&errorResponse); err == nil {
			return nil, nil, 0, fmt.Errorf("search request failed with status: %s, error: %v", res.Status(), errorResponse)
		}
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

// countKeywordMatches 计算关键词在文本中的精确匹配次数
func countKeywordMatches(text, keyword string) int {
	if text == "" || keyword == "" {
		return 0
	}

	// 转换为小写进行匹配
	lowerText := strings.ToLower(text)
	lowerKeyword := strings.ToLower(keyword)

	// 使用词边界匹配，确保精确匹配
	// 例如："英语" 不会匹配 "英无语"
	count := 0
	start := 0

	for {
		pos := strings.Index(lowerText[start:], lowerKeyword)
		if pos == -1 {
			break
		}

		actualPos := start + pos

		// 检查词边界
		isWordBoundary := true

		// 检查前一个字符
		if actualPos > 0 {
			prevChar := lowerText[actualPos-1]
			if isWordChar(prevChar) {
				isWordBoundary = false
			}
		}

		// 检查后一个字符
		if actualPos+len(lowerKeyword) < len(lowerText) {
			nextChar := lowerText[actualPos+len(lowerKeyword)]
			if isWordChar(nextChar) {
				isWordBoundary = false
			}
		}

		if isWordBoundary {
			count++
		}

		start = actualPos + len(lowerKeyword)
	}

	return count
}

// isWordChar 检查字符是否为单词字符
func isWordChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') || c == '_' || c >= 128 // 包含中文字符
}

// calculateRelevanceScore 计算相关性评分
func calculateRelevanceScore(conversationDoc *models.ConversationDocument, keyword string) float64 {
	score := 0.0

	// 计算标题匹配
	title := conversationDoc.Title
	if title == "" {
		title = conversationDoc.SourceTitle
	}

	titleMatches := countRelevanceMatches(title, keyword)
	score += float64(titleMatches) * 10.0 // 标题匹配权重最高

	// 计算消息内容匹配
	messageMatches := 0
	for _, msg := range conversationDoc.Messages {
		content := msg.Content
		if content == "" {
			content = msg.SourceContent
		}
		messageMatches += countRelevanceMatches(content, keyword)
	}
	score += float64(messageMatches) * 5.0 // 消息匹配权重中等

	// 计算标签匹配
	tagMatches := 0
	for _, tag := range conversationDoc.Tags {
		tagMatches += countRelevanceMatches(tag.Name, keyword)
	}
	score += float64(tagMatches) * 8.0 // 标签匹配权重较高

	return score
}

// countRelevanceMatches 计算相关性匹配次数（更宽松的匹配）
func countRelevanceMatches(text, keyword string) int {
	if text == "" || keyword == "" {
		return 0
	}

	// 转换为小写进行匹配
	lowerText := strings.ToLower(text)
	lowerKeyword := strings.ToLower(keyword)

	// 直接包含匹配
	if strings.Contains(lowerText, lowerKeyword) {
		// 计算出现次数
		count := 0
		start := 0
		for {
			pos := strings.Index(lowerText[start:], lowerKeyword)
			if pos == -1 {
				break
			}
			count++
			start = start + pos + len(lowerKeyword)
		}
		return count
	}

	// 对于短关键词，也使用词边界匹配
	if len([]rune(keyword)) <= 3 {
		return countKeywordMatches(text, keyword)
	}

	return 0
}

// hasExactMatch 检查对话是否包含相关匹配的关键词
func (r *ElasticsearchRepositoryImpl) hasExactMatch(doc *models.ConversationDocument, keyword string) bool {
	// 检查标题
	title := doc.Title
	if title == "" {
		title = doc.SourceTitle
	}

	if r.containsKeyword(title, keyword) {
		return true
	}

	// 检查消息内容
	for _, msg := range doc.Messages {
		content := msg.Content
		if content == "" {
			content = msg.SourceContent
		}
		if r.containsKeyword(content, keyword) {
			return true
		}
	}

	// 检查标签
	for _, tag := range doc.Tags {
		if r.containsKeyword(tag.Name, keyword) {
			return true
		}
	}

	return false
}

// containsKeyword 检查文本是否包含关键词（更宽松的匹配）
func (r *ElasticsearchRepositoryImpl) containsKeyword(text, keyword string) bool {
	if text == "" || keyword == "" {
		return false
	}

	// 转换为小写进行匹配
	lowerText := strings.ToLower(text)
	lowerKeyword := strings.ToLower(keyword)

	// 直接包含匹配
	if strings.Contains(lowerText, lowerKeyword) {
		return true
	}

	// 对于中文，也检查词边界匹配
	if len([]rune(keyword)) <= 3 { // 短关键词使用词边界匹配
		return countKeywordMatches(text, keyword) > 0
	}

	return false
}

// sortByRelevance 按相关性评分排序对话
func (r *ElasticsearchRepositoryImpl) sortByRelevance(docs []*models.ConversationDocument, keyword string) {
	// 使用简单的冒泡排序，按相关性评分降序排列
	n := len(docs)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			score1 := calculateRelevanceScore(docs[j], keyword)
			score2 := calculateRelevanceScore(docs[j+1], keyword)

			if score1 < score2 {
				// 交换位置
				docs[j], docs[j+1] = docs[j+1], docs[j]
			}
		}
	}
}
