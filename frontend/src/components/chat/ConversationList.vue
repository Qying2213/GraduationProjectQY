<template>
  <div class="conversation-list">
    <div
      v-for="conversation in conversations"
      :key="conversation.id"
      class="conversation-item"
      :class="{ 'is-selected': conversation.id === selectedId }"
      @click="handleSelect(conversation)"
    >
      <!-- Avatar -->
      <div class="conversation-avatar">
        <el-avatar :size="48" :style="{ background: getAvatarColor(conversation.participant.id) }">
          {{ conversation.participant.name?.charAt(0) || conversation.participant.username?.charAt(0) || '?' }}
        </el-avatar>
        <!-- Online status indicator -->
        <span
          v-if="conversation.participant.is_online"
          class="online-indicator"
        />
      </div>

      <!-- Content -->
      <div class="conversation-content">
        <div class="conversation-header">
          <span class="participant-name">
            {{ conversation.participant.name || conversation.participant.username }}
          </span>
          <span class="last-time">
            {{ formatLastTime(conversation.last_message_at) }}
          </span>
        </div>
        <div class="conversation-preview">
          <span class="last-message">
            {{ getLastMessagePreview(conversation.last_message) }}
          </span>
          <!-- Unread badge -->
          <el-badge
            v-if="conversation.unread_count > 0"
            :value="conversation.unread_count > 99 ? '99+' : conversation.unread_count"
            class="unread-badge"
          />
        </div>
      </div>
    </div>

    <!-- Empty state -->
    <div v-if="conversations.length === 0" class="empty-state">
      <el-empty description="暂无会话" :image-size="80" />
    </div>
  </div>
</template>

<script setup lang="ts">
/**
 * ConversationList.vue - 会话列表组件
 * Requirements: 9.1 (Conversation list), 9.2 (Last message preview), 9.3 (Unread count badge)
 * 
 * Displays list of conversations with:
 * - Other user's avatar and name
 * - Last message preview
 * - Unread count badge
 * - Last message time
 * - Online status indicator
 */
import type { ConversationWithDetails, ChatMessage } from '@/api/chat'

defineProps<{
  /** List of conversations to display */
  conversations: ConversationWithDetails[]
  /** Currently selected conversation ID */
  selectedId?: number
}>()

const emit = defineEmits<{
  /** Emitted when a conversation is selected */
  (e: 'select', conversation: ConversationWithDetails): void
}>()

/**
 * Handle conversation selection
 * @param conversation - The selected conversation
 */
const handleSelect = (conversation: ConversationWithDetails) => {
  emit('select', conversation)
}

/**
 * Get avatar background color based on user ID
 * @param userId - User ID for consistent color
 * @returns CSS color string
 */
const getAvatarColor = (userId: number): string => {
  const colors = [
    '#0ea5e9', // sky
    '#8b5cf6', // violet
    '#ec4899', // pink
    '#f97316', // orange
    '#10b981', // emerald
    '#6366f1', // indigo
    '#14b8a6', // teal
    '#f59e0b', // amber
  ]
  return colors[userId % colors.length]
}

/**
 * Format last message time for display
 * @param timestamp - ISO timestamp string or undefined
 * @returns Formatted time string
 */
const formatLastTime = (timestamp?: string): string => {
  if (!timestamp) return ''
  
  const date = new Date(timestamp)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMins = Math.floor(diffMs / 60000)
  const diffDays = Math.floor(diffMs / 86400000)
  
  // Within 1 minute
  if (diffMins < 1) {
    return '刚刚'
  }
  
  // Within 1 hour
  if (diffMins < 60) {
    return `${diffMins}分钟前`
  }
  
  // Today
  if (date.toDateString() === now.toDateString()) {
    return date.toLocaleTimeString('zh-CN', {
      hour: '2-digit',
      minute: '2-digit'
    })
  }
  
  // Yesterday
  const yesterday = new Date(now)
  yesterday.setDate(yesterday.getDate() - 1)
  if (date.toDateString() === yesterday.toDateString()) {
    return '昨天'
  }
  
  // Within a week
  if (diffDays < 7) {
    const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
    return weekdays[date.getDay()]
  }
  
  // Older
  return date.toLocaleDateString('zh-CN', {
    month: '2-digit',
    day: '2-digit'
  })
}

/**
 * Get preview text for last message
 * @param message - Last message or undefined
 * @returns Preview text
 */
const getLastMessagePreview = (message?: ChatMessage): string => {
  if (!message) return '暂无消息'
  
  // Handle different message types
  switch (message.message_type) {
    case 'image':
      return '[图片]'
    case 'file':
      return '[文件]'
    default:
      // Truncate long text messages
      const maxLength = 30
      if (message.content.length > maxLength) {
        return message.content.substring(0, maxLength) + '...'
      }
      return message.content
  }
}
</script>

<style scoped lang="scss">
.conversation-list {
  height: 100%;
  overflow-y: auto;
  background: #fff;
}

.conversation-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.2s;
  border-bottom: 1px solid #f1f5f9;

  &:hover {
    background: #f8fafc;
  }

  &.is-selected {
    background: #e0f2fe;
    
    &:hover {
      background: #e0f2fe;
    }
  }
}

.conversation-avatar {
  position: relative;
  flex-shrink: 0;
  margin-right: 12px;

  .online-indicator {
    position: absolute;
    bottom: 2px;
    right: 2px;
    width: 12px;
    height: 12px;
    background: #10b981;
    border: 2px solid #fff;
    border-radius: 50%;
  }
}

.conversation-content {
  flex: 1;
  min-width: 0;
  overflow: hidden;
}

.conversation-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.participant-name {
  font-weight: 500;
  font-size: 15px;
  color: #1e293b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.last-time {
  font-size: 12px;
  color: #94a3b8;
  flex-shrink: 0;
  margin-left: 8px;
}

.conversation-preview {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.last-message {
  font-size: 13px;
  color: #64748b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

.unread-badge {
  flex-shrink: 0;
  margin-left: 8px;

  :deep(.el-badge__content) {
    background: #ef4444;
    border: none;
  }
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
  color: #94a3b8;
}
</style>
