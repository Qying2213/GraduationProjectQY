// recommendation-service 负责人岗匹配、语义推荐和 RAG 检索。
// 前端通过 API Gateway 调用它；resume-service 在 AI 评估链路中会通过内部接口调用它。
package main

import (
	"common/middleware"
	"fmt"
	"log"
	"net/http"
	"os"
	"recommendation-service/handlers"

	"common/database"

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
	// 加载 .env 文件
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// 打印 Embedding 配置状态。未配置时会使用 mock embedding，这样本地演示
	// 不会因为外部 API Key 缺失而完全不可用。
	arkAPIKey := os.Getenv("ARK_API_KEY")
	volcModelID := os.Getenv("VOLC_MODEL_ID")
	if arkAPIKey != "" {
		log.Printf("Volcengine Embedding configured: model=%s, api_key=%s...", volcModelID, arkAPIKey[:min(10, len(arkAPIKey))])
	} else {
		log.Println("Warning: ARK_API_KEY not set, embedding will use mock mode")
	}

	// 数据库连接（与其他服务保持一致的配置方式）
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
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected")

	// RAG 索引表在服务启动时自动保证存在，降低本地部署和答辩演示的准备成本。
	if err := ensureRAGTables(db); err != nil {
		log.Printf("Warning: Failed to ensure RAG tables: %v", err)
	}

	// 初始化 Redis 缓存
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	if err := database.InitRedis(database.RedisConfig{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       0,
	}); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v (caching disabled)", err)
	} else {
		log.Println("Redis connected, caching enabled")
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowedOrigins := map[string]bool{
			"http://localhost:3000": true,
			"http://localhost:5173": true,
			"http://127.0.0.1:3000": true,
			"http://127.0.0.1:5173": true,
		}
		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}
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
	api.Use(middleware.JWTAuth(), middleware.RoleAuth("admin", "hr", "hr_manager", "recruiter"))
	{
		api.POST("/jobs-for-talent", recommendHandler.RecommendJobsForTalent)
		api.POST("/talents-for-job", recommendHandler.RecommendTalentsForJob)
		api.GET("/stats", recommendHandler.GetRecommendationStats)
		api.POST("/batch", recommendHandler.BatchRecommend)
		api.POST("/attribution-report", recommendHandler.GenerateAttributionReport)
		api.POST("/semantic-match", recommendHandler.SemanticMatch)

		// RAG 相关接口：给后台管理页面使用，需要登录和 HR/管理员权限。
		api.POST("/rag/query", recommendHandler.RAGQuery)
		api.POST("/rag/index-talent", recommendHandler.IndexTalent)
		api.POST("/rag/index-job", recommendHandler.IndexJob)
		api.POST("/rag/index-resume", recommendHandler.IndexResume)
		api.POST("/rag/index-all", recommendHandler.IndexAll)
		api.POST("/rag/match", recommendHandler.RAGMatch)
	}

	internal := r.Group("/internal/recommendations")
	internal.Use(func(c *gin.Context) {
		internalAPIKey := getEnv("INTERNAL_API_KEY", "")
		if internalAPIKey != "" && c.GetHeader("X-Internal-Token") != internalAPIKey {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "error": "invalid internal token"})
			c.Abort()
			return
		}
		c.Next()
	})
	{
		// 内部接口给 resume-service 调用。使用 X-Internal-Token 可以把服务间
		// 调用和用户 JWT 区分开，避免把后台用户权限直接暴露给内部链路。
		internal.POST("/rag/query", recommendHandler.RAGQuery)
		internal.POST("/rag/index-talent", recommendHandler.IndexTalent)
		internal.POST("/rag/index-job", recommendHandler.IndexJob)
		internal.POST("/rag/index-resume", recommendHandler.IndexResume)
		internal.POST("/rag/index-all", recommendHandler.IndexAll)
		internal.POST("/rag/match", recommendHandler.RAGMatch)
	}

	log.Println("Recommendation service is running on :8087")
	if err := r.Run(":8087"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// ensureRAGTables 创建项目使用的轻量级向量索引表。
// 当前向量以 JSONB 存储，便于普通 PostgreSQL 环境直接部署；相似度在 Go 中计算，
// 对毕业设计演示规模已经足够。
func ensureRAGTables(db *gorm.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS talent_embeddings (
			id SERIAL PRIMARY KEY,
			talent_id INTEGER REFERENCES talents(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			embedding JSONB NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS job_embeddings (
			id SERIAL PRIMARY KEY,
			job_id INTEGER REFERENCES jobs(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			embedding JSONB NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS resume_embeddings (
			id SERIAL PRIMARY KEY,
			resume_id INTEGER REFERENCES resumes(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			embedding JSONB NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`ALTER TABLE talent_embeddings ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP`,
		`ALTER TABLE job_embeddings ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP`,
		`ALTER TABLE resume_embeddings ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_talent_embeddings_talent_id ON talent_embeddings(talent_id)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_job_embeddings_job_id ON job_embeddings(job_id)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_resume_embeddings_resume_id ON resume_embeddings(resume_id)`,
	}

	for _, statement := range statements {
		if err := db.Exec(statement).Error; err != nil {
			return err
		}
	}
	return nil
}
