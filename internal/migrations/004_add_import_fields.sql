-- +goose Up
-- +goose StatementBegin
-- Add import-related fields to existing tables
-- This migration adds fields to support data import functionality
-- Add import-related fields to conversations table
ALTER TABLE conversations
ADD COLUMN IF NOT EXISTS import_job_id UUID,
    ADD COLUMN IF NOT EXISTS original_id VARCHAR(255),
    ADD COLUMN IF NOT EXISTS import_metadata JSONB;
-- Add import-related fields to messages table  
ALTER TABLE messages
ADD COLUMN IF NOT EXISTS original_id VARCHAR(255),
    ADD COLUMN IF NOT EXISTS import_metadata JSONB;
-- Create import_jobs table to track import operations
CREATE TABLE IF NOT EXISTS import_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    provider VARCHAR(50) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    -- pending, processing, completed, failed
    success_count INTEGER DEFAULT 0,
    error_count INTEGER DEFAULT 0,
    error_details TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE
);
-- Create import_details table for detailed import tracking
CREATE TABLE IF NOT EXISTS import_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES import_jobs(id),
    conversation_id UUID REFERENCES conversations(id),
    message_id UUID REFERENCES messages(id),
    original_data JSONB,
    status VARCHAR(20) NOT NULL,
    -- success, failed, skipped
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- Add indexes for better performance
CREATE INDEX IF NOT EXISTS idx_conversations_import_job_id ON conversations(import_job_id);
CREATE INDEX IF NOT EXISTS idx_conversations_original_id ON conversations(original_id);
CREATE INDEX IF NOT EXISTS idx_messages_original_id ON messages(original_id);
CREATE INDEX IF NOT EXISTS idx_import_jobs_user_id ON import_jobs(user_id);
CREATE INDEX IF NOT EXISTS idx_import_jobs_status ON import_jobs(status);
CREATE INDEX IF NOT EXISTS idx_import_details_job_id ON import_details(job_id);
-- Add trigger for import_jobs table
CREATE TRIGGER update_import_jobs_updated_at BEFORE
UPDATE ON import_jobs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd