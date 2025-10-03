package importer

import (
	"context"
	"fmt"
	"os"
	"time"

	"chat-assistant-backend/internal/config"
	"chat-assistant-backend/internal/importer/parsers"
	"chat-assistant-backend/internal/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Importer 核心导入器
type Importer struct {
	config      *config.Config
	loader      *Loader
	validator   *Validator
	transformer *Transformer
}

// ImportResult 导入结果
type ImportResult struct {
	Platform          string   `json:"platform"`
	ConversationCount int      `json:"conversation_count"`
	MessageCount      int      `json:"message_count"`
	SuccessCount      int      `json:"success_count"`
	ErrorCount        int      `json:"error_count"`
	Errors            []string `json:"errors,omitempty"`
	Duration          string   `json:"duration"`
}

// NewImporter 创建导入器
func NewImporter(cfg *config.Config) *Importer {
	return &Importer{
		config:      cfg,
		loader:      NewLoader(cfg),
		validator:   NewValidator(),
		transformer: NewTransformer(),
	}
}

// Import 执行导入
func (i *Importer) Import(filePath, platform, userIDStr string, dryRun bool) (*ImportResult, error) {
	startTime := time.Now()
	log := logger.GetLogger()

	// 解析用户ID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	log.Info("Starting import process",
		zap.String("file", filePath),
		zap.String("platform", platform),
		zap.String("user_id", userID.String()),
		zap.Bool("dry_run", dryRun),
	)

	// 获取解析器
	parser, err := parsers.GetParser(platform)
	if err != nil {
		return nil, fmt.Errorf("failed to get parser: %w", err)
	}

	// 读取文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 解析数据
	standardData, err := parser.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse data: %w", err)
	}

	// 验证数据
	if err := i.validator.Validate(standardData); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 转换数据
	conversations, messages, err := i.transformer.Transform(standardData, userID, platform)
	if err != nil {
		return nil, fmt.Errorf("transformation failed: %w", err)
	}

	result := &ImportResult{
		Platform:          platform,
		ConversationCount: len(conversations),
		MessageCount:      len(messages),
		SuccessCount:      len(conversations),
		ErrorCount:        0,
		Duration:          time.Since(startTime).String(),
	}

	// 如果不是dry run，写入数据库
	if !dryRun {
		if err := i.loader.Load(context.Background(), conversations, messages); err != nil {
			result.ErrorCount = 1
			result.Errors = append(result.Errors, err.Error())
			return result, fmt.Errorf("failed to load data: %w", err)
		}
	}

	log.Info("Import completed",
		zap.String("platform", platform),
		zap.Int("conversations", len(conversations)),
		zap.Int("messages", len(messages)),
		zap.String("duration", result.Duration),
	)

	return result, nil
}
