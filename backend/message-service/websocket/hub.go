package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

// Message WebSocket消息
type Message struct {
	Type   string      `json:"type"`
	UserID uint        `json:"user_id,omitempty"`
	Data   interface{} `json:"data"`
}

// Hub WebSocket连接管理中心
type Hub struct {
	// 已注册的客户端
	clients map[*Client]bool

	// 用户ID到客户端的映射
	userClients map[uint][]*Client

	// 广播消息通道
	broadcast chan *Message

	// 定向消息通道
	unicast chan *Message

	// 注册请求
	register chan *Client

	// 注销请求
	unregister chan *Client

	// 互斥锁
	mu sync.RWMutex
}

// NewHub 创建新的Hub
func NewHub() *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		userClients: make(map[uint][]*Client),
		broadcast:   make(chan *Message, 256),
		unicast:     make(chan *Message, 256),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
	}
}

// Run 运行Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			if client.UserID > 0 {
				h.userClients[client.UserID] = append(h.userClients[client.UserID], client)
			}
			h.mu.Unlock()
			log.Printf("Client registered: user_id=%d, total=%d", client.UserID, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				// 从用户映射中移除
				if client.UserID > 0 {
					clients := h.userClients[client.UserID]
					for i, c := range clients {
						if c == client {
							h.userClients[client.UserID] = append(clients[:i], clients[i+1:]...)
							break
						}
					}
					if len(h.userClients[client.UserID]) == 0 {
						delete(h.userClients, client.UserID)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("Client unregistered: user_id=%d, total=%d", client.UserID, len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			data, _ := json.Marshal(message)
			for client := range h.clients {
				select {
				case client.send <- data:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()

		case message := <-h.unicast:
			h.mu.RLock()
			if clients, ok := h.userClients[message.UserID]; ok {
				data, _ := json.Marshal(message)
				for _, client := range clients {
					select {
					case client.send <- data:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast 广播消息给所有客户端
func (h *Hub) Broadcast(msgType string, data interface{}) {
	h.broadcast <- &Message{
		Type: msgType,
		Data: data,
	}
}

// SendToUser 发送消息给指定用户
func (h *Hub) SendToUser(userID uint, msgType string, data interface{}) {
	h.unicast <- &Message{
		Type:   msgType,
		UserID: userID,
		Data:   data,
	}
}

// GetOnlineUsers 获取在线用户数
func (h *Hub) GetOnlineUsers() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.userClients)
}

// IsUserOnline 检查用户是否在线
func (h *Hub) IsUserOnline(userID uint) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.userClients[userID]
	return ok
}
