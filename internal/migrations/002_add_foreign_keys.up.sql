-- +goose Up
-- +goose StatementBegin
-- Add foreign key constraints
ALTER TABLE conversations
ADD CONSTRAINT fk_conversations_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE messages
ADD CONSTRAINT fk_messages_conversation_id FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE;
-- +goose StatementEnd