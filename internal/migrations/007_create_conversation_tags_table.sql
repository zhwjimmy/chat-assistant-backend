-- +goose Up
-- +goose StatementBegin
-- Create conversation_tags junction table for many-to-many relationship
CREATE TABLE conversation_tags (
    conversation_id UUID NOT NULL,
    tag_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (conversation_id, tag_id),
    FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);
-- Create indexes for better query performance
CREATE INDEX idx_conversation_tags_conversation_id ON conversation_tags(conversation_id);
CREATE INDEX idx_conversation_tags_tag_id ON conversation_tags(tag_id);
-- Add table and column comments
COMMENT ON TABLE conversation_tags IS '对话标签关联表，实现多对多关系';
COMMENT ON COLUMN conversation_tags.conversation_id IS '对话ID';
COMMENT ON COLUMN conversation_tags.tag_id IS '标签ID';
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- Drop conversation_tags table
DROP TABLE IF EXISTS conversation_tags CASCADE;
-- +goose StatementEnd