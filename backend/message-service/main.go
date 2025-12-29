package main

import (
	"log"
	"message-service/handlers"
	"message-service/models"

	"common/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=qinyang dbname=talent_platform port=5432 sslmode=disable TimeZone=Asia/Shanghai"
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
	{
		api.POST("", messageHandler.SendMessage)
		api.GET("", messageHandler.GetMessages)
		api.GET("/unread-count", messageHandler.GetUnreadCount)
		api.PUT("/:id/read", messageHandler.MarkAsRead)
		api.DELETE("/:id", messageHandler.DeleteMessage)
	}

	log.Println("Message service is running on :8085")
	if err := r.Run(":8085"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
