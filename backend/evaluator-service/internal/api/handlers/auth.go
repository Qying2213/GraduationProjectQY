package handlers

import (
	"net/http"

	"evaluator-service/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

// LoginRequest 登录请求
type LoginRequest struct {
	CorpCode string `json:"corp_code" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 用户登录
func (h *Handlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "corp_code, username, password 不能为空"})
		return
	}

	result, err := h.authSvc.Login(c.Request.Context(), req.CorpCode, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetProfile 获取当前用户信息
func (h *Handlers) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	user, err := h.authSvc.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            user.ID,
		"corp_code":     user.CorpCode,
		"username":      user.Username,
		"last_login_at": user.LastLoginAt,
	})
}

// Logout 用户登出
func (h *Handlers) Logout(c *gin.Context) {
	// JWT 无状态，客户端删除 token 即可
	c.JSON(http.StatusOK, gin.H{"message": "已登出"})
}
