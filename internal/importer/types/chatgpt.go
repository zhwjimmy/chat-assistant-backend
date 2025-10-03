package types

// ChatGPTExportData ChatGPT导出数据结构（简略版本）
// 实际使用时需要根据ChatGPT的真实导出格式进行调整
type ChatGPTExportData struct {
	Conversations map[string]ChatGPTConversation `json:"conversations"`
	Mapping       map[string]string              `json:"mapping"`
	CurrentModel  string                         `json:"current_model"`
}

// ChatGPTConversation ChatGPT对话结构（简略版本）
type ChatGPTConversation struct {
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	CreateTime  float64                `json:"create_time"`
	UpdateTime  float64                `json:"update_time"`
	Mapping     map[string]interface{} `json:"mapping"`
	CurrentNode string                 `json:"current_node"`
}

// ChatGPTMessage ChatGPT消息结构（简略版本）
type ChatGPTMessage struct {
	ID       string                 `json:"id"`
	Role     string                 `json:"role"`
	Content  ChatGPTContent         `json:"content"`
	Created  float64                `json:"created"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ChatGPTContent ChatGPT消息内容结构
type ChatGPTContent struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}
