package models

import (
	"time"

	"gorm.io/gorm"
)

// EvaluationResult 是 AI 评估完成后持久化的业务结果记录。
//
// 它同时保存结构化字段和原始 JSON：结构化字段方便管理页筛选和排序，原始 JSON
// 保留完整模型输出，便于后续审计、报告渲染和不同 prompt/版本效果对比。
type EvaluationResult struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联信息：把一次 AI 评估挂回简历、人才和职位，形成招聘业务闭环。
	ResumeID uint  `json:"resume_id" gorm:"index"`
	TalentID *uint `json:"talent_id" gorm:"index"`
	JobID    *uint `json:"job_id" gorm:"index"`

	// 简历信息：冗余文件名和路径，避免列表页频繁 join resume 表。
	ResumeName string `json:"resume_name" gorm:"size:100"`
	ResumeFile string `json:"resume_file" gorm:"size:500"`

	// 解析结果：从 LLM/解析流程抽取出的候选人基础画像。
	ParsedName       string `json:"parsed_name" gorm:"size:50"`
	ParsedPhone      string `json:"parsed_phone" gorm:"size:20"`
	ParsedEmail      string `json:"parsed_email" gorm:"size:100"`
	ParsedEducation  string `json:"parsed_education" gorm:"size:50"`
	ParsedSchool     string `json:"parsed_school" gorm:"size:100"`
	ParsedExperience string `json:"parsed_experience" gorm:"size:50"`
	ParsedLocation   string `json:"parsed_location" gorm:"size:50"`
	ParsedSkills     string `json:"parsed_skills" gorm:"type:text"` // JSON数组
	ParsedReport     string `json:"parsed_report" gorm:"type:text"` // 完整结构化报告 JSON
	RawResult        string `json:"raw_result" gorm:"type:text"`    // Coze 原始返回 JSON

	// 匹配结果：用于 HR 快速判断候选人与职位的契合度。
	MatchScore   float64 `json:"match_score"`
	MatchLevel   string  `json:"match_level" gorm:"size:20"`     // high, medium, low
	MatchDetails string  `json:"match_details" gorm:"type:text"` // JSON数组

	// 归因报告：支撑“为什么推荐/为什么扣分”的可解释性展示。
	ReportSummary        string `json:"report_summary" gorm:"type:text"`
	ReportStrengths      string `json:"report_strengths" gorm:"type:text"` // JSON数组
	ReportGaps           string `json:"report_gaps" gorm:"type:text"`      // JSON数组
	ReportRecommendation string `json:"report_recommendation" gorm:"type:text"`
	ReportDimensions     string `json:"report_dimensions" gorm:"type:text"` // JSON数组

	// 风控信息：把简历时间线、学历、经历等风险点沉淀为可追踪字段。
	RiskScore int    `json:"risk_score"`
	RiskItems string `json:"risk_items" gorm:"type:text"` // JSON数组

	// 状态：区分评估处理中、成功和失败，方便前端轮询与历史审计。
	Status   string `json:"status" gorm:"size:20;default:'completed'"` // pending, completed, failed
	ErrorMsg string `json:"error_msg" gorm:"type:text"`
	EvalType string `json:"eval_type" gorm:"size:20"` // parse, match, full
}

// TableName 指定表名
func (EvaluationResult) TableName() string {
	return "evaluation_results"
}
