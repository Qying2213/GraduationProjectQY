package models

import (
	"time"

	"gorm.io/gorm"
)

// AIProcessLog 记录一次简历评估中可观察的 AI 执行链路。
//
// 只有最终分数不足以支撑毕业设计演示和问题排查。该模型会保存序列化后的
// ProcessTrace JSON，让前端展示 OCR、Embedding、RAG、LLM 是否启用、成功、跳过或失败。
type AIProcessLog struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	EvaluationID *uint  `json:"evaluation_id" gorm:"index"`     // 可为空，兼容早期只按 resume_id 记录的链路日志
	ResumeID     uint   `json:"resume_id" gorm:"index"`         // 简历是链路追踪的业务主线
	Status       string `json:"status" gorm:"size:20;index"`    // completed, failed
	ProcessTrace string `json:"process_trace" gorm:"type:text"` // OCR/Embedding/RAG/LLM 的结构化 JSON
	ErrorMsg     string `json:"error_msg" gorm:"type:text"`     // 失败时给前端和日志排查使用
}

// TableName 指定表名
func (AIProcessLog) TableName() string {
	return "ai_process_logs"
}
