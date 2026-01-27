import request from '@/utils/request'
import type { Message, ApiResponse } from '@/types'

export const messageApi = {
    // 发送消息
    send(data: Partial<Message>) {
        return request.post<ApiResponse<Message>>('/messages', data)
    },

    // 获取消息列表（用户 ID 从 JWT 中获取）
    list(params?: { page?: number; page_size?: number; type?: string; is_read?: string }) {
        return request.get<ApiResponse>('/messages', { params })
    },

    // 标记为已读
    markAsRead(id: number) {
        return request.put<ApiResponse>(`/messages/${id}/read`)
    },

    // 获取未读数量（用户 ID 从 JWT 中获取）
    getUnreadCount() {
        return request.get<ApiResponse>('/messages/unread-count')
    },

    // 删除消息
    delete(id: number) {
        return request.delete<ApiResponse>(`/messages/${id}`)
    },

    // 获取消息统计
    getStats() {
        return request.get<ApiResponse>('/messages/stats')
    }
}
