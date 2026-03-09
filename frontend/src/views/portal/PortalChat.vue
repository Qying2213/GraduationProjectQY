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
 * Requirements: 8.1 (WebSocket connection), 8.2 (Real-time messaging), 8.3 (Display messages)
 * 
 * Two-column layout with:
 * - Left: Conversation list
 * - Right: Chat window with messages and input
 * 
 * Features:
 * - Load conversations on mount
 * - Select conversation to view messages
 * - Send messages via API and WebSocket
 * - Receive real-time messages via WebSocket
 * - Mark messages as read when conversation is selected
 * - Handle WebSocket connection/reconnection
 */
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh, Warning, ChatDotRound } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/user'
import { useWebSocket } from '@/utils/websocket'
import { chatApi, type ConversationWithDetails, type ChatMessage } from '@/api/chat'
import { ConversationList, ChatWindow, ChatInput } from '@/components/chat'

// ============================================================================
// Store and WebSocket Setup
// ============================================================================

const userStore = useUserStore()
const { connected: wsConnected, subscribe, send, connect } = useWebSocket()

// ============================================================================
// State
// ============================================================================

/** List of conversations */
const conversations = ref<ConversationWithDetails[]>([])

/** Currently selected conversation ID */
const selectedConversationId = ref<number | null>(null)

/** Messages in the selected conversation */
const messages = ref<ChatMessage[]>([])

/** Total unread message count */
const totalUnreadCount = ref(0)

/** Loading states */
const loadingConversations = ref(false)
const loadingMessages = ref(false)

/** Pagination for messages */
const messagePage = ref(1)
const messagePageSize = ref(20)
const hasMoreMessages = ref(true)

/** Reference to ChatWindow component */
const chatWindowRef = ref<InstanceType<typeof ChatWindow> | null>(null)

// ============================================================================
// Computed Properties
// ============================================================================

/** Current user ID */
const currentUserId = computed(() => userStore.user?.id || 0)

/** Currently selected conversation object */
const selectedConversation = computed(() => {
  if (!selectedConversationId.value) return null
  return conversations.value.find(c => c.id === selectedConversationId.value) || null
})

// ============================================================================
// API Methods
// ============================================================================

/**
 * Load conversation list
 * Requirements: 9.1 (Conversation list sorted by last message time)
 */
const loadConversations = async () => {
  loadingConversations.value = true
  try {
    const res = await chatApi.getConversations()
    if (res.data?.code === 0 && res.data.data) {
      conversations.value = res.data.data.conversations || []
      // Calculate total unread count
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
 * Load messages for a conversation
 * Requirements: 8.5 (Chat message persistence), 8.6 (Message pagination)
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
        // Prepend older messages (they come in reverse order from API)
        messages.value = [...newMessages.reverse(), ...messages.value]
      } else {
        // Initial load - messages come newest first, reverse for display
        messages.value = newMessages.reverse()
      }

      // Check if there are more messages
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
 * Load more messages (infinite scroll)
 * Requirements: 8.6 (Load older messages with pagination)
 */
const loadMoreMessages = async () => {
  if (!selectedConversationId.value || loadingMessages.value || !hasMoreMessages.value) {
    return
  }

  messagePage.value++
  await loadMessages(selectedConversationId.value, true)
}

/**
 * Mark conversation as read
 * Requirements: 9.4 (Mark messages as read when conversation is selected)
 */
const markConversationAsRead = async (conversationId: number) => {
  try {
    await chatApi.markAsRead(conversationId)
    
    // Update local unread count
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
 * Send a message
 * Requirements: 8.2 (Deliver message in real-time via WebSocket)
 */
const handleSendMessage = async (content: string) => {
  if (!selectedConversationId.value || !content.trim()) return

  try {
    // Send via API (which will also broadcast via WebSocket)
    const res = await chatApi.sendMessage(selectedConversationId.value, content.trim())
    
    if (res.data?.code === 0 && res.data.data) {
      const newMessage = res.data.data
      
      // WebSocket may arrive before this HTTP response; dedupe by ID.
      const exists = messages.value.some(m => m.id === newMessage.id)
      if (!exists) {
        messages.value.push(newMessage)
      }
      
      // Update conversation's last message
      updateConversationLastMessage(selectedConversationId.value, newMessage)
      
      // Scroll to bottom
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
// WebSocket Handlers
// ============================================================================

/**
 * Handle incoming chat message via WebSocket
 * Requirements: 8.3 (Display message immediately in conversation view)
 */
const handleWebSocketChatMessage = (wsMessage: any) => {
  const { conversation_id, message } = wsMessage.data || {}
  
  if (!message) return

  // If this message is for the currently selected conversation, add it
  if (conversation_id === selectedConversationId.value) {
    // Check if message already exists (avoid duplicates)
    const exists = messages.value.some(m => m.id === message.id)
    if (!exists) {
      messages.value.push(message)
      
      // Scroll to bottom for new messages
      nextTick(() => {
        chatWindowRef.value?.scrollToBottom()
      })
      
      // Mark as read since we're viewing this conversation
      markConversationAsRead(conversation_id)
    }
  } else {
    // Update unread count for other conversations
    const conversation = conversations.value.find(c => c.id === conversation_id)
    if (conversation) {
      conversation.unread_count = (conversation.unread_count || 0) + 1
      totalUnreadCount.value++
    }
  }

  // Update conversation's last message
  updateConversationLastMessage(conversation_id, message)
}

/**
 * Handle message read status update via WebSocket
 * Requirements: 9.4 (Mark messages as read)
 */
const handleWebSocketReadStatus = (wsMessage: any) => {
  const { conversation_id } = wsMessage.data || {}
  
  if (conversation_id === selectedConversationId.value) {
    // Mark all messages as read in the current conversation
    messages.value.forEach(msg => {
      if (msg.sender_id === currentUserId.value) {
        msg.is_read = true
      }
    })
  }
}

// ============================================================================
// Helper Methods
// ============================================================================

/**
 * Handle conversation selection
 */
const handleSelectConversation = async (conversation: ConversationWithDetails) => {
  if (selectedConversationId.value === conversation.id) return

  selectedConversationId.value = conversation.id
  messages.value = []
  
  // Load messages for the selected conversation
  await loadMessages(conversation.id)
  
  // Mark as read
  if (conversation.unread_count > 0) {
    await markConversationAsRead(conversation.id)
  }
}

/**
 * Update conversation's last message in the list
 */
const updateConversationLastMessage = (conversationId: number, message: ChatMessage) => {
  const conversation = conversations.value.find(c => c.id === conversationId)
  if (conversation) {
    conversation.last_message = message
    conversation.last_message_at = message.created_at
    
    // Re-sort conversations by last message time
    conversations.value.sort((a, b) => {
      const timeA = a.last_message_at ? new Date(a.last_message_at).getTime() : 0
      const timeB = b.last_message_at ? new Date(b.last_message_at).getTime() : 0
      return timeB - timeA
    })
  }
}

/**
 * Get avatar background color based on user ID
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
// Lifecycle Hooks
// ============================================================================

onMounted(async () => {
  // Connect to WebSocket with auth token
  if (userStore.token) {
    connect(userStore.token)
  }

  // Subscribe to WebSocket events
  subscribe('chat', handleWebSocketChatMessage)
  subscribe('chat_read', handleWebSocketReadStatus)

  // Load initial data
  await loadConversations()
})

onUnmounted(() => {
  // WebSocket cleanup is handled by the useWebSocket composable
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
