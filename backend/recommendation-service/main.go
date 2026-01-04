package main

import (
	"log"
	"os"
	"recommendation-service/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 连接数据库
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=talent_platform port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
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

	recommendHandler := handlers.NewRecommendationHandler(db)

	api := r.Group("/api/v1/recommendations")
	{
		api.POST("/jobs-for-talent", recommendHandler.RecommendJobsForTalent)
		api.POST("/talents-for-job", recommendHandler.RecommendTalentsForJob)
		api.GET("/stats", recommendHandler.GetRecommendationStats)
	}

	log.Println("Recommendation service is running on :8087")
	if err := r.Run(":8087"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
