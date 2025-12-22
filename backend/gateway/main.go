package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// ==================== 统计处理器 ====================

// StatsHandler 统计处理器
type StatsHandler struct{}

// NewStatsHandler 创建统计处理器
func NewStatsHandler() *StatsHandler {
	return &StatsHandler{}
}

// GetDashboardStats 获取仪表板统计
func (h *StatsHandler) GetDashboardStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"total_talents":      1560,
			"total_jobs":         48,
			"total_applications": 326,
			"total_interviews":   89,
			"match_rate":         89.5,
			"talent_trend":       12.5,
			"job_trend":          8.2,
			"application_trend":  -3.1,
		},
	})
}

// GetRecruitmentFunnel 获取招聘漏斗
func (h *StatsHandler) GetRecruitmentFunnel(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"resumes":     1256,
			"screened":    565,
			"interviewed": 328,
			"passed":      156,
			"hired":       45,
		},
	})
}

// GetChannelStats 获取渠道统计
func (h *StatsHandler) GetChannelStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": []gin.H{
			{"name": "官网投递", "count": 456, "rate": 36},
			{"name": "猎聘网", "count": 312, "rate": 25},
			{"name": "BOSS直聘", "count": 234, "rate": 19},
			{"name": "内部推荐", "count": 156, "rate": 12},
			{"name": "其他渠道", "count": 98, "rate": 8},
		},
	})
}

// GetDepartmentProgress 获取部门招聘进度
func (h *StatsHandler) GetDepartmentProgress(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": []gin.H{
			{"department": "技术部", "target": 20, "hired": 15, "progress": 75},
			{"department": "产品部", "target": 8, "hired": 6, "progress": 75},
			{"department": "设计部", "target": 5, "hired": 5, "progress": 100},
			{"department": "市场部", "target": 10, "hired": 4, "progress": 40},
			{"department": "运营部", "target": 6, "hired": 3, "progress": 50},
		},
	})
}

// GetInterviewerRank 获取面试官排行
func (h *StatsHandler) GetInterviewerRank(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": []gin.H{
			{"name": "陈总监", "department": "技术部", "interviews": 45, "pass_rate": 68, "avg_score": 4.5},
			{"name": "刘经理", "department": "产品部", "interviews": 38, "pass_rate": 72, "avg_score": 4.3},
			{"name": "周主管", "department": "技术部", "interviews": 32, "pass_rate": 65, "avg_score": 4.2},
			{"name": "王总监", "department": "设计部", "interviews": 28, "pass_rate": 78, "avg_score": 4.6},
			{"name": "HR小李", "department": "人力资源", "interviews": 56, "pass_rate": 82, "avg_score": 4.4},
		},
	})
}

// GetTrendData 获取趋势数据
func (h *StatsHandler) GetTrendData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": []gin.H{
			{"date": "2024-01", "resumes": 120, "interviews": 45, "hired": 8},
			{"date": "2024-02", "resumes": 132, "interviews": 52, "hired": 12},
			{"date": "2024-03", "resumes": 101, "interviews": 38, "hired": 6},
			{"date": "2024-04", "resumes": 134, "interviews": 48, "hired": 10},
			{"date": "2024-05", "resumes": 90, "interviews": 35, "hired": 5},
			{"date": "2024-06", "resumes": 230, "interviews": 85, "hired": 18},
			{"date": "2024-07", "resumes": 210, "interviews": 78, "hired": 15},
			{"date": "2024-08", "resumes": 182, "interviews": 68, "hired": 12},
			{"date": "2024-09", "resumes": 191, "interviews": 72, "hired": 14},
			{"date": "2024-10", "resumes": 234, "interviews": 88, "hired": 16},
			{"date": "2024-11", "resumes": 290, "interviews": 108, "hired": 22},
			{"date": "2024-12", "resumes": 330, "interviews": 125, "hired": 28},
		},
	})
}

// GetJobRank 获取职位热度排行
func (h *StatsHandler) GetJobRank(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": []gin.H{
			{"title": "Go开发工程师", "count": 198},
			{"title": "前端工程师", "count": 156},
			{"title": "产品经理", "count": 134},
			{"title": "UI设计师", "count": 112},
			{"title": "数据分析师", "count": 89},
		},
	})
}

// ==================== 服务注册 ====================

// ServiceRegistry 服务注册表
var serviceRegistry = map[string]string{
	"user":           "http://localhost:8081",
	"talent":         "http://localhost:8082",
	"job":            "http://localhost:8083",
	"resume":         "http://localhost:8084",
	"recommendation": "http://localhost:8085",
	"message":        "http://localhost:8086",
	"interview":      "http://localhost:8087",
}

// ReverseProxy 反向代理处理器
func ReverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, err := url.Parse(target)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid service URL"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
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

func main() {
	r := gin.Default()

	// CORS中间件
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

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "healthy",
			"services": serviceRegistry,
		})
	})

	// API版本1路由
	api := r.Group("/api/v1")

	// 用户服务路由
	api.Any("/register", ReverseProxy(serviceRegistry["user"]))
	api.Any("/login", ReverseProxy(serviceRegistry["user"]))
	api.Any("/profile", ReverseProxy(serviceRegistry["user"]))
	api.Any("/users", ReverseProxy(serviceRegistry["user"]))
	api.Any("/users/*path", ReverseProxy(serviceRegistry["user"]))

	// 人才服务路由
	api.Any("/talents", ReverseProxy(serviceRegistry["talent"]))
	api.Any("/talents/*path", ReverseProxy(serviceRegistry["talent"]))

	// 职位服务路由
	api.Any("/jobs", ReverseProxy(serviceRegistry["job"]))
	api.Any("/jobs/*path", ReverseProxy(serviceRegistry["job"]))

	// 简历服务路由
	api.Any("/resumes", ReverseProxy(serviceRegistry["resume"]))
	api.Any("/resumes/*path", ReverseProxy(serviceRegistry["resume"]))
	api.Any("/applications", ReverseProxy(serviceRegistry["resume"]))
	api.Any("/applications/*path", ReverseProxy(serviceRegistry["resume"]))

	// 推荐服务路由
	api.Any("/recommendations", ReverseProxy(serviceRegistry["recommendation"]))
	api.Any("/recommendations/*path", ReverseProxy(serviceRegistry["recommendation"]))

	// 消息服务路由
	api.Any("/messages", ReverseProxy(serviceRegistry["message"]))
	api.Any("/messages/*path", ReverseProxy(serviceRegistry["message"]))

	// 面试服务路由
	api.Any("/interviews", ReverseProxy(serviceRegistry["interview"]))
	api.Any("/interviews/*path", ReverseProxy(serviceRegistry["interview"]))

	// 统计服务路由 (网关内置)
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

	log.Println("API Gateway is running on :8080")
	log.Println("Registered services:", serviceRegistry)
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start gateway:", err)
	}
}
