package chatgpt

import (
	"encoding/json"
	"fmt"

	"chat-assistant-backend/internal/importer/types"
)

// Parser ChatGPT解析器
type Parser struct{}

// NewParser 创建ChatGPT解析器
func NewParser() *Parser {
	return &Parser{}
}

// Platform 返回平台名称
func (p *Parser) Platform() string {
	return "chatgpt"
}

// Parse 解析ChatGPT导出数据
func (p *Parser) Parse(data []byte) (*types.StandardFormat, error) {
	// 简略实现 - 实际需要根据ChatGPT的真实导出格式调整
	var chatgptData ChatGPTExportData
	if err := json.Unmarshal(data, &chatgptData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ChatGPT data: %w", err)
	}

	// 转换为标准化格式
	standardData := &types.StandardFormat{
		Conversations: make([]*types.StandardConversation, 0),
	}

	// 简略转换逻辑 - 实际需要根据真实格式调整
	for _, conv := range chatgptData.Conversations {
		stdConv := &types.StandardConversation{
			ID:       conv.ID,
			Title:    conv.Title,
			Provider: "chatgpt",
			Model:    "gpt-4", // 默认模型，实际应该从数据中获取
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

// ChatGPTExportData ChatGPT导出数据结构（简略版本）
type ChatGPTExportData struct {
	Conversations []ChatGPTConversation `json:"conversations"`
}

// ChatGPTConversation ChatGPT对话结构（简略版本）
type ChatGPTConversation struct {
	ID       string           `json:"id"`
	Title    string           `json:"title"`
	Messages []ChatGPTMessage `json:"messages"`
}

// ChatGPTMessage ChatGPT消息结构（简略版本）
type ChatGPTMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
