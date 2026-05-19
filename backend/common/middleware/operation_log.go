// 操作日志中间件负责把用户行为记录成可写入 Elasticsearch 的结构化日志。
// 对毕业设计来说，它提供了审计能力：创建职位、更新面试、触发 AI 评估等 HR
// 操作都可以按用户、路径、状态码和耗时追踪。
package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"

	"common/elasticsearch"

	"github.com/gin-gonic/gin"
)

// responseWriter 在开启响应体记录时捕获返回内容，同时仍然把数据正常写回客户端。
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// OperationLogConfig 控制操作日志记录哪些请求/响应细节。
// 文件上传会被跳过，因为记录 multipart 内容既影响性能，也可能暴露简历文件。
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

// OperationLog 在每次请求结束后生成一条结构化操作日志。
// 最终写入采用异步方式，避免 Elasticsearch 抖动直接拖慢业务接口。
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

		// 读取请求体 - 跳过 multipart/form-data（文件上传）。读完普通 JSON 后
		// 必须把 Body 放回去，否则后续 handler 将读不到请求内容。
		var requestBody string
		contentType := c.GetHeader("Content-Type")
		isMultipart := strings.HasPrefix(contentType, "multipart/form-data")

		if config.LogRequestBody && c.Request.Body != nil && !isMultipart {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			logBodyBytes := bodyBytes
			if len(logBodyBytes) > config.MaxBodySize {
				logBodyBytes = logBodyBytes[:config.MaxBodySize]
			}
			requestBody = string(logBodyBytes)
			// 重新设置完整请求体，避免后续处理器读取到被截断的 JSON
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		} else if isMultipart {
			requestBody = "[multipart/form-data - skipped]"
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

		// 确定日志级别：根据 HTTP 状态把操作分为 info/warn/error，方便日志页筛选。
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

		// 异步记录日志，避免日志系统抖动拖慢业务接口。
		logService.LogAsync(log)
	}
}

// parseActionAndModule 根据 method/path 推断中文操作类型和业务模块，用于操作日志页面展示。
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
