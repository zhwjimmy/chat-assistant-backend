package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"chat-assistant-backend/internal/models"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
)

// ElasticsearchIndexer 定义 Elasticsearch 索引操作接口
type ElasticsearchIndexer interface {
	// 索引 conversation 文档
	IndexConversation(doc *models.ConversationDocument) error

	// 向 conversation 添加 message
	AddMessageToConversation(conversationID uuid.UUID, message models.MessageDocument) error

	// 更新 conversation 中的 message
	UpdateMessageInConversation(conversationID uuid.UUID, message models.MessageDocument) error

	// 从 conversation 中删除 message
	RemoveMessageFromConversation(conversationID uuid.UUID, messageID uuid.UUID) error

	// 删除整个 conversation
	DeleteConversation(conversationID uuid.UUID) error

	// 批量索引 conversations
	BulkIndexConversations(docs []*models.ConversationDocument) error

	// 更新 conversation 基本信息（不包含 messages）
	UpdateConversation(doc *models.ConversationDocument) error

	// 检查 conversation 是否存在
	ConversationExists(conversationID uuid.UUID) (bool, error)
}

// ElasticsearchIndexerImpl 默认的索引器实现
type ElasticsearchIndexerImpl struct {
	esClient  *es.Client
	indexName string
}

// NewElasticsearchIndexer 创建新的索引器
func NewElasticsearchIndexer(esClient *es.Client, indexName string) ElasticsearchIndexer {
	return &ElasticsearchIndexerImpl{
		esClient:  esClient,
		indexName: indexName,
	}
}

// IndexConversation 索引 conversation 文档
func (i *ElasticsearchIndexerImpl) IndexConversation(doc *models.ConversationDocument) error {
	ctx := context.Background()

	// 序列化文档
	docBytes, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal conversation document: %w", err)
	}

	// 创建索引请求
	req := esapi.IndexRequest{
		Index:      i.indexName,
		DocumentID: doc.ID.String(),
		Body:       bytes.NewReader(docBytes),
		Refresh:    "true",
	}

	// 执行请求
	res, err := req.Do(ctx, i.esClient)
	if err != nil {
		return fmt.Errorf("failed to index conversation: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("index request failed with status: %s", res.Status())
	}

	return nil
}

// AddMessageToConversation 向 conversation 添加 message
func (i *ElasticsearchIndexerImpl) AddMessageToConversation(conversationID uuid.UUID, message models.MessageDocument) error {
	ctx := context.Background()

	// 构建脚本，向 messages 数组添加新消息
	script := `
		if (ctx._source.messages == null) {
			ctx._source.messages = []
		}
		ctx._source.messages.add(params.message)
	`

	// 序列化消息
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	var messageData map[string]interface{}
	if err := json.Unmarshal(messageBytes, &messageData); err != nil {
		return fmt.Errorf("failed to unmarshal message data: %w", err)
	}

	// 构建更新请求
	updateBody := map[string]interface{}{
		"script": map[string]interface{}{
			"source": script,
			"params": map[string]interface{}{
				"message": messageData,
			},
		},
	}

	updateBytes, err := json.Marshal(updateBody)
	if err != nil {
		return fmt.Errorf("failed to marshal update body: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      i.indexName,
		DocumentID: conversationID.String(),
		Body:       bytes.NewReader(updateBytes),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, i.esClient)
	if err != nil {
		return fmt.Errorf("failed to add message to conversation: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("update request failed with status: %s", res.Status())
	}

	return nil
}

// UpdateMessageInConversation 更新 conversation 中的 message
func (i *ElasticsearchIndexerImpl) UpdateMessageInConversation(conversationID uuid.UUID, message models.MessageDocument) error {
	ctx := context.Background()

	// 构建脚本，更新 messages 数组中的特定消息
	script := `
		if (ctx._source.messages != null) {
			for (int i = 0; i < ctx._source.messages.size(); i++) {
				if (ctx._source.messages[i].id == params.messageId) {
					ctx._source.messages[i] = params.message
					break
				}
			}
		}
	`

	// 序列化消息
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	var messageData map[string]interface{}
	if err := json.Unmarshal(messageBytes, &messageData); err != nil {
		return fmt.Errorf("failed to unmarshal message data: %w", err)
	}

	// 构建更新请求
	updateBody := map[string]interface{}{
		"script": map[string]interface{}{
			"source": script,
			"params": map[string]interface{}{
				"messageId": message.ID.String(),
				"message":   messageData,
			},
		},
	}

	updateBytes, err := json.Marshal(updateBody)
	if err != nil {
		return fmt.Errorf("failed to marshal update body: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      i.indexName,
		DocumentID: conversationID.String(),
		Body:       bytes.NewReader(updateBytes),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, i.esClient)
	if err != nil {
		return fmt.Errorf("failed to update message in conversation: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("update request failed with status: %s", res.Status())
	}

	return nil
}

// RemoveMessageFromConversation 从 conversation 中删除 message
func (i *ElasticsearchIndexerImpl) RemoveMessageFromConversation(conversationID uuid.UUID, messageID uuid.UUID) error {
	ctx := context.Background()

	// 构建脚本，从 messages 数组中删除特定消息
	script := `
		if (ctx._source.messages != null) {
			ctx._source.messages.removeIf(msg -> msg.id == params.messageId)
		}
	`

	// 构建更新请求
	updateBody := map[string]interface{}{
		"script": map[string]interface{}{
			"source": script,
			"params": map[string]interface{}{
				"messageId": messageID.String(),
			},
		},
	}

	updateBytes, err := json.Marshal(updateBody)
	if err != nil {
		return fmt.Errorf("failed to marshal update body: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      i.indexName,
		DocumentID: conversationID.String(),
		Body:       bytes.NewReader(updateBytes),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, i.esClient)
	if err != nil {
		return fmt.Errorf("failed to remove message from conversation: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("update request failed with status: %s", res.Status())
	}

	return nil
}

// DeleteConversation 删除整个 conversation
func (i *ElasticsearchIndexerImpl) DeleteConversation(conversationID uuid.UUID) error {
	ctx := context.Background()

	req := esapi.DeleteRequest{
		Index:      i.indexName,
		DocumentID: conversationID.String(),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, i.esClient)
	if err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("delete request failed with status: %s", res.Status())
	}

	return nil
}

// BulkIndexConversations 批量索引 conversations
func (i *ElasticsearchIndexerImpl) BulkIndexConversations(docs []*models.ConversationDocument) error {
	if len(docs) == 0 {
		return nil
	}

	ctx := context.Background()

	// 构建批量请求体
	var bulkBody strings.Builder
	for _, doc := range docs {
		// 添加索引操作元数据
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": i.indexName,
				"_id":    doc.ID.String(),
			},
		}
		metaBytes, _ := json.Marshal(meta)
		bulkBody.Write(metaBytes)
		bulkBody.WriteString("\n")

		// 添加文档数据
		docBytes, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("failed to marshal conversation document: %w", err)
		}
		bulkBody.Write(docBytes)
		bulkBody.WriteString("\n")
	}

	// 执行批量请求
	req := esapi.BulkRequest{
		Body:    strings.NewReader(bulkBody.String()),
		Refresh: "true",
	}

	res, err := req.Do(ctx, i.esClient)
	if err != nil {
		return fmt.Errorf("failed to bulk index conversations: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk request failed with status: %s", res.Status())
	}

	return nil
}

// UpdateConversation 更新 conversation 基本信息（不包含 messages）
func (i *ElasticsearchIndexerImpl) UpdateConversation(doc *models.ConversationDocument) error {
	ctx := context.Background()

	// 构建更新文档，排除 messages 字段
	updateDoc := map[string]interface{}{
		"id":           doc.ID,
		"user_id":      doc.UserID,
		"title":        doc.Title,
		"provider":     doc.Provider,
		"model":        doc.Model,
		"source_id":    doc.SourceID,
		"source_title": doc.SourceTitle,
		"created_at":   doc.CreatedAt,
		"updated_at":   doc.UpdatedAt,
	}

	updateBytes, err := json.Marshal(updateDoc)
	if err != nil {
		return fmt.Errorf("failed to marshal update document: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      i.indexName,
		DocumentID: doc.ID.String(),
		Body:       bytes.NewReader(updateBytes),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, i.esClient)
	if err != nil {
		return fmt.Errorf("failed to update conversation: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("update request failed with status: %s", res.Status())
	}

	return nil
}

// ConversationExists 检查 conversation 是否存在
func (i *ElasticsearchIndexerImpl) ConversationExists(conversationID uuid.UUID) (bool, error) {
	ctx := context.Background()

	req := esapi.ExistsRequest{
		Index:      i.indexName,
		DocumentID: conversationID.String(),
	}

	res, err := req.Do(ctx, i.esClient)
	if err != nil {
		return false, fmt.Errorf("failed to check conversation existence: %w", err)
	}
	defer res.Body.Close()

	return res.StatusCode == 200, nil
}
