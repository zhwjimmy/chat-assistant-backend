-- +goose Down
-- +goose StatementBegin
-- Drop triggers
DROP TRIGGER IF EXISTS update_messages_updated_at ON messages;
DROP TRIGGER IF EXISTS update_conversations_updated_at ON conversations;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();
-- Drop indexes
DROP INDEX IF EXISTS idx_messages_role;
DROP INDEX IF EXISTS idx_messages_deleted_at;
DROP INDEX IF EXISTS idx_messages_conversation_id;
DROP INDEX IF EXISTS idx_conversations_provider;
DROP INDEX IF EXISTS idx_conversations_deleted_at;
DROP INDEX IF EXISTS idx_conversations_user_id;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_username;
-- Drop tables
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS conversations;
DROP TABLE IF EXISTS users;
-- Drop UUID extension (optional, as it might be used by other tables)
-- DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd