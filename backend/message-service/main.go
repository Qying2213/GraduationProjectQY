package main

import (
	"log"
	"message-service/handlers"
	"message-service/models"
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

	if err := db.AutoMigrate(&models.Message{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.SimpleOperationLog("message-service"))

	messageHandler := handlers.NewMessageHandler(db)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "message-service", "port": 8085})
	})

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

	log.Println("Message service is running on :8085")
	if err := r.Run(":8085"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
