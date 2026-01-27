import { defineStore } from 'pinia'
import { ref } from 'vue'
import { messageApi } from '@/api/message'

export const useMessageStore = defineStore('message', () => {
  const unreadCount = ref(0)

  // 获取未读消息数（用户 ID 从 JWT 中获取）
  const fetchUnreadCount = async () => {
    try {
      const res = await messageApi.getUnreadCount()
      if (res.data.code === 0 && res.data.data) {
        unreadCount.value = res.data.data.unread_count || 0
      }
    } catch (error) {
      console.error('获取未读消息数失败:', error)
    }
  }

  // 减少未读数
  const decreaseUnreadCount = (count: number = 1) => {
    unreadCount.value = Math.max(0, unreadCount.value - count)
  }

  // 重置未读数
  const resetUnreadCount = () => {
    unreadCount.value = 0
  }

  return {
    unreadCount,
    fetchUnreadCount,
    decreaseUnreadCount,
    resetUnreadCount
  }
})
