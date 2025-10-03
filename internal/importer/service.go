package importer

import (
	"chat-assistant-backend/internal/config"
	"chat-assistant-backend/internal/importer/parsers"
)

// Service 导入服务
type Service struct {
	config   *config.Config
	importer *Importer
}

// NewService 创建导入服务
func NewService(cfg *config.Config) *Service {
	return &Service{
		config:   cfg,
		importer: NewImporter(cfg),
	}
}

// Import 执行导入
func (s *Service) Import(filePath, platform, userID string, dryRun bool) (*ImportResult, error) {
	return s.importer.Import(filePath, platform, userID, dryRun)
}

// GetSupportedPlatforms 获取支持的平台列表
func (s *Service) GetSupportedPlatforms() []string {
	return parsers.GetSupportedPlatforms()
}
