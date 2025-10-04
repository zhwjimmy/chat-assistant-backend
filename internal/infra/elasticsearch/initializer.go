package elasticsearch

import (
	"context"
	"fmt"
	"time"

	"chat-assistant-backend/internal/repositories"
)

// Initializer 负责初始化 Elasticsearch 索引
type Initializer struct {
	client  *Client
	indexer repositories.ElasticsearchIndexer
}

// NewInitializer 创建新的初始化器
func NewInitializer(client *Client, indexer repositories.ElasticsearchIndexer) *Initializer {
	return &Initializer{
		client:  client,
		indexer: indexer,
	}
}

// Initialize 初始化所有必要的索引
func (i *Initializer) Initialize(ctx context.Context) error {
	cfg := i.client.GetConfig()

	// 等待 Elasticsearch 可用
	healthChecker := NewHealthChecker(i.client)
	if err := healthChecker.WaitForHealthy(ctx, 60*time.Second); err != nil {
		return fmt.Errorf("elasticsearch is not healthy: %w", err)
	}

	// 创建 conversation 索引
	if err := i.createConversationIndex(ctx, cfg.Index.Conversations); err != nil {
		return fmt.Errorf("failed to create conversation index: %w", err)
	}

	// 创建 message 索引（如果需要独立索引）
	if err := i.createMessageIndex(ctx, cfg.Index.Messages); err != nil {
		return fmt.Errorf("failed to create message index: %w", err)
	}

	return nil
}

// createConversationIndex 创建 conversation 索引
func (i *Initializer) createConversationIndex(ctx context.Context, indexName string) error {
	// 检查索引是否已存在
	exists, err := i.client.IndexExists(ctx, indexName)
	if err != nil {
		return fmt.Errorf("failed to check if conversation index exists: %w", err)
	}

	if exists {
		// 索引已存在，可以选择更新映射或跳过
		return nil
	}

	// 创建索引
	mapping := ConversationMapping()
	if err := i.client.CreateIndex(ctx, indexName, mapping); err != nil {
		return fmt.Errorf("failed to create conversation index: %w", err)
	}

	return nil
}

// createMessageIndex 创建 message 索引
func (i *Initializer) createMessageIndex(ctx context.Context, indexName string) error {
	// 检查索引是否已存在
	exists, err := i.client.IndexExists(ctx, indexName)
	if err != nil {
		return fmt.Errorf("failed to check if message index exists: %w", err)
	}

	if exists {
		// 索引已存在，可以选择更新映射或跳过
		return nil
	}

	// 创建索引
	mapping := MessageMapping()
	if err := i.client.CreateIndex(ctx, indexName, mapping); err != nil {
		return fmt.Errorf("failed to create message index: %w", err)
	}

	return nil
}

// RecreateIndexes 重新创建所有索引（会删除现有数据）
func (i *Initializer) RecreateIndexes(ctx context.Context) error {
	cfg := i.client.GetConfig()

	// 删除现有索引
	if err := i.client.DeleteIndex(ctx, cfg.Index.Conversations); err != nil {
		// 忽略索引不存在的错误
	}

	if err := i.client.DeleteIndex(ctx, cfg.Index.Messages); err != nil {
		// 忽略索引不存在的错误
	}

	// 重新创建索引
	return i.Initialize(ctx)
}

// GetIndexStatus 获取索引状态信息
func (i *Initializer) GetIndexStatus(ctx context.Context) (map[string]interface{}, error) {
	cfg := i.client.GetConfig()
	status := make(map[string]interface{})

	// 检查 conversation 索引状态
	convExists, err := i.client.IndexExists(ctx, cfg.Index.Conversations)
	if err != nil {
		return nil, fmt.Errorf("failed to check conversation index status: %w", err)
	}
	status["conversation_index_exists"] = convExists

	// 检查 message 索引状态
	msgExists, err := i.client.IndexExists(ctx, cfg.Index.Messages)
	if err != nil {
		return nil, fmt.Errorf("failed to check message index status: %w", err)
	}
	status["message_index_exists"] = msgExists

	// 获取集群健康状态
	health, err := i.client.ClusterHealth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster health: %w", err)
	}
	status["cluster_health"] = health

	return status, nil
}

// ConversationMapping 返回 conversation 索引的映射定义
func ConversationMapping() string {
	return `{
		"mappings": {
			"properties": {
				"id": {
					"type": "keyword"
				},
				"user_id": {
					"type": "keyword"
				},
				"title": {
					"type": "text",
					"analyzer": "standard",
					"fields": {
						"keyword": {
							"type": "keyword"
						},
						"exact": {
							"type": "text",
							"analyzer": "keyword"
						}
					}
				},
				"provider": {
					"type": "keyword"
				},
				"model": {
					"type": "keyword"
				},
				"source_id": {
					"type": "keyword"
				},
				"source_title": {
					"type": "text",
					"analyzer": "standard",
					"fields": {
						"exact": {
							"type": "text",
							"analyzer": "keyword"
						}
					}
				},
				"created_at": {
					"type": "date"
				},
				"updated_at": {
					"type": "date"
				},
				"messages": {
					"type": "nested",
					"properties": {
						"id": {
							"type": "keyword"
						},
						"conversation_id": {
							"type": "keyword"
						},
						"role": {
							"type": "keyword"
						},
						"content": {
							"type": "text",
							"analyzer": "standard",
							"fields": {
								"exact": {
									"type": "text",
									"analyzer": "keyword"
								}
							}
						},
						"source_id": {
							"type": "keyword"
						},
						"source_content": {
							"type": "text",
							"analyzer": "standard",
							"fields": {
								"exact": {
									"type": "text",
									"analyzer": "keyword"
								}
							}
						},
						"created_at": {
							"type": "date"
						},
						"updated_at": {
							"type": "date"
						}
					}
				}
			}
		},
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0,
			"analysis": {
				"analyzer": {
					"standard": {
						"type": "standard",
						"stopwords": "_english_"
					}
				}
			}
		}
	}`
}

// MessageMapping 返回 message 索引的映射定义（独立索引方案）
func MessageMapping() string {
	return `{
		"mappings": {
			"properties": {
				"id": {
					"type": "keyword"
				},
				"conversation_id": {
					"type": "keyword"
				},
				"role": {
					"type": "keyword"
				},
				"content": {
					"type": "text",
					"analyzer": "standard"
				},
				"source_id": {
					"type": "keyword"
				},
				"source_content": {
					"type": "text",
					"analyzer": "standard"
				},
				"created_at": {
					"type": "date"
				},
				"updated_at": {
					"type": "date"
				}
			}
		},
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0,
			"analysis": {
				"analyzer": {
					"standard": {
						"type": "standard",
						"stopwords": "_english_"
					}
				}
			}
		}
	}`
}
