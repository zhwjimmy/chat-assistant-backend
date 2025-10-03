-- +goose Down
-- +goose StatementBegin
-- Drop check constraints
ALTER TABLE messages DROP CONSTRAINT IF EXISTS chk_messages_role;
ALTER TABLE conversations DROP CONSTRAINT IF EXISTS chk_conversations_provider;
-- Drop additional indexes
DROP INDEX IF EXISTS idx_messages_created_at;
DROP INDEX IF EXISTS idx_conversations_created_at;
-- +goose StatementEnd