package main

import (
	"log"
	"os"
	"user-service/handlers"

	"common/middleware"

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

	// user-service 是身份与用户资料服务，负责注册、登录、个人资料和后台用户列表。
	// 其他服务只消费 JWT 中的用户信息，不直接承担账号密码校验。
	r := gin.Default()

	// CORS中间件
	r.Use(middleware.CORS())

	// 操作日志中间件（ES）
	r.Use(middleware.SimpleOperationLog("user-service"))

	// 初始化处理器
	userHandler := handlers.NewUserHandler(db)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "user-service",
			"port":    8081,
		})
	})

	// 公开路由：注册和登录不需要 JWT，登录成功后才签发 token。
	public := r.Group("/api/v1")
	{
		public.POST("/register", userHandler.Register)
		public.POST("/login", userHandler.Login)
	}

	// 个人资料路由：只允许当前登录用户查看和更新自己的资料。
	auth := r.Group("/api/v1")
	auth.Use(middleware.JWTAuth())
	{
		auth.GET("/profile", userHandler.GetProfile)
		auth.PUT("/profile", userHandler.UpdateProfile)
	}

	// 后台用户列表：仅管理员或招聘侧角色可查看，用于系统管理页面。
	adminOrHR := r.Group("/api/v1")
	adminOrHR.Use(middleware.JWTAuth(), middleware.RoleAuth("admin", "hr", "hr_manager", "recruiter"))
	{
		adminOrHR.GET("/users", userHandler.ListUsers)
	}

	log.Println("User service is running on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
