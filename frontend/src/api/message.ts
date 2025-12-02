import request from '@/utils/request'
import type { Message, ApiResponse } from '@/types'

export const messageApi = {
    // 发送消息
    send(data: Partial<Message>) {
        return request.post<ApiResponse<Message>>('/messages', data)
    },

    // 获取消息列表
    list(params: any) {
        return request.get<ApiResponse>('/messages', { params })
    },

    // 标记为已读
    markAsRead(id: number) {
        return request.put<ApiResponse>(`/messages/${id}/read`)
    },

    // 获取未读数量
    getUnreadCount(userId: number) {
        return request.get<ApiResponse>('/messages/unread-count', {
            params: { user_id: userId }
        })
    },

    // 删除消息
    delete(id: number) {
        return request.delete<ApiResponse>(`/messages/${id}`)
    }
}
