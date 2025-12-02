package main

import (
	"log"
	"message-service/handlers"
	"message-service/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=talent_platform port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	if err := db.AutoMigrate(&models.Message{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	messageHandler := handlers.NewMessageHandler(db)

	api := r.Group("/api/v1/messages")
	{
		api.POST("", messageHandler.SendMessage)
		api.GET("", messageHandler.GetMessages)
		api.GET("/unread-count", messageHandler.GetUnreadCount)
		api.PUT("/:id/read", messageHandler.MarkAsRead)
		api.DELETE("/:id", messageHandler.DeleteMessage)
	}

	log.Println("Message service is running on :8086")
	if err := r.Run(":8086"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
