package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"chat-assistant-backend/internal/config"
	"chat-assistant-backend/internal/infra/elasticsearch"
	"chat-assistant-backend/internal/repositories"
)

func main() {
	command := flag.String("command", "status", "Command: status, init, recreate, health")
	flag.Parse()

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 创建 ES 客户端
	esConfig := &elasticsearch.Config{
		Hosts:    cfg.Elasticsearch.Hosts,
		Username: cfg.Elasticsearch.Username,
		Password: cfg.Elasticsearch.Password,
		Timeout:  cfg.Elasticsearch.Timeout,
		Index: elasticsearch.IndexConfig{
			Conversations: cfg.Elasticsearch.Index.Conversations,
			Messages:      cfg.Elasticsearch.Index.Messages,
		},
	}

	client, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v", err)
	}

	// 创建索引器
	indexer := repositories.NewElasticsearchIndexer(client.GetClient(), cfg.Elasticsearch.Index.Conversations)

	// 创建初始化器
	initializer := elasticsearch.NewInitializer(client, indexer)

	ctx := context.Background()

	// 执行命令
	switch *command {
	case "status":
		if err := showStatus(ctx, initializer); err != nil {
			log.Fatalf("Failed to get status: %v", err)
		}
	case "init":
		if err := initializeIndexes(ctx, initializer); err != nil {
			log.Fatalf("Failed to initialize indexes: %v", err)
		}
		fmt.Println("Indexes initialized successfully")
	case "recreate":
		if err := recreateIndexes(ctx, initializer); err != nil {
			log.Fatalf("Failed to recreate indexes: %v", err)
		}
		fmt.Println("Indexes recreated successfully")
	case "health":
		if err := showHealth(ctx, client); err != nil {
			log.Fatalf("Failed to get health: %v", err)
		}
	default:
		fmt.Printf("Unknown command: %s\n", *command)
		fmt.Println("Available commands: status, init, recreate, health")
		os.Exit(1)
	}
}

func showStatus(ctx context.Context, initializer *elasticsearch.Initializer) error {
	status, err := initializer.GetIndexStatus(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Elasticsearch Index Status:")
	fmt.Printf("  Conversation Index Exists: %v\n", status["conversation_index_exists"])
	fmt.Printf("  Message Index Exists: %v\n", status["message_index_exists"])

	if health, ok := status["cluster_health"].(map[string]interface{}); ok {
		fmt.Printf("  Cluster Status: %v\n", health["status"])
		fmt.Printf("  Number of Nodes: %v\n", health["number_of_nodes"])
	}

	return nil
}

func initializeIndexes(ctx context.Context, initializer *elasticsearch.Initializer) error {
	return initializer.Initialize(ctx)
}

func recreateIndexes(ctx context.Context, initializer *elasticsearch.Initializer) error {
	return initializer.RecreateIndexes(ctx)
}

func showHealth(ctx context.Context, client *elasticsearch.Client) error {
	healthChecker := elasticsearch.NewHealthChecker(client)
	status := healthChecker.Check(ctx)

	fmt.Printf("Elasticsearch Health Status: %s\n", status.Status)
	if status.Error != "" {
		fmt.Printf("Error: %s\n", status.Error)
	}

	if status.Details != nil {
		fmt.Println("Details:")
		for key, value := range status.Details {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	return nil
}
