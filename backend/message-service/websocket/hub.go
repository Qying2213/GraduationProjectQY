// Package websocket 实现 message-service 的实时消息层。
//
// 网关会把 /api/v1/ws 隧道转发到这里。Hub 负责连接注册、在线状态、会话订阅和
// 消息分发，Client 负责单个 WebSocket 连接的读写循环。
package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

// ============================================================================
// 聊天消息类型常量
// 对应实时消息投递、已读回执、在线状态和输入状态等前端交互。
// ============================================================================

const (
	// MsgTypeChat 表示普通聊天消息，通过 WebSocket 实时投递。
	MsgTypeChat = "chat"

	// MsgTypeChatRead 表示消息已读通知。
	MsgTypeChatRead = "chat_read"

	// MsgTypeOnlineStatus 表示用户在线状态变化。
	MsgTypeOnlineStatus = "online_status"

	// MsgTypeTyping 表示正在输入状态。
	MsgTypeTyping = "typing"

	// MsgTypeSystem 表示系统通知。
	MsgTypeSystem = "system"

	// MsgTypeNotification 表示通用通知。
	MsgTypeNotification = "notification"
)

// ============================================================================
// WebSocket 消息结构
// ============================================================================

// Message WebSocket消息
type Message struct {
	Type   string      `json:"type"`
	UserID uint        `json:"user_id,omitempty"`
	Data   interface{} `json:"data"`
}

// ChatWebSocketMessage 表示通过 WebSocket 推送的聊天消息。
type ChatWebSocketMessage struct {
	Type           string `json:"type"`
	ConversationID uint   `json:"conversation_id"`
	Message        struct {
		ID        uint   `json:"id"`
		SenderID  uint   `json:"sender_id"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
	} `json:"message"`
}

// ChatReadMessage 表示会话消息已读通知。
type ChatReadMessage struct {
	Type           string `json:"type"`
	ConversationID uint   `json:"conversation_id"`
	ReaderID       uint   `json:"reader_id"`
	MarkedCount    int    `json:"marked_count"`
}

// OnlineStatusMessage 表示用户在线状态更新。
type OnlineStatusMessage struct {
	Type     string `json:"type"`
	UserID   uint   `json:"user_id"`
	IsOnline bool   `json:"is_online"`
}

// TypingMessage 表示正在输入提示。
type TypingMessage struct {
	Type           string `json:"type"`
	ConversationID uint   `json:"conversation_id"`
	UserID         uint   `json:"user_id"`
	IsTyping       bool   `json:"is_typing"`
}

// UserOnlineInfo 保存单个用户的在线状态和最后活跃时间。
type UserOnlineInfo struct {
	UserID     uint
	IsOnline   bool
	LastSeenAt time.Time
}

// Hub WebSocket 连接管理中心。
//
// 它相当于进程内的消息分发器：HTTP handler 把事件投递到 channel，Hub 再把事件
// 推送给在线客户端；离线消息持久化仍然由数据库层的聊天/消息 handler 负责。
type Hub struct {
	// 已注册的客户端
	clients map[*Client]bool

	// 用户ID到客户端的映射
	userClients map[uint][]*Client

	// 用户在线状态（userID -> 在线信息）
	onlineStatus map[uint]*UserOnlineInfo

	// 用户订阅的会话（conversationID -> userIDs），用于向会话成员广播聊天事件。
	conversationSubscribers map[uint]map[uint]bool

	// 广播消息通道
	broadcast chan *Message

	// 定向消息通道
	unicast chan *Message

	// 聊天消息通道，用于实时投递会话消息。
	chatMessage chan *ChatWebSocketMessage

	// 已读通知通道
	chatRead chan *ChatReadMessage

	// 在线状态通道，用于推送在线/离线变化。
	onlineStatusChan chan *OnlineStatusMessage

	// 正在输入通道
	typing chan *TypingMessage

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
		clients:                 make(map[*Client]bool),
		userClients:             make(map[uint][]*Client),
		onlineStatus:            make(map[uint]*UserOnlineInfo),
		conversationSubscribers: make(map[uint]map[uint]bool),
		broadcast:               make(chan *Message, 256),
		unicast:                 make(chan *Message, 256),
		chatMessage:             make(chan *ChatWebSocketMessage, 256),
		chatRead:                make(chan *ChatReadMessage, 256),
		onlineStatusChan:        make(chan *OnlineStatusMessage, 256),
		typing:                  make(chan *TypingMessage, 256),
		register:                make(chan *Client),
		unregister:              make(chan *Client),
	}
}

// Run 是 Hub 的事件循环。连接注册、注销、广播和定向推送都通过这个 select 循环
// 统一处理，避免并发修改连接状态导致行为不可控。
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.handleRegister(client)

		case client := <-h.unregister:
			h.handleUnregister(client)

		case message := <-h.broadcast:
			h.handleBroadcast(message)

		case message := <-h.unicast:
			h.handleUnicast(message)

		case chatMsg := <-h.chatMessage:
			h.handleChatMessage(chatMsg)

		case readMsg := <-h.chatRead:
			h.handleChatRead(readMsg)

		case statusMsg := <-h.onlineStatusChan:
			h.handleOnlineStatus(statusMsg)

		case typingMsg := <-h.typing:
			h.handleTyping(typingMsg)
		}
	}
}

// handleRegister 处理客户端注册并标记用户在线。
// 一个用户可能同时有多个连接，例如多个浏览器标签页或多台设备。
func (h *Hub) handleRegister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true
	if client.UserID > 0 {
		h.userClients[client.UserID] = append(h.userClients[client.UserID], client)

		// 更新在线状态。
		h.onlineStatus[client.UserID] = &UserOnlineInfo{
			UserID:     client.UserID,
			IsOnline:   true,
			LastSeenAt: time.Now(),
		}

		log.Printf("Client registered: user_id=%d, total=%d", client.UserID, len(h.clients))

		// 异步广播在线状态，避免阻塞注册流程。
		go h.broadcastOnlineStatus(client.UserID, true)
	}
}

// handleUnregister 处理客户端注销，并在用户没有其他连接时标记为离线。
func (h *Hub) handleUnregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

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

			// 如果该用户没有其他连接，则标记为离线。
			if len(h.userClients[client.UserID]) == 0 {
				delete(h.userClients, client.UserID)

				// 更新在线状态。
				if info, exists := h.onlineStatus[client.UserID]; exists {
					info.IsOnline = false
					info.LastSeenAt = time.Now()
				}

				log.Printf("Client unregistered: user_id=%d, total=%d", client.UserID, len(h.clients))

				// 异步广播离线状态，避免阻塞注销流程。
				go h.broadcastOnlineStatus(client.UserID, false)
			}
		}
	}
}

// handleBroadcast 向所有已连接客户端广播消息。
func (h *Hub) handleBroadcast(message *Message) {
	h.mu.Lock()
	defer h.mu.Unlock()

	data, _ := json.Marshal(message)
	for client := range h.clients {
		select {
		case client.send <- data:
		default:
			delete(h.clients, client)
			close(client.send)
		}
	}
}

// handleUnicast 向指定用户的所有在线连接发送消息。
func (h *Hub) handleUnicast(message *Message) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.userClients[message.UserID]; ok {
		data, _ := json.Marshal(message)
		for _, client := range clients {
			select {
			case client.send <- data:
			default:
				delete(h.clients, client)
				close(client.send)
			}
		}
	}
}

// handleChatMessage 把聊天消息实时推送给订阅该会话的在线用户。
func (h *Hub) handleChatMessage(chatMsg *ChatWebSocketMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, err := json.Marshal(chatMsg)
	if err != nil {
		log.Printf("Error marshaling chat message: %v", err)
		return
	}

	// 发送给该会话的所有订阅者。
	if subscribers, ok := h.conversationSubscribers[chatMsg.ConversationID]; ok {
		for userID := range subscribers {
			if clients, ok := h.userClients[userID]; ok {
				for _, client := range clients {
					select {
					case client.send <- data:
					default:
						log.Printf("Failed to send chat message to user %d", userID)
					}
				}
			}
		}
	}
}

// handleChatRead 把已读通知推送给会话订阅者。
func (h *Hub) handleChatRead(readMsg *ChatReadMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, err := json.Marshal(readMsg)
	if err != nil {
		log.Printf("Error marshaling chat read message: %v", err)
		return
	}

	// 发送给该会话的所有订阅者。
	if subscribers, ok := h.conversationSubscribers[readMsg.ConversationID]; ok {
		for userID := range subscribers {
			if clients, ok := h.userClients[userID]; ok {
				for _, client := range clients {
					select {
					case client.send <- data:
					default:
						log.Printf("Failed to send chat read notification to user %d", userID)
					}
				}
			}
		}
	}
}

// handleOnlineStatus 广播用户在线状态变化。
func (h *Hub) handleOnlineStatus(statusMsg *OnlineStatusMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, err := json.Marshal(statusMsg)
	if err != nil {
		log.Printf("Error marshaling online status message: %v", err)
		return
	}

	// 广播给所有已连接客户端。
	for client := range h.clients {
		select {
		case client.send <- data:
		default:
			log.Printf("Failed to send online status to client")
		}
	}
}

// handleTyping 把“正在输入”状态推送给同一会话的其他成员。
func (h *Hub) handleTyping(typingMsg *TypingMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, err := json.Marshal(typingMsg)
	if err != nil {
		log.Printf("Error marshaling typing message: %v", err)
		return
	}

	// 发送给该会话除发送者之外的所有订阅者。
	if subscribers, ok := h.conversationSubscribers[typingMsg.ConversationID]; ok {
		for userID := range subscribers {
			// 不把输入状态回发给发送者自己。
			if userID == typingMsg.UserID {
				continue
			}
			if clients, ok := h.userClients[userID]; ok {
				for _, client := range clients {
					select {
					case client.send <- data:
					default:
						log.Printf("Failed to send typing indicator to user %d", userID)
					}
				}
			}
		}
	}
}

// broadcastOnlineStatus 广播指定用户的在线/离线变化。
func (h *Hub) broadcastOnlineStatus(userID uint, isOnline bool) {
	h.onlineStatusChan <- &OnlineStatusMessage{
		Type:     MsgTypeOnlineStatus,
		UserID:   userID,
		IsOnline: isOnline,
	}
}

// Broadcast 广播消息给所有客户端
func (h *Hub) Broadcast(msgType string, data interface{}) {
	h.broadcast <- &Message{
		Type: msgType,
		Data: data,
	}
}

// SendToUser 发送消息给指定用户的在线连接。
// 离线消息的持久化由聊天 handler 和数据库负责，这里只处理在线投递。
func (h *Hub) SendToUser(userID uint, msgType string, data interface{}) {
	h.unicast <- &Message{
		Type:   msgType,
		UserID: userID,
		Data:   data,
	}
}

// SendChatMessage 向会话参与者发送聊天消息。
func (h *Hub) SendChatMessage(conversationID uint, messageID uint, senderID uint, content string, createdAt string) {
	msg := &ChatWebSocketMessage{
		Type:           MsgTypeChat,
		ConversationID: conversationID,
	}
	msg.Message.ID = messageID
	msg.Message.SenderID = senderID
	msg.Message.Content = content
	msg.Message.CreatedAt = createdAt

	h.chatMessage <- msg
}

// SendChatReadNotification 向会话参与者发送已读通知。
func (h *Hub) SendChatReadNotification(conversationID uint, readerID uint, markedCount int) {
	h.chatRead <- &ChatReadMessage{
		Type:           MsgTypeChatRead,
		ConversationID: conversationID,
		ReaderID:       readerID,
		MarkedCount:    markedCount,
	}
}

// SendTypingIndicator 向会话参与者发送“正在输入”状态。
func (h *Hub) SendTypingIndicator(conversationID uint, userID uint, isTyping bool) {
	h.typing <- &TypingMessage{
		Type:           MsgTypeTyping,
		ConversationID: conversationID,
		UserID:         userID,
		IsTyping:       isTyping,
	}
}

// SubscribeToConversation 让用户订阅某个会话的实时更新。
func (h *Hub) SubscribeToConversation(conversationID uint, userID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.conversationSubscribers[conversationID] == nil {
		h.conversationSubscribers[conversationID] = make(map[uint]bool)
	}
	h.conversationSubscribers[conversationID][userID] = true
	log.Printf("User %d subscribed to conversation %d", userID, conversationID)
}

// UnsubscribeFromConversation 取消用户对某个会话的实时订阅。
func (h *Hub) UnsubscribeFromConversation(conversationID uint, userID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if subscribers, ok := h.conversationSubscribers[conversationID]; ok {
		delete(subscribers, userID)
		if len(subscribers) == 0 {
			delete(h.conversationSubscribers, conversationID)
		}
	}
	log.Printf("User %d unsubscribed from conversation %d", userID, conversationID)
}

// GetOnlineUsers 获取在线用户数
func (h *Hub) GetOnlineUsers() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.userClients)
}

// IsUserOnline 检查用户是否在线。
func (h *Hub) IsUserOnline(userID uint) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.userClients[userID]
	return ok
}

// GetUserOnlineStatus 返回指定用户的在线状态信息。
func (h *Hub) GetUserOnlineStatus(userID uint) *UserOnlineInfo {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if info, ok := h.onlineStatus[userID]; ok {
		return info
	}
	return &UserOnlineInfo{
		UserID:   userID,
		IsOnline: false,
	}
}

// GetOnlineUserIDs 返回所有当前在线用户 ID。
func (h *Hub) GetOnlineUserIDs() []uint {
	h.mu.RLock()
	defer h.mu.RUnlock()

	userIDs := make([]uint, 0, len(h.userClients))
	for userID := range h.userClients {
		userIDs = append(userIDs, userID)
	}
	return userIDs
}

// GetConversationSubscribers 返回订阅指定会话的用户 ID 列表。
func (h *Hub) GetConversationSubscribers(conversationID uint) []uint {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if subscribers, ok := h.conversationSubscribers[conversationID]; ok {
		userIDs := make([]uint, 0, len(subscribers))
		for userID := range subscribers {
			userIDs = append(userIDs, userID)
		}
		return userIDs
	}
	return []uint{}
}
