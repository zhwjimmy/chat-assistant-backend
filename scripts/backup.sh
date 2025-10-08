#!/bin/bash

# Database backup script
set -e

# Read configuration from environment variables
CONTAINER_NAME="${DB_CONTAINER_NAME:-chat-assistant-postgres}"
DB_NAME="${DB_NAME:-chat_assistant}"
DB_USER="${DB_USER:-postgres}"
BACKUP_DIR="${BACKUP_DIR:-./tmp}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/chat_assistant_${TIMESTAMP}.dump.gz"

# Create backup directory
mkdir -p "$BACKUP_DIR"

echo "Creating database backup..."
echo "Container: $CONTAINER_NAME"
echo "Database: $DB_NAME"
echo "Backup file: $BACKUP_FILE"

# Perform backup
docker exec "$CONTAINER_NAME" pg_dump -U "$DB_USER" -d "$DB_NAME" -Fc | gzip > "$BACKUP_FILE"

echo "Backup completed: $BACKUP_FILE"
