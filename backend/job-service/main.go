package main

import (
	"job-service/handlers"
	"job-service/models"
	"log"

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

	if err := db.AutoMigrate(&models.Job{}); err != nil {
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

	jobHandler := handlers.NewJobHandler(db)

	api := r.Group("/api/v1/jobs")
	{
		api.POST("", jobHandler.CreateJob)
		api.GET("", jobHandler.ListJobs)
		api.GET("/stats", jobHandler.GetJobStats)
		api.GET("/:id", jobHandler.GetJob)
		api.PUT("/:id", jobHandler.UpdateJob)
		api.DELETE("/:id", jobHandler.DeleteJob)
	}

	log.Println("Job service is running on :8083")
	if err := r.Run(":8083"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
