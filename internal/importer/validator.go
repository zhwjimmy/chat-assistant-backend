package importer

import (
	"fmt"

	"chat-assistant-backend/internal/importer/types"
)

// Validator 数据验证器
type Validator struct{}

// NewValidator 创建验证器
func NewValidator() *Validator {
	return &Validator{}
}

// Validate 验证标准化数据
func (v *Validator) Validate(data *types.StandardFormat) error {
	if data == nil {
		return fmt.Errorf("data is nil")
	}

	if len(data.Conversations) == 0 {
		return fmt.Errorf("no conversations found")
	}

	// 验证每个对话
	for i, conv := range data.Conversations {
		if err := v.validateConversation(conv, i); err != nil {
			return fmt.Errorf("conversation %d validation failed: %w", i, err)
		}
	}

	return nil
}

// validateConversation 验证单个对话
func (v *Validator) validateConversation(conv *types.StandardConversation, index int) error {
	if conv == nil {
		return fmt.Errorf("conversation is nil")
	}

	if conv.ID == "" {
		return fmt.Errorf("conversation ID is empty")
	}

	if conv.Title == "" {
		// return fmt.Errorf("conversation title is empty")
	}

	if conv.Provider == "" {
		return fmt.Errorf("conversation provider is empty")
	}

	// 验证消息
	for j, msg := range conv.Messages {
		if err := v.validateMessage(msg, j); err != nil {
			return fmt.Errorf("message %d validation failed: %w", j, err)
		}
	}

	return nil
}

// validateMessage 验证单个消息
func (v *Validator) validateMessage(msg *types.StandardMessage, index int) error {
	if msg == nil {
		return fmt.Errorf("message is nil")
	}

	if msg.Role == "" {
		return fmt.Errorf("message role is empty")
	}

	if msg.Content == "" {
		// return fmt.Errorf("message content is empty")
	}

	// 验证角色
	validRoles := map[string]bool{
		"user":      true,
		"assistant": true,
		"system":    true,
	}

	if !validRoles[msg.Role] {
		return fmt.Errorf("invalid message role: %s", msg.Role)
	}

	return nil
}
