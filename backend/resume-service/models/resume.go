package models

import (
	"time"

	"gorm.io/gorm"
)

// Resume 表示系统中的一份简历记录。
// 这张表同时承担两类职责：
// 1. 保存上传文件本身的元信息（文件名、路径、大小、类型）
// 2. 保存后续处理结果（OCR 文本、结构化解析结果、匹配分）
type Resume struct {
	ID            uint           `gorm:"primarykey" json:"id"`                    // 主键 ID
	CreatedAt     time.Time      `json:"created_at"`                              // 创建时间，通常对应上传或创建记录的时间
	UpdatedAt     time.Time      `json:"updated_at"`                              // 最近一次更新时间
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`                          // 软删除时间；json:"-" 表示接口不向前端暴露
	TalentID      *uint          `json:"talent_id"`                               // 关联的人才 ID，可为空，表示简历先上传、后绑定候选人
	JobID         *uint          `json:"job_id"`                                  // 关联的职位 ID，可为空，表示简历尚未用于某个具体岗位
	FileName      string         `gorm:"size:255" json:"file_name"`               // 原始文件名，便于页面展示和下载
	FilePath      string         `gorm:"size:500" json:"file_path"`               // 服务端磁盘上的真实路径，仅后端内部使用
	FileURL       string         `gorm:"size:500" json:"file_url"`                // 对外访问路径，前端通常通过这个地址预览或下载
	FileSize      int64          `json:"file_size"`                               // 文件大小，单位字节
	FileType      string         `gorm:"size:20" json:"file_type"`                // 文件扩展名，如 .pdf / .doc / .docx
	ParsedData    string         `gorm:"type:text" json:"parsed_data"`            // 结构化解析结果，通常以 JSON 字符串方式存储
	ExtractedText string         `gorm:"type:text" json:"extracted_text"`         // OCR 或文本抽取得到的原始文本内容
	MatchScore    int            `json:"match_score"`                             // 简历与职位匹配后的分数，供排序或筛选使用
	Status        string         `gorm:"size:20;default:'pending'" json:"status"` // 简历处理状态：pending / parsed / active / archived
}

// Application 表示一次职位投递记录。
// 阅读这张表时要重点区分两个概念：
// 1. Stage: 招聘流程走到哪一步，例如 screening / interview / offer
// 2. Status: 当前处理结果或动作状态，例如 pending / reviewed / rejected
// 这两个字段同时存在时，Stage 更偏“流程位置”，Status 更偏“当前结论”。
type Application struct {
	ID          uint           `gorm:"primarykey" json:"id"`                    // 投递记录主键
	CreatedAt   time.Time      `json:"created_at"`                              // 投递创建时间
	UpdatedAt   time.Time      `json:"updated_at"`                              // 投递最近更新时间
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                          // 软删除时间，便于支持“撤回申请”而不是物理删除
	JobID       uint           `json:"job_id"`                                  // 被投递的职位 ID
	TalentID    uint           `json:"talent_id"`                               // 投递者对应的人才 ID
	ResumeID    uint           `json:"resume_id"`                               // 本次投递所使用的简历 ID
	Stage       string         `gorm:"size:50;default:'applied'" json:"stage"`  // 招聘流程阶段：applied / screening / interview / offer / hired / rejected
	Status      string         `gorm:"size:20;default:'pending'" json:"status"` // 当前处理状态：pending / reviewed / interview / rejected / accepted
	Source      string         `gorm:"size:50" json:"source"`                   // 投递来源渠道，例如官网、内推、Boss、智联
	CoverLetter string         `gorm:"type:text" json:"cover_letter"`           // 求职者投递时附带的自荐信或补充说明
	Notes       string         `gorm:"type:text" json:"notes"`                  // HR 内部备注，也会被用来记录状态流转历史
}
