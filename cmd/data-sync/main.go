package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"chat-assistant-backend/internal/config"
	"chat-assistant-backend/internal/infra/elasticsearch"
	"chat-assistant-backend/internal/repositories"
	"chat-assistant-backend/internal/services"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 命令行参数
	var (
		dryRun = flag.Bool("dry-run", false, "试运行，不实际同步")
		help   = flag.Bool("help", false, "显示帮助信息")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	db, err := initializeDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化 ES 客户端
	esClient, err := initializeElasticsearch(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Elasticsearch: %v", err)
	}

	// 创建 repositories
	conversationRepo := repositories.NewConversationRepository(db)
	indexer := repositories.NewElasticsearchIndexer(esClient.GetClient(), cfg.Elasticsearch.Index.Conversations)

	// 创建同步服务
	syncService := services.NewSyncService(conversationRepo, indexer)

	// 执行同步
	if *dryRun {
		log.Println("Dry run mode - fetching sample data...")
		conversations, err := conversationRepo.FindAll()
		if err != nil {
			log.Fatalf("Failed to fetch conversations: %v", err)
		}
		log.Printf("Dry run: Found %d conversations to sync", len(conversations))

		// 显示前几个 conversation 的示例
		if len(conversations) > 0 {
			log.Printf("Sample conversation: ID=%s, Title=%s, Messages=%d",
				conversations[0].ID, conversations[0].Title, len(conversations[0].Messages))
		}
		log.Println("Dry run completed - no data was actually synced")
	} else {
		log.Println("Starting data sync...")
		if err := syncService.SyncAll(); err != nil {
			log.Fatalf("Sync failed: %v", err)
		}
		log.Println("Data sync completed successfully")
	}
}

func showHelp() {
	fmt.Println("Data Sync Tool - 同步数据库数据到 Elasticsearch")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  data-sync [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -dry-run")
	fmt.Println("       试运行，不实际同步")
	fmt.Println("  -help")
	fmt.Println("       显示帮助信息")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  data-sync                    # 同步所有数据")
	fmt.Println("  data-sync -dry-run          # 试运行")
}

func initializeDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.Database.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 测试连接
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established")
	return db, nil
}

func initializeElasticsearch(cfg *config.Config) (*elasticsearch.Client, error) {
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
		return nil, fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	// 测试连接
	ctx := context.Background()
	if err := client.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping Elasticsearch: %w", err)
	}

	log.Println("Elasticsearch connection established")
	return client, nil
}
