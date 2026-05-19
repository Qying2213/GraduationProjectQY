package models

import (
	"time"
)

// Conversation 表示两个用户之间的一条聊天会话。
// 会话只保存双方参与人和最后一条消息索引，具体消息内容在 chat_messages 表中。
type Conversation struct {
	ID            uint       `gorm:"primarykey" json:"id"`
	ParticipantA  uint       `gorm:"column:participant_a;not null" json:"participant_a"` // 用户A ID
	ParticipantB  uint       `gorm:"column:participant_b;not null" json:"participant_b"` // 用户B ID
	LastMessageID *uint      `gorm:"column:last_message_id" json:"last_message_id"`      // 最后一条消息ID
	LastMessageAt *time.Time `gorm:"column:last_message_at" json:"last_message_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// GORM 关联：用于列表页预加载最后一条消息，或查询会话下的消息集合。
	LastMessage *ChatMessage  `gorm:"foreignKey:LastMessageID" json:"last_message,omitempty"`
	Messages    []ChatMessage `gorm:"foreignKey:ConversationID" json:"messages,omitempty"`
}

// TableName 指定 Conversation 对应的数据库表名。
func (Conversation) TableName() string {
	return "conversations"
}

// ChatMessage 表示会话中的单条聊天消息。
// 消息落库后，即使接收方离线，也可以在下次打开会话时加载出来。
type ChatMessage struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	ConversationID uint      `gorm:"column:conversation_id;not null" json:"conversation_id"`
	SenderID       uint      `gorm:"column:sender_id;not null" json:"sender_id"`
	Content        string    `gorm:"type:text;not null" json:"content"`
	MessageType    string    `gorm:"column:message_type;size:20;default:'text'" json:"message_type"` // text, image, file
	IsRead         bool      `gorm:"column:is_read;default:false" json:"is_read"`
	CreatedAt      time.Time `json:"created_at"`

	// GORM 关联：回到所属会话。
	Conversation *Conversation `gorm:"foreignKey:ConversationID" json:"conversation,omitempty"`
}

// TableName 指定 ChatMessage 对应的数据库表名。
func (ChatMessage) TableName() string {
	return "chat_messages"
}

// ============================================================================
// API 响应辅助结构
// ============================================================================

// ConversationWithDetails 是会话列表接口返回的增强结构。
// 它在会话基础数据上补充对方用户、最后消息、未读数等前端展示字段。
type ConversationWithDetails struct {
	ID            uint         `json:"id"`
	Participant   *UserInfo    `json:"participant"`  // 会话另一方用户
	LastMessage   *ChatMessage `json:"last_message"` // 会话最后一条消息
	UnreadCount   int          `json:"unread_count"` // 当前用户未读消息数
	LastMessageAt *time.Time   `json:"last_message_at"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

// UserInfo 是会话页面展示参与人所需的基础用户信息。
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar,omitempty"`
	Role     string `json:"role,omitempty"`
	IsOnline bool   `json:"is_online"`
}

// ConversationListResponse 是会话列表接口响应结构。
type ConversationListResponse struct {
	Conversations []ConversationWithDetails `json:"conversations"`
	Total         int64                     `json:"total"`
}

// MessageListResponse 是会话消息分页接口响应结构。
type MessageListResponse struct {
	Messages []ChatMessage `json:"messages"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

// SendMessageRequest 是发送聊天消息的请求结构。
type SendMessageRequest struct {
	Content     string `json:"content" binding:"required"`
	MessageType string `json:"message_type"` // 未传时默认为 text
}

// CreateConversationRequest 是创建会话的请求结构。
type CreateConversationRequest struct {
	ParticipantID uint `json:"participant_id" binding:"required"`
}

// CreateConversationResponse 是创建/获取会话接口的响应结构。
type CreateConversationResponse struct {
	Conversation ConversationWithDetails `json:"conversation"`
	IsNew        bool                    `json:"is_new"` // true 表示新建，false 表示复用已有会话
}

// UnreadCountResponse 是总未读数接口响应结构。
type UnreadCountResponse struct {
	TotalUnread int `json:"total_unread"`
}
