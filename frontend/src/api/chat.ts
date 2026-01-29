import request from '@/utils/request'
import type { ApiResponse } from '@/types'

// ============================================================================
// TypeScript Interfaces for Chat API
// Requirements: 8.2 (Real-time Messaging), 9.1 (Conversation Management)
// ============================================================================

/**
 * Basic user information for conversation display
 */
export interface UserInfo {
  id: number
  username: string
  name: string
  avatar?: string
  role?: string
  is_online: boolean
}

/**
 * Chat message in a conversation
 * Requirements: 8.5 (Chat Message Persistence)
 */
export interface ChatMessage {
  id: number
  conversation_id: number
  sender_id: number
  content: string
  message_type: 'text' | 'image' | 'file'
  is_read: boolean
  created_at: string
}

/**
 * Conversation between two users
 * Requirements: 9.1 (Conversation List)
 */
export interface Conversation {
  id: number
  participant_a: number
  participant_b: number
  last_message_id?: number
  last_message_at?: string
  created_at: string
  updated_at: string
}

/**
 * Conversation with additional computed fields for display
 * Requirements: 9.1 (Conversation List), 9.2 (Last Message Preview), 9.3 (Unread Count)
 */
export interface ConversationWithDetails {
  id: number
  participant: UserInfo
  last_message?: ChatMessage
  unread_count: number
  last_message_at?: string
  created_at: string
  updated_at: string
}

/**
 * Response structure for conversation list API
 */
export interface ConversationListResponse {
  conversations: ConversationWithDetails[]
  total: number
}

/**
 * Response structure for message list API
 */
export interface MessageListResponse {
  messages: ChatMessage[]
  total: number
  page: number
  page_size: number
}

/**
 * Response structure for creating a conversation
 */
export interface CreateConversationResponse {
  conversation: ConversationWithDetails
  is_new: boolean
}

/**
 * Response structure for unread count API
 */
export interface UnreadCountResponse {
  total_unread: number
}

// ============================================================================
// Chat API Methods
// Requirements: 8.2 (Real-time Messaging), 9.1 (Conversation Management)
// ============================================================================

export const chatApi = {
  /**
   * Get conversation list for current user
   * Requirements: 9.1 (Conversation List), 9.2 (Last Message Preview), 9.3 (Unread Count)
   * @param params - Optional pagination parameters
   * @returns List of conversations with details
   */
  getConversations(params?: { page?: number; page_size?: number }) {
    return request.get<ApiResponse<ConversationListResponse>>('/conversations', { params })
  },

  /**
   * Create or get existing conversation with a user
   * Requirements: 9.1 (Conversation Management)
   * @param userId - The ID of the user to start conversation with
   * @returns The conversation (new or existing)
   */
  createOrGetConversation(userId: number) {
    return request.post<ApiResponse<CreateConversationResponse>>('/conversations', {
      participant_id: userId
    })
  },

  /**
   * Get total unread message count for current user
   * Requirements: 10.1 (Notification Badge), 10.2 (Real-time Unread Count)
   * @returns Total unread message count
   */
  getTotalUnreadCount() {
    return request.get<ApiResponse<UnreadCountResponse>>('/conversations/unread-count')
  },

  /**
   * Get messages in a conversation with pagination
   * Requirements: 8.5 (Chat Message Persistence), 8.6 (Message Pagination)
   * @param conversationId - The conversation ID
   * @param params - Optional pagination parameters
   * @returns List of messages
   */
  getMessages(conversationId: number, params?: { page?: number; page_size?: number }) {
    return request.get<ApiResponse<MessageListResponse>>(
      `/conversations/${conversationId}/messages`,
      { params }
    )
  },

  /**
   * Send a message in a conversation
   * Requirements: 8.2 (Real-time Messaging), 8.5 (Chat Message Persistence)
   * @param conversationId - The conversation ID
   * @param content - Message content
   * @param messageType - Message type (default: 'text')
   * @returns The sent message
   */
  sendMessage(conversationId: number, content: string, messageType: string = 'text') {
    return request.post<ApiResponse<ChatMessage>>(
      `/conversations/${conversationId}/messages`,
      {
        content,
        message_type: messageType
      }
    )
  },

  /**
   * Mark all messages in a conversation as read
   * Requirements: 9.4 (Mark Messages as Read), 10.6 (Unread Count Update)
   * @param conversationId - The conversation ID
   * @returns Success response
   */
  markAsRead(conversationId: number) {
    return request.put<ApiResponse<void>>(`/conversations/${conversationId}/read`)
  }
}
