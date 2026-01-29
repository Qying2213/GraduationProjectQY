package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

// ============================================================================
// Chat Message Type Constants
// Requirements: 8.2 (Real-time message delivery), 9.6 (Online status indicator)
// ============================================================================

const (
	// MsgTypeChat is for chat messages
	// Requirements: 8.2 - Deliver messages in real-time via WebSocket
	MsgTypeChat = "chat"

	// MsgTypeChatRead is for message read notifications
	// Requirements: 9.4 - Mark messages as read
	MsgTypeChatRead = "chat_read"

	// MsgTypeOnlineStatus is for online status updates
	// Requirements: 9.6 - Display online status indicator
	MsgTypeOnlineStatus = "online_status"

	// MsgTypeTyping is for typing indicator
	MsgTypeTyping = "typing"

	// MsgTypeSystem is for system notifications (existing)
	MsgTypeSystem = "system"

	// MsgTypeNotification is for general notifications (existing)
	MsgTypeNotification = "notification"
)

// ============================================================================
// WebSocket Message Structures
// ============================================================================

// Message WebSocket消息
type Message struct {
	Type   string      `json:"type"`
	UserID uint        `json:"user_id,omitempty"`
	Data   interface{} `json:"data"`
}

// ChatWebSocketMessage represents a chat message sent via WebSocket
// Requirements: 8.2 - Real-time message delivery
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

// ChatReadMessage represents a message read notification
// Requirements: 9.4 - Mark messages as read
type ChatReadMessage struct {
	Type           string `json:"type"`
	ConversationID uint   `json:"conversation_id"`
	ReaderID       uint   `json:"reader_id"`
	MarkedCount    int    `json:"marked_count"`
}

// OnlineStatusMessage represents an online status update
// Requirements: 9.6 - Display online status indicator
type OnlineStatusMessage struct {
	Type     string `json:"type"`
	UserID   uint   `json:"user_id"`
	IsOnline bool   `json:"is_online"`
}

// TypingMessage represents a typing indicator
type TypingMessage struct {
	Type           string `json:"type"`
	ConversationID uint   `json:"conversation_id"`
	UserID         uint   `json:"user_id"`
	IsTyping       bool   `json:"is_typing"`
}

// UserOnlineInfo stores online status information for a user
type UserOnlineInfo struct {
	UserID     uint
	IsOnline   bool
	LastSeenAt time.Time
}

// Hub WebSocket连接管理中心
// Requirements: 8.2 (Real-time message delivery), 8.4 (Offline message delivery), 9.6 (Online status)
type Hub struct {
	// 已注册的客户端
	clients map[*Client]bool

	// 用户ID到客户端的映射
	userClients map[uint][]*Client

	// 用户在线状态 (userID -> online info)
	// Requirements: 9.6 - Display online status indicator
	onlineStatus map[uint]*UserOnlineInfo

	// 用户订阅的会话 (conversationID -> userIDs)
	// Used for broadcasting chat messages to conversation participants
	conversationSubscribers map[uint]map[uint]bool

	// 广播消息通道
	broadcast chan *Message

	// 定向消息通道
	unicast chan *Message

	// 聊天消息通道
	// Requirements: 8.2 - Real-time message delivery
	chatMessage chan *ChatWebSocketMessage

	// 已读通知通道
	chatRead chan *ChatReadMessage

	// 在线状态通道
	// Requirements: 9.6 - Online status indicator
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

// Run 运行Hub
// Requirements: 8.2 (Real-time message delivery), 9.6 (Online status)
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

// handleRegister handles client registration
// Requirements: 9.6 - Update online status when user connects
func (h *Hub) handleRegister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true
	if client.UserID > 0 {
		h.userClients[client.UserID] = append(h.userClients[client.UserID], client)

		// Update online status
		h.onlineStatus[client.UserID] = &UserOnlineInfo{
			UserID:     client.UserID,
			IsOnline:   true,
			LastSeenAt: time.Now(),
		}

		log.Printf("Client registered: user_id=%d, total=%d", client.UserID, len(h.clients))

		// Broadcast online status to other users (non-blocking)
		go h.broadcastOnlineStatus(client.UserID, true)
	}
}

// handleUnregister handles client unregistration
// Requirements: 9.6 - Update online status when user disconnects
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

			// If no more clients for this user, mark as offline
			if len(h.userClients[client.UserID]) == 0 {
				delete(h.userClients, client.UserID)

				// Update online status
				if info, exists := h.onlineStatus[client.UserID]; exists {
					info.IsOnline = false
					info.LastSeenAt = time.Now()
				}

				log.Printf("Client unregistered: user_id=%d, total=%d", client.UserID, len(h.clients))

				// Broadcast offline status to other users (non-blocking)
				go h.broadcastOnlineStatus(client.UserID, false)
			}
		}
	}
}

// handleBroadcast handles broadcast messages to all clients
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

// handleUnicast handles unicast messages to specific user
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

// handleChatMessage handles chat messages
// Requirements: 8.2 - Deliver messages in real-time via WebSocket
func (h *Hub) handleChatMessage(chatMsg *ChatWebSocketMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, err := json.Marshal(chatMsg)
	if err != nil {
		log.Printf("Error marshaling chat message: %v", err)
		return
	}

	// Send to all subscribers of this conversation
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

// handleChatRead handles chat read notifications
// Requirements: 9.4 - Mark messages as read
func (h *Hub) handleChatRead(readMsg *ChatReadMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, err := json.Marshal(readMsg)
	if err != nil {
		log.Printf("Error marshaling chat read message: %v", err)
		return
	}

	// Send to all subscribers of this conversation
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

// handleOnlineStatus handles online status updates
// Requirements: 9.6 - Display online status indicator
func (h *Hub) handleOnlineStatus(statusMsg *OnlineStatusMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, err := json.Marshal(statusMsg)
	if err != nil {
		log.Printf("Error marshaling online status message: %v", err)
		return
	}

	// Broadcast to all connected clients
	for client := range h.clients {
		select {
		case client.send <- data:
		default:
			log.Printf("Failed to send online status to client")
		}
	}
}

// handleTyping handles typing indicator messages
func (h *Hub) handleTyping(typingMsg *TypingMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, err := json.Marshal(typingMsg)
	if err != nil {
		log.Printf("Error marshaling typing message: %v", err)
		return
	}

	// Send to all subscribers of this conversation except the sender
	if subscribers, ok := h.conversationSubscribers[typingMsg.ConversationID]; ok {
		for userID := range subscribers {
			// Don't send typing indicator back to the sender
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

// broadcastOnlineStatus broadcasts online status change to all connected clients
// Requirements: 9.6 - Display online status indicator
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

// SendToUser 发送消息给指定用户
// Requirements: 8.4 - Store message and deliver when recipient connects (offline delivery)
func (h *Hub) SendToUser(userID uint, msgType string, data interface{}) {
	h.unicast <- &Message{
		Type:   msgType,
		UserID: userID,
		Data:   data,
	}
}

// SendChatMessage sends a chat message to conversation participants
// Requirements: 8.2 - Deliver messages in real-time via WebSocket
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

// SendChatReadNotification sends a read notification to conversation participants
// Requirements: 9.4 - Mark messages as read
func (h *Hub) SendChatReadNotification(conversationID uint, readerID uint, markedCount int) {
	h.chatRead <- &ChatReadMessage{
		Type:           MsgTypeChatRead,
		ConversationID: conversationID,
		ReaderID:       readerID,
		MarkedCount:    markedCount,
	}
}

// SendTypingIndicator sends a typing indicator to conversation participants
func (h *Hub) SendTypingIndicator(conversationID uint, userID uint, isTyping bool) {
	h.typing <- &TypingMessage{
		Type:           MsgTypeTyping,
		ConversationID: conversationID,
		UserID:         userID,
		IsTyping:       isTyping,
	}
}

// SubscribeToConversation subscribes a user to a conversation for real-time updates
// Requirements: 8.2 - Real-time message delivery
func (h *Hub) SubscribeToConversation(conversationID uint, userID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.conversationSubscribers[conversationID] == nil {
		h.conversationSubscribers[conversationID] = make(map[uint]bool)
	}
	h.conversationSubscribers[conversationID][userID] = true
	log.Printf("User %d subscribed to conversation %d", userID, conversationID)
}

// UnsubscribeFromConversation unsubscribes a user from a conversation
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

// IsUserOnline 检查用户是否在线
// Requirements: 9.6 - Display online status indicator
func (h *Hub) IsUserOnline(userID uint) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.userClients[userID]
	return ok
}

// GetUserOnlineStatus returns the online status info for a user
// Requirements: 9.6 - Display online status indicator
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

// GetOnlineUserIDs returns a list of all online user IDs
// Requirements: 9.6 - Display online status indicator
func (h *Hub) GetOnlineUserIDs() []uint {
	h.mu.RLock()
	defer h.mu.RUnlock()

	userIDs := make([]uint, 0, len(h.userClients))
	for userID := range h.userClients {
		userIDs = append(userIDs, userID)
	}
	return userIDs
}

// GetConversationSubscribers returns the list of user IDs subscribed to a conversation
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
