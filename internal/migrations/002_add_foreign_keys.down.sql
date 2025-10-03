-- +goose Down
-- +goose StatementBegin
-- Drop foreign key constraints
ALTER TABLE messages DROP CONSTRAINT IF EXISTS fk_messages_conversation_id;
ALTER TABLE conversations DROP CONSTRAINT IF EXISTS fk_conversations_user_id;
-- +goose StatementEnd