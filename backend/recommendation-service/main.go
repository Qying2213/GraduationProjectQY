package main

import (
	"log"
	"net/http"
	"os"
	"recommendation-service/handlers"

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
	// 数据库连接（与其他服务保持一致的配置方式）
	dbHost := getEnv("DB_HOST", "localhost")
	dbUser := getEnv("DB_USER", "qinyang")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "talent_platform")
	dbPort := getEnv("DB_PORT", "5432")

	dsn := "host=" + dbHost + " user=" + dbUser + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=Asia/Shanghai"
	if dbPassword != "" {
		dsn = "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=Asia/Shanghai"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
	} else {
		log.Println("Database connected")
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

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "recommendation-service", "port": 8087})
	})

	recommendHandler := handlers.NewRecommendationHandler(db)

	api := r.Group("/api/v1/recommendations")
	{
		api.POST("/jobs-for-talent", recommendHandler.RecommendJobsForTalent)
		api.POST("/talents-for-job", recommendHandler.RecommendTalentsForJob)
		api.GET("/stats", recommendHandler.GetRecommendationStats)
		api.POST("/batch", recommendHandler.BatchRecommend)
		api.POST("/attribution-report", recommendHandler.GenerateAttributionReport)
		api.POST("/semantic-match", recommendHandler.SemanticMatch)
	}

	log.Println("Recommendation service is running on :8087")
	if err := r.Run(":8087"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
