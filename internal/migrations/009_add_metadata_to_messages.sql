-- +goose Up
-- +goose StatementBegin
-- Add metadata field to messages table
ALTER TABLE messages
ADD COLUMN metadata TEXT;
-- Add column comment
COMMENT ON COLUMN messages.metadata IS '可选元信息';
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- Remove metadata field from messages table
ALTER TABLE messages DROP COLUMN IF EXISTS metadata;
-- +goose StatementEnd