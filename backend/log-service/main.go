package main

import (
	"common/middleware"
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
	// log-service 是操作审计服务，负责查询各业务服务写入 Elasticsearch 的操作日志。
	// 它不处理业务写入，只提供日志检索、统计和清理能力。
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

	// 初始化日志查询处理器，内部会连接 Elasticsearch 日志服务。
	logHandler := handlers.NewLogHandler()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "log-service",
			"port":    8088,
		})
	})

	// 日志查看页面（无需登录），方便本地开发和演示快速查看。
	r.GET("/", logHandler.LogViewerPage)
	r.GET("/logs", logHandler.LogViewerPage)

	// 日志 API：后台审计能力，需要管理员或招聘侧角色。
	api := r.Group("/api/v1/logs")
	api.Use(middleware.JWTAuth(), middleware.RoleAuth("admin", "hr", "hr_manager", "recruiter"))
	{
		api.GET("", logHandler.QueryLogs)            // 查询日志
		api.GET("/stats", logHandler.GetStats)       // 获取统计
		api.GET("/services", logHandler.GetServices) // 获取服务列表
		api.GET("/actions", logHandler.GetActions)   // 获取操作类型
		api.DELETE("/cleanup", logHandler.Cleanup)   // 清理旧日志
	}

	port := getEnv("LOG_SERVICE_PORT", "8088")
	log.Printf("Log service is running on :%s", port)
	log.Printf("日志查看页面: http://localhost:%s/logs", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
