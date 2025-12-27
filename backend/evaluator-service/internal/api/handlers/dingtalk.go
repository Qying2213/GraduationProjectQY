package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"evaluator-service/internal/api/middleware"
	"evaluator-service/internal/dingtalk"
	"evaluator-service/internal/logging"
	"evaluator-service/internal/models"

	"github.com/gin-gonic/gin"
)

// GetDingTalkConfig 获取钉钉配置（兼容旧接口）
func (h *Handlers) GetDingTalkConfig(c *gin.Context) {
	userID := middleware.GetUserID(c)
	config, err := h.dtRepo.GetByUser(userID)
	if err != nil {
		fail(c, err)
		return
	}
	if config == nil {
		config = &models.DingTalkConfig{
			Name:      "默认机器人",
			PushTime:  "09:00",
			PushLimit: 10,
			Enabled:   false,
		}
	}
	ok(c, config)
}

// ListDingTalkConfigs 获取所有钉钉配置
func (h *Handlers) ListDingTalkConfigs(c *gin.Context) {
	userID := middleware.GetUserID(c)
	configs, err := h.dtRepo.ListByUser(userID)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, gin.H{"configs": configs, "total": len(configs)})
}

// GetDingTalkConfigByID 根据ID获取配置
func (h *Handlers) GetDingTalkConfigByID(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")
	var idUint uint
	if _, err := fmt.Sscanf(id, "%d", &idUint); err != nil {
		c.JSON(400, gin.H{"error": "无效的ID"})
		return
	}

	config, err := h.dtRepo.GetByIDAndUser(idUint, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	if config == nil {
		c.JSON(404, gin.H{"error": "配置不存在"})
		return
	}
	ok(c, config)
}

// DeleteDingTalkConfig 删除配置
func (h *Handlers) DeleteDingTalkConfig(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")
	var idUint uint
	if _, err := fmt.Sscanf(id, "%d", &idUint); err != nil {
		c.JSON(400, gin.H{"error": "无效的ID"})
		return
	}

	if err := h.dtRepo.DeleteByUser(idUint, userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	ok(c, gin.H{"success": true, "message": "删除成功"})
}

// UpsertDingTalkConfig 更新钉钉配置
func (h *Handlers) UpsertDingTalkConfig(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var config models.DingTalkConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		bad(c, err)
		return
	}

	if err := h.dtRepo.UpsertByUser(&config, userID); err != nil {
		if err.Error() == "无权访问该资源" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		fail(c, err)
		return
	}

	// 重启钉钉服务以应用新配置
	if h.dtService != nil {
		if err := h.dtService.Restart(c.Request.Context()); err != nil {
			h.log.Error("restart dingtalk service failed", logging.Err(err))
		}
	}

	ok(c, gin.H{"success": true, "config": config})
}

// TestDingTalkPush 测试钉钉推送
func (h *Handlers) TestDingTalkPush(c *gin.Context) {
	userID := middleware.GetUserID(c)
	config, err := h.dtRepo.GetByUser(userID)
	if err != nil {
		fail(c, err)
		return
	}

	if config == nil || config.Webhook == "" {
		c.JSON(400, gin.H{"error": "钉钉配置不存在或Webhook未配置"})
		return
	}

	h.log.Info("test dingtalk push triggered",
		logging.KV("webhook", config.Webhook),
		logging.KV("has_secret", config.Secret != ""))

	client := NewDingTalkClient(config, h.log)

	atUserIDs := []string{}
	if config.AtUserIDs != "" {
		for _, id := range strings.Split(config.AtUserIDs, ",") {
			atUserIDs = append(atUserIDs, strings.TrimSpace(id))
		}
	}

	testContent := fmt.Sprintf("## 🧪 测试消息\n\n这是一条来自简历评估系统的测试消息\n\n- 发送时间: %s\n- 配置状态: %s\n- 推送时间: %s\n- 推送数量: %d人\n\n✅ 如果您看到这条消息，说明钉钉机器人配置成功！",
		time.Now().Format("2006-01-02 15:04:05"),
		map[bool]string{true: "已启用", false: "未启用"}[config.Enabled],
		config.PushTime,
		config.PushLimit)

	err = client.SendMarkdownMessage("测试消息", testContent, atUserIDs, false)
	if err != nil {
		h.log.Error("test push failed", logging.Err(err))
		c.JSON(500, gin.H{"error": "发送失败: " + err.Error()})
		return
	}

	h.log.Info("test push sent successfully")
	ok(c, gin.H{"success": true, "message": "测试消息已发送，请查看钉钉群"})
}

// PushNow 立即推送候选人
func (h *Handlers) PushNow(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req struct {
		ConfigID uint `json:"config_id"`
	}
	_ = c.ShouldBindJSON(&req)

	var config *models.DingTalkConfig
	var err error

	if req.ConfigID > 0 {
		config, err = h.dtRepo.GetByIDAndUser(req.ConfigID, userID)
	} else {
		config, err = h.dtRepo.GetByUser(userID)
	}

	if err != nil {
		if err.Error() == "无权访问该资源" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		fail(c, err)
		return
	}

	if config == nil || config.Webhook == "" {
		c.JSON(400, gin.H{"error": "钉钉配置不存在或Webhook未配置"})
		return
	}

	h.log.Info("manual push triggered",
		logging.KV("config_id", config.ID),
		logging.KV("config_name", config.Name))

	// 查询用户未通知的候选人
	candidates, err := h.repo.FindUnnotifiedByUser(config.PushLimit, userID)
	if err != nil {
		fail(c, err)
		return
	}

	if len(candidates) == 0 {
		ok(c, gin.H{"success": true, "message": "没有待推送的候选人", "count": 0})
		return
	}

	if err := h.dtService.PushCandidates(candidates, config); err != nil {
		h.log.Error("push failed", logging.Err(err))
		c.JSON(500, gin.H{"error": "推送失败: " + err.Error()})
		return
	}

	h.log.Info("manual push completed",
		logging.KV("count", len(candidates)))

	ok(c, gin.H{
		"success": true,
		"message": fmt.Sprintf("已推送 %d 位候选人", len(candidates)),
		"count":   len(candidates),
	})
}

// buildCandidateListMarkdown 构建候选人列表Markdown
func buildCandidateListMarkdown(candidates []models.Candidate, atUserIds []string) string {
	var sb strings.Builder

	sb.WriteString("## 📋 待面试候选人列表\n\n")
	sb.WriteString(fmt.Sprintf("> 共 **%d** 位候选人待处理\n\n", len(candidates)))

	for i, c := range candidates {
		gradeEmoji := getGradeEmoji(c.Grade)
		sb.WriteString(fmt.Sprintf("**%d. %s** %s\n", i+1, c.Name, gradeEmoji))
		sb.WriteString(fmt.Sprintf("- 评分: **%.1f** 分 | 评级: **%s**\n", c.TotalScore, c.Grade))
		sb.WriteString(fmt.Sprintf("- JD匹配: %d%% | 状态: %s\n", c.JDMatch, c.Status))
		sb.WriteString(fmt.Sprintf("- 建议: %s\n", c.Recommendation))
		sb.WriteString("\n")
	}

	sb.WriteString("---\n")
	sb.WriteString("💡 **回复候选人序号（如：1）查看详细信息**\n\n")

	if len(atUserIds) > 0 {
		sb.WriteString("\n")
		for _, userId := range atUserIds {
			if userId != "" {
				sb.WriteString(fmt.Sprintf("@%s ", userId))
			}
		}
	}

	return sb.String()
}

func getGradeEmoji(grade string) string {
	switch grade {
	case "A":
		return "🌟"
	case "B":
		return "✨"
	case "C":
		return "⭐"
	case "D":
		return "💫"
	default:
		return "📄"
	}
}

// NewDingTalkClient 创建钉钉客户端的辅助函数
func NewDingTalkClient(config *models.DingTalkConfig, log *logging.Logger) *dingtalk.Client {
	return dingtalk.NewClient(config, log)
}
