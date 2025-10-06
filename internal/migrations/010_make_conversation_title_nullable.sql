-- +goose Up
-- +goose StatementBegin
-- Make conversation title field nullable
ALTER TABLE conversations
ALTER COLUMN title DROP NOT NULL;
-- Update column comment
COMMENT ON COLUMN conversations.title IS '对话标题 (可选)';
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- Make conversation title field not null again
-- First, update any NULL values to empty string
UPDATE conversations
SET title = ''
WHERE title IS NULL;
-- Then add NOT NULL constraint back
ALTER TABLE conversations
ALTER COLUMN title
SET NOT NULL;
-- Update column comment back
COMMENT ON COLUMN conversations.title IS '对话标题';
-- +goose StatementEnd