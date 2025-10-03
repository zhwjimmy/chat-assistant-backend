-- +goose Down
-- +goose StatementBegin
-- Drop messages table and related objects
DROP TRIGGER IF EXISTS update_messages_updated_at ON messages;
DROP TABLE IF EXISTS messages CASCADE;
-- +goose StatementEnd