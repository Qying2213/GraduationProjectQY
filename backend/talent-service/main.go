package main

import (
	"log"
	"os"
	"talent-service/handlers"

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
	// 加载 .env 文件，确保单独启动服务时也能拿到与其他服务一致的 JWT/DB 配置。
	godotenv.Load("../.env")
	godotenv.Load(".env")
	godotenv.Load("../../.env")

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

	// talents 表由 databaseSQL/schema.sql 管理，避免运行时 AutoMigrate
	// 在当前 GORM/Postgres 驱动组合下触发数组字段迁移兼容问题。

	// talent-service 是人才库服务，负责人才档案的创建、筛选、搜索、统计和维护。
	// 只有招聘侧角色可以访问，候选人自助信息会通过 user/resume 服务间接关联到这里。
	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.SimpleOperationLog("talent-service"))

	talentHandler := handlers.NewTalentHandler(db)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "talent-service", "port": 8086})
	})

	api := r.Group("/api/v1/talents")
	api.Use(middleware.JWTAuth(), middleware.RoleAuth("admin", "hr", "hr_manager", "recruiter"))
	{
		// 人才库核心接口：围绕“人才档案”做 CRUD、条件筛选和统计看板数据。
		api.POST("", talentHandler.CreateTalent)
		api.GET("", talentHandler.ListTalents)
		api.GET("/stats", talentHandler.GetTalentStats)
		api.GET("/search", talentHandler.SearchTalents)
		api.GET("/:id", talentHandler.GetTalent)
		api.PUT("/:id", talentHandler.UpdateTalent)
		api.DELETE("/:id", talentHandler.DeleteTalent)
	}

	log.Println("Talent service is running on :8086")
	if err := r.Run(":8086"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
