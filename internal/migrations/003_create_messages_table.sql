-- +goose Up
-- +goose StatementBegin
-- Create messages table
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL,
    role VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    source_id VARCHAR(255) NOT NULL,
    source_content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
-- Create indexes for messages table
CREATE INDEX idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX idx_messages_deleted_at ON messages(deleted_at);
CREATE INDEX idx_messages_role ON messages(role);
CREATE INDEX idx_messages_source_id ON messages(source_id);
-- Create trigger for automatic timestamp updates
CREATE TRIGGER update_messages_updated_at BEFORE
UPDATE ON messages FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- Add foreign key constraints
ALTER TABLE messages
ADD CONSTRAINT fk_messages_conversation_id FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE;
-- Add table and column comments
COMMENT ON TABLE messages IS '消息表';
COMMENT ON COLUMN messages.conversation_id IS '对话ID';
COMMENT ON COLUMN messages.role IS '消息角色 (user, assistant, system)';
COMMENT ON COLUMN messages.content IS '消息内容';
COMMENT ON COLUMN messages.source_id IS '原始数据中的ID，用于关联导入内容';
COMMENT ON COLUMN messages.source_content IS '原始数据中的内容，用于对比和调试';
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- Drop messages table and related objects
DROP TRIGGER IF EXISTS update_messages_updated_at ON messages;
-- Drop foreign key constraint first
ALTER TABLE messages DROP CONSTRAINT IF EXISTS fk_messages_conversation_id;
DROP TABLE IF EXISTS messages CASCADE;
-- +goose StatementEnd