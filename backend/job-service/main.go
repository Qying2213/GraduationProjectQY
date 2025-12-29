package main

import (
	"job-service/handlers"
	"job-service/models"
	"log"

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

	if err := db.AutoMigrate(&models.Job{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.SimpleOperationLog("job-service"))

	jobHandler := handlers.NewJobHandler(db)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "job-service", "port": 8082})
	})

	api := r.Group("/api/v1/jobs")
	{
		api.POST("", jobHandler.CreateJob)
		api.GET("", jobHandler.ListJobs)
		api.GET("/stats", jobHandler.GetJobStats)
		api.GET("/:id", jobHandler.GetJob)
		api.PUT("/:id", jobHandler.UpdateJob)
		api.DELETE("/:id", jobHandler.DeleteJob)
	}

	log.Println("Job service is running on :8082")
	if err := r.Run(":8082"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
