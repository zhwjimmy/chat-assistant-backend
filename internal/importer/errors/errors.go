package errors

import (
	"fmt"
)

// ImportError 导入错误
type ImportError struct {
	Type       string `json:"type"`        // conversation, message
	OriginalID string `json:"original_id"` // 原始ID
	Message    string `json:"message"`     // 错误信息
}

// Error 实现error接口
func (e *ImportError) Error() string {
	return fmt.Sprintf("%s (%s): %s", e.Type, e.OriginalID, e.Message)
}

// NewImportError 创建导入错误
func NewImportError(errorType, originalID, message string) *ImportError {
	return &ImportError{
		Type:       errorType,
		OriginalID: originalID,
		Message:    message,
	}
}
