#!/bin/bash

# Database restore script
set -e

# Read configuration from environment variables
CONTAINER_NAME="${DB_CONTAINER_NAME:-chat-assistant-postgres}"
DB_NAME="${DB_NAME:-chat_assistant}"
DB_USER="${DB_USER:-postgres}"
BACKUP_DIR="${BACKUP_DIR:-./tmp}"

# Get backup file
if [ $# -eq 0 ]; then
    # Use latest backup if no file specified
    BACKUP_FILE=$(ls -t "$BACKUP_DIR"/*.dump.gz 2>/dev/null | head -1)
    if [ -z "$BACKUP_FILE" ]; then
        echo "No backup files found in $BACKUP_DIR"
        exit 1
    fi
    echo "Using latest backup: $BACKUP_FILE"
else
    BACKUP_FILE="$1"
fi

# Validate backup file
if [ ! -f "$BACKUP_FILE" ]; then
    echo "Backup file not found: $BACKUP_FILE"
    exit 1
fi

echo "Restoring database..."
echo "Container: $CONTAINER_NAME"
echo "Database: $DB_NAME"
echo "Backup file: $BACKUP_FILE"

# Drop and recreate database
docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -c "DROP DATABASE IF EXISTS $DB_NAME;"
docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -c "CREATE DATABASE $DB_NAME;"

# Restore backup
gunzip -c "$BACKUP_FILE" | docker exec -i "$CONTAINER_NAME" pg_restore -U "$DB_USER" -d "$DB_NAME"

echo "Database restored successfully!"
