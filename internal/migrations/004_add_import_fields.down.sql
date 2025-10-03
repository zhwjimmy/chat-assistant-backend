-- +goose Down
-- +goose StatementBegin
-- Drop trigger for import_jobs table
DROP TRIGGER IF EXISTS update_import_jobs_updated_at ON import_jobs;
-- Drop indexes
DROP INDEX IF EXISTS idx_import_details_job_id;
DROP INDEX IF EXISTS idx_import_jobs_status;
DROP INDEX IF EXISTS idx_import_jobs_user_id;
DROP INDEX IF EXISTS idx_messages_original_id;
DROP INDEX IF EXISTS idx_conversations_original_id;
DROP INDEX IF EXISTS idx_conversations_import_job_id;
-- Drop tables
DROP TABLE IF EXISTS import_details;
DROP TABLE IF EXISTS import_jobs;
-- Remove import-related columns from existing tables
ALTER TABLE messages DROP COLUMN IF EXISTS import_metadata,
    DROP COLUMN IF EXISTS original_id;
ALTER TABLE conversations DROP COLUMN IF EXISTS import_metadata,
    DROP COLUMN IF EXISTS original_id,
    DROP COLUMN IF EXISTS import_job_id;
-- +goose StatementEnd