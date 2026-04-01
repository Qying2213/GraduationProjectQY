package model

import (
	"time"

	"gorm.io/gorm"
)

type ForumPost struct {
	ID           uint           `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	BoardID      uint           `gorm:"index;not null"`
	UserID       uint           `gorm:"index;not null"`
	Title        string         `gorm:"size:200;not null"`
	Content      string         `gorm:"type:text;not null"`
	Status       string         `gorm:"size:20;default:'published';index"`
	IsPinned     bool           `gorm:"default:false;index"`
	IsLocked     bool           `gorm:"default:false"`
	ViewCount    int64          `gorm:"default:0"`
	LikeCount    int64          `gorm:"default:0"`
	CommentCount int64          `gorm:"default:0"`
}

func (ForumPost) TableName() string {
	return "forum_posts"
}

type ForumBoard struct {
	ID          uint           `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string         `gorm:"size:100;not null;uniqueIndex"`
	Description string         `gorm:"size:255"`
	Status      string         `gorm:"size:20;default:'active';index"`
	SortOrder   int            `gorm:"default:0;index"`
}

func (ForumBoard) TableName() string {
	return "forum_boards"
}
