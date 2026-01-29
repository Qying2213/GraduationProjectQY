package handlers

import (
	"message-service/models"
	"message-service/websocket"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ChatHandler handles chat-related HTTP requests
// Requirements: 9.1 (Conversation List), 9.2 (Last Message Preview), 9.3 (Unread Count)
type ChatHandler struct {
	DB  *gorm.DB
	Hub *websocket.Hub
}

// NewChatHandler creates a new ChatHandler instance
func NewChatHandler(db *gorm.DB, hub *websocket.Hub) *ChatHandler {
	return &ChatHandler{DB: db, Hub: hub}
}

// User represents the user model for joining with conversations
// This is a local struct to avoid importing from user-service
type User struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	Username string `json:"username"`
	RealName string `gorm:"column:real_name" json:"real_name"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}

// GetConversations retrieves the list of conversations for the current user
// GET /conversations
// Requirements: 9.1 (sorted by last message time), 9.2 (last message preview), 9.3 (unread count)
// Property 17: Conversation Sorting Order - conversations sorted by last_message_at DESC
func (h *ChatHandler) GetConversations(c *gin.Context) {
	// Get current user ID from JWT
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

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// Query conversations where current user is a participant
	// Sorted by last_message_at DESC (Requirement 9.1)
	var conversations []models.Conversation
	query := h.DB.Model(&models.Conversation{}).
		Where("participant_a = ? OR participant_b = ?", userID, userID).
		Order("last_message_at DESC NULLS LAST")

	// Get total count
	var total int64
	query.Count(&total)

	// Get paginated results with last message preloaded
	if err := query.Preload("LastMessage").
		Offset(offset).
		Limit(pageSize).
		Find(&conversations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "获取会话列表失败"})
		return
	}

	// Build response with participant info and unread count
	conversationDetails := make([]models.ConversationWithDetails, 0, len(conversations))
	for _, conv := range conversations {
		// Determine the other participant
		var otherUserID uint
		if conv.ParticipantA == userID {
			otherUserID = conv.ParticipantB
		} else {
			otherUserID = conv.ParticipantA
		}

		// Get participant info
		var participant User
		if err := h.DB.First(&participant, otherUserID).Error; err != nil {
			// If user not found, use placeholder
			participant = User{
				ID:       otherUserID,
				Username: "Unknown",
				RealName: "未知用户",
			}
		}

		// Count unread messages (messages sent by other user that are not read)
		// Property 18: Unread Count Accuracy - count messages where is_read=false and sender_id != current user
		var unreadCount int64
		h.DB.Model(&models.ChatMessage{}).
			Where("conversation_id = ? AND sender_id != ? AND is_read = ?", conv.ID, userID, false).
			Count(&unreadCount)

		// Build participant info
		participantInfo := &models.UserInfo{
			ID:       participant.ID,
			Username: participant.Username,
			Name:     participant.RealName,
			Avatar:   participant.Avatar,
			Role:     participant.Role,
			IsOnline: h.Hub != nil && h.Hub.IsUserOnline(participant.ID), // Check online status via Hub
		}

		// Use real_name if available, otherwise use username
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

// CreateOrGetConversation creates a new conversation or returns an existing one
// POST /conversations
// Requirements: 9.1, 9.2, 9.3
func (h *ChatHandler) CreateOrGetConversation(c *gin.Context) {
	// Get current user ID from JWT
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

	// Parse request body
	var req models.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "参数错误: participant_id 是必填项"})
		return
	}

	// Validate participant_id
	if req.ParticipantID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "参数错误: participant_id 不能为空"})
		return
	}

	// Cannot create conversation with self
	if req.ParticipantID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "不能与自己创建会话"})
		return
	}

	// Check if participant exists
	var participant User
	if err := h.DB.First(&participant, req.ParticipantID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "error": "用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "查询用户失败"})
		return
	}

	// Normalize participant order (smaller ID first) to ensure uniqueness
	participantA := userID
	participantB := req.ParticipantID
	if participantA > participantB {
		participantA, participantB = participantB, participantA
	}

	// Try to find existing conversation
	var conversation models.Conversation
	err := h.DB.Where("participant_a = ? AND participant_b = ?", participantA, participantB).
		Preload("LastMessage").
		First(&conversation).Error

	isNew := false
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new conversation
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

	// Count unread messages
	var unreadCount int64
	h.DB.Model(&models.ChatMessage{}).
		Where("conversation_id = ? AND sender_id != ? AND is_read = ?", conversation.ID, userID, false).
		Count(&unreadCount)

	// Build participant info
	participantInfo := &models.UserInfo{
		ID:       participant.ID,
		Username: participant.Username,
		Name:     participant.RealName,
		Avatar:   participant.Avatar,
		Role:     participant.Role,
		IsOnline: h.Hub != nil && h.Hub.IsUserOnline(participant.ID), // Check online status via Hub
	}

	// Use real_name if available, otherwise use username
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

// GetTotalUnreadCount returns the total unread message count for the current user
// GET /conversations/unread-count
// Requirements: 10.1, 10.2
func (h *ChatHandler) GetTotalUnreadCount(c *gin.Context) {
	// Get current user ID from JWT
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

	// Get all conversation IDs where user is a participant
	var conversationIDs []uint
	h.DB.Model(&models.Conversation{}).
		Where("participant_a = ? OR participant_b = ?", userID, userID).
		Pluck("id", &conversationIDs)

	// Count total unread messages across all conversations
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

// GetMessages retrieves messages for a conversation with pagination
// GET /conversations/:id/messages?page=1&page_size=20
// Requirements: 8.5 (Chat Message Persistence), 8.6 (Load older messages with pagination)
// Property 15: Chat Message Persistence Round-Trip - messages are persisted and retrievable
func (h *ChatHandler) GetMessages(c *gin.Context) {
	// Get current user ID from JWT
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

	// Parse conversation ID from URL
	conversationIDStr := c.Param("id")
	conversationID, err := strconv.ParseUint(conversationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "无效的会话ID"})
		return
	}

	// Verify conversation exists and user is a participant
	var conversation models.Conversation
	if err := h.DB.First(&conversation, conversationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "error": "会话不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "查询会话失败"})
		return
	}

	// Check if user is a participant in this conversation
	if conversation.ParticipantA != userID && conversation.ParticipantB != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "error": "权限不足"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// Query messages for this conversation
	// Order by created_at ASC (oldest first) for chat display
	var messages []models.ChatMessage
	query := h.DB.Model(&models.ChatMessage{}).
		Where("conversation_id = ?", conversationID).
		Order("created_at ASC")

	// Get total count
	var total int64
	query.Count(&total)

	// Get paginated results
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

// SendMessage sends a new message in a conversation
// POST /conversations/:id/messages
// Requirements: 8.5 (Chat Message Persistence)
// Property 15: Chat Message Persistence Round-Trip - messages are persisted and retrievable
func (h *ChatHandler) SendMessage(c *gin.Context) {
	// Get current user ID from JWT
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

	// Parse conversation ID from URL
	conversationIDStr := c.Param("id")
	conversationID, err := strconv.ParseUint(conversationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "无效的会话ID"})
		return
	}

	// Verify conversation exists and user is a participant
	var conversation models.Conversation
	if err := h.DB.First(&conversation, conversationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "error": "会话不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "查询会话失败"})
		return
	}

	// Check if user is a participant in this conversation
	if conversation.ParticipantA != userID && conversation.ParticipantB != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "error": "权限不足"})
		return
	}

	// Parse request body
	var req models.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "参数错误: content 是必填项"})
		return
	}

	// Validate content is not empty
	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "消息内容不能为空"})
		return
	}

	// Set default message type if not provided
	messageType := req.MessageType
	if messageType == "" {
		messageType = "text"
	}

	// Create the message
	message := models.ChatMessage{
		ConversationID: uint(conversationID),
		SenderID:       userID,
		Content:        req.Content,
		MessageType:    messageType,
		IsRead:         false,
	}

	// Use transaction to create message and update conversation
	tx := h.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "创建消息失败"})
		return
	}

	// Create the message
	if err := tx.Create(&message).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "创建消息失败"})
		return
	}

	// Update conversation's last_message_id and last_message_at
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

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "创建消息失败"})
		return
	}

	// Broadcast message via WebSocket for real-time delivery
	// Requirements: 8.2 - Deliver messages in real-time via WebSocket
	// Requirements: 8.4 - Store message and deliver when recipient connects (offline delivery)
	if h.Hub != nil {
		// Determine the recipient (the other participant)
		var recipientID uint
		if conversation.ParticipantA == userID {
			recipientID = conversation.ParticipantB
		} else {
			recipientID = conversation.ParticipantA
		}

		// Send to recipient via WebSocket if online
		h.Hub.SendChatMessage(
			uint(conversationID),
			message.ID,
			userID,
			message.Content,
			message.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		)

		// Also send to sender (for multi-device sync)
		h.Hub.SendToUser(userID, "chat", map[string]interface{}{
			"conversation_id": conversationID,
			"message":         message,
		})

		// Send to recipient
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

// MarkAsRead marks all unread messages in a conversation as read
// PUT /conversations/:id/read
// Requirements: 9.4 (Mark all messages as read when opening conversation)
// Property 19: Mark As Read Behavior - all messages where sender_id != current user are marked as read
func (h *ChatHandler) MarkAsRead(c *gin.Context) {
	// Get current user ID from JWT
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

	// Parse conversation ID from URL
	conversationIDStr := c.Param("id")
	conversationID, err := strconv.ParseUint(conversationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "无效的会话ID"})
		return
	}

	// Verify conversation exists and user is a participant
	var conversation models.Conversation
	if err := h.DB.First(&conversation, conversationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "error": "会话不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "查询会话失败"})
		return
	}

	// Check if user is a participant in this conversation
	if conversation.ParticipantA != userID && conversation.ParticipantB != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "error": "权限不足"})
		return
	}

	// Count unread messages before marking as read (for response)
	var unreadCount int64
	h.DB.Model(&models.ChatMessage{}).
		Where("conversation_id = ? AND sender_id != ? AND is_read = ?", conversationID, userID, false).
		Count(&unreadCount)

	// Mark all unread messages from other user as read
	// Property 19: Only mark messages where sender_id != current user
	result := h.DB.Model(&models.ChatMessage{}).
		Where("conversation_id = ? AND sender_id != ? AND is_read = ?", conversationID, userID, false).
		Update("is_read", true)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "标记已读失败"})
		return
	}

	// Broadcast read notification via WebSocket
	// Requirements: 9.4 - Mark messages as read
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
