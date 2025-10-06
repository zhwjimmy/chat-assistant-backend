-- +goose Up
-- +goose StatementBegin
-- Make message content field nullable
ALTER TABLE messages
ALTER COLUMN content DROP NOT NULL;
-- Update column comment
COMMENT ON COLUMN messages.content IS '消息内容 (可选)';
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- Make message content field not null again
-- First, update any NULL values to empty string
UPDATE messages
SET content = ''
WHERE content IS NULL;
-- Then add NOT NULL constraint back
ALTER TABLE messages
ALTER COLUMN content
SET NOT NULL;
-- Update column comment back
COMMENT ON COLUMN messages.content IS '消息内容';
-- +goose StatementEnd