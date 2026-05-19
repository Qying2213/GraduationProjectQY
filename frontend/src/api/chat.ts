import request from '@/utils/request'
import type { ApiResponse } from '@/types'

// ============================================================================
// 聊天 API 类型定义
// 这里的类型与 message-service 的会话、消息和未读数接口保持一致。
// ============================================================================

/**
 * 会话展示用的用户基础信息。
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
 * 会话中的单条聊天消息。
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
 * 两个用户之间的一条会话。
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
 * 会话列表展示用的增强结构，包含对方用户、最后消息和未读数。
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
 * 会话列表接口响应结构。
 */
export interface ConversationListResponse {
  conversations: ConversationWithDetails[]
  total: number
}

/**
 * 消息列表接口响应结构。
 */
export interface MessageListResponse {
  messages: ChatMessage[]
  total: number
  page: number
  page_size: number
}

/**
 * 创建/获取会话接口响应结构。
 */
export interface CreateConversationResponse {
  conversation: ConversationWithDetails
  is_new: boolean
}

/**
 * 总未读数接口响应结构。
 */
export interface UnreadCountResponse {
  total_unread: number
}

// ============================================================================
// 聊天 API 方法
// 页面层通过这些方法调用网关下的 message-service 会话接口。
// ============================================================================

export const chatApi = {
  /**
   * 获取当前用户的会话列表。
   * @param params 可选分页参数
   * 返回：带详情的会话列表
   */
  getConversations(params?: { page?: number; page_size?: number }) {
    return request.get<ApiResponse<ConversationListResponse>>('/conversations', { params })
  },

  /**
   * 与指定用户创建或获取已有会话。
   * @param userId 对方用户 ID
   * 返回：新建或已存在的会话
   */
  createOrGetConversation(userId: number) {
    return request.post<ApiResponse<CreateConversationResponse>>('/conversations', {
      participant_id: userId
    })
  },

  /**
   * 获取当前用户所有会话的总未读数。
   * 返回：总未读消息数
   */
  getTotalUnreadCount() {
    return request.get<ApiResponse<UnreadCountResponse>>('/conversations/unread-count')
  },

  /**
   * 分页获取指定会话的消息列表。
   * @param conversationId 会话 ID
   * @param params 可选分页参数
   * 返回：消息列表
   */
  getMessages(conversationId: number, params?: { page?: number; page_size?: number }) {
    return request.get<ApiResponse<MessageListResponse>>(
      `/conversations/${conversationId}/messages`,
      { params }
    )
  },

  /**
   * 在指定会话中发送消息。
   * @param conversationId 会话 ID
   * @param content 消息内容
   * @param messageType 消息类型，默认 text
   * 返回：已发送的消息
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
   * 将指定会话中的未读消息标记为已读。
   * @param conversationId 会话 ID
   * 返回：成功响应
   */
  markAsRead(conversationId: number) {
    return request.put<ApiResponse<void>>(`/conversations/${conversationId}/read`)
  }
}
