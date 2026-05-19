package models

import (
	"time"

	"gorm.io/gorm"
)

// Message 是站内信/系统通知模型。
// 它用于面试通知、系统提醒、offer/反馈等非聊天会话类消息。
type Message struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	SenderID   *uint          `gorm:"column:sender_id" json:"sender_id"`
	ReceiverID uint           `gorm:"column:receiver_id;not null" json:"receiver_id"`
	Title      string         `gorm:"size:200" json:"title"`
	Content    string         `gorm:"type:text" json:"content"`
	Type       string         `gorm:"size:20;default:'system'" json:"type"` // system, chat, interview, feedback, offer, reminder
	IsRead     bool           `gorm:"default:false" json:"is_read"`
	ReadAt     *time.Time     `json:"read_at,omitempty"`
}
