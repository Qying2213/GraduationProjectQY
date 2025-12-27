package main

import (
	"log"
	"os"
	"resume-service/handlers"
	"resume-service/models"

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
	// 数据库连接（支持环境变量配置）
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
	aiHandler := handlers.NewAIEvaluateHandler(db)

	api := r.Group("/api/v1")
	{
		// Resume routes
		resumes := api.Group("/resumes")
		{
			resumes.POST("", resumeHandler.UploadResume)
			resumes.POST("/upload", resumeHandler.UploadResumeFile)
			resumes.GET("", resumeHandler.ListResumes)
			resumes.GET("/evaluation", resumeHandler.ListResumesForEvaluation) // 用于自动评估系统
			resumes.GET("/:id", resumeHandler.GetResume)
			resumes.GET("/:id/download", resumeHandler.DownloadResume)
			resumes.DELETE("/:id", resumeHandler.DeleteResume)
			resumes.PUT("/:id/status", resumeHandler.UpdateResumeStatus) // 更新简历状态
			resumes.POST("/parse", resumeHandler.ParseResume)
			resumes.POST("/match", resumeHandler.MatchResumeToJob)
		}

		// AI Evaluation routes
		ai := api.Group("/ai")
		{
			ai.GET("/config", aiHandler.CheckAIConfig)
			ai.POST("/evaluate", aiHandler.EvaluateByResumeID)
			ai.POST("/evaluate/upload", aiHandler.EvaluateUploadedFile)
			ai.POST("/evaluate/batch", aiHandler.BatchEvaluate)
			ai.GET("/evaluate/:id/result", aiHandler.GetEvaluationResult)
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
