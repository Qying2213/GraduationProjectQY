package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Talent struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Name            string         `gorm:"size:100;not null" json:"name"`
	Email           string         `gorm:"size:100;not null" json:"email"`
	Phone           string         `gorm:"size:20" json:"phone"`
	Skills          pq.StringArray `gorm:"type:text[]" json:"skills"`
	Experience      int            `json:"experience"`
	Education       string         `gorm:"size:50" json:"education"`
	Status          string         `gorm:"size:20;default:'active'" json:"status"`
	Tags            pq.StringArray `gorm:"type:text[]" json:"tags"`
	UserID          *uint          `json:"user_id,omitempty"`
	Location        string         `gorm:"size:100" json:"location"`
	Salary          string         `gorm:"size:50" json:"salary"`
	Summary         string         `gorm:"type:text" json:"summary"`
	Gender          string         `gorm:"size:10" json:"gender"`
	Age             int            `json:"age"`
	CurrentCompany  string         `gorm:"size:100" json:"current_company"`
	CurrentPosition string         `gorm:"size:100" json:"current_position"`
	Source          string         `gorm:"size:50" json:"source"`
	ResumeID        *uint          `json:"resume_id,omitempty"`
}
