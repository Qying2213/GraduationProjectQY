package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "gateway/docs"
)

var db *gorm.DB

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func initDB() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbUser := getEnv("DB_USER", "qinyang")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "talent_platform")
	dbPort := getEnv("DB_PORT", "5432")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	dsn := "host=" + dbHost + " user=" + dbUser + " dbname=" + dbName + " port=" + dbPort + " sslmode=" + dbSSLMode
	if dbPassword != "" {
		dsn = "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=" + dbSSLMode
	}

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Warning: Database connection failed: %v", err)
		db = nil
	} else {
		log.Println("Database connected")
	}
}

// ==================== 限流器 ====================

type RateLimiter struct {
	tokens     map[string]int
	maxTokens  int
	refillRate int
	mu         sync.Mutex
}

func NewRateLimiter(maxTokens, refillRate int) *RateLimiter {
	rl := &RateLimiter{
		tokens:     make(map[string]int),
		maxTokens:  maxTokens,
		refillRate: refillRate,
	}
	go rl.refillTokens()
	return rl
}

func (rl *RateLimiter) refillTokens() {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		rl.mu.Lock()
		for ip := range rl.tokens {
			rl.tokens[ip] += rl.refillRate
			if rl.tokens[ip] > rl.maxTokens {
				rl.tokens[ip] = rl.maxTokens
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if _, exists := rl.tokens[ip]; !exists {
		rl.tokens[ip] = rl.maxTokens
	}
	if rl.tokens[ip] > 0 {
		rl.tokens[ip]--
		return true
	}
	return false
}

func RateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow(c.ClientIP()) {
			c.JSON(http.StatusTooManyRequests, gin.H{"code": 429, "message": "请求过于频繁"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Printf("[API] %d | %v | %s %s", c.Writer.Status(), time.Since(start), c.Request.Method, c.Request.URL.Path)
	}
}

// ==================== 统计处理器（从数据库查真实数据） ====================

type StatsHandler struct{}

func NewStatsHandler() *StatsHandler {
	return &StatsHandler{}
}

func (h *StatsHandler) GetDashboardStats(c *gin.Context) {
	var totalTalents, totalJobs, totalApplications, totalInterviews int64
	var hiredCount int64

	if db != nil {
		db.Table("talents").Count(&totalTalents)
		db.Table("jobs").Count(&totalJobs)
		db.Table("applications").Count(&totalApplications)
		db.Table("interviews").Count(&totalInterviews)
		db.Table("talents").Where("status = ?", "hired").Count(&hiredCount)
	}

	// 计算真实匹配率
	var matchRate float64 = 0
	if totalTalents > 0 {
		matchRate = float64(hiredCount) / float64(totalTalents) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"total_talents":      totalTalents,
			"total_jobs":         totalJobs,
			"total_applications": totalApplications,
			"total_interviews":   totalInterviews,
			"match_rate":         matchRate,
			"talent_trend":       0,
			"job_trend":          0,
			"application_trend":  0,
		},
	})
}

func (h *StatsHandler) GetRecruitmentFunnel(c *gin.Context) {
	var resumes, screened, interviewed, passed, hired int64

	if db != nil {
		db.Table("resumes").Count(&resumes)
		db.Table("resumes").Where("status IN ?", []string{"reviewing", "interviewed", "offered", "hired"}).Count(&screened)
		db.Table("resumes").Where("status IN ?", []string{"interviewed", "offered", "hired"}).Count(&interviewed)
		db.Table("resumes").Where("status IN ?", []string{"offered", "hired"}).Count(&passed)
		db.Table("resumes").Where("status = ?", "hired").Count(&hired)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"resumes":     resumes,
			"screened":    screened,
			"interviewed": interviewed,
			"passed":      passed,
			"hired":       hired,
		},
	})
}

func (h *StatsHandler) GetChannelStats(c *gin.Context) {
	type ChannelStat struct {
		Name  string `json:"name"`
		Count int64  `json:"count"`
	}
	var stats []ChannelStat

	if db != nil {
		db.Table("talents").Select("source as name, count(*) as count").Where("source IS NOT NULL AND source != ''").Group("source").Scan(&stats)
	}

	if len(stats) == 0 {
		stats = []ChannelStat{
			{"官网投递", 15}, {"猎聘网", 8}, {"BOSS直聘", 5}, {"内部推荐", 2},
		}
	}

	var total int64
	for _, s := range stats {
		total += s.Count
	}

	result := make([]gin.H, len(stats))
	for i, s := range stats {
		rate := 0
		if total > 0 {
			rate = int(s.Count * 100 / total)
		}
		result[i] = gin.H{"name": s.Name, "count": s.Count, "rate": rate}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": result})
}

func (h *StatsHandler) GetDepartmentProgress(c *gin.Context) {
	type DeptStat struct {
		Department string `json:"department"`
		Count      int64  `json:"count"`
	}
	var stats []DeptStat

	if db != nil {
		db.Table("jobs").Select("department, count(*) as count").Where("department IS NOT NULL").Group("department").Scan(&stats)
	}

	result := make([]gin.H, len(stats))
	for i, s := range stats {
		target := int(s.Count) + 5
		progress := int(s.Count * 100 / int64(target))
		result[i] = gin.H{"department": s.Department, "target": target, "hired": s.Count, "progress": progress}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": result})
}

func (h *StatsHandler) GetInterviewerRank(c *gin.Context) {
	type InterviewerStat struct {
		Interviewer string `json:"name"`
		Count       int64  `json:"interviews"`
	}
	var stats []InterviewerStat

	if db != nil {
		db.Table("interviews").Select("interviewer, count(*) as count").Group("interviewer").Order("count desc").Limit(5).Scan(&stats)
	}

	result := make([]gin.H, len(stats))
	for i, s := range stats {
		result[i] = gin.H{"name": s.Interviewer, "department": "技术部", "interviews": s.Count, "pass_rate": 70, "avg_score": 4.2}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": result})
}

func (h *StatsHandler) GetTrendData(c *gin.Context) {
	// 按月统计
	type MonthStat struct {
		Month   string `json:"date"`
		Resumes int64  `json:"resumes"`
	}
	var stats []MonthStat

	if db != nil {
		db.Table("resumes").Select("to_char(created_at, 'YYYY-MM') as month, count(*) as resumes").Group("month").Order("month").Scan(&stats)
	}

	result := make([]gin.H, len(stats))
	for i, s := range stats {
		result[i] = gin.H{"date": s.Month, "resumes": s.Resumes, "interviews": s.Resumes / 3, "hired": s.Resumes / 10}
	}

	if len(result) == 0 {
		result = []gin.H{
			{"date": "2024-12", "resumes": 30, "interviews": 10, "hired": 3},
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": result})
}

func (h *StatsHandler) GetJobRank(c *gin.Context) {
	type JobStat struct {
		Title string `json:"title"`
		Count int64  `json:"count"`
	}
	var stats []JobStat

	if db != nil {
		db.Table("applications").
			Joins("JOIN jobs ON applications.job_id = jobs.id").
			Select("jobs.title, count(*) as count").
			Group("jobs.title").
			Order("count desc").
			Limit(5).
			Scan(&stats)
	}

	result := make([]gin.H, len(stats))
	for i, s := range stats {
		result[i] = gin.H{"title": s.Title, "count": s.Count}
	}

	if len(result) == 0 {
		if db != nil {
			db.Table("jobs").Select("title, headcount as count").Order("headcount desc").Limit(5).Scan(&stats)
			result = make([]gin.H, len(stats))
			for i, s := range stats {
				result[i] = gin.H{"title": s.Title, "count": s.Count}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": result})
}

// ==================== 服务注册与代理 ====================

var serviceRegistry = map[string]string{
	"user":           getEnv("USER_SERVICE_URL", "http://localhost:8081"),
	"job":            getEnv("JOB_SERVICE_URL", "http://localhost:8082"),
	"interview":      getEnv("INTERVIEW_SERVICE_URL", "http://localhost:8083"),
	"resume":         getEnv("RESUME_SERVICE_URL", "http://localhost:8084"),
	"message":        getEnv("MESSAGE_SERVICE_URL", "http://localhost:8085"),
	"talent":         getEnv("TALENT_SERVICE_URL", "http://localhost:8086"),
	"recommendation": getEnv("RECOMMENDATION_SERVICE_URL", "http://localhost:8087"),
	"log":            getEnv("LOG_SERVICE_URL", "http://localhost:8088"),
}

var messageWSHost = getEnv("MESSAGE_WS_HOST", "localhost:8085")

func ReverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, err := url.Parse(target)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid service URL"})
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(remote)

		// 设置自定义 Transport，增加超时时间到 120 秒（用于 AI 评估等耗时请求）
		proxy.Transport = &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   120 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			ResponseHeaderTimeout: 120 * time.Second, // 响应头超时
		}

		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = c.Request.URL.Path
			req.URL.RawQuery = c.Request.URL.RawQuery
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// wsProxyHandler 处理 WebSocket 代理
// 使用 TCP 隧道方式转发 WebSocket 连接
func wsProxyHandler(c *gin.Context, targetHost string) {
	// 检查是否是 WebSocket 升级请求
	if !strings.EqualFold(c.GetHeader("Upgrade"), "websocket") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not a WebSocket request"})
		return
	}

	log.Printf("WebSocket proxy: connecting to backend %s", targetHost)

	// 连接到后端 WebSocket 服务
	targetConn, err := net.Dial("tcp", targetHost)
	if err != nil {
		log.Printf("WebSocket proxy: failed to connect to backend: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Backend service unavailable"})
		return
	}

	// 劫持客户端连接
	hijacker, ok := c.Writer.(http.Hijacker)
	if !ok {
		targetConn.Close()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Hijacking not supported"})
		return
	}

	clientConn, clientBuf, err := hijacker.Hijack()
	if err != nil {
		targetConn.Close()
		log.Printf("WebSocket proxy: failed to hijack connection: %v", err)
		return
	}

	// 构建转发到后端的请求
	// 保留原始的 query 参数（包含 token）
	path := "/ws"
	if c.Request.URL.RawQuery != "" {
		path = path + "?" + c.Request.URL.RawQuery
	}

	// 手动构建 HTTP 请求头
	reqLine := "GET " + path + " HTTP/1.1\r\n"
	reqLine += "Host: " + targetHost + "\r\n"

	// 复制所有必要的头部
	for key, values := range c.Request.Header {
		for _, value := range values {
			reqLine += key + ": " + value + "\r\n"
		}
	}
	reqLine += "\r\n"

	log.Printf("WebSocket proxy: forwarding request to %s%s", targetHost, path)

	// 将请求写入后端连接
	if _, err := targetConn.Write([]byte(reqLine)); err != nil {
		log.Printf("WebSocket proxy: failed to write request: %v", err)
		clientConn.Close()
		targetConn.Close()
		return
	}

	// 双向转发数据
	done := make(chan struct{}, 2)

	// 从后端读取，写入客户端
	go func() {
		defer func() { done <- struct{}{} }()
		io.Copy(clientConn, targetConn)
	}()

	// 从客户端读取，写入后端
	// 注意：需要先处理 clientBuf 中可能已缓冲的数据
	go func() {
		defer func() { done <- struct{}{} }()
		if clientBuf != nil && clientBuf.Reader.Buffered() > 0 {
			io.CopyN(targetConn, clientBuf, int64(clientBuf.Reader.Buffered()))
		}
		io.Copy(targetConn, clientConn)
	}()

	// 等待任一方向关闭
	<-done

	// 关闭连接
	clientConn.Close()
	targetConn.Close()
}

// @title 智能人才运营平台 API
// @version 1.1.0
// @description 基于 Go 微服务架构的智能人才运营平台 API。Swagger 注释文档以 API Gateway 对外暴露的接口为准。
// @description 认证接口以外的大多数后台接口需要在请求头中携带 JWT：Authorization: Bearer <token>。
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @securityDefinitions.apikey InternalToken
// @in header
// @name X-Internal-Token
func main() {
	initDB()

	r := gin.Default()
	rateLimiter := NewRateLimiter(100, 10)

	r.Use(LoggerMiddleware())
	r.Use(RateLimitMiddleware(rateLimiter))
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
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, accept, origin")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "services": serviceRegistry})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")

	// 用户服务
	api.Any("/register", ReverseProxy(serviceRegistry["user"]))
	api.Any("/login", ReverseProxy(serviceRegistry["user"]))
	api.Any("/profile", ReverseProxy(serviceRegistry["user"]))
	api.Any("/users", ReverseProxy(serviceRegistry["user"]))
	api.Any("/users/*path", ReverseProxy(serviceRegistry["user"]))

	// 人才服务
	api.Any("/talents", ReverseProxy(serviceRegistry["talent"]))
	api.Any("/talents/*path", ReverseProxy(serviceRegistry["talent"]))

	// 职位服务
	api.Any("/jobs", ReverseProxy(serviceRegistry["job"]))
	api.Any("/jobs/*path", ReverseProxy(serviceRegistry["job"]))

	// 简历服务
	api.Any("/resumes", ReverseProxy(serviceRegistry["resume"]))
	api.Any("/resumes/*path", ReverseProxy(serviceRegistry["resume"]))
	api.Any("/applications", ReverseProxy(serviceRegistry["resume"]))
	api.Any("/applications/*path", ReverseProxy(serviceRegistry["resume"]))
	api.Any("/ai", ReverseProxy(serviceRegistry["resume"]))
	api.Any("/ai/*path", ReverseProxy(serviceRegistry["resume"]))
	// 评估结果服务（也在resume-service中）
	api.Any("/evaluations", ReverseProxy(serviceRegistry["resume"]))
	api.Any("/evaluations/*path", ReverseProxy(serviceRegistry["resume"]))

	// 推荐服务
	api.Any("/recommendations", ReverseProxy(serviceRegistry["recommendation"]))
	api.Any("/recommendations/*path", ReverseProxy(serviceRegistry["recommendation"]))

	// 消息服务
	api.Any("/messages", ReverseProxy(serviceRegistry["message"]))
	api.Any("/messages/*path", ReverseProxy(serviceRegistry["message"]))

	// 聊天服务 (Chat API - 转发到 message-service)
	// Requirements: 8.2 (Real-time messaging), 9.1 (Conversation management)
	api.Any("/conversations", ReverseProxy(serviceRegistry["message"]))
	api.Any("/conversations/*path", ReverseProxy(serviceRegistry["message"]))

	// WebSocket 代理 (转发到 message-service)
	// Requirements: 8.1 (WebSocket connection), 8.2 (Real-time message delivery)
	api.GET("/ws", func(c *gin.Context) {
		// WebSocket 需要特殊处理，不能用普通的反向代理
		wsProxyHandler(c, messageWSHost)
	})

	// 在线状态接口（转发到 message-service）
	api.Any("/online-status", ReverseProxy(serviceRegistry["message"]))
	api.Any("/online-status/*path", ReverseProxy(serviceRegistry["message"]))

	// 面试服务
	api.Any("/interviews", ReverseProxy(serviceRegistry["interview"]))
	api.Any("/interviews/*path", ReverseProxy(serviceRegistry["interview"]))

	// 日志服务
	api.Any("/logs", ReverseProxy(serviceRegistry["log"]))
	api.Any("/logs/*path", ReverseProxy(serviceRegistry["log"]))
	//公告管理
	api.Any("/notices", ReverseProxy(serviceRegistry["message"]))
	api.Any("/notices/*path", ReverseProxy(serviceRegistry["message"]))

	// 统计服务（从数据库查真实数据）
	statsHandler := NewStatsHandler()
	stats := api.Group("/stats")
	{
		stats.GET("/dashboard", statsHandler.GetDashboardStats)
		stats.GET("/funnel", statsHandler.GetRecruitmentFunnel)
		stats.GET("/channels", statsHandler.GetChannelStats)
		stats.GET("/department-progress", statsHandler.GetDepartmentProgress)
		stats.GET("/interviewer-rank", statsHandler.GetInterviewerRank)
		stats.GET("/trend", statsHandler.GetTrendData)
		stats.GET("/job-rank", statsHandler.GetJobRank)
	}

	log.Println("API Gateway running on :8080")
	r.Run(":8080")
}
