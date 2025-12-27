package handlers

import (
	"context"
	"net/http"
	"time"

	"evaluator-service/internal/api/middleware"
	"evaluator-service/internal/database"
	"evaluator-service/internal/logging"
	"evaluator-service/internal/repository"
	"evaluator-service/internal/service"
	"evaluator-service/internal/utils"

	"github.com/gin-gonic/gin"
)

// getPositionService 获取 PositionService 实例
func (h *Handlers) getPositionService() *service.PositionService {
	repo := repository.NewPositionRepository(database.DB)
	return service.NewPositionService(repo, h.log)
}

// SyncPositions 同步岗位数据
func (h *Handlers) SyncPositions(c *gin.Context) {
	userID := middleware.GetUserID(c)

	// 获取用户凭据
	cred, err := h.credRepo.GetByOrgAndUser("motern", userID)
	if err != nil {
		fail(c, err)
		return
	}
	if cred == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "凭据未找到，请重新登录", "code": "CREDENTIALS_NOT_FOUND"})
		return
	}

	// 解密密码
	encKey := getEncryptionKey(h.cfg.Credentials.EncKey)
	password, decErr := utils.DecryptAESGCM(encKey, cred.PasswordCipher)
	if decErr != nil {
		h.log.Error("decrypt credential failed", logging.KV("org", cred.Org), logging.Err(decErr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "凭据解密失败，请重新登录"})
		return
	}

	// 注入给 Python 的环境变量
	env := map[string]string{
		"WT_USERNAME": cred.Account,
		"WT_PASSWORD": password,
	}

	// 调用 PositionService 同步岗位
	positionSvc := h.getPositionService()
	ctx, cancel := context.WithTimeout(c.Request.Context(), 300*time.Second)
	defer cancel()

	syncedCount, err := positionSvc.SyncPositions(ctx, userID, env)
	if err != nil {
		h.log.Error("Position sync failed", logging.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "岗位同步失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"synced_count": syncedCount,
		"message":      "岗位同步完成",
	})
}

// GetPositions 获取岗位列表
func (h *Handlers) GetPositions(c *gin.Context) {
	userID := middleware.GetUserID(c)

	positionSvc := h.getPositionService()
	positions, err := positionSvc.GetAllByUser(userID)
	if err != nil {
		h.log.Error("Failed to get positions", logging.Err(err))
		fail(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"positions": positions,
	})
}
