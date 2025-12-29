package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"

	"common/elasticsearch"

	"github.com/gin-gonic/gin"
)

// 响应体写入器
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// OperationLogConfig 日志中间件配置
type OperationLogConfig struct {
	ServiceName     string
	SkipPaths       []string // 跳过的路径
	LogRequestBody  bool     // 是否记录请求体
	LogResponseBody bool     // 是否记录响应体
	MaxBodySize     int      // 最大记录的body大小
}

// DefaultOperationLogConfig 默认配置
func DefaultOperationLogConfig(serviceName string) *OperationLogConfig {
	return &OperationLogConfig{
		ServiceName:     serviceName,
		SkipPaths:       []string{"/health", "/metrics", "/favicon.ico"},
		LogRequestBody:  true,
		LogResponseBody: false,
		MaxBodySize:     4096,
	}
}

// OperationLog 操作日志中间件
func OperationLog(config *OperationLogConfig) gin.HandlerFunc {
	logService := elasticsearch.NewLogService(config.ServiceName)

	return func(c *gin.Context) {
		// 检查是否跳过
		for _, path := range config.SkipPaths {
			if strings.HasPrefix(c.Request.URL.Path, path) {
				c.Next()
				return
			}
		}

		start := time.Now()

		// 读取请求体
		var requestBody string
		if config.LogRequestBody && c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			if len(bodyBytes) > config.MaxBodySize {
				bodyBytes = bodyBytes[:config.MaxBodySize]
			}
			requestBody = string(bodyBytes)
			// 重新设置请求体
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 包装响应写入器
		var responseBody string
		if config.LogResponseBody {
			rw := &responseWriter{body: bytes.NewBuffer(nil), ResponseWriter: c.Writer}
			c.Writer = rw
			defer func() {
				body := rw.body.Bytes()
				if len(body) > config.MaxBodySize {
					body = body[:config.MaxBodySize]
				}
				responseBody = string(body)
			}()
		}

		// 处理请求
		c.Next()

		// 计算耗时
		duration := time.Since(start).Milliseconds()

		// 获取用户信息
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")

		// 确定日志级别
		level := "info"
		if c.Writer.Status() >= 500 {
			level = "error"
		} else if c.Writer.Status() >= 400 {
			level = "warn"
		}

		// 获取错误信息
		var errorMsg string
		if len(c.Errors) > 0 {
			errorMsg = c.Errors.String()
		}

		// 解析操作类型和模块
		action, module := parseActionAndModule(c.Request.Method, c.Request.URL.Path)

		// 构建日志
		log := &elasticsearch.OperationLog{
			Timestamp:    time.Now(),
			Service:      config.ServiceName,
			IP:           c.ClientIP(),
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			Query:        c.Request.URL.RawQuery,
			StatusCode:   c.Writer.Status(),
			Duration:     duration,
			RequestBody:  requestBody,
			ResponseBody: responseBody,
			UserAgent:    c.Request.UserAgent(),
			Action:       action,
			Module:       module,
			Level:        level,
			ErrorMsg:     errorMsg,
		}

		if userID != nil {
			if id, ok := userID.(uint); ok {
				log.UserID = id
			}
		}
		if username != nil {
			if name, ok := username.(string); ok {
				log.Username = name
			}
		}

		// 异步记录日志
		logService.LogAsync(log)
	}
}

// parseActionAndModule 解析操作类型和模块
func parseActionAndModule(method, path string) (action, module string) {
	// 解析模块
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 3 {
		module = parts[2] // /api/v1/users -> users
	} else if len(parts) >= 1 {
		module = parts[0]
	}

	// 解析操作类型
	switch method {
	case "GET":
		if strings.Contains(path, "/") && !strings.HasSuffix(path, "s") {
			action = "查看"
		} else {
			action = "查询"
		}
	case "POST":
		if strings.Contains(path, "login") {
			action = "登录"
		} else if strings.Contains(path, "register") {
			action = "注册"
		} else if strings.Contains(path, "upload") {
			action = "上传"
		} else {
			action = "创建"
		}
	case "PUT", "PATCH":
		action = "更新"
	case "DELETE":
		action = "删除"
	default:
		action = method
	}

	return action, module
}

// SimpleOperationLog 简化版日志中间件
func SimpleOperationLog(serviceName string) gin.HandlerFunc {
	return OperationLog(DefaultOperationLogConfig(serviceName))
}
