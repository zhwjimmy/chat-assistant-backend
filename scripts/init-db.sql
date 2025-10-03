-- Initialize chat_assistant database
-- This script runs when the PostgreSQL container starts for the first time

-- Create the database if it doesn't exist (it should already exist from POSTGRES_DB)
-- But we'll add some initial setup here

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create a simple test table to verify the database is working
CREATE TABLE IF NOT EXISTS health_check (
    id SERIAL PRIMARY KEY,
    message TEXT NOT NULL DEFAULT 'Database is healthy',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Insert a test record
INSERT INTO health_check (message) VALUES ('Database initialized successfully') ON CONFLICT DO NOTHING;

-- Grant permissions
GRANT ALL PRIVILEGES ON DATABASE chat_assistant TO postgres;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO postgres;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO postgres;
