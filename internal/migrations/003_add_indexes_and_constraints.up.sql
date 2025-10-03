-- +goose Up
-- +goose StatementBegin
-- Add additional indexes for better performance
CREATE INDEX IF NOT EXISTS idx_conversations_created_at ON conversations(created_at);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);
-- Add check constraints for data validation
ALTER TABLE conversations
ADD CONSTRAINT chk_conversations_provider CHECK (
        provider IN ('openai', 'gemini', 'claude', 'local', 'other')
    );
ALTER TABLE messages
ADD CONSTRAINT chk_messages_role CHECK (role IN ('user', 'assistant', 'system'));
-- Add not null constraints for required fields
ALTER TABLE conversations
ALTER COLUMN user_id
SET NOT NULL,
    ALTER COLUMN title
SET NOT NULL,
    ALTER COLUMN provider
SET NOT NULL;
ALTER TABLE messages
ALTER COLUMN conversation_id
SET NOT NULL,
    ALTER COLUMN role
SET NOT NULL,
    ALTER COLUMN content
SET NOT NULL;
-- +goose StatementEnd