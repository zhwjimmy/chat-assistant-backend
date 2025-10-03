package parsers

import (
	chatgptParser "chat-assistant-backend/internal/importer/parsers/chatgpt"
	claudeParser "chat-assistant-backend/internal/importer/parsers/claude"
	geminiParser "chat-assistant-backend/internal/importer/parsers/gemini"
)

// RegisterAll 注册所有解析器
func RegisterAll() {
	Register(chatgptParser.NewParser())
	Register(claudeParser.NewParser())
	Register(geminiParser.NewParser())
}

// RegisterChatGPT 注册ChatGPT解析器
func RegisterChatGPT() {
	Register(chatgptParser.NewParser())
}

// RegisterClaude 注册Claude解析器
func RegisterClaude() {
	Register(claudeParser.NewParser())
}

// RegisterGemini 注册Gemini解析器
func RegisterGemini() {
	Register(geminiParser.NewParser())
}
