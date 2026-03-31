package models

import (
	"time"

	"gorm.io/gorm"
)

type Notice struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Title     string         `gorm:"size:200;not null" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	Status    string         `gorm:"size:20;default:'draft'" json:"status"` // draft, published
	CreatedBy uint           `json:"created_by"`
	IsPinned  bool           `gorm:"default:false;index" json:"is_pinned"`
	Priority  string         `gorm:"size:20;default:'normal';index" json:"priority"` // normal, high, urgent
}
