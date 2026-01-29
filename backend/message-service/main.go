package main

import (
	"log"
	"message-service/handlers"
	"message-service/models"
	"message-service/websocket"
	"os"

	"common/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// 数据库连接（支持环境变量配置）
	dbHost := getEnv("DB_HOST", "localhost")
	dbUser := getEnv("DB_USER", "qinyang")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "talent_platform")
	dbPort := getEnv("DB_PORT", "5432")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	dsn := "host=" + dbHost + " user=" + dbUser + " dbname=" + dbName + " port=" + dbPort + " sslmode=" + dbSSLMode + " TimeZone=Asia/Shanghai"
	if dbPassword != "" {
		dsn = "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=" + dbSSLMode + " TimeZone=Asia/Shanghai"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(&models.Message{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Auto-migrate chat models (conversations and chat_messages)
	if err := db.AutoMigrate(&models.Conversation{}, &models.ChatMessage{}); err != nil {
		log.Fatal("Failed to migrate chat tables:", err)
	}

	// Initialize WebSocket Hub
	// Requirements: 8.2 (Real-time message delivery), 9.6 (Online status)
	hub := websocket.NewHub()
	go hub.Run()

	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.SimpleOperationLog("message-service"))

	messageHandler := handlers.NewMessageHandler(db)
	chatHandler := handlers.NewChatHandler(db, hub)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "message-service", "port": 8085})
	})

	// WebSocket endpoint for real-time chat
	// Requirements: 8.1 (WebSocket connection), 8.2 (Real-time message delivery)
	r.GET("/ws", middleware.JWTAuth(), func(c *gin.Context) {
		jwtUserID, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{"error": "未授权"})
			return
		}
		userID, ok := jwtUserID.(uint)
		if !ok {
			c.JSON(500, gin.H{"error": "invalid user_id type"})
			return
		}
		websocket.ServeWs(hub, c.Writer, c.Request, userID)
	})

	// Online status endpoint
	// Requirements: 9.6 (Display online status indicator)
	r.GET("/api/v1/online-status", middleware.JWTAuth(), func(c *gin.Context) {
		onlineUsers := hub.GetOnlineUserIDs()
		c.JSON(200, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"online_users": onlineUsers,
				"total":        len(onlineUsers),
			},
		})
	})

	// Check specific user online status
	// Requirements: 9.6 (Display online status indicator)
	r.GET("/api/v1/online-status/:user_id", middleware.JWTAuth(), func(c *gin.Context) {
		var userID uint
		if _, err := c.Params.Get("user_id"); err {
			c.JSON(400, gin.H{"error": "invalid user_id"})
			return
		}
		// Parse user_id from URL
		if n, err := parseUint(c.Param("user_id")); err == nil {
			userID = n
		} else {
			c.JSON(400, gin.H{"error": "invalid user_id"})
			return
		}

		isOnline := hub.IsUserOnline(userID)
		statusInfo := hub.GetUserOnlineStatus(userID)

		c.JSON(200, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"user_id":      userID,
				"is_online":    isOnline,
				"last_seen_at": statusInfo.LastSeenAt,
			},
		})
	})

	// System messages API (existing)
	api := r.Group("/api/v1/messages")
	api.Use(middleware.JWTAuth())
	{
		api.POST("", messageHandler.SendMessage)
		api.GET("", messageHandler.GetMessages)
		api.GET("/stats", messageHandler.GetMessageStats)
		api.GET("/unread-count", messageHandler.GetUnreadCount)
		api.PUT("/:id/read", messageHandler.MarkAsRead)
		api.DELETE("/:id", messageHandler.DeleteMessage)
	}

	// Chat conversations API (new)
	// Requirements: 9.1 (Conversation List), 9.2 (Last Message Preview), 9.3 (Unread Count)
	chatAPI := r.Group("/api/v1/conversations")
	chatAPI.Use(middleware.JWTAuth())
	{
		// GET /conversations - 获取会话列表
		// Returns conversations sorted by last_message_at DESC with unread count and last message
		chatAPI.GET("", chatHandler.GetConversations)

		// POST /conversations - 创建/获取会话
		// Creates a new conversation or returns existing one with participant
		chatAPI.POST("", chatHandler.CreateOrGetConversation)

		// GET /conversations/unread-count - 获取总未读数
		// Returns total unread message count across all conversations
		chatAPI.GET("/unread-count", chatHandler.GetTotalUnreadCount)

		// GET /conversations/:id/messages - 获取会话消息列表（分页）
		// Returns messages for a conversation with pagination (oldest first for chat display)
		// Requirements: 8.5 (Chat Message Persistence), 8.6 (Load older messages with pagination)
		chatAPI.GET("/:id/messages", chatHandler.GetMessages)

		// POST /conversations/:id/messages - 发送消息
		// Creates a new message in the conversation and updates last_message_id/last_message_at
		// Requirements: 8.5 (Chat Message Persistence), 8.2 (Real-time delivery via WebSocket)
		chatAPI.POST("/:id/messages", chatHandler.SendMessage)

		// PUT /conversations/:id/read - 标记会话已读
		// Marks all unread messages from other user as read
		// Requirements: 9.4 (Mark all messages as read when opening conversation)
		chatAPI.PUT("/:id/read", chatHandler.MarkAsRead)
	}

	log.Println("Message service is running on :8085")
	if err := r.Run(":8085"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// parseUint parses a string to uint
func parseUint(s string) (uint, error) {
	var n uint
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, nil
		}
		n = n*10 + uint(c-'0')
	}
	return n, nil
}
