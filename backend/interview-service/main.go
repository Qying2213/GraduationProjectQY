package main

import (
	"fmt"
	"interview-service/handlers"
	"interview-service/models"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 数据库配置
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "talent_platform")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(&models.Interview{}, &models.InterviewFeedback{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 初始化处理器
	interviewHandler := handlers.NewInterviewHandler(db)

	// 创建路由
	r := gin.Default()

	// CORS 中间件
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

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "interview-service"})
	})

	// API 路由
	api := r.Group("/api/v1")
	{
		interviews := api.Group("/interviews")
		{
			interviews.POST("", interviewHandler.CreateInterview)
			interviews.GET("", interviewHandler.ListInterviews)
			interviews.GET("/stats", interviewHandler.GetInterviewStats)
			interviews.GET("/today", interviewHandler.GetTodayInterviews)
			interviews.GET("/interviewer/:interviewer_id", interviewHandler.GetInterviewerSchedule)
			interviews.GET("/:id", interviewHandler.GetInterview)
			interviews.PUT("/:id", interviewHandler.UpdateInterview)
			interviews.DELETE("/:id", interviewHandler.DeleteInterview)
			interviews.POST("/:id/cancel", interviewHandler.CancelInterview)
			interviews.POST("/:id/complete", interviewHandler.CompleteInterview)
		}
	}

	// 启动服务
	port := getEnv("PORT", "8087")
	log.Printf("Interview Service is running on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
