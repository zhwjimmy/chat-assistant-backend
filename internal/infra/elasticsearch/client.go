package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"chat-assistant-backend/internal/config"
	"chat-assistant-backend/internal/repositories"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// Client wraps the Elasticsearch client with additional functionality
type Client struct {
	es  *elasticsearch.Client
	cfg *Config
}

// NewClient creates a new Elasticsearch client
func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// Build client configuration
	esConfig := elasticsearch.Config{
		Addresses: cfg.Hosts,
	}

	// Add authentication if provided
	if cfg.Username != "" && cfg.Password != "" {
		esConfig.Username = cfg.Username
		esConfig.Password = cfg.Password
	}

	// Create Elasticsearch client
	es, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	client := &Client{
		es:  es,
		cfg: cfg,
	}

	// Test connection
	if err := client.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to Elasticsearch: %w", err)
	}

	return client, nil
}

// GetClient returns the underlying Elasticsearch client
func (c *Client) GetClient() *elasticsearch.Client {
	return c.es
}

// GetConfig returns the client configuration
func (c *Client) GetConfig() *Config {
	return c.cfg
}

// Ping tests the connection to Elasticsearch
func (c *Client) Ping(ctx context.Context) error {
	req := esapi.PingRequest{
		Pretty: false,
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("ping request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("ping failed with status: %s", res.Status())
	}

	return nil
}

// Info returns cluster information
func (c *Client) Info(ctx context.Context) (map[string]interface{}, error) {
	req := esapi.InfoRequest{
		Pretty: false,
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return nil, fmt.Errorf("info request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("info request failed with status: %s", res.Status())
	}

	var info map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to decode info response: %w", err)
	}

	return info, nil
}

// ClusterHealth returns cluster health status
func (c *Client) ClusterHealth(ctx context.Context) (map[string]interface{}, error) {
	req := esapi.ClusterHealthRequest{
		Pretty: false,
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return nil, fmt.Errorf("cluster health request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("cluster health request failed with status: %s", res.Status())
	}

	var health map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&health); err != nil {
		return nil, fmt.Errorf("failed to decode cluster health response: %w", err)
	}

	return health, nil
}

// CreateIndex creates an index with the given name and mapping
func (c *Client) CreateIndex(ctx context.Context, indexName string, mapping string) error {
	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  strings.NewReader(mapping),
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("create index request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("create index failed with status: %s", res.Status())
	}

	return nil
}

// DeleteIndex deletes an index
func (c *Client) DeleteIndex(ctx context.Context, indexName string) error {
	req := esapi.IndicesDeleteRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("delete index request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("delete index failed with status: %s", res.Status())
	}

	return nil
}

// IndexExists checks if an index exists
func (c *Client) IndexExists(ctx context.Context, indexName string) (bool, error) {
	req := esapi.IndicesExistsRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return false, fmt.Errorf("index exists request failed: %w", err)
	}
	defer res.Body.Close()

	return res.StatusCode == 200, nil
}

// HealthChecker provides health check functionality for Elasticsearch
type HealthChecker struct {
	client *Client
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(client *Client) *HealthChecker {
	return &HealthChecker{
		client: client,
	}
}

// HealthStatus represents the health status of Elasticsearch
type HealthStatus struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

// Check performs a health check on Elasticsearch
func (h *HealthChecker) Check(ctx context.Context) *HealthStatus {
	status := &HealthStatus{
		Timestamp: time.Now(),
	}

	// Test basic connectivity
	if err := h.client.Ping(ctx); err != nil {
		status.Status = "unhealthy"
		status.Error = fmt.Sprintf("ping failed: %v", err)
		return status
	}

	// Get cluster health
	health, err := h.client.ClusterHealth(ctx)
	if err != nil {
		status.Status = "degraded"
		status.Error = fmt.Sprintf("cluster health check failed: %v", err)
		status.Details = map[string]interface{}{
			"ping": "ok",
		}
		return status
	}

	// Extract cluster status
	clusterStatus, ok := health["status"].(string)
	if !ok {
		status.Status = "unknown"
		status.Error = "unable to determine cluster status"
		status.Details = health
		return status
	}

	// Map Elasticsearch cluster status to our health status
	switch clusterStatus {
	case "green":
		status.Status = "healthy"
	case "yellow":
		status.Status = "degraded"
	case "red":
		status.Status = "unhealthy"
	default:
		status.Status = "unknown"
	}

	status.Details = health
	return status
}

// IsHealthy returns true if Elasticsearch is healthy
func (h *HealthChecker) IsHealthy(ctx context.Context) bool {
	status := h.Check(ctx)
	return status.Status == "healthy"
}

// WaitForHealthy waits for Elasticsearch to become healthy
func (h *HealthChecker) WaitForHealthy(ctx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for Elasticsearch to become healthy: %w", ctx.Err())
		case <-ticker.C:
			if h.IsHealthy(ctx) {
				return nil
			}
		}
	}
}

// NewElasticsearchClientFromConfig creates a new Elasticsearch client from config
func NewElasticsearchClientFromConfig(cfg *config.Config) (*Client, error) {
	esConfig := &Config{
		Hosts:    cfg.Elasticsearch.Hosts,
		Username: cfg.Elasticsearch.Username,
		Password: cfg.Elasticsearch.Password,
		Timeout:  cfg.Elasticsearch.Timeout,
		Index: IndexConfig{
			Conversations: cfg.Elasticsearch.Index.Conversations,
			Messages:      cfg.Elasticsearch.Index.Messages,
		},
	}

	return NewClient(esConfig)
}

// NewElasticsearchIndexerFromClient creates a new Elasticsearch indexer from client
func NewElasticsearchIndexerFromClient(esClient *Client, cfg *config.Config) repositories.ElasticsearchIndexer {
	return repositories.NewElasticsearchIndexer(esClient.GetClient(), cfg.Elasticsearch.Index.Conversations)
}

// NewElasticsearchClient extracts the underlying Elasticsearch client
func NewElasticsearchClient(client *Client) *elasticsearch.Client {
	return client.GetClient()
}

// NewElasticsearchIndexName provides the conversations index name
func NewElasticsearchIndexName(cfg *config.Config) string {
	return cfg.Elasticsearch.Index.Conversations
}
