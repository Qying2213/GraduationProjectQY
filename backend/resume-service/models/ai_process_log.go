package models

import (
	"time"

	"gorm.io/gorm"
)

// AIProcessLog 记录 AI 评估链路过程（OCR -> Embedding -> RAG -> LLM）
type AIProcessLog struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	EvaluationID *uint  `json:"evaluation_id" gorm:"index"`
	ResumeID     uint   `json:"resume_id" gorm:"index"`
	Status       string `json:"status" gorm:"size:20;index"` // completed, failed
	ProcessTrace string `json:"process_trace" gorm:"type:text"`
	ErrorMsg     string `json:"error_msg" gorm:"type:text"`
}

// TableName 指定表名
func (AIProcessLog) TableName() string {
	return "ai_process_logs"
}
