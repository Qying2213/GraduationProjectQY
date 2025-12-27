package service

import (
	"errors"
	"time"

	"evaluator-service/internal/config"
	"evaluator-service/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

// UserClaims JWT 载荷中的用户信息
type UserClaims struct {
	UserID   uint   `json:"user_id"`
	CorpCode string `json:"corp_code"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTService JWT 服务
type JWTService struct {
	secretKey     []byte
	expireHours   int
}

// NewJWTService 创建 JWT 服务实例
func NewJWTService(cfg *config.Config) *JWTService {
	// 使用凭据加密密钥作为 JWT 签名密钥，如果未配置则使用默认值
	secretKey := cfg.Credentials.EncKey
	if secretKey == "" {
		secretKey = "resume-evaluator-default-jwt-secret-key-32b"
	}
	return &JWTService{
		secretKey:   []byte(secretKey),
		expireHours: 24, // 默认 24 小时过期
	}
}

// GenerateToken 为用户生成 JWT token
func (s *JWTService) GenerateToken(user *models.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(s.expireHours) * time.Hour)
	
	claims := UserClaims{
		UserID:   user.ID,
		CorpCode: user.CorpCode,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "resume-evaluator",
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", time.Time{}, err
	}
	
	return tokenString, expiresAt, nil
}

// ValidateToken 验证 JWT token 并返回用户信息
func (s *JWTService) ValidateToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, errors.New("invalid token")
}
