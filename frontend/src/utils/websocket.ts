import { ref, onUnmounted } from 'vue'
import { ElNotification } from 'element-plus'

export interface WebSocketMessage {
  type: string
  user_id?: number
  data: any
}

export type MessageHandler = (message: WebSocketMessage) => void

// WebSocketService 封装前端实时通信。
// 它连接 API Gateway 暴露的 /api/v1/ws，负责自动重连、消息分发和默认通知提示。
class WebSocketService {
  private ws: WebSocket | null = null
  private url: string
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectDelay = 3000
  private handlers: Map<string, Set<MessageHandler>> = new Map()
  private isConnected = ref(false)

  constructor() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    this.url = `${protocol}//${window.location.host}/api/v1/ws`
  }

  connect(token?: string) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      return
    }

    const url = token ? `${this.url}?token=${token}` : this.url
    this.ws = new WebSocket(url)

    this.ws.onopen = () => {
      console.log('WebSocket connected')
      this.isConnected.value = true
      this.reconnectAttempts = 0
    }

    this.ws.onmessage = (event) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data)
        this.handleMessage(message)
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error)
      }
    }

    this.ws.onclose = (event) => {
      console.log('WebSocket closed:', event.code, event.reason)
      this.isConnected.value = false
      this.attemptReconnect()
    }

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error)
    }
  }

  private attemptReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++
      console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`)
      // 重连时重新读取 token，避免使用已经退出登录或过期的旧连接信息。
      const token = localStorage.getItem('token')
      setTimeout(() => this.connect(token || undefined), this.reconnectDelay)
    } else {
      console.log('Max reconnect attempts reached')
    }
  }

  private handleMessage(message: WebSocketMessage) {
    // 调用注册的处理器
    const handlers = this.handlers.get(message.type)
    if (handlers) {
      handlers.forEach(handler => handler(message))
    }

    // 调用通用处理器
    const allHandlers = this.handlers.get('*')
    if (allHandlers) {
      allHandlers.forEach(handler => handler(message))
    }

    // 默认通知处理，页面也可以通过 subscribe 注册自己的业务处理器。
    this.showNotification(message)
  }

  // showNotification 将常见实时事件转换成 Element Plus 通知。
  // 复杂页面仍可订阅原始消息，自行决定如何刷新列表或变更状态。
  private showNotification(message: WebSocketMessage) {
    switch (message.type) {
      case 'new_message':
        ElNotification({
          title: '新消息',
          message: message.data.title || '您有一条新消息',
          type: 'info',
          duration: 5000,
        })
        break

      case 'interview_reminder':
        ElNotification({
          title: '面试提醒',
          message: message.data.message || '您有一场面试即将开始',
          type: 'warning',
          duration: 10000,
        })
        break

      case 'application_update':
        ElNotification({
          title: '申请状态更新',
          message: message.data.message || '您的申请状态已更新',
          type: 'success',
          duration: 5000,
        })
        break

      case 'system_notice':
        ElNotification({
          title: '系统通知',
          message: message.data.message,
          type: 'info',
          duration: 8000,
        })
        break
    }
  }

  on(type: string, handler: MessageHandler) {
    if (!this.handlers.has(type)) {
      this.handlers.set(type, new Set())
    }
    this.handlers.get(type)!.add(handler)
  }

  off(type: string, handler: MessageHandler) {
    const handlers = this.handlers.get(type)
    if (handlers) {
      handlers.delete(handler)
    }
  }

  send(type: string, data: any) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({ type, data }))
    } else {
      console.warn('WebSocket is not connected')
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
  }

  get connected() {
    return this.isConnected
  }
}

// 单例：整个前端应用共享同一条 WebSocket 连接，避免多个页面重复建立连接。
export const wsService = new WebSocketService()

// Vue Composable：组件通过 useWebSocket 订阅消息，卸载时自动取消订阅。
export function useWebSocket() {
  const connected = wsService.connected

  const subscribe = (type: string, handler: MessageHandler) => {
    wsService.on(type, handler)
    onUnmounted(() => {
      wsService.off(type, handler)
    })
  }

  const send = (type: string, data: any) => {
    wsService.send(type, data)
  }

  return {
    connected,
    subscribe,
    send,
    connect: (token?: string) => wsService.connect(token),
    disconnect: () => wsService.disconnect(),
  }
}
