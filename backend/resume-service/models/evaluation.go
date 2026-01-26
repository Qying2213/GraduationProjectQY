package models

import (
	"time"

	"gorm.io/gorm"
)

// EvaluationResult AI评估结果
type EvaluationResult struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联信息
	ResumeID uint  `json:"resume_id" gorm:"index"`
	TalentID *uint `json:"talent_id" gorm:"index"`
	JobID    *uint `json:"job_id" gorm:"index"`

	// 简历信息
	ResumeName string `json:"resume_name" gorm:"size:100"`
	ResumeFile string `json:"resume_file" gorm:"size:500"`

	// 解析结果
	ParsedName       string `json:"parsed_name" gorm:"size:50"`
	ParsedPhone      string `json:"parsed_phone" gorm:"size:20"`
	ParsedEmail      string `json:"parsed_email" gorm:"size:100"`
	ParsedEducation  string `json:"parsed_education" gorm:"size:50"`
	ParsedExperience string `json:"parsed_experience" gorm:"size:50"`
	ParsedLocation   string `json:"parsed_location" gorm:"size:50"`
	ParsedSkills     string `json:"parsed_skills" gorm:"type:text"` // JSON数组

	// 匹配结果
	MatchScore   float64 `json:"match_score"`
	MatchLevel   string  `json:"match_level" gorm:"size:20"`     // high, medium, low
	MatchDetails string  `json:"match_details" gorm:"type:text"` // JSON数组

	// 归因报告
	ReportSummary        string `json:"report_summary" gorm:"type:text"`
	ReportStrengths      string `json:"report_strengths" gorm:"type:text"` // JSON数组
	ReportGaps           string `json:"report_gaps" gorm:"type:text"`      // JSON数组
	ReportRecommendation string `json:"report_recommendation" gorm:"type:text"`
	ReportDimensions     string `json:"report_dimensions" gorm:"type:text"` // JSON数组

	// 风控信息
	RiskScore int    `json:"risk_score"`
	RiskItems string `json:"risk_items" gorm:"type:text"` // JSON数组

	// 状态
	Status   string `json:"status" gorm:"size:20;default:'completed'"` // pending, completed, failed
	ErrorMsg string `json:"error_msg" gorm:"type:text"`
	EvalType string `json:"eval_type" gorm:"size:20"` // parse, match, full
}

// TableName 指定表名
func (EvaluationResult) TableName() string {
	return "evaluation_results"
}
