package types

// ClaudeExportData Claude导出数据结构（简略版本）
// 实际使用时需要根据Claude的真实导出格式进行调整
type ClaudeExportData struct {
	Conversations []ClaudeConversation   `json:"conversations"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// ClaudeConversation Claude对话结构（简略版本）
type ClaudeConversation struct {
	ID        string                 `json:"id"`
	Title     string                 `json:"title"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
	Messages  []ClaudeMessage        `json:"messages"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ClaudeMessage Claude消息结构（简略版本）
type ClaudeMessage struct {
	ID        string                 `json:"id"`
	Role      string                 `json:"role"`
	Content   string                 `json:"content"`
	CreatedAt string                 `json:"created_at"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}
