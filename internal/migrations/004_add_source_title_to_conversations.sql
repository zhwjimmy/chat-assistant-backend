-- +goose Up
-- +goose StatementBegin
-- Add source_title column to conversations table
ALTER TABLE conversations
ADD COLUMN source_title VARCHAR(500);
-- Add comment for the new column
COMMENT ON COLUMN conversations.source_title IS '原始数据中的标题，用于对比和调试';
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- Remove source_title column from conversations table
ALTER TABLE conversations DROP COLUMN IF EXISTS source_title;
-- +goose StatementEnd