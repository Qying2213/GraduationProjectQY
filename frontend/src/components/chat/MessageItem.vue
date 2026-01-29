<template>
  <div
    class="message-item"
    :class="{ 'is-self': isSelf }"
  >
    <div class="message-content">
      <div class="message-bubble">
        <span class="message-text">{{ message.content }}</span>
      </div>
      <div class="message-meta">
        <span class="message-time">{{ formatTime(message.created_at) }}</span>
        <span v-if="isSelf" class="message-status">
          <el-icon v-if="message.is_read" class="read-icon"><Check /></el-icon>
          <el-icon v-else class="unread-icon"><Clock /></el-icon>
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
/**
 * MessageItem.vue - 消息项组件
 * Requirements: 8.3 (Display messages), 9.2 (Message display)
 * 
 * Displays a single chat message with different styles for sent vs received messages.
 * Shows message content, time, and read status.
 */
import { Check, Clock } from '@element-plus/icons-vue'
import type { ChatMessage } from '@/api/chat'

defineProps<{
  /** The chat message to display */
  message: ChatMessage
  /** Whether this message was sent by the current user */
  isSelf: boolean
}>()

/**
 * Format timestamp to display time
 * @param timestamp - ISO timestamp string
 * @returns Formatted time string
 */
const formatTime = (timestamp: string): string => {
  const date = new Date(timestamp)
  const now = new Date()
  const isToday = date.toDateString() === now.toDateString()
  
  const timeStr = date.toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit'
  })
  
  if (isToday) {
    return timeStr
  }
  
  // Check if yesterday
  const yesterday = new Date(now)
  yesterday.setDate(yesterday.getDate() - 1)
  if (date.toDateString() === yesterday.toDateString()) {
    return `昨天 ${timeStr}`
  }
  
  // Show date for older messages
  return date.toLocaleDateString('zh-CN', {
    month: '2-digit',
    day: '2-digit'
  }) + ' ' + timeStr
}
</script>

<style scoped lang="scss">
.message-item {
  display: flex;
  margin-bottom: 16px;
  padding: 0 16px;

  &.is-self {
    justify-content: flex-end;

    .message-bubble {
      background: #0ea5e9;
      color: #fff;
      border-radius: 16px 16px 4px 16px;
    }

    .message-meta {
      justify-content: flex-end;
    }
  }

  &:not(.is-self) {
    justify-content: flex-start;

    .message-bubble {
      background: #f1f5f9;
      color: #1e293b;
      border-radius: 16px 16px 16px 4px;
    }

    .message-meta {
      justify-content: flex-start;
    }
  }
}

.message-content {
  max-width: 70%;
  display: flex;
  flex-direction: column;
}

.message-bubble {
  padding: 10px 14px;
  word-break: break-word;
  line-height: 1.5;
  font-size: 14px;
}

.message-text {
  white-space: pre-wrap;
}

.message-meta {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-top: 4px;
  padding: 0 4px;
}

.message-time {
  font-size: 12px;
  color: #94a3b8;
}

.message-status {
  display: flex;
  align-items: center;
  
  .read-icon {
    color: #0ea5e9;
    font-size: 14px;
  }
  
  .unread-icon {
    color: #94a3b8;
    font-size: 14px;
  }
}
</style>
