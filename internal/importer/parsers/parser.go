package parsers

import (
	"fmt"

	"chat-assistant-backend/internal/importer/types"
)

// Parser 解析器接口
type Parser interface {
	Parse(data []byte) (*types.StandardFormat, error)
	Platform() string
}

// Registry 解析器注册中心
type Registry struct {
	parsers map[string]Parser
}

var registry = &Registry{
	parsers: make(map[string]Parser),
}

// Register 注册解析器
func Register(parser Parser) {
	registry.parsers[parser.Platform()] = parser
}

// GetParser 获取解析器
func GetParser(platform string) (Parser, error) {
	parser, exists := registry.parsers[platform]
	if !exists {
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
	return parser, nil
}

// GetSupportedPlatforms 获取支持的平台列表
func GetSupportedPlatforms() []string {
	platforms := make([]string, 0, len(registry.parsers))
	for platform := range registry.parsers {
		platforms = append(platforms, platform)
	}
	return platforms
}
