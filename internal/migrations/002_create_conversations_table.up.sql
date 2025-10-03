-- +goose Up
-- +goose StatementBegin
-- Create conversations table
CREATE TABLE conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    title VARCHAR(500) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    model VARCHAR(50),
    source_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
-- Create indexes for conversations table
CREATE INDEX idx_conversations_user_id ON conversations(user_id);
CREATE INDEX idx_conversations_deleted_at ON conversations(deleted_at);
CREATE INDEX idx_conversations_provider ON conversations(provider);
CREATE INDEX idx_conversations_source_id ON conversations(source_id);
-- Create trigger for automatic timestamp updates
CREATE TRIGGER update_conversations_updated_at BEFORE
UPDATE ON conversations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- Add foreign key constraints
ALTER TABLE conversations
ADD CONSTRAINT fk_conversations_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
-- Add table and column comments
COMMENT ON TABLE conversations IS '对话表';
COMMENT ON COLUMN conversations.user_id IS '用户ID';
COMMENT ON COLUMN conversations.title IS '对话标题';
COMMENT ON COLUMN conversations.provider IS 'AI提供商 (openai, gemini, local等)';
COMMENT ON COLUMN conversations.model IS 'AI模型 (gpt-4, gemini-pro, llama-3等)';
COMMENT ON COLUMN conversations.source_id IS '原始数据中的ID，用于关联导入内容';
-- +goose StatementEnd