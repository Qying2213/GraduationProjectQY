package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// 写入等待时间
	writeWait = 10 * time.Second

	// 读取pong消息的等待时间
	pongWait = 60 * time.Second

	// 发送ping消息的间隔
	pingPeriod = (pongWait * 9) / 10

	// 最大消息大小（聊天消息内容可能比普通控制消息更长）
	maxMessageSize = 4096
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		// 如果没有 Origin 头（如某些客户端或测试），允许连接
		if origin == "" {
			return true
		}
		allowedOrigins := map[string]bool{
			"http://localhost:3000": true,
			"http://localhost:5173": true,
			"http://127.0.0.1:3000": true,
			"http://127.0.0.1:5173": true,
		}
		return allowedOrigins[origin]
	},
}

// ClientMessage 表示前端通过 WebSocket 发来的控制消息。
type ClientMessage struct {
	Type           string `json:"type"`
	ConversationID uint   `json:"conversation_id,omitempty"`
	Content        string `json:"content,omitempty"`
	IsTyping       bool   `json:"is_typing,omitempty"`
}

// Client WebSocket客户端
type Client struct {
	hub *Hub

	// WebSocket连接
	conn *websocket.Conn

	// 发送消息的缓冲通道
	send chan []byte

	// 用户ID
	UserID uint

	// 可选的消息回调，用于处理业务侧自定义的入站消息。
	OnMessage func(client *Client, message *ClientMessage)
}

// readPump 从 WebSocket 连接读取消息。
// 读循环负责解析前端发来的订阅、取消订阅和输入状态等控制消息。
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		log.Printf("Received message from user %d: %s", c.UserID, message)

		// 解析前端发来的消息。
		var clientMsg ClientMessage
		if err := json.Unmarshal(message, &clientMsg); err != nil {
			log.Printf("Error parsing client message: %v", err)
			continue
		}

		// 根据消息类型处理订阅、取消订阅和输入状态。
		switch clientMsg.Type {
		case MsgTypeTyping:
			// 处理正在输入状态。
			if clientMsg.ConversationID > 0 {
				c.hub.SendTypingIndicator(clientMsg.ConversationID, c.UserID, clientMsg.IsTyping)
			}

		case "subscribe":
			// 订阅会话，接收该会话的实时消息。
			if clientMsg.ConversationID > 0 {
				c.hub.SubscribeToConversation(clientMsg.ConversationID, c.UserID)
			}

		case "unsubscribe":
			// 取消会话订阅。
			if clientMsg.ConversationID > 0 {
				c.hub.UnsubscribeFromConversation(clientMsg.ConversationID, c.UserID)
			}

		default:
			// 如果业务层注册了回调，则交给业务层继续处理。
			if c.OnMessage != nil {
				c.OnMessage(c, &clientMsg)
			}
		}
	}
}

// writePump 向 WebSocket 连接写入消息。
// 写循环集中处理服务端推送、消息批量刷出和心跳 ping。
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 将队列中的消息一起发送
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs 处理 WebSocket 升级请求，并把连接注册到 Hub。
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, userID uint) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		UserID: userID,
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
