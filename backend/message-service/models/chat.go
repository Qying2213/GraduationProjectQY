package models

import (
	"time"
)

// Conversation represents a chat conversation between two users
// Table: conversations
// Requirements: 8.5 (Chat Message Persistence), 9.1 (Conversation List)
type Conversation struct {
	ID            uint       `gorm:"primarykey" json:"id"`
	ParticipantA  uint       `gorm:"column:participant_a;not null" json:"participant_a"` // 用户A ID
	ParticipantB  uint       `gorm:"column:participant_b;not null" json:"participant_b"` // 用户B ID
	LastMessageID *uint      `gorm:"column:last_message_id" json:"last_message_id"`      // 最后一条消息ID
	LastMessageAt *time.Time `gorm:"column:last_message_at" json:"last_message_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// GORM associations
	LastMessage *ChatMessage  `gorm:"foreignKey:LastMessageID" json:"last_message,omitempty"`
	Messages    []ChatMessage `gorm:"foreignKey:ConversationID" json:"messages,omitempty"`
}

// TableName specifies the table name for Conversation
func (Conversation) TableName() string {
	return "conversations"
}

// ChatMessage represents a single message in a conversation
// Table: chat_messages
// Requirements: 8.5 (Chat Message Persistence)
type ChatMessage struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	ConversationID uint      `gorm:"column:conversation_id;not null" json:"conversation_id"`
	SenderID       uint      `gorm:"column:sender_id;not null" json:"sender_id"`
	Content        string    `gorm:"type:text;not null" json:"content"`
	MessageType    string    `gorm:"column:message_type;size:20;default:'text'" json:"message_type"` // text, image, file
	IsRead         bool      `gorm:"column:is_read;default:false" json:"is_read"`
	CreatedAt      time.Time `json:"created_at"`

	// GORM associations
	Conversation *Conversation `gorm:"foreignKey:ConversationID" json:"conversation,omitempty"`
}

// TableName specifies the table name for ChatMessage
func (ChatMessage) TableName() string {
	return "chat_messages"
}

// ============================================================================
// Helper structs for API responses
// ============================================================================

// ConversationWithDetails is used for API responses with additional computed fields
// Requirements: 9.1 (Conversation List), 9.2 (Last Message Preview), 9.3 (Unread Count)
type ConversationWithDetails struct {
	ID            uint         `json:"id"`
	Participant   *UserInfo    `json:"participant"`  // The other participant (not current user)
	LastMessage   *ChatMessage `json:"last_message"` // Last message in conversation
	UnreadCount   int          `json:"unread_count"` // Number of unread messages
	LastMessageAt *time.Time   `json:"last_message_at"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

// UserInfo represents basic user information for conversation display
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar,omitempty"`
	Role     string `json:"role,omitempty"`
	IsOnline bool   `json:"is_online"`
}

// ConversationListResponse is the response structure for conversation list API
type ConversationListResponse struct {
	Conversations []ConversationWithDetails `json:"conversations"`
	Total         int64                     `json:"total"`
}

// MessageListResponse is the response structure for message list API
type MessageListResponse struct {
	Messages []ChatMessage `json:"messages"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

// SendMessageRequest is the request structure for sending a message
type SendMessageRequest struct {
	Content     string `json:"content" binding:"required"`
	MessageType string `json:"message_type"` // defaults to "text" if not provided
}

// CreateConversationRequest is the request structure for creating a conversation
type CreateConversationRequest struct {
	ParticipantID uint `json:"participant_id" binding:"required"`
}

// CreateConversationResponse is the response structure for creating a conversation
type CreateConversationResponse struct {
	Conversation ConversationWithDetails `json:"conversation"`
	IsNew        bool                    `json:"is_new"` // true if newly created, false if existing
}

// UnreadCountResponse is the response structure for unread count API
type UnreadCountResponse struct {
	TotalUnread int `json:"total_unread"`
}
