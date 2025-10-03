package types

// GeminiExportData Gemini导出数据结构（简略版本）
// 实际使用时需要根据Gemini的真实导出格式进行调整
type GeminiExportData struct {
	Conversations []GeminiConversation   `json:"conversations"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// GeminiConversation Gemini对话结构（简略版本）
type GeminiConversation struct {
	ID        string                 `json:"id"`
	Title     string                 `json:"title"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
	Messages  []GeminiMessage        `json:"messages"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// GeminiMessage Gemini消息结构（简略版本）
type GeminiMessage struct {
	ID        string                 `json:"id"`
	Role      string                 `json:"role"`
	Content   string                 `json:"content"`
	CreatedAt string                 `json:"created_at"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}
