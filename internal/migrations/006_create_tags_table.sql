-- +goose Up
-- +goose StatementBegin
-- Create tags table
CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(500) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
-- Create indexes for tags table
CREATE INDEX idx_tags_name ON tags(name);
CREATE INDEX idx_tags_deleted_at ON tags(deleted_at);
-- Create trigger for automatic timestamp updates
CREATE TRIGGER update_tags_updated_at BEFORE
UPDATE ON tags FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- Add table and column comments
COMMENT ON TABLE tags IS '标签表';
COMMENT ON COLUMN tags.name IS '标签名称';
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- Drop tags table and related objects
DROP TRIGGER IF EXISTS update_tags_updated_at ON tags;
DROP TABLE IF EXISTS tags CASCADE;
-- +goose StatementEnd