package main

import (
	"log"
	"os"
	"resume-service/handlers"
	"resume-service/models"

	"common/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	// 加载 .env 文件 - 从 backend 目录加载
	godotenv.Load("../.env")    // backend/.env (从 resume-service 目录)
	godotenv.Load(".env")       // 当前目录
	godotenv.Load("../../.env") // 项目根目录
	// 数据库连接（支持环境变量配置）
	dbHost := getEnv("DB_HOST", "localhost")
	dbUser := getEnv("DB_USER", "qinyang")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "talent_platform")
	dbPort := getEnv("DB_PORT", "5432")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	dsn := "host=" + dbHost + " user=" + dbUser + " dbname=" + dbName + " port=" + dbPort + " sslmode=" + dbSSLMode + " TimeZone=Asia/Shanghai"
	if dbPassword != "" {
		dsn = "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=" + dbSSLMode + " TimeZone=Asia/Shanghai"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	if err := db.AutoMigrate(&models.Resume{}, &models.Application{}, &models.EvaluationResult{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.SimpleOperationLog("resume-service"))

	resumeHandler := handlers.NewResumeHandler(db)
	aiHandler := handlers.NewAIEvaluateHandler(db)
	aiParseHandler := handlers.NewAIParseHandler(db)
	evalHandler := handlers.NewEvaluationHandler(db)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "resume-service", "port": 8084})
	})

	api := r.Group("/api/v1")
	{
		// Resume routes
		resumes := api.Group("/resumes")
		{
			resumes.POST("", resumeHandler.UploadResume)
			resumes.POST("/upload", resumeHandler.UploadResumeFile)
			resumes.GET("", resumeHandler.ListResumes)
			resumes.GET("/evaluation", resumeHandler.ListResumesForEvaluation) // 用于自动评估系统
			resumes.GET("/file/:filename", resumeHandler.ServeResumeFile)      // 提供文件访问
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
			ai.GET("/current-task", aiHandler.GetCurrentTask)
			ai.POST("/evaluate", aiHandler.EvaluateByResumeID)
			ai.POST("/evaluate/upload", aiHandler.EvaluateUploadedFile)
			ai.POST("/evaluate/batch", aiHandler.BatchEvaluate)
			ai.GET("/evaluate/:id/result", aiHandler.GetEvaluationResult)
			ai.POST("/parse", aiParseHandler.AIParseResume) // AI智能解析
			ai.POST("/ocr", aiParseHandler.OCRExtract)      // OCR文本提取
		}

		// Evaluation Results routes (评估结果管理)
		evaluations := api.Group("/evaluations")
		{
			evaluations.GET("", evalHandler.ListEvaluations)
			evaluations.GET("/stats", evalHandler.GetEvaluationStats)
			evaluations.GET("/:id", evalHandler.GetEvaluation)
			evaluations.DELETE("/:id", evalHandler.DeleteEvaluation)
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
