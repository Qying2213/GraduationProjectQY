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
 *
 * 展示单个会话的消息历史，支持向上滚动加载更多、新消息自动滚到底部和日期分隔线。
 */
import { ref, watch, nextTick, onMounted } from 'vue'
import { Loading } from '@element-plus/icons-vue'
import MessageItem from './MessageItem.vue'
import type { ChatMessage } from '@/api/chat'

const props = defineProps<{
  /** 当前会话 ID */
  conversationId: number
  /** 需要展示的消息列表 */
  messages: ChatMessage[]
  /** 是否正在加载消息 */
  loading: boolean
  /** 当前用户 ID，用于判断消息左右对齐 */
  currentUserId: number
}>()

const emit = defineEmits<{
  /** 用户滚动到顶部并需要加载更多历史消息时触发 */
  (e: 'loadMore'): void
}>()

/** 消息容器 DOM 引用 */
const messagesContainerRef = ref<HTMLElement | null>(null)

/** 是否应该自动滚动到底部 */
const shouldAutoScroll = ref(true)

/** 加载历史消息前的滚动高度，用于维持滚动位置 */
const previousScrollHeight = ref(0)

/**
 * 将消息容器滚动到底部。
 */
const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainerRef.value) {
      messagesContainerRef.value.scrollTop = messagesContainerRef.value.scrollHeight
    }
  })
}

/**
 * 处理滚动事件，用于触发向上加载更多。
 */
const handleScroll = () => {
  if (!messagesContainerRef.value) return
  
  const { scrollTop, scrollHeight, clientHeight } = messagesContainerRef.value
  
  // 用户接近底部时才自动滚动，避免打断正在查看历史消息的操作。
  shouldAutoScroll.value = scrollHeight - scrollTop - clientHeight < 100
  
  // 滚动到顶部附近时加载更多历史消息。
  if (scrollTop < 50 && !props.loading && props.messages.length > 0) {
    previousScrollHeight.value = scrollHeight
    emit('loadMore')
  }
}

/**
 * 加载历史消息后维持原有阅读位置。
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
 * 判断某条消息前是否需要显示日期分隔线。
 * @param message 当前消息
 * @param index 消息在数组中的位置
 * 返回：是否显示日期分隔线
 */
const shouldShowDateSeparator = (message: ChatMessage, index: number): boolean => {
  if (index === 0) return true
  
  const currentDate = new Date(message.created_at).toDateString()
  const previousDate = new Date(props.messages[index - 1].created_at).toDateString()
  
  return currentDate !== previousDate
}

/**
 * 格式化日期分隔线文本。
 * @param timestamp ISO 时间字符串
 * 返回：日期展示文本
 */
const formatDateSeparator = (timestamp: string): string => {
  const date = new Date(timestamp)
  const now = new Date()
  
  // 今天
  if (date.toDateString() === now.toDateString()) {
    return '今天'
  }
  
  // 昨天
  const yesterday = new Date(now)
  yesterday.setDate(yesterday.getDate() - 1)
  if (date.toDateString() === yesterday.toDateString()) {
    return '昨天'
  }
  
  // 今年
  if (date.getFullYear() === now.getFullYear()) {
    return date.toLocaleDateString('zh-CN', {
      month: 'long',
      day: 'numeric'
    })
  }
  
  // 其他年份
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
