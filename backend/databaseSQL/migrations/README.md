# Database Migrations

This directory contains SQL migration scripts for the 智能人才运营平台 database.

## Migration Files

| File | Description | Requirements |
|------|-------------|--------------|
| `001_create_chat_tables.sql` | Creates conversations and chat_messages tables for chat functionality | 8.5 |

## Running Migrations

### Prerequisites

- PostgreSQL database server running
- Database connection configured (see `backend/.env`)

### Using psql

```bash
# Connect to your database and run the migration
psql -h localhost -U postgres -d talent_platform -f backend/database/migrations/001_create_chat_tables.sql
```

### Using Docker

If running PostgreSQL in Docker:

```bash
# Copy migration file to container and execute
docker cp backend/database/migrations/001_create_chat_tables.sql postgres_container:/tmp/
docker exec -it postgres_container psql -U postgres -d talent_platform -f /tmp/001_create_chat_tables.sql
```

### Using Go Application

The migrations can also be executed programmatically using GORM's AutoMigrate or raw SQL execution.

## Migration Design Principles

1. **Idempotent**: All migrations use `IF NOT EXISTS` clauses, making them safe to run multiple times
2. **Documented**: Each migration includes comments explaining the purpose and requirements
3. **Indexed**: Performance indexes are created for common query patterns
4. **Ordered**: Migration files are numbered (001, 002, etc.) to ensure correct execution order

## Tables Created

### conversations
- Stores chat conversations between two users
- Unique constraint on participant pair to prevent duplicate conversations
- Tracks last message for efficient conversation list display

### chat_messages
- Stores individual messages within conversations
- Supports different message types (text, image, file)
- Tracks read status for unread count functionality

## Indexes

| Index | Table | Purpose |
|-------|-------|---------|
| `idx_conversations_participants` | conversations | Fast lookup by participant pair |
| `idx_conversations_last_message` | conversations | Sort conversations by recent activity |
| `idx_chat_messages_conversation` | chat_messages | Fetch messages with pagination |
| `idx_chat_messages_unread` | chat_messages | Count/query unread messages efficiently |

## Rollback

To rollback this migration (remove the tables):

```sql
-- WARNING: This will delete all chat data!
DROP TABLE IF EXISTS chat_messages CASCADE;
DROP TABLE IF EXISTS conversations CASCADE;
```
