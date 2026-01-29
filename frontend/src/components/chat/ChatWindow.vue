<template>
  <div class="chat-window">
    <!-- Loading indicator for initial load -->
    <div v-if="loading && messages.length === 0" class="loading-container">
      <el-icon class="loading-icon"><Loading /></el-icon>
      <span>加载中...</span>
    </div>

    <!-- Messages container -->
    <div
      v-else
      ref="messagesContainerRef"
      class="messages-container"
      @scroll="handleScroll"
    >
      <!-- Load more indicator -->
      <div v-if="loading && messages.length > 0" class="load-more-indicator">
        <el-icon class="loading-icon"><Loading /></el-icon>
        <span>加载更多消息...</span>
      </div>

      <!-- Empty state -->
      <div v-if="!loading && messages.length === 0" class="empty-messages">
        <el-empty description="暂无消息，开始聊天吧" :image-size="100" />
      </div>

      <!-- Messages list -->
      <div class="messages-list">
        <template v-for="(message, index) in messages" :key="message.id">
          <!-- Date separator -->
          <div
            v-if="shouldShowDateSeparator(message, index)"
            class="date-separator"
          >
            <span>{{ formatDateSeparator(message.created_at) }}</span>
          </div>
          
          <!-- Message item -->
          <MessageItem
            :message="message"
            :is-self="message.sender_id === currentUserId"
          />
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
/**
 * ChatWindow.vue - 聊天窗口组件
 * Requirements: 8.3 (Display messages), 8.6 (Infinite scroll pagination)
 * 
 * Displays chat messages for a conversation with:
 * - Message history with infinite scroll (load more on scroll up)
 * - Auto-scroll to bottom on new messages
 * - Date separators between messages
 */
import { ref, watch, nextTick, onMounted } from 'vue'
import { Loading } from '@element-plus/icons-vue'
import MessageItem from './MessageItem.vue'
import type { ChatMessage } from '@/api/chat'

const props = defineProps<{
  /** Current conversation ID */
  conversationId: number
  /** List of messages to display */
  messages: ChatMessage[]
  /** Whether messages are being loaded */
  loading: boolean
  /** Current user ID for determining message alignment */
  currentUserId: number
}>()

const emit = defineEmits<{
  /** Emitted when user scrolls to top to load more messages */
  (e: 'loadMore'): void
}>()

/** Reference to the messages container element */
const messagesContainerRef = ref<HTMLElement | null>(null)

/** Flag to track if we should auto-scroll */
const shouldAutoScroll = ref(true)

/** Previous scroll height for maintaining position after loading more */
const previousScrollHeight = ref(0)

/**
 * Scroll to the bottom of the messages container
 */
const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainerRef.value) {
      messagesContainerRef.value.scrollTop = messagesContainerRef.value.scrollHeight
    }
  })
}

/**
 * Handle scroll event for infinite scroll
 */
const handleScroll = () => {
  if (!messagesContainerRef.value) return
  
  const { scrollTop, scrollHeight, clientHeight } = messagesContainerRef.value
  
  // Check if user is near bottom (within 100px)
  shouldAutoScroll.value = scrollHeight - scrollTop - clientHeight < 100
  
  // Load more when scrolled to top (within 50px)
  if (scrollTop < 50 && !props.loading && props.messages.length > 0) {
    previousScrollHeight.value = scrollHeight
    emit('loadMore')
  }
}

/**
 * Maintain scroll position after loading more messages
 */
const maintainScrollPosition = () => {
  nextTick(() => {
    if (messagesContainerRef.value && previousScrollHeight.value > 0) {
      const newScrollHeight = messagesContainerRef.value.scrollHeight
      const scrollDiff = newScrollHeight - previousScrollHeight.value
      messagesContainerRef.value.scrollTop = scrollDiff
      previousScrollHeight.value = 0
    }
  })
}

/**
 * Check if date separator should be shown before a message
 * @param message - Current message
 * @param index - Message index in array
 * @returns Whether to show date separator
 */
const shouldShowDateSeparator = (message: ChatMessage, index: number): boolean => {
  if (index === 0) return true
  
  const currentDate = new Date(message.created_at).toDateString()
  const previousDate = new Date(props.messages[index - 1].created_at).toDateString()
  
  return currentDate !== previousDate
}

/**
 * Format date for separator display
 * @param timestamp - ISO timestamp string
 * @returns Formatted date string
 */
const formatDateSeparator = (timestamp: string): string => {
  const date = new Date(timestamp)
  const now = new Date()
  
  // Today
  if (date.toDateString() === now.toDateString()) {
    return '今天'
  }
  
  // Yesterday
  const yesterday = new Date(now)
  yesterday.setDate(yesterday.getDate() - 1)
  if (date.toDateString() === yesterday.toDateString()) {
    return '昨天'
  }
  
  // This year
  if (date.getFullYear() === now.getFullYear()) {
    return date.toLocaleDateString('zh-CN', {
      month: 'long',
      day: 'numeric'
    })
  }
  
  // Other years
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

// Watch for new messages and auto-scroll if appropriate
watch(
  () => props.messages.length,
  (newLength, oldLength) => {
    if (newLength > oldLength) {
      // New messages added
      if (shouldAutoScroll.value) {
        scrollToBottom()
      } else if (previousScrollHeight.value > 0) {
        // Messages loaded at top, maintain position
        maintainScrollPosition()
      }
    }
  }
)

// Watch for conversation change and scroll to bottom
watch(
  () => props.conversationId,
  () => {
    shouldAutoScroll.value = true
    nextTick(() => {
      scrollToBottom()
    })
  }
)

// Initial scroll to bottom
onMounted(() => {
  scrollToBottom()
})

// Expose scrollToBottom for parent component
defineExpose({
  scrollToBottom
})
</script>

<style scoped lang="scss">
.chat-window {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #f8fafc;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #64748b;
  gap: 12px;

  .loading-icon {
    font-size: 32px;
    animation: spin 1s linear infinite;
  }
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 16px 0;
}

.load-more-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 12px;
  color: #64748b;
  gap: 8px;

  .loading-icon {
    font-size: 16px;
    animation: spin 1s linear infinite;
  }
}

.empty-messages {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 300px;
}

.messages-list {
  display: flex;
  flex-direction: column;
}

.date-separator {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px 0;

  span {
    background: #e2e8f0;
    color: #64748b;
    font-size: 12px;
    padding: 4px 12px;
    border-radius: 12px;
  }
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
