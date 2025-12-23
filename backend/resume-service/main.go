package main

import (
	"log"
	"resume-service/handlers"
	"resume-service/models"

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

	if err := db.AutoMigrate(&models.Resume{}, &models.Application{}); err != nil {
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

	resumeHandler := handlers.NewResumeHandler(db)

	api := r.Group("/api/v1")
	{
		// Resume routes
		resumes := api.Group("/resumes")
		{
			resumes.POST("", resumeHandler.UploadResume)
			resumes.GET("", resumeHandler.ListResumes)
			resumes.GET("/:id", resumeHandler.GetResume)
			resumes.DELETE("/:id", resumeHandler.DeleteResume)
			resumes.POST("/parse", resumeHandler.ParseResume)
			resumes.POST("/match", resumeHandler.MatchResumeToJob)
		}

		// Application routes
		applications := api.Group("/applications")
		{
			applications.POST("", resumeHandler.CreateApplication)
			applications.GET("", resumeHandler.ListApplications)
			applications.PUT("/:id", resumeHandler.UpdateApplication)
		}
	}

	log.Println("Resume service is running on :8084")
	if err := r.Run(":8084"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
