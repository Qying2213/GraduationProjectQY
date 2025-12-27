package models

import "time"

type Candidate struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UserID           uint      `json:"user_id" gorm:"index"` // 关联用户，数据隔离
	ApplyID          string    `json:"apply_id" gorm:"size:100;index:idx_user_apply,unique,priority:2"` // 招聘系统申请ID，与UserID组合唯一
	Name             string    `json:"name" gorm:"index;size:100"`
	Filename         string    `json:"filename" gorm:"size:500"`
	PDFPath          string    `json:"pdf_path" gorm:"size:500"`
	TotalScore       float64   `json:"total_score"`
	Grade            string    `json:"grade" gorm:"size:10"`
	JDMatch          int       `json:"jd_match"`
	AgeScore         int       `json:"age_score"`
	ExperienceScore  int       `json:"experience_score"`
	EducationScore   int       `json:"education_score"`
	CompanyScore     int       `json:"company_score"`
	TechScore        int       `json:"tech_score"`
	ProjectScore     int       `json:"project_score"`
	AgeReason        string    `json:"age_reason" gorm:"type:text"`
	ExperienceReason string    `json:"experience_reason" gorm:"type:text"`
	EducationReason  string    `json:"education_reason" gorm:"type:text"`
	CompanyReason    string    `json:"company_reason" gorm:"type:text"`
	TechReason       string    `json:"tech_reason" gorm:"type:text"`
	ProjectReason    string    `json:"project_reason" gorm:"type:text"`
	Recommendation   string    `json:"recommendation" gorm:"size:100"`
	ReportMarkdown   string    `json:"report_markdown" gorm:"type:text"`
	ResumeMarkdown   string    `json:"resume_markdown" gorm:"type:text"`
	CozeReportJSON   string    `json:"coze_report_json" gorm:"type:text"` // 完整的 Coze 评估报告 JSON
	Status           string    `json:"status" gorm:"size:20;default:'待面试'"`
	Notes            string    `json:"notes" gorm:"type:text"`
	// 钉钉通知相关字段
	NotifyCount   int        `json:"notify_count" gorm:"default:0"`
	FirstNotifyAt *time.Time `json:"first_notify_at"`
	LastNotifyAt  *time.Time `json:"last_notify_at"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type JDMatchResult struct {
	Score         int      `json:"score"`
	MatchedSkills []string `json:"matched_skills"`
	MissingSkills []string `json:"missing_skills"`
	Summary       string   `json:"summary"`
}

type RequirementResult struct {
	EducationPass   bool     `json:"education_pass"`
	EducationDetail string   `json:"education_detail"`
	ExperiencePass  bool     `json:"experience_pass"`
	ExperienceYears *float64 `json:"experience_years"`
	BlacklistPass   bool     `json:"blacklist_pass"`
	BlacklistHits   []string `json:"blacklist_hits"`
	OverallPass     bool     `json:"overall_pass"`
}

type ScoringResult struct {
	TotalScore       float64 `json:"total_score"`
	Grade            string  `json:"grade"`
	AgeScore         int     `json:"age_score"`
	ExperienceScore  int     `json:"experience_score"`
	EducationScore   int     `json:"education_score"`
	CompanyScore     int     `json:"company_score"`
	TechScore        int     `json:"tech_score"`
	ProjectScore     int     `json:"project_score"`
	AgeReason        string  `json:"age_reason"`
	ExperienceReason string  `json:"experience_reason"`
	EducationReason  string  `json:"education_reason"`
	CompanyReason    string  `json:"company_reason"`
	TechReason       string  `json:"tech_reason"`
	ProjectReason    string  `json:"project_reason"`
}

type EvaluationResult struct {
	JDMatch        JDMatchResult     `json:"jd_match"`
	Requirement    RequirementResult `json:"requirement"`
	Recommendation string            `json:"recommendation"`
}

type InterviewQuestion struct {
	Question       string `json:"question"`
	Category       string `json:"category"`
	Purpose        string `json:"purpose"`
	RedFlag        string `json:"red_flag"`
	ExpectedAnswer string `json:"expected_answer"`
}
