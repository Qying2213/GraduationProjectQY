package models

import (
	"time"

	"gorm.io/gorm"
)

// Application 表示一次职位投递记录。
// job-service 查询职位投递列表时会读取这张表，用于 HR 查看某个岗位下的候选人。
type Application struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	JobID       uint           `json:"job_id"`
	TalentID    uint           `json:"talent_id"`
	ResumeID    uint           `json:"resume_id"`
	Status      string         `gorm:"size:20;default:'pending'" json:"status"` // pending, viewed, interview, offer, rejected
	CoverLetter string         `gorm:"type:text" json:"cover_letter"`
	Notes       string         `gorm:"type:text" json:"notes"`
}

// Talent 是职位投递查询中用于关联候选人信息的轻量结构。
// 这里没有复用 talent-service 模型，是为了避免跨服务代码包依赖。
type Talent struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Name            string         `gorm:"size:100;not null" json:"name"`
	Email           string         `gorm:"size:100;not null" json:"email"`
	Phone           string         `gorm:"size:20" json:"phone"`
	Summary         string         `gorm:"type:text" json:"summary"`
	Location        string         `gorm:"size:100" json:"location"`
	CurrentCompany  string         `gorm:"size:100" json:"current_company"`
	CurrentPosition string         `gorm:"size:100" json:"current_position"`
}

// ApplicationWithCandidate 是 HR 查看职位投递列表时的聚合返回结构。
type ApplicationWithCandidate struct {
	ID             uint      `json:"id"`
	JobID          uint      `json:"job_id"`
	TalentID       uint      `json:"talent_id"`
	ResumeID       uint      `json:"resume_id"`
	Status         string    `json:"status"`
	CoverLetter    string    `json:"cover_letter"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	CandidateName  string    `json:"candidate_name"`
	CandidateEmail string    `json:"candidate_email"`
	CandidatePhone string    `json:"candidate_phone"`
	ResumeSummary  string    `json:"resume_summary"`
}
