package handlers

import (
	"message-service/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NoticeHandler 负责公告管理接口。
// 公告与个人消息不同，它面向后台统一发布，可按状态、优先级和置顶规则展示。
type NoticeHandler struct {
	DB *gorm.DB
}

func NewNoticeHandler(db *gorm.DB) *NoticeHandler {
	return &NoticeHandler{DB: db}
}

// ListNotices 查询公告列表。
// 支持关键词、发布状态、优先级和是否置顶筛选，排序时置顶和高优先级公告优先展示。
func (h *NoticeHandler) ListNotices(c *gin.Context) {
	keyword := strings.TrimSpace(c.Query("keyword"))
	status := strings.TrimSpace(c.Query("status"))
	priority := strings.TrimSpace(c.Query("priority"))
	page, _ := strconv.Atoi(c.Copy().DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	isPinned := strings.TrimSpace(c.Query("is_pinned"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	query := h.DB.Model(&models.Notice{})

	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("title ILIKE ? OR content ILIKE ?", like, like)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if priority != "" {
		if !isValidPriority(priority) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "priority must be normal, high or urgent",
			})
			return
		}
		query = query.Where("priority = ?", priority)
	}

	if isPinned != "" {
		pinned, err := strconv.ParseBool(isPinned)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "invalid is_pinned value",
			})
			return
		}
		query = query.Where("is_pinned = ?", pinned)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to count notices",
		})
		return
	}
	var notices []models.Notice
	if err := query.Order("is_pinned DESC").Order("CASE priority WHEN 'urgent' THEN 3 WHEN 'high' THEN 2 ELSE 1 END DESC").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&notices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to fetch notices",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"notices":   notices,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})

}

// CreateNotice 创建公告。
// 默认创建草稿，只有 draft/published 两种状态，优先级会统一归一化。
func (h *NoticeHandler) CreateNotice(c *gin.Context) {
	var req struct {
		Title    string `json:"title" binding:"required"`
		Content  string `json:"content"`
		Status   string `json:"status"`
		IsPinned bool   `json:"is_pinned"`
		Priority string `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	if req.Status == "" {
		req.Status = "draft"
	}
	if req.Status != "draft" && req.Status != "published" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "status must be draft or published",
		})
		return
	}
	priority, ok := normalizePriority(req.Priority)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "priority must be normal, high or urgent",
		})
		return
	}

	notice := models.Notice{
		Title:    strings.TrimSpace(req.Title),
		Content:  strings.TrimSpace(req.Content),
		Status:   req.Status,
		IsPinned: req.IsPinned,
		Priority: priority,
	}
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(uint); ok {
			notice.CreatedBy = uid
		}
	}
	if err := h.DB.Create(&notice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create notice",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Notice created successfully",
		"data":    notice,
	})

}

// GetNotice 获取单条公告详情。
func (h *NoticeHandler) GetNotice(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid notice id",
		})
		return
	}

	var notice models.Notice
	if err := h.DB.First(&notice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Notice not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    notice,
	})
}
func (h *NoticeHandler) UpdateNotice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid notice id",
		})
		return
	}

	var req struct {
		Title    string `json:"title" binding:"required"`
		Content  string `json:"content"`
		Status   string `json:"status"`
		IsPinned bool   `json:"is_pinned"`
		Priority string `json:"priority"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	if req.Status != "draft" && req.Status != "published" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "status must be draft or published",
		})
		return
	}
	var notice models.Notice
	if err := h.DB.First(&notice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Notice not found",
		})
		return
	}

	notice.Title = strings.TrimSpace(req.Title)
	notice.Content = strings.TrimSpace(req.Content)
	notice.Status = req.Status
	notice.IsPinned = req.IsPinned

	priorityInput := strings.TrimSpace(req.Priority)
	if priorityInput == "" {
		priorityInput = notice.Priority
	}
	priority, ok := normalizePriority(priorityInput)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "priority must be normal, high or urgent",
		})
		return
	}
	notice.Priority = priority
	if err := h.DB.Save(&notice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to update notice",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Notice updated successfully",
		"data":    notice,
	})

}
func (h *NoticeHandler) DeleteNotice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid notice id",
		})
		return
	}

	var notice models.Notice
	if err := h.DB.First(&notice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Notice not found",
		})
		return
	}

	if err := h.DB.Delete(&notice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to delete notice",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Notice deleted successfully",
	})
}

const (
	noticePriorityNormal = "normal"
	noticePriorityHigh   = "high"
	noticePriorityUrgent = "urgent"
)

func isValidPriority(priority string) bool {
	switch priority {
	case noticePriorityNormal, noticePriorityHigh, noticePriorityUrgent:
		return true
	default:
		return false
	}
}
func normalizePriority(priority string) (string, bool) {
	v := strings.TrimSpace(priority)
	if v == "" {
		return noticePriorityNormal, true
	}
	if isValidPriority(v) {
		return v, true
	}
	return "", false
}
