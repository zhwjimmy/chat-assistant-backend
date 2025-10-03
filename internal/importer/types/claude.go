package types

// ClaudeExportData Claude导出数据结构
type ClaudeExportData []ClaudeConversation

// ClaudeConversation Claude对话结构
type ClaudeConversation struct {
	UUID         string                 `json:"uuid"`
	Name         string                 `json:"name"`
	Summary      string                 `json:"summary"`
	CreatedAt    string                 `json:"created_at"`
	UpdatedAt    string                 `json:"updated_at"`
	Account      ClaudeAccount          `json:"account"`
	ChatMessages []ClaudeMessage        `json:"chat_messages"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// ClaudeAccount Claude账户信息
type ClaudeAccount struct {
	UUID string `json:"uuid"`
}

// ClaudeMessage Claude消息结构
type ClaudeMessage struct {
	UUID        string                 `json:"uuid"`
	Text        string                 `json:"text"`
	Content     []ClaudeContent        `json:"content"`
	Sender      string                 `json:"sender"`
	CreatedAt   string                 `json:"created_at"` // 2025-09-22T09:17:21.803710Z
	UpdatedAt   string                 `json:"updated_at"`
	Attachments []interface{}          `json:"attachments"`
	Files       []interface{}          `json:"files"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ClaudeContent Claude消息内容结构
type ClaudeContent struct {
	StartTimestamp string                 `json:"start_timestamp"` // 2025-09-22T09:17:21.803710Z
	StopTimestamp  string                 `json:"stop_timestamp"`
	Flags          interface{}            `json:"flags"`
	Type           string                 `json:"type"`
	Text           string                 `json:"text"`
	Citations      []interface{}          `json:"citations"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}
