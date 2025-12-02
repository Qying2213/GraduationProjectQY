package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	FromID    uint           `json:"from_id"`
	ToID      uint           `json:"to_id"`
	Title     string         `gorm:"size:200" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	Type      string         `gorm:"size:20;default:'system'" json:"type"` // system, user, notification
	IsRead    bool           `gorm:"default:false" json:"is_read"`
	RelatedID uint           `json:"related_id,omitempty"` // 关联的对象ID（如职位ID、申请ID）
}
