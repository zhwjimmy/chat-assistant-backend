-- +goose Up
-- +goose StatementBegin
-- Add metadata field to conversations table
ALTER TABLE conversations
ADD COLUMN metadata TEXT;
-- Add column comment
COMMENT ON COLUMN conversations.metadata IS '可选元信息';
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- Remove metadata field from conversations table
ALTER TABLE conversations DROP COLUMN IF EXISTS metadata;
-- +goose StatementEnd