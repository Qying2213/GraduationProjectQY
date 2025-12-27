package middleware

import (
	"strings"

	"evaluator-service/internal/service"

	"github.com/gin-gonic/gin"
)

// Auth 认证中间件，验证 JWT token 并将用户信息存入 context
func Auth(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization header 提取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "未授权访问"})
			return
		}

		// 检查 Bearer 前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"error": "无效的授权格式"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "token 不能为空"})
			return
		}

		// 验证 token
		claims, err := authSvc.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "token 无效或已过期"})
			return
		}

		// 将用户信息存入 context
		c.Set("user_id", claims.UserID)
		c.Set("corp_code", claims.CorpCode)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// GetUserID 从 context 获取当前用户 ID
func GetUserID(c *gin.Context) uint {
	if v, exists := c.Get("user_id"); exists {
		if id, ok := v.(uint); ok {
			return id
		}
	}
	return 0
}

// GetCorpCode 从 context 获取当前用户的组织代码
func GetCorpCode(c *gin.Context) string {
	if v, exists := c.Get("corp_code"); exists {
		if code, ok := v.(string); ok {
			return code
		}
	}
	return ""
}

// GetUsername 从 context 获取当前用户名
func GetUsername(c *gin.Context) string {
	if v, exists := c.Get("username"); exists {
		if name, ok := v.(string); ok {
			return name
		}
	}
	return ""
}
