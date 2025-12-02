package models

import (
	"time"

	"gorm.io/gorm"
)

type Resume struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	TalentID   uint           `json:"talent_id"`
	FileName   string         `gorm:"size:255" json:"file_name"`
	FileURL    string         `gorm:"size:500" json:"file_url"`
	FileSize   int64          `json:"file_size"`
	ParsedData string         `gorm:"type:text" json:"parsed_data"` // JSON格式存储解析后的数据
	Status     string         `gorm:"size:20;default:'active'" json:"status"` // active, archived
}

type Application struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	JobID     uint           `json:"job_id"`
	TalentID  uint           `json:"talent_id"`
	ResumeID  uint           `json:"resume_id"`
	Status    string         `gorm:"size:20;default:'pending'" json:"status"` // pending, reviewed, interview, rejected, accepted
	CoverLetter string       `gorm:"type:text" json:"cover_letter"`
	Notes     string         `gorm:"type:text" json:"notes"`
}
