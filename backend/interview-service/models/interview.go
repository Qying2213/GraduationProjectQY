package models

import (
	"time"

	"gorm.io/gorm"
)

// InterviewType 面试类型
type InterviewType string

const (
	InterviewTypeInitial InterviewType = "initial" // 初试
	InterviewTypeSecond  InterviewType = "second"  // 复试
	InterviewTypeFinal   InterviewType = "final"   // 终面
	InterviewTypeHR      InterviewType = "hr"      // HR面试
)

// InterviewMethod 面试方式
type InterviewMethod string

const (
	InterviewMethodOnsite InterviewMethod = "onsite" // 现场面试
	InterviewMethodVideo  InterviewMethod = "video"  // 视频面试
	InterviewMethodPhone  InterviewMethod = "phone"  // 电话面试
)

// InterviewStatus 面试状态
type InterviewStatus string

const (
	InterviewStatusScheduled InterviewStatus = "scheduled" // 已安排
	InterviewStatusCompleted InterviewStatus = "completed" // 已完成
	InterviewStatusCancelled InterviewStatus = "cancelled" // 已取消
	InterviewStatusNoShow    InterviewStatus = "no_show"   // 爽约
)

// Interview 是一次面试安排记录。
// 它把候选人、职位、面试官、时间、方式和状态放在一起，支撑日程管理和面试流程流转。
type Interview struct {
	ID            uint            `json:"id" gorm:"primaryKey"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	CandidateID   uint            `json:"candidate_id" gorm:"index;not null"`
	CandidateName string          `json:"candidate_name" gorm:"size:100;not null"`
	PositionID    uint            `json:"position_id" gorm:"index;not null"`
	Position      string          `json:"position" gorm:"size:200;not null"`
	Type          InterviewType   `json:"type" gorm:"size:20;not null;default:'initial'"`
	Date          string          `json:"date" gorm:"size:20;not null"` // YYYY-MM-DD
	Time          string          `json:"time" gorm:"size:10;not null"` // HH:mm
	Duration      int             `json:"duration" gorm:"default:60"`   // 分钟
	InterviewerID uint            `json:"interviewer_id" gorm:"index;not null"`
	Interviewer   string          `json:"interviewer" gorm:"size:100;not null"`
	Method        InterviewMethod `json:"method" gorm:"size:20;not null;default:'onsite'"`
	Location      string          `json:"location" gorm:"size:500"` // 地点或会议链接
	Status        InterviewStatus `json:"status" gorm:"size:20;not null;default:'scheduled'"`
	Notes         string          `json:"notes" gorm:"type:text"`     // 备注
	Feedback      string          `json:"feedback" gorm:"type:text"`  // 面试反馈
	Rating        int             `json:"rating" gorm:"default:0"`    // 评分 1-5
	CreatedBy     uint            `json:"created_by" gorm:"not null"` // 创建人
}

// InterviewFeedback 是面试官提交的结构化反馈。
// 它与 Interview 分表保存，方便多维度记录优势、不足、评分和通过建议。
type InterviewFeedback struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	InterviewID    uint      `json:"interview_id" gorm:"index;not null"`
	InterviewerID  uint      `json:"interviewer_id" gorm:"index;not null"`
	Rating         int       `json:"rating" gorm:"not null"` // 1-5
	Strengths      string    `json:"strengths" gorm:"type:text"`
	Weaknesses     string    `json:"weaknesses" gorm:"type:text"`
	Comments       string    `json:"comments" gorm:"type:text"`
	Recommendation string    `json:"recommendation" gorm:"size:50"` // pass, fail, pending
}

// InterviewScheduleRequest 是创建面试时前端提交的请求结构。
type InterviewScheduleRequest struct {
	CandidateID   uint   `json:"candidate_id" binding:"required"`
	CandidateName string `json:"candidate_name" binding:"required"`
	PositionID    uint   `json:"position_id" binding:"required"`
	Position      string `json:"position" binding:"required"`
	Type          string `json:"type" binding:"required"`
	Date          string `json:"date" binding:"required"`
	Time          string `json:"time" binding:"required"`
	Duration      int    `json:"duration"`
	InterviewerID uint   `json:"interviewer_id" binding:"required"`
	Interviewer   string `json:"interviewer" binding:"required"`
	Method        string `json:"method"`
	Location      string `json:"location"`
	Notes         string `json:"notes"`
}

// InterviewUpdateRequest 面试更新请求
type InterviewUpdateRequest struct {
	Type          string `json:"type"`
	Date          string `json:"date"`
	Time          string `json:"time"`
	Duration      int    `json:"duration"`
	InterviewerID uint   `json:"interviewer_id"`
	Interviewer   string `json:"interviewer"`
	Method        string `json:"method"`
	Location      string `json:"location"`
	Status        string `json:"status"`
	Notes         string `json:"notes"`
	Feedback      string `json:"feedback"`
	Rating        int    `json:"rating"`
}

// InterviewListResponse 面试列表响应
type InterviewListResponse struct {
	Interviews []Interview `json:"interviews"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
}

// InterviewStats 面试统计
type InterviewStats struct {
	TotalInterviews     int64 `json:"total_interviews"`
	ScheduledInterviews int64 `json:"scheduled_interviews"`
	CompletedInterviews int64 `json:"completed_interviews"`
	CancelledInterviews int64 `json:"cancelled_interviews"`
	TodayInterviews     int64 `json:"today_interviews"`
	WeekInterviews      int64 `json:"week_interviews"`
}

// BeforeCreate 创建前钩子
func (i *Interview) BeforeCreate(tx *gorm.DB) error {
	if i.Duration == 0 {
		i.Duration = 60
	}
	if i.Type == "" {
		i.Type = InterviewTypeInitial
	}
	if i.Method == "" {
		i.Method = InterviewMethodOnsite
	}
	if i.Status == "" {
		i.Status = InterviewStatusScheduled
	}
	return nil
}

// TableName 设置表名
func (Interview) TableName() string {
	return "interviews"
}

func (InterviewFeedback) TableName() string {
	return "interview_feedbacks"
}

// GetScheduledAt 获取面试时间
func (i *Interview) GetScheduledAt() (time.Time, error) {
	return time.Parse("2006-01-02 15:04", i.Date+" "+i.Time)
}
