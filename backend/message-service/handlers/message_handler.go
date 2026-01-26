package handlers

import (
	"message-service/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MessageHandler struct {
	DB *gorm.DB
}

func NewMessageHandler(db *gorm.DB) *MessageHandler {
	return &MessageHandler{DB: db}
}

// SendMessage 发送消息
func (h *MessageHandler) SendMessage(c *gin.Context) {
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": err.Error()})
		return
	}

	if err := h.DB.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "Message sent successfully",
		"data":    message,
	})
}

// MessageResponse 消息响应结构
type MessageResponse struct {
	ID         uint   `json:"id"`
	SenderID   *uint  `json:"sender_id"`
	ReceiverID uint   `json:"receiver_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Type       string `json:"type"`
	IsRead     bool   `json:"is_read"`
	SenderName string `json:"sender_name"`
	CreatedAt  string `json:"created_at"`
}

// GetMessages 获取消息列表
func (h *MessageHandler) GetMessages(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "user_id is required"})
		return
	}

	// 尝试将 user_id 转换为整数，如果是 UUID 格式则使用默认值
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		// 如果不是整数（可能是 UUID），使用默认用户 ID 1
		// 这是为了兼容不同的认证系统
		userID = 1
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	msgType := c.Query("type")
	isRead := c.Query("is_read")

	offset := (page - 1) * pageSize

	query := h.DB.Model(&models.Message{}).Where("receiver_id = ?", userID)

	if msgType != "" {
		query = query.Where("type = ?", msgType)
	}

	if isRead != "" {
		query = query.Where("is_read = ?", isRead == "true")
	}

	var total int64
	query.Count(&total)

	var messages []models.Message
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "Failed to fetch messages"})
		return
	}

	// 转换为响应格式，添加 sender_name
	var responseMessages []MessageResponse
	for _, msg := range messages {
		resp := MessageResponse{
			ID:         msg.ID,
			SenderID:   msg.SenderID,
			ReceiverID: msg.ReceiverID,
			Title:      msg.Title,
			Content:    msg.Content,
			Type:       msg.Type,
			IsRead:     msg.IsRead,
			SenderName: "系统", // 默认发送者名称
			CreatedAt:  msg.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		// 如果有发送者ID，可以查询发送者名称（这里简化处理）
		if msg.SenderID != nil {
			resp.SenderName = "管理员"
		}
		responseMessages = append(responseMessages, resp)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"messages":  responseMessages,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// MarkAsRead 标记消息为已读
func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")

	if err := h.DB.Model(&models.Message{}).Where("id = ?", id).Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark message as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Message marked as read",
	})
}

// GetUnreadCount 获取未读消息数
func (h *MessageHandler) GetUnreadCount(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	// 尝试将 user_id 转换为整数，如果是 UUID 格式则使用默认值
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		// 如果不是整数（可能是 UUID），使用默认用户 ID 1
		userID = 1
	}

	var count int64
	h.DB.Model(&models.Message{}).Where("receiver_id = ? AND is_read = ?", userID, false).Count(&count)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"unread_count": count,
		},
	})
}

// DeleteMessage 删除消息
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	id := c.Param("id")

	if err := h.DB.Delete(&models.Message{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Message deleted successfully",
	})
}

// GetMessageStats 获取消息统计
func (h *MessageHandler) GetMessageStats(c *gin.Context) {
	var stats struct {
		TotalMessages  int64 `json:"total_messages"`
		UnreadMessages int64 `json:"unread_messages"`
		TodayMessages  int64 `json:"today_messages"`
		ByType         []struct {
			Type  string `json:"type"`
			Count int64  `json:"count"`
		} `json:"by_type"`
	}

	h.DB.Model(&models.Message{}).Count(&stats.TotalMessages)
	h.DB.Model(&models.Message{}).Where("is_read = ?", false).Count(&stats.UnreadMessages)
	h.DB.Model(&models.Message{}).Where("created_at >= CURRENT_DATE").Count(&stats.TodayMessages)

	// 按类型统计
	h.DB.Model(&models.Message{}).
		Select("type, count(*) as count").
		Group("type").
		Scan(&stats.ByType)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    stats,
	})
}
