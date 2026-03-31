package main

import (
	"log"
	"os"
	"resume-service/handlers"

	"common/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// getEnv 统一读取环境变量。
// 这样本地开发时可以直接使用默认值，部署时再由外部环境覆盖。
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// main 是 resume-service 的组装入口。
// 新人第一次读这个服务时，建议先看这里：
// 1. 服务依赖了哪些 handler
// 2. 暴露了哪些路由分组
// 3. 哪些接口要求 JWT 认证
func main() {
	// 加载 .env 文件 - 从 backend 目录加载
	// 这里按“离当前进程最近 -> 项目根目录”的顺序尝试，方便在不同目录下启动服务。
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

	// AutoMigrate 已禁用 - 表结构通过 SQL 脚本管理
	// 如需迁移，请运行: psql -d talent_platform -f backend/database/schema.sql
	// if err := db.AutoMigrate(&models.Resume{}, &models.Application{}, &models.EvaluationResult{}); err != nil {
	// 	log.Fatal("Failed to migrate database:", err)
	// }

	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.SimpleOperationLog("resume-service"))

	// 这里集中初始化本服务会用到的所有业务 handler。
	// 可以把它理解成“接口层”的依赖注入位置。
	resumeHandler := handlers.NewResumeHandler(db)
	aiHandler := handlers.NewAIEvaluateHandler(db)
	aiParseHandler := handlers.NewAIParseHandler(db)
	evalHandler := handlers.NewEvaluationHandler(db)
	onlineResumeHandler := handlers.NewOnlineResumeHandler(db)
	riskHandler := handlers.NewRiskCheckHandler()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "resume-service", "port": 8084})
	})

	api := r.Group("/api/v1")
	{
		// Resume routes
		// 这组接口负责“简历文件生命周期”：
		// 上传、查询、下载、删除、状态更新、文本解析、岗位匹配。
		resumes := api.Group("/resumes")
		{
			resumes.POST("", resumeHandler.UploadResume)
			resumes.POST("/upload", middleware.JWTAuth(), resumeHandler.UploadResumeFile)
			resumes.GET("", resumeHandler.ListResumes)
			resumes.GET("/evaluation", resumeHandler.ListResumesForEvaluation) // 用于自动评估系统
			resumes.GET("/file/:filename", resumeHandler.ServeResumeFile)      // 提供文件访问
			resumes.GET("/:id", resumeHandler.GetResume)
			resumes.GET("/:id/download", resumeHandler.DownloadResume)
			resumes.DELETE("/:id", resumeHandler.DeleteResume)
			resumes.PUT("/:id/status", resumeHandler.UpdateResumeStatus) // 更新简历状态
			resumes.POST("/parse", resumeHandler.ParseResume)
			resumes.POST("/match", resumeHandler.MatchResumeToJob)

			// 风控检查接口
			resumes.POST("/risk-check", riskHandler.CheckResumeRisk)
			resumes.POST("/risk-check/time-conflict", riskHandler.CheckTimeConflict)
			resumes.POST("/risk-check/education", riskHandler.CheckEducationFraud)

			// 在线简历接口 (需要JWT认证)
			// Requirements: 4.3 (持久化), 4.6 (字段验证)
			resumes.GET("/online", middleware.JWTAuth(), onlineResumeHandler.GetOnlineResume)
			resumes.PUT("/online", middleware.JWTAuth(), onlineResumeHandler.SaveOnlineResume)
		}

		// AI Evaluation routes
		// 这组接口是 AI 能力入口，负责触发评估、OCR、解析和查询评估结果。
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
		// 这组接口偏“结果视图”，用于后台页面查看和清理 AI 评估记录。
		evaluations := api.Group("/evaluations")
		{
			evaluations.GET("", evalHandler.ListEvaluations)
			evaluations.GET("/stats", evalHandler.GetEvaluationStats)
			evaluations.GET("/:id/process", evalHandler.GetEvaluationProcess)
			evaluations.GET("/:id", evalHandler.GetEvaluation)
			evaluations.DELETE("/:id", evalHandler.DeleteEvaluation)
		}

		// Application routes (需要JWT认证)
		// 这组接口负责“职位投递生命周期”。
		// 因为会涉及当前登录用户，所以整组接口都挂了 JWTAuth。
		applications := api.Group("/applications")
		applications.Use(middleware.JWTAuth())
		{
			applications.POST("", resumeHandler.CreateApplication)
			applications.GET("", resumeHandler.ListApplications)
			applications.PUT("/:id", resumeHandler.UpdateApplication)
			applications.DELETE("/:id", resumeHandler.DeleteApplication)
		}
	}

	log.Println("Resume service is running on :8084")
	if err := r.Run(":8084"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
