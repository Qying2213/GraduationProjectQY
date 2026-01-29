package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// BasicInfo 基本信息结构
type BasicInfo struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Location string `json:"location"`
	Avatar   string `json:"avatar,omitempty"`
	Gender   string `json:"gender,omitempty"`
	Age      int    `json:"age,omitempty"`
	Summary  string `json:"summary,omitempty"` // 个人简介
}

// WorkExperience 工作经历结构
type WorkExperience struct {
	Company     string `json:"company"`
	Position    string `json:"position"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"` // 空表示至今
	IsCurrent   bool   `json:"is_current"`
	Description string `json:"description,omitempty"`
	Location    string `json:"location,omitempty"`
}

// Education 教育经历结构
type Education struct {
	School     string `json:"school"`
	Degree     string `json:"degree"` // 学历：本科、硕士、博士等
	Major      string `json:"major"`  // 专业
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date,omitempty"`
	IsCurrent  bool   `json:"is_current"`
	GPA        string `json:"gpa,omitempty"`
	Activities string `json:"activities,omitempty"` // 校园活动
}

// WorkExperienceList 工作经历列表（用于GORM JSON存储）
type WorkExperienceList []WorkExperience

// EducationList 教育经历列表（用于GORM JSON存储）
type EducationList []Education

// SkillList 技能列表（用于GORM JSON存储）
type SkillList []string

// Value 实现 driver.Valuer 接口 - WorkExperienceList
func (w WorkExperienceList) Value() (driver.Value, error) {
	if w == nil {
		return "[]", nil
	}
	return json.Marshal(w)
}

// Scan 实现 sql.Scanner 接口 - WorkExperienceList
func (w *WorkExperienceList) Scan(value interface{}) error {
	if value == nil {
		*w = WorkExperienceList{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, w)
}

// Value 实现 driver.Valuer 接口 - EducationList
func (e EducationList) Value() (driver.Value, error) {
	if e == nil {
		return "[]", nil
	}
	return json.Marshal(e)
}

// Scan 实现 sql.Scanner 接口 - EducationList
func (e *EducationList) Scan(value interface{}) error {
	if value == nil {
		*e = EducationList{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, e)
}

// Value 实现 driver.Valuer 接口 - SkillList
func (s SkillList) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	return json.Marshal(s)
}

// Scan 实现 sql.Scanner 接口 - SkillList
func (s *SkillList) Scan(value interface{}) error {
	if value == nil {
		*s = SkillList{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, s)
}

// OnlineResume 在线简历模型
// 存储用户在线编辑的简历数据
type OnlineResume struct {
	ID             uint               `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	DeletedAt      gorm.DeletedAt     `gorm:"index" json:"-"`
	UserID         uint               `gorm:"uniqueIndex;not null" json:"user_id"` // 关联用户，一个用户只有一份在线简历
	TalentID       *uint              `gorm:"index" json:"talent_id,omitempty"`    // 可选关联人才档案
	Name           string             `gorm:"size:100" json:"name"`
	Phone          string             `gorm:"size:20" json:"phone"`
	Email          string             `gorm:"size:100" json:"email"`
	Location       string             `gorm:"size:100" json:"location,omitempty"`
	Avatar         string             `gorm:"size:500" json:"avatar,omitempty"`
	Gender         string             `gorm:"size:10" json:"gender,omitempty"`
	Age            int                `json:"age,omitempty"`
	Summary        string             `gorm:"type:text" json:"summary,omitempty"`
	WorkExperience WorkExperienceList `gorm:"type:jsonb" json:"work_experience"`
	Education      EducationList      `gorm:"type:jsonb" json:"education"`
	Skills         SkillList          `gorm:"type:jsonb" json:"skills"`
	IsComplete     bool               `gorm:"default:false" json:"is_complete"` // 简历是否完整
}

// TableName 指定表名
func (OnlineResume) TableName() string {
	return "online_resumes"
}

// GetBasicInfo 获取基本信息结构
func (o *OnlineResume) GetBasicInfo() BasicInfo {
	return BasicInfo{
		Name:     o.Name,
		Phone:    o.Phone,
		Email:    o.Email,
		Location: o.Location,
		Avatar:   o.Avatar,
		Gender:   o.Gender,
		Age:      o.Age,
		Summary:  o.Summary,
	}
}

// SetBasicInfo 设置基本信息
func (o *OnlineResume) SetBasicInfo(info BasicInfo) {
	o.Name = info.Name
	o.Phone = info.Phone
	o.Email = info.Email
	o.Location = info.Location
	o.Avatar = info.Avatar
	o.Gender = info.Gender
	o.Age = info.Age
	o.Summary = info.Summary
}

// Validate 验证必填字段
// Requirements 4.6: 验证必填字段（name, phone, email）
func (o *OnlineResume) Validate() error {
	if o.Name == "" {
		return errors.New("姓名不能为空")
	}
	if o.Phone == "" {
		return errors.New("手机号不能为空")
	}
	if o.Email == "" {
		return errors.New("邮箱不能为空")
	}
	return nil
}

// UpdateCompleteness 更新简历完整度状态
func (o *OnlineResume) UpdateCompleteness() {
	// 简历完整需要：基本信息完整 + 至少一条工作经历或教育经历 + 至少一个技能
	hasBasicInfo := o.Name != "" && o.Phone != "" && o.Email != ""
	hasExperience := len(o.WorkExperience) > 0 || len(o.Education) > 0
	hasSkills := len(o.Skills) > 0
	o.IsComplete = hasBasicInfo && hasExperience && hasSkills
}

// ========== API Request/Response 结构 ==========

// OnlineResumeRequest 保存在线简历请求
type OnlineResumeRequest struct {
	BasicInfo      BasicInfo        `json:"basic_info"`
	WorkExperience []WorkExperience `json:"work_experience"`
	Education      []Education      `json:"education"`
	Skills         []string         `json:"skills"`
}

// OnlineResumeResponse 在线简历响应
type OnlineResumeResponse struct {
	ID             uint             `json:"id"`
	BasicInfo      BasicInfo        `json:"basic_info"`
	WorkExperience []WorkExperience `json:"work_experience"`
	Education      []Education      `json:"education"`
	Skills         []string         `json:"skills"`
	IsComplete     bool             `json:"is_complete"`
	UpdatedAt      time.Time        `json:"updated_at"`
}

// ToResponse 将模型转换为响应结构
func (o *OnlineResume) ToResponse() OnlineResumeResponse {
	workExp := make([]WorkExperience, len(o.WorkExperience))
	copy(workExp, o.WorkExperience)

	edu := make([]Education, len(o.Education))
	copy(edu, o.Education)

	skills := make([]string, len(o.Skills))
	copy(skills, o.Skills)

	return OnlineResumeResponse{
		ID:             o.ID,
		BasicInfo:      o.GetBasicInfo(),
		WorkExperience: workExp,
		Education:      edu,
		Skills:         skills,
		IsComplete:     o.IsComplete,
		UpdatedAt:      o.UpdatedAt,
	}
}

// FromRequest 从请求结构更新模型
func (o *OnlineResume) FromRequest(req OnlineResumeRequest) {
	o.SetBasicInfo(req.BasicInfo)
	o.WorkExperience = WorkExperienceList(req.WorkExperience)
	o.Education = EducationList(req.Education)
	o.Skills = SkillList(req.Skills)
	o.UpdateCompleteness()
}
