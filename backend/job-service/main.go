package main

import (
	"fmt"
	"log"
	"os"

	"common/database"
	"common/middleware"
	"job-service/handlers"
	"job-service/models"

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
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	dsn := "host=" + dbHost + " user=" + dbUser + " dbname=" + dbName + " port=" + dbPort + " sslmode=" + dbSSLMode + " TimeZone=Asia/Shanghai"
	if dbPassword != "" {
		dsn = "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=" + dbSSLMode + " TimeZone=Asia/Shanghai"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// job-service 拥有职位主数据，启动时保证 jobs 表结构可用。
	if err := db.AutoMigrate(&models.Job{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Redis 是职位详情缓存和热门职位排行的增强依赖；连接失败时仅关闭缓存能力。
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

	r.Use(middleware.CORS())
	r.Use(middleware.SimpleOperationLog("job-service"))

	jobHandler := handlers.NewJobHandler(db)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "job-service", "port": 8082})
	})

	// 公开职位接口：候选人和未登录用户也可以浏览职位列表、详情和热门职位。
	api := r.Group("/api/v1/jobs")
	{
		api.GET("", jobHandler.ListJobs)
		api.GET("/stats", jobHandler.GetJobStats)
		api.GET("/hot", jobHandler.GetHotJobs)
		api.GET("/:id", jobHandler.GetJob)
	}

	// 职位管理接口：创建、编辑、删除和查看投递列表，需要招聘侧角色。
	protected := r.Group("/api/v1/jobs")
	protected.Use(middleware.JWTAuth(), middleware.RoleAuth("admin", "hr", "hr_manager", "recruiter"))
	{
		protected.POST("", jobHandler.CreateJob)
		protected.PUT("/:id", jobHandler.UpdateJob)
		protected.DELETE("/:id", jobHandler.DeleteJob)
		protected.GET("/:id/applications", jobHandler.GetJobApplications)
	}

	log.Println("Job service is running on :8082")
	if err := r.Run(":8082"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
