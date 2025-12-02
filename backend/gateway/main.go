package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// ServiceRegistry 服务注册表
var serviceRegistry = map[string]string{
	"user":           "http://localhost:8081",
	"talent":         "http://localhost:8082",
	"job":            "http://localhost:8083",
	"resume":         "http://localhost:8084",
	"recommendation": "http://localhost:8085",
	"message":        "http://localhost:8086",
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

	log.Println("API Gateway is running on :8080")
	log.Println("Registered services:", serviceRegistry)
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start gateway:", err)
	}
}
