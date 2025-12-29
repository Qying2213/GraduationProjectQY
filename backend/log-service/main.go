package main

import (
	"log"
	"os"

	"log-service/handlers"

	"github.com/gin-gonic/gin"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	r := gin.Default()

	// CORS中间件
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

	// 初始化处理器
	logHandler := handlers.NewLogHandler()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "log-service",
			"port":    8088,
		})
	})

	// 日志API
	api := r.Group("/api/v1/logs")
	{
		api.GET("", logHandler.QueryLogs)            // 查询日志
		api.GET("/stats", logHandler.GetStats)       // 获取统计
		api.GET("/services", logHandler.GetServices) // 获取服务列表
		api.GET("/actions", logHandler.GetActions)   // 获取操作类型
		api.DELETE("/cleanup", logHandler.Cleanup)   // 清理旧日志
	}

	port := getEnv("LOG_SERVICE_PORT", "8088")
	log.Printf("Log service is running on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
