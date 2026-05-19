package handlers

import (
	"message-service/models"
	"message-service/websocket"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ChatHandler 处理聊天会话相关 HTTP 接口。
// 它负责会话列表、未读数、消息分页、发送消息和 WebSocket 实时推送。
type ChatHandler struct {
	DB  *gorm.DB
	Hub *websocket.Hub
}

// NewChatHandler 创建聊天处理器，并注入 WebSocket Hub 用于在线消息推送。
func NewChatHandler(db *gorm.DB, hub *websocket.Hub) *ChatHandler {
	return &ChatHandler{DB: db, Hub: hub}
}

// User 是聊天查询中用于关联 users 表的本地结构。
// message-service 不直接依赖 user-service 的代码包，只通过共享数据库读取必要字段。
type User struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	Username string `json:"username"`
	RealName string `gorm:"column:real_name" json:"real_name"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
}

// TableName 指定本地 User 结构对应 users 表。
func (User) TableName() string {
	return "users"
}

// GetConversations 获取当前用户的会话列表。
// 会话按最后一条消息时间倒序排列，并附带对方用户信息、最后消息和未读数。
func (h *ChatHandler) GetConversations(c *gin.Context) {
	// 从 JWT 上下文获取当前用户 ID。
	jwtUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "error": "未授权，请先登录"})
		return
	}

	userID, ok := jwtUserID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "invalid user_id type"})
		return
	}

	// 解析分页参数。
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// 查询当前用户参与的会话，并按最后消息时间倒序排序。
	var conversations []models.Conversation
	query := h.DB.Model(&models.Conversation{}).
		Where("participant_a = ? OR participant_b = ?", userID, userID).
		Order("last_message_at DESC NULLS LAST")

	// 获取总数。
	var total int64
	query.Count(&total)

	// 分页查询，并预加载最后一条消息。
	if err := query.Preload("LastMessage").
		Offset(offset).
		Limit(pageSize).
		Find(&conversations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "获取会话列表失败"})
		return
	}

	// 组装对方用户信息、在线状态和未读数。
	conversationDetails := make([]models.ConversationWithDetails, 0, len(conversations))
	for _, conv := range conversations {
		// 判断会话中的另一方用户。
		var otherUserID uint
		if conv.ParticipantA == userID {
			otherUserID = conv.ParticipantB
		} else {
			otherUserID = conv.ParticipantA
		}

		// 查询对方用户基础信息。
		var participant User
		if err := h.DB.First(&participant, otherUserID).Error; err != nil {
			// 用户不存在时使用占位信息，避免整个会话列表失败。
			participant = User{
				ID:       otherUserID,
				Username: "Unknown",
				RealName: "未知用户",
			}
		}

		// 只统计对方发送且当前用户未读的消息，保证未读数准确。
		var unreadCount int64
		h.DB.Model(&models.ChatMessage{}).
			Where("conversation_id = ? AND sender_id != ? AND is_read = ?", conv.ID, userID, false).
			Count(&unreadCount)

		// 构建前端需要的会话参与人信息。
		participantInfo := &models.UserInfo{
			ID:       participant.ID,
			Username: participant.Username,
			Name:     participant.RealName,
			Avatar:   participant.Avatar,
			Role:     participant.Role,
			IsOnline: h.Hub != nil && h.Hub.IsUserOnline(participant.ID), // 通过 Hub 获取在线状态
		}

		// 优先展示真实姓名，没有时回退到用户名。
		if participantInfo.Name == "" {
			participantInfo.Name = participant.Username
		}

		detail := models.ConversationWithDetails{
			ID:            conv.ID,
			Participant:   participantInfo,
			LastMessage:   conv.LastMessage,
			UnreadCount:   int(unreadCount),
			LastMessageAt: conv.LastMessageAt,
			CreatedAt:     conv.CreatedAt,
			UpdatedAt:     conv.UpdatedAt,
		}

		conversationDetails = append(conversationDetails, detail)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": models.ConversationListResponse{
			Conversations: conversationDetails,
			Total:         total,
		},
	})
}

// CreateOrGetConversation 创建或获取一对一会话。
// 如果两个用户之间已经存在会话，则直接返回旧会话，避免重复创建。
func (h *ChatHandler) CreateOrGetConversation(c *gin.Context) {
	// 从 JWT 上下文获取当前用户 ID。
	jwtUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "error": "未授权，请先登录"})
		return
	}

	userID, ok := jwtUserID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "invalid user_id type"})
		return
	}

	// 解析请求体。
	var req models.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "参数错误: participant_id 是必填项"})
		return
	}

	// 校验会话参与人。
	if req.ParticipantID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "参数错误: participant_id 不能为空"})
		return
	}

	// 不允许和自己创建会话。
	if req.ParticipantID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "不能与自己创建会话"})
		return
	}

	// 检查对方用户是否存在。
	var participant User
	if err := h.DB.First(&participant, req.ParticipantID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "error": "用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "查询用户失败"})
		return
	}

	// 固定参与者顺序，保证同一对用户只会生成一个会话。
	participantA := userID
	participantB := req.ParticipantID
	if participantA > participantB {
		participantA, participantB = participantB, participantA
	}

	// 优先查找已存在的会话。
	var conversation models.Conversation
	err := h.DB.Where("participant_a = ? AND participant_b = ?", participantA, participantB).
		Preload("LastMessage").
		First(&conversation).Error

	isNew := false
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 不存在时创建新会话。
			conversation = models.Conversation{
				ParticipantA: participantA,
				ParticipantB: participantB,
			}
			if err := h.DB.Create(&conversation).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "创建会话失败"})
				return
			}
			isNew = true
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "查询会话失败"})
			return
		}
	}

	// 统计当前会话未读消息。
	var unreadCount int64
	h.DB.Model(&models.ChatMessage{}).
		Where("conversation_id = ? AND sender_id != ? AND is_read = ?", conversation.ID, userID, false).
		Count(&unreadCount)

	// 构建对方参与人信息。
	participantInfo := &models.UserInfo{
		ID:       participant.ID,
		Username: participant.Username,
		Name:     participant.RealName,
		Avatar:   participant.Avatar,
		Role:     participant.Role,
		IsOnline: h.Hub != nil && h.Hub.IsUserOnline(participant.ID), // Check online status via Hub
	}

	// 优先展示真实姓名，没有时回退到用户名。
	if participantInfo.Name == "" {
		participantInfo.Name = participant.Username
	}

	detail := models.ConversationWithDetails{
		ID:            conversation.ID,
		Participant:   participantInfo,
		LastMessage:   conversation.LastMessage,
		UnreadCount:   int(unreadCount),
		LastMessageAt: conversation.LastMessageAt,
		CreatedAt:     conversation.CreatedAt,
		UpdatedAt:     conversation.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": models.CreateConversationResponse{
			Conversation: detail,
			IsNew:        isNew,
		},
	})
}

// GetTotalUnreadCount 获取当前用户所有会话的总未读数。
func (h *ChatHandler) GetTotalUnreadCount(c *gin.Context) {
	// 从 JWT 上下文获取当前用户 ID。
	jwtUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "error": "未授权，请先登录"})
		return
	}

	userID, ok := jwtUserID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "invalid user_id type"})
		return
	}

	// 找到当前用户参与的所有会话 ID。
	var conversationIDs []uint
	h.DB.Model(&models.Conversation{}).
		Where("participant_a = ? OR participant_b = ?", userID, userID).
		Pluck("id", &conversationIDs)

	// 统计所有会话中由对方发送且尚未阅读的消息。
	var totalUnread int64
	if len(conversationIDs) > 0 {
		h.DB.Model(&models.ChatMessage{}).
			Where("conversation_id IN ? AND sender_id != ? AND is_read = ?", conversationIDs, userID, false).
			Count(&totalUnread)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": models.UnreadCountResponse{
			TotalUnread: int(totalUnread),
		},
	})
}

// GetMessages 分页获取指定会话的聊天记录。
// 只有会话参与者可以读取消息，返回结果按时间正序用于聊天窗口展示。
func (h *ChatHandler) GetMessages(c *gin.Context) {
	// 从 JWT 上下文获取当前用户 ID。
	jwtUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "error": "未授权，请先登录"})
		return
	}

	userID, ok := jwtUserID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "invalid user_id type"})
		return
	}

	// 从 URL 解析会话 ID。
	conversationIDStr := c.Param("id")
	conversationID, err := strconv.ParseUint(conversationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "无效的会话ID"})
		return
	}

	// 校验会话存在。
	var conversation models.Conversation
	if err := h.DB.First(&conversation, conversationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "error": "会话不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "查询会话失败"})
		return
	}

	// 校验当前用户确实是会话参与者。
	if conversation.ParticipantA != userID && conversation.ParticipantB != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "error": "权限不足"})
		return
	}

	// 解析分页参数。
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// 查询会话消息，按创建时间正序返回，便于聊天窗口直接展示。
	var messages []models.ChatMessage
	query := h.DB.Model(&models.ChatMessage{}).
		Where("conversation_id = ?", conversationID).
		Order("created_at ASC")

	// 获取总数。
	var total int64
	query.Count(&total)

	// 获取分页结果。
	if err := query.Offset(offset).Limit(pageSize).Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "获取消息列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": models.MessageListResponse{
			Messages: messages,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// SendMessage 在指定会话中发送新消息。
// 消息先持久化到数据库，再通过 WebSocket 推送给双方在线连接。
func (h *ChatHandler) SendMessage(c *gin.Context) {
	// 从 JWT 上下文获取当前用户 ID。
	jwtUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "error": "未授权，请先登录"})
		return
	}

	userID, ok := jwtUserID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "invalid user_id type"})
		return
	}

	// 从 URL 解析会话 ID。
	conversationIDStr := c.Param("id")
	conversationID, err := strconv.ParseUint(conversationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "无效的会话ID"})
		return
	}

	// 校验会话存在。
	var conversation models.Conversation
	if err := h.DB.First(&conversation, conversationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "error": "会话不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "查询会话失败"})
		return
	}

	// 校验当前用户确实是会话参与者。
	if conversation.ParticipantA != userID && conversation.ParticipantB != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "error": "权限不足"})
		return
	}

	// 解析请求体。
	var req models.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "参数错误: content 是必填项"})
		return
	}

	// 校验消息内容不能为空。
	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "消息内容不能为空"})
		return
	}

	// 未传消息类型时默认按文本消息处理。
	messageType := req.MessageType
	if messageType == "" {
		messageType = "text"
	}

	// 构建待保存的消息记录。
	message := models.ChatMessage{
		ConversationID: uint(conversationID),
		SenderID:       userID,
		Content:        req.Content,
		MessageType:    messageType,
		IsRead:         false,
	}

	// 使用事务同时保存消息并更新会话最后消息，保证列表预览一致。
	tx := h.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "创建消息失败"})
		return
	}

	// 保存消息。
	if err := tx.Create(&message).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "创建消息失败"})
		return
	}

	// 更新会话最后消息和最后消息时间。
	now := message.CreatedAt
	if err := tx.Model(&conversation).Updates(map[string]interface{}{
		"last_message_id": message.ID,
		"last_message_at": now,
		"updated_at":      now,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "更新会话失败"})
		return
	}

	// 提交事务。
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "创建消息失败"})
		return
	}

	// 通过 WebSocket 实时推送消息；离线用户后续仍可从数据库读取。
	if h.Hub != nil {
		// 计算接收方，即会话中的另一位参与者。
		var recipientID uint
		if conversation.ParticipantA == userID {
			recipientID = conversation.ParticipantB
		} else {
			recipientID = conversation.ParticipantA
		}

		// 推送到会话订阅者。
		h.Hub.SendChatMessage(
			uint(conversationID),
			message.ID,
			userID,
			message.Content,
			message.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		)

		// 同步给发送方的其他在线设备。
		h.Hub.SendToUser(userID, "chat", map[string]interface{}{
			"conversation_id": conversationID,
			"message":         message,
		})

		// 推送给接收方在线连接。
		h.Hub.SendToUser(recipientID, "chat", map[string]interface{}{
			"conversation_id": conversationID,
			"message":         message,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    message,
	})
}

// MarkAsRead 将指定会话中对方发来的未读消息标记为已读。
func (h *ChatHandler) MarkAsRead(c *gin.Context) {
	// 从 JWT 上下文获取当前用户 ID。
	jwtUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "error": "未授权，请先登录"})
		return
	}

	userID, ok := jwtUserID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "invalid user_id type"})
		return
	}

	// 从 URL 解析会话 ID。
	conversationIDStr := c.Param("id")
	conversationID, err := strconv.ParseUint(conversationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "无效的会话ID"})
		return
	}

	// 校验会话存在。
	var conversation models.Conversation
	if err := h.DB.First(&conversation, conversationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "error": "会话不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "查询会话失败"})
		return
	}

	// 校验当前用户确实是会话参与者。
	if conversation.ParticipantA != userID && conversation.ParticipantB != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "error": "权限不足"})
		return
	}

	// 先统计未读数量，用于响应和已读通知。
	var unreadCount int64
	h.DB.Model(&models.ChatMessage{}).
		Where("conversation_id = ? AND sender_id != ? AND is_read = ?", conversationID, userID, false).
		Count(&unreadCount)

	// 只标记对方发送的未读消息，避免把自己发送的消息误算为已读处理对象。
	result := h.DB.Model(&models.ChatMessage{}).
		Where("conversation_id = ? AND sender_id != ? AND is_read = ?", conversationID, userID, false).
		Update("is_read", true)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "标记已读失败"})
		return
	}

	// 通过 WebSocket 广播已读通知。
	if h.Hub != nil && unreadCount > 0 {
		h.Hub.SendChatReadNotification(uint(conversationID), userID, int(unreadCount))
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"marked_count": unreadCount,
		},
	})
}
