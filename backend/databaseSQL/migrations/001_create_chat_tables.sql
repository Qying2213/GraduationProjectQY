-- Migration: 001_create_chat_tables.sql
-- Description: Create conversations and chat_messages tables for chat functionality
-- Requirements: 8.5 (Chat Message Persistence)
-- 
-- This migration is idempotent - safe to run multiple times
-- All CREATE statements use IF NOT EXISTS to prevent errors on re-run

-- ============================================================================
-- 会话表 (Conversations Table)
-- ============================================================================
-- Stores chat conversations between two users (candidate and HR)
-- Each conversation is unique per participant pair

CREATE TABLE IF NOT EXISTS conversations (
    id SERIAL PRIMARY KEY,
    participant_a INTEGER REFERENCES users(id) NOT NULL,
    participant_b INTEGER REFERENCES users(id) NOT NULL,
    last_message_id INTEGER,
    last_message_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(participant_a, participant_b)
);

-- Add comment for documentation
COMMENT ON TABLE conversations IS '聊天会话表 - 存储用户之间的聊天会话';
COMMENT ON COLUMN conversations.participant_a IS '参与者A的用户ID';
COMMENT ON COLUMN conversations.participant_b IS '参与者B的用户ID';
COMMENT ON COLUMN conversations.last_message_id IS '最后一条消息的ID';
COMMENT ON COLUMN conversations.last_message_at IS '最后一条消息的时间';

-- ============================================================================
-- 聊天消息表 (Chat Messages Table)
-- ============================================================================
-- Stores individual chat messages within conversations

CREATE TABLE IF NOT EXISTS chat_messages (
    id SERIAL PRIMARY KEY,
    conversation_id INTEGER REFERENCES conversations(id) NOT NULL,
    sender_id INTEGER REFERENCES users(id) NOT NULL,
    content TEXT NOT NULL,
    message_type VARCHAR(20) DEFAULT 'text',
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add comment for documentation
COMMENT ON TABLE chat_messages IS '聊天消息表 - 存储会话中的消息';
COMMENT ON COLUMN chat_messages.conversation_id IS '所属会话ID';
COMMENT ON COLUMN chat_messages.sender_id IS '发送者用户ID';
COMMENT ON COLUMN chat_messages.content IS '消息内容';
COMMENT ON COLUMN chat_messages.message_type IS '消息类型: text, image, file';
COMMENT ON COLUMN chat_messages.is_read IS '是否已读';

-- ============================================================================
-- 索引 (Indexes)
-- ============================================================================
-- Performance indexes for common query patterns

-- Index for finding conversations by participants
-- Used when: Looking up conversation between two users
CREATE INDEX IF NOT EXISTS idx_conversations_participants 
    ON conversations(participant_a, participant_b);

-- Index for sorting conversations by last message time
-- Used when: Displaying conversation list sorted by recent activity
CREATE INDEX IF NOT EXISTS idx_conversations_last_message 
    ON conversations(last_message_at DESC);

-- Index for fetching messages in a conversation ordered by time
-- Used when: Loading chat history with pagination
CREATE INDEX IF NOT EXISTS idx_chat_messages_conversation 
    ON chat_messages(conversation_id, created_at DESC);

-- Partial index for unread messages (only indexes unread messages)
-- Used when: Counting unread messages, marking messages as read
CREATE INDEX IF NOT EXISTS idx_chat_messages_unread 
    ON chat_messages(conversation_id, is_read) 
    WHERE is_read = FALSE;

-- ============================================================================
-- Add foreign key for last_message_id after chat_messages table exists
-- ============================================================================
-- Note: This is done separately because chat_messages references conversations
-- and we need both tables to exist first

DO $$
BEGIN
    -- Check if the foreign key constraint already exists
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_conversations_last_message'
        AND table_name = 'conversations'
    ) THEN
        -- Add foreign key constraint for last_message_id
        ALTER TABLE conversations 
        ADD CONSTRAINT fk_conversations_last_message 
        FOREIGN KEY (last_message_id) REFERENCES chat_messages(id);
    END IF;
END $$;

-- ============================================================================
-- Migration complete
-- ============================================================================
