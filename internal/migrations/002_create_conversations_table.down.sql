-- +goose Down
-- +goose StatementBegin
-- Drop conversations table and related objects
DROP TRIGGER IF EXISTS update_conversations_updated_at ON conversations;
DROP TABLE IF EXISTS conversations CASCADE;
-- +goose StatementEnd