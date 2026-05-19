package main

import (
	"fmt"
	"log"
	"message-service/handlers"
	"message-service/models"
	"message-service/websocket"
	"net/http"
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

	// 自动迁移消息、通知、会话和聊天消息表，保证本地启动时数据库结构可用。
	if err := db.AutoMigrate(&models.Message{}, &models.Notice{}, &models.Conversation{}, &models.ChatMessage{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 初始化 WebSocket Hub，统一维护在线连接、实时聊天投递和在线状态。
	hub := websocket.NewHub()
	go hub.Run()

	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.SimpleOperationLog("message-service"))

	messageHandler := handlers.NewMessageHandler(db)
	chatHandler := handlers.NewChatHandler(db, hub)
	noticeHandler := handlers.NewNoticeHandler(db)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "message-service", "port": 8085})
	})

	// WebSocket 实时聊天入口。JWTAuth 会从 Header 或 query token 中解析用户身份。
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

	// 在线用户列表接口，供前端展示当前在线用户总览。
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

	// 查询指定用户在线状态，供会话列表或用户卡片显示在线标识。
	r.GET("/api/v1/online-status/:user_id", middleware.JWTAuth(), func(c *gin.Context) {
		var userID uint
		if _, err := c.Params.Get("user_id"); err {
			c.JSON(400, gin.H{"error": "invalid user_id"})
			return
		}
		// 从 URL 中解析 user_id。
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

	// 站内信接口：处理系统消息、未读数、标记已读和删除等功能。
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

	// 服务间通知内部接口：其他服务可通过 X-Internal-Token 写入站内信。
	r.POST("/internal/messages", func(c *gin.Context) {
		internalAPIKey := getEnv("INTERNAL_API_KEY", "")
		if internalAPIKey != "" && c.GetHeader("X-Internal-Token") != internalAPIKey {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "error": "invalid internal token"})
			return
		}
		messageHandler.SendMessage(c)
	})

	// 聊天会话接口：提供会话列表、未读数、消息分页、发送消息和已读状态。
	chatAPI := r.Group("/api/v1/conversations")
	chatAPI.Use(middleware.JWTAuth())
	{
		// GET /conversations - 获取会话列表，按最新消息时间排序并带未读数。
		chatAPI.GET("", chatHandler.GetConversations)

		// POST /conversations - 创建/获取会话
		// 如果双方会话已存在则复用，否则创建新会话。
		chatAPI.POST("", chatHandler.CreateOrGetConversation)

		// GET /conversations/unread-count - 获取总未读数
		// 汇总当前用户所有会话的未读消息数量。
		chatAPI.GET("/unread-count", chatHandler.GetTotalUnreadCount)

		// GET /conversations/:id/messages - 获取会话消息列表（分页）
		// 分页返回会话消息，按聊天展示需要保持时间顺序。
		chatAPI.GET("/:id/messages", chatHandler.GetMessages)

		// POST /conversations/:id/messages - 发送消息
		// 持久化新消息，并通过 WebSocket 推送给会话订阅者。
		chatAPI.POST("/:id/messages", chatHandler.SendMessage)

		// PUT /conversations/:id/read - 标记会话已读
		// 将当前会话中对方发送的未读消息批量标记为已读。
		chatAPI.PUT("/:id/read", chatHandler.MarkAsRead)
	}
	noticeAPI := r.Group("/api/v1/notices")
	noticeAPI.Use(middleware.JWTAuth(), middleware.RoleAuth("admin", "hr", "hr_manager", "recruiter"))
	{
		noticeAPI.GET("", noticeHandler.ListNotices)
		noticeAPI.GET("/:id", noticeHandler.GetNotice)
		noticeAPI.POST("", noticeHandler.CreateNotice)
		noticeAPI.PUT("/:id", noticeHandler.UpdateNotice)
		noticeAPI.DELETE("/:id", noticeHandler.DeleteNotice)

	}

	log.Println("Message service is running on :8085")
	if err := r.Run(":8085"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// parseUint 将路径参数中的数字字符串解析为 uint。
func parseUint(s string) (uint, error) {
	if s == "" {
		return 0, fmt.Errorf("empty input")
	}

	var n uint
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid character: %c", c)
		}
		n = n*10 + uint(c-'0')
	}
	return n, nil
}
