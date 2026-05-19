// Package middleware 提供多个 Go 微服务共用的 Gin 中间件。
// JWT/RBAC 是系统最重要的安全边界：网关只负责转发请求，真正的身份校验和
// 角色权限判断仍然由各业务服务自己完成。
package middleware

import (
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret     []byte
	jwtSecretOnce sync.Once
)

// getJWTSecret 使用 sync.Once 缓存 JWT 密钥。
// 这样整个进程生命周期内只读取一次环境变量，避免每次请求都重复解析。
func getJWTSecret() []byte {
	jwtSecretOnce.Do(func() {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "default_jwt_secret_for_development" // 开发环境默认值
		}
		jwtSecret = []byte(secret)
	})
	return jwtSecret
}

// Claims 是项目内 JWT 的自定义载荷。
// 除了标准注册声明外，还把 user_id / username / role 放进去，
// 这样后续业务 handler 就可以直接从 context 里拿到这些信息。
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT token。
// token 中携带 user_id、username 和 role，下游服务可以直接根据这些信息做权限
// 判断，不需要每个请求都回查 user-service。
func GenerateToken(userID uint, username, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

// ParseToken 解析 JWT token，并使用共享密钥校验签名。
// Swagger 或接口测试中看到的 “Invalid token”，通常就是这里解析或验签失败。
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// JWTAuth JWT 认证中间件。
//
// 普通 REST 请求从 Authorization 头取 token，WebSocket 握手可从 query 参数取
// token。校验通过后会把用户信息写入 Gin context，供业务 handler 和 RoleAuth 使用。
// 读取 token 的优先级：
// 1. Authorization: Bearer <token>
// 2. URL 查询参数 token（主要给 WebSocket 握手使用）
// 校验通过后，会把 user_id / username / role 写入 Gin context。
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 优先从 Header 获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}

		// 如果 Header 中没有，尝试从 URL 参数获取（用于 WebSocket 连接）
		if tokenString == "" {
			tokenString = c.Query("token")
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}

		claims, err := ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// RoleAuth 角色权限中间件。
// 它实现管理员、HR、招聘者等角色的粗粒度 RBAC 控制，例如创建职位、管理人才、
// 查看操作日志等后台能力。
// 它依赖 JWTAuth 先把 role 写进 context，因此通常应当和 JWTAuth 组合使用。
func RoleAuth(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role not found"})
			c.Abort()
			return
		}

		roleStr := role.(string)
		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}
