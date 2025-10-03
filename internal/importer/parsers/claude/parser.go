package claude

import (
	"encoding/json"
	"fmt"
	"time"

	"chat-assistant-backend/internal/importer/types"
)

// Parser Claude解析器
type Parser struct{}

// NewParser 创建Claude解析器
func NewParser() *Parser {
	return &Parser{}
}

// Platform 返回平台名称
func (p *Parser) Platform() string {
	return "claude"
}

// Parse 解析Claude导出数据
func (p *Parser) Parse(data []byte) (*types.StandardFormat, error) {
	var claudeData types.ClaudeExportData
	if err := json.Unmarshal(data, &claudeData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Claude data: %w", err)
	}

	// 转换为标准化格式
	standardData := &types.StandardFormat{
		Conversations: make([]*types.StandardConversation, 0),
	}

	// 转换对话数据
	for _, conv := range claudeData {
		// 解析时间
		createdAt, _ := time.Parse(time.RFC3339, conv.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339, conv.UpdatedAt)

		stdConv := &types.StandardConversation{
			ID:        conv.UUID,
			Title:     conv.Name,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			Provider:  "claude",
			Model:     "claude-3", // 默认模型，实际应该从数据中获取
			Messages:  make([]*types.StandardMessage, 0),
			Metadata: map[string]interface{}{
				"summary": conv.Summary,
				"account": conv.Account,
			},
		}

		// 转换消息数据
		for _, msg := range conv.ChatMessages {
			// 解析消息时间
			msgCreatedAt, _ := time.Parse(time.RFC3339, msg.CreatedAt)
			msgUpdatedAt, _ := time.Parse(time.RFC3339, msg.UpdatedAt)

			// 确定角色
			role := "user"
			if msg.Sender == "assistant" {
				role = "assistant"
			}

			// 提取消息内容
			content := msg.Text
			if content == "" && len(msg.Content) > 0 {
				content = msg.Content[0].Text
			}

			stdMsg := &types.StandardMessage{
				ID:        msg.UUID,
				Role:      role,
				Content:   content,
				CreatedAt: msgCreatedAt,
				Metadata: map[string]interface{}{
					"updated_at":  msgUpdatedAt,
					"attachments": msg.Attachments,
					"files":       msg.Files,
					"content":     msg.Content,
				},
			}
			stdConv.Messages = append(stdConv.Messages, stdMsg)
		}

		standardData.Conversations = append(standardData.Conversations, stdConv)
	}

	return standardData, nil
}
