-- +goose Up
-- +goose StatementBegin
-- Add unique constraints to prevent duplicate imports
-- Unique constraint for conversations: user_id + source_id should be unique
ALTER TABLE conversations
ADD CONSTRAINT uk_conversations_user_source UNIQUE (user_id, source_id);
-- Unique constraint for messages: conversation_id + source_id should be unique
ALTER TABLE messages
ADD CONSTRAINT uk_messages_conversation_source UNIQUE (conversation_id, source_id);
-- Add comments
COMMENT ON CONSTRAINT uk_conversations_user_source ON conversations IS '确保同一用户的相同source_id不会重复导入';
COMMENT ON CONSTRAINT uk_messages_conversation_source ON messages IS '确保同一对话中的相同source_id消息不会重复导入';
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- Remove unique constraints
ALTER TABLE conversations DROP CONSTRAINT IF EXISTS uk_conversations_user_source;
ALTER TABLE messages DROP CONSTRAINT IF EXISTS uk_messages_conversation_source;
-- +goose StatementEnd