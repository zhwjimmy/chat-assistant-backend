package gemini

import (
	"encoding/json"
	"fmt"

	"chat-assistant-backend/internal/importer/types"
)

// Parser Gemini解析器
type Parser struct{}

// NewParser 创建Gemini解析器
func NewParser() *Parser {
	return &Parser{}
}

// Platform 返回平台名称
func (p *Parser) Platform() string {
	return "gemini"
}

// Parse 解析Gemini导出数据
func (p *Parser) Parse(data []byte) (*types.StandardFormat, error) {
	// 简略实现 - 实际需要根据Gemini的真实导出格式调整
	var geminiData GeminiExportData
	if err := json.Unmarshal(data, &geminiData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Gemini data: %w", err)
	}

	// 转换为标准化格式
	standardData := &types.StandardFormat{
		Conversations: make([]*types.StandardConversation, 0),
	}

	// 简略转换逻辑 - 实际需要根据真实格式调整
	for _, conv := range geminiData.Conversations {
		stdConv := &types.StandardConversation{
			ID:       conv.ID,
			Title:    conv.Title,
			Provider: "gemini",
			Model:    "gemini-pro", // 默认模型，实际应该从数据中获取
			Messages: make([]*types.StandardMessage, 0),
		}

		// 简略消息转换
		for _, msg := range conv.Messages {
			stdMsg := &types.StandardMessage{
				Role:    msg.Role,
				Content: msg.Content,
			}
			stdConv.Messages = append(stdConv.Messages, stdMsg)
		}

		standardData.Conversations = append(standardData.Conversations, stdConv)
	}

	return standardData, nil
}

// GeminiExportData Gemini导出数据结构（简略版本）
type GeminiExportData struct {
	Conversations []GeminiConversation `json:"conversations"`
}

// GeminiConversation Gemini对话结构（简略版本）
type GeminiConversation struct {
	ID       string          `json:"id"`
	Title    string          `json:"title"`
	Messages []GeminiMessage `json:"messages"`
}

// GeminiMessage Gemini消息结构（简略版本）
type GeminiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
