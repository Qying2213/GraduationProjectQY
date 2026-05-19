<template>
  <div class="portal-chat">
    <div class="chat-container">
      <!-- Left Panel: Conversation List -->
      <div class="conversation-panel">
        <div class="panel-header">
          <h2>消息</h2>
          <el-badge :value="totalUnreadCount" :max="99" :hidden="totalUnreadCount === 0">
            <el-icon class="refresh-icon" @click="loadConversations"><Refresh /></el-icon>
          </el-badge>
        </div>
        
        <!-- Connection Status -->
        <div v-if="!wsConnected" class="connection-status">
          <el-icon class="status-icon"><Warning /></el-icon>
          <span>连接中...</span>
        </div>
        
        <!-- Conversation List -->
        <div class="conversation-list-wrapper" v-loading="loadingConversations">
          <ConversationList
            :conversations="conversations"
            :selected-id="selectedConversationId"
            @select="handleSelectConversation"
          />
        </div>
      </div>

      <!-- Right Panel: Chat Window -->
      <div class="chat-panel">
        <template v-if="selectedConversation">
          <!-- Chat Header -->
          <div class="chat-header">
            <div class="participant-info">
              <el-avatar 
                :size="40" 
                :style="{ background: getAvatarColor(selectedConversation.participant.id) }"
              >
                {{ selectedConversation.participant.name?.charAt(0) || selectedConversation.participant.username?.charAt(0) || '?' }}
              </el-avatar>
              <div class="participant-details">
                <span class="participant-name">
                  {{ selectedConversation.participant.name || selectedConversation.participant.username }}
                </span>
                <span class="participant-status" :class="{ online: selectedConversation.participant.is_online }">
                  {{ selectedConversation.participant.is_online ? '在线' : '离线' }}
                </span>
              </div>
            </div>
          </div>

          <!-- Chat Window Component -->
          <ChatWindow
            ref="chatWindowRef"
            :conversation-id="selectedConversationId!"
            :messages="messages"
            :loading="loadingMessages"
            :current-user-id="currentUserId"
            @load-more="loadMoreMessages"
          />

          <!-- Chat Input Component -->
          <ChatInput
            :disabled="!wsConnected"
            @send="handleSendMessage"
          />
        </template>

        <!-- Empty State: No Conversation Selected -->
        <div v-else class="empty-chat">
          <el-empty description="选择一个会话开始聊天">
            <template #image>
              <el-icon class="empty-icon"><ChatDotRound /></el-icon>
            </template>
          </el-empty>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
/**
 * PortalChat.vue - 求职者端聊天页面
 *
 * 两栏布局：左侧会话列表，右侧聊天窗口和输入框。
 * 页面负责加载会话、切换消息、发送消息、接收 WebSocket 实时推送，并在选中会话后标记已读。
 */
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh, Warning, ChatDotRound } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/user'
import { useWebSocket } from '@/utils/websocket'
import { chatApi, type ConversationWithDetails, type ChatMessage } from '@/api/chat'
import { ConversationList, ChatWindow, ChatInput } from '@/components/chat'

// ============================================================================
// Store 和 WebSocket 初始化
// ============================================================================

const userStore = useUserStore()
const { connected: wsConnected, subscribe, send, connect } = useWebSocket()

// ============================================================================
// 页面状态
// ============================================================================

/** 会话列表 */
const conversations = ref<ConversationWithDetails[]>([])

/** 当前选中的会话 ID */
const selectedConversationId = ref<number | null>(null)

/** 当前会话中的消息 */
const messages = ref<ChatMessage[]>([])

/** 总未读消息数 */
const totalUnreadCount = ref(0)

/** 加载状态 */
const loadingConversations = ref(false)
const loadingMessages = ref(false)

/** 消息分页状态 */
const messagePage = ref(1)
const messagePageSize = ref(20)
const hasMoreMessages = ref(true)

/** ChatWindow 组件引用 */
const chatWindowRef = ref<InstanceType<typeof ChatWindow> | null>(null)

// ============================================================================
// 计算属性
// ============================================================================

/** 当前用户 ID */
const currentUserId = computed(() => userStore.user?.id || 0)

/** 当前选中的会话对象 */
const selectedConversation = computed(() => {
  if (!selectedConversationId.value) return null
  return conversations.value.find(c => c.id === selectedConversationId.value) || null
})

// ============================================================================
// API 方法
// ============================================================================

/**
 * 加载会话列表，后端已按最后消息时间排序。
 */
const loadConversations = async () => {
  loadingConversations.value = true
  try {
    const res = await chatApi.getConversations()
    if (res.data?.code === 0 && res.data.data) {
      conversations.value = res.data.data.conversations || []
      // 计算总未读数。
      totalUnreadCount.value = conversations.value.reduce(
        (sum, conv) => sum + (conv.unread_count || 0),
        0
      )
    }
  } catch (error) {
    console.error('Failed to load conversations:', error)
    ElMessage.error('加载会话列表失败')
  } finally {
    loadingConversations.value = false
  }
}

/**
 * 加载指定会话的消息，append=true 时用于向上翻页加载历史消息。
 */
const loadMessages = async (conversationId: number, append = false) => {
  if (!append) {
    loadingMessages.value = true
    messagePage.value = 1
    hasMoreMessages.value = true
  }

  try {
    const res = await chatApi.getMessages(conversationId, {
      page: messagePage.value,
      page_size: messagePageSize.value
    })

    if (res.data?.code === 0 && res.data.data) {
      const newMessages = res.data.data.messages || []
      
      if (append) {
        // 加载历史消息时插入到当前列表前面。
        messages.value = [...newMessages.reverse(), ...messages.value]
      } else {
        // 首次加载时反转为聊天窗口需要的正序展示。
        messages.value = newMessages.reverse()
      }

      // 判断是否还有更多历史消息。
      hasMoreMessages.value = newMessages.length >= messagePageSize.value
    }
  } catch (error) {
    console.error('Failed to load messages:', error)
    ElMessage.error('加载消息失败')
  } finally {
    loadingMessages.value = false
  }
}

/**
 * 加载更多历史消息，用于聊天窗口向上滚动分页。
 */
const loadMoreMessages = async () => {
  if (!selectedConversationId.value || loadingMessages.value || !hasMoreMessages.value) {
    return
  }

  messagePage.value++
  await loadMessages(selectedConversationId.value, true)
}

/**
 * 将当前会话标记为已读。
 */
const markConversationAsRead = async (conversationId: number) => {
  try {
    await chatApi.markAsRead(conversationId)
    
    // 同步更新本地未读数，避免等待下一次列表刷新。
    const conversation = conversations.value.find(c => c.id === conversationId)
    if (conversation && conversation.unread_count > 0) {
      totalUnreadCount.value = Math.max(0, totalUnreadCount.value - conversation.unread_count)
      conversation.unread_count = 0
    }
  } catch (error) {
    console.error('Failed to mark as read:', error)
  }
}

/**
 * 发送消息。
 */
const handleSendMessage = async (content: string) => {
  if (!selectedConversationId.value || !content.trim()) return

  try {
    // 通过 API 保存消息，后端会同时通过 WebSocket 广播。
    const res = await chatApi.sendMessage(selectedConversationId.value, content.trim())
    
    if (res.data?.code === 0 && res.data.data) {
      const newMessage = res.data.data
      
      // WebSocket 可能早于 HTTP 响应到达，因此按消息 ID 去重。
      const exists = messages.value.some(m => m.id === newMessage.id)
      if (!exists) {
        messages.value.push(newMessage)
      }
      
      // 更新会话最后一条消息。
      updateConversationLastMessage(selectedConversationId.value, newMessage)
      
      // 发送成功后滚动到底部。
      nextTick(() => {
        chatWindowRef.value?.scrollToBottom()
      })
    }
  } catch (error) {
    console.error('Failed to send message:', error)
    ElMessage.error('发送消息失败')
  }
}

// ============================================================================
// WebSocket 处理器
// ============================================================================

/**
 * 处理 WebSocket 推来的新聊天消息。
 */
const handleWebSocketChatMessage = (wsMessage: any) => {
  const { conversation_id, message } = wsMessage.data || {}
  
  if (!message) return

  // 如果消息属于当前会话，则直接追加到聊天窗口。
  if (conversation_id === selectedConversationId.value) {
    // 按消息 ID 去重，避免 HTTP 响应和 WebSocket 推送重复插入。
    const exists = messages.value.some(m => m.id === message.id)
    if (!exists) {
      messages.value.push(message)
      
      // 新消息进入当前会话后滚动到底部。
      nextTick(() => {
        chatWindowRef.value?.scrollToBottom()
      })
      
      // 当前正在查看该会话，因此立即标记为已读。
      markConversationAsRead(conversation_id)
    }
  } else {
    // 其他会话的新消息只更新未读数。
    const conversation = conversations.value.find(c => c.id === conversation_id)
    if (conversation) {
      conversation.unread_count = (conversation.unread_count || 0) + 1
      totalUnreadCount.value++
    }
  }

  // 更新会话最后一条消息。
  updateConversationLastMessage(conversation_id, message)
}

/**
 * 处理 WebSocket 推来的已读状态变更。
 */
const handleWebSocketReadStatus = (wsMessage: any) => {
  const { conversation_id } = wsMessage.data || {}
  
  if (conversation_id === selectedConversationId.value) {
    // 将当前会话中自己发送的消息标记为对方已读。
    messages.value.forEach(msg => {
      if (msg.sender_id === currentUserId.value) {
        msg.is_read = true
      }
    })
  }
}

// ============================================================================
// 辅助方法
// ============================================================================

/**
 * 处理会话选择。
 */
const handleSelectConversation = async (conversation: ConversationWithDetails) => {
  if (selectedConversationId.value === conversation.id) return

  selectedConversationId.value = conversation.id
  messages.value = []
  
  // 加载当前选中会话的消息。
  await loadMessages(conversation.id)
  
  // 选中后标记为已读。
  if (conversation.unread_count > 0) {
    await markConversationAsRead(conversation.id)
  }
}

/**
 * 更新会话列表中的最后一条消息。
 */
const updateConversationLastMessage = (conversationId: number, message: ChatMessage) => {
  const conversation = conversations.value.find(c => c.id === conversationId)
  if (conversation) {
    conversation.last_message = message
    conversation.last_message_at = message.created_at
    
    // 按最后消息时间重新排序会话。
    conversations.value.sort((a, b) => {
      const timeA = a.last_message_at ? new Date(a.last_message_at).getTime() : 0
      const timeB = b.last_message_at ? new Date(b.last_message_at).getTime() : 0
      return timeB - timeA
    })
  }
}

/**
 * 根据用户 ID 生成头像背景色。
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

// ============================================================================
// 生命周期钩子
// ============================================================================

onMounted(async () => {
  // 使用登录 token 连接 WebSocket。
  if (userStore.token) {
    connect(userStore.token)
  }

  // 订阅 WebSocket 事件。
  subscribe('chat', handleWebSocketChatMessage)
  subscribe('chat_read', handleWebSocketReadStatus)

  // 加载初始数据。
  await loadConversations()
})

onUnmounted(() => {
  // WebSocket 订阅清理由 useWebSocket 组合式函数处理。
})

// Watch for WebSocket connection changes
watch(wsConnected, (connected) => {
  if (connected) {
    // Reload conversations when reconnected
    loadConversations()
  }
})
</script>

<style scoped lang="scss">
.portal-chat {
  height: calc(100vh - 120px);
  padding: 24px;
  background: #f8fafc;

  .chat-container {
    display: flex;
    height: 100%;
    max-width: 1200px;
    margin: 0 auto;
    background: #fff;
    border-radius: 16px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
    overflow: hidden;
  }

  // Left Panel - Conversation List
  .conversation-panel {
    width: 320px;
    min-width: 280px;
    border-right: 1px solid #e2e8f0;
    display: flex;
    flex-direction: column;
    background: #fff;

    .panel-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 20px;
      border-bottom: 1px solid #e2e8f0;

      h2 {
        font-size: 20px;
        font-weight: 600;
        color: #1e293b;
        margin: 0;
      }

      .refresh-icon {
        font-size: 20px;
        color: #64748b;
        cursor: pointer;
        transition: color 0.2s, transform 0.2s;

        &:hover {
          color: #0ea5e9;
          transform: rotate(180deg);
        }
      }
    }

    .connection-status {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 8px;
      padding: 8px 16px;
      background: #fef3c7;
      color: #92400e;
      font-size: 13px;

      .status-icon {
        font-size: 16px;
      }
    }

    .conversation-list-wrapper {
      flex: 1;
      overflow-y: auto;
    }
  }

  // Right Panel - Chat Window
  .chat-panel {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-width: 0;
    background: #f8fafc;

    .chat-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 16px 20px;
      background: #fff;
      border-bottom: 1px solid #e2e8f0;

      .participant-info {
        display: flex;
        align-items: center;
        gap: 12px;

        .participant-details {
          display: flex;
          flex-direction: column;

          .participant-name {
            font-size: 16px;
            font-weight: 600;
            color: #1e293b;
          }

          .participant-status {
            font-size: 12px;
            color: #94a3b8;

            &.online {
              color: #10b981;
            }
          }
        }
      }
    }

    .empty-chat {
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #f8fafc;

      .empty-icon {
        font-size: 80px;
        color: #cbd5e1;
      }

      :deep(.el-empty__description) {
        color: #94a3b8;
      }
    }
  }
}

// Responsive adjustments
@media (max-width: 768px) {
  .portal-chat {
    padding: 0;
    height: calc(100vh - 60px);

    .chat-container {
      border-radius: 0;
      flex-direction: column;
    }

    .conversation-panel {
      width: 100%;
      height: 40%;
      border-right: none;
      border-bottom: 1px solid #e2e8f0;
    }

    .chat-panel {
      height: 60%;
    }
  }
}
</style>
