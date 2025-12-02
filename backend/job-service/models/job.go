package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Job struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Title        string         `gorm:"size:200;not null" json:"title"`
	Description  string         `gorm:"type:text" json:"description"`
	Requirements pq.StringArray `gorm:"type:text[]" json:"requirements"`
	Salary       string         `gorm:"size:100" json:"salary"`
	Location     string         `gorm:"size:100" json:"location"`
	Type         string         `gorm:"size:20;default:'full-time'" json:"type"` // full-time, part-time, contract, internship
	Status       string         `gorm:"size:20;default:'open'" json:"status"`    // open, closed, filled
	CreatedBy    uint           `json:"created_by"`
	Department   string         `gorm:"size:100" json:"department"`
	Level        string         `gorm:"size:50" json:"level"` // junior, mid, senior, lead
	Skills       pq.StringArray `gorm:"type:text[]" json:"skills"`
	Benefits     pq.StringArray `gorm:"type:text[]" json:"benefits"`
}
