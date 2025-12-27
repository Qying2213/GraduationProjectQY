package models

import "time"

// Position 岗位信息
type Position struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UserID           uint      `json:"user_id" gorm:"index:idx_user_post,unique,priority:1"`
	PostID           string    `json:"post_id" gorm:"size:50;index:idx_user_post,unique,priority:2"`
	PostName         string    `json:"post_name" gorm:"size:200"`
	RecruitType      string    `json:"recruit_type" gorm:"size:20"`
	ServiceCondition string    `json:"service_condition" gorm:"type:text"`
	WorkContent      string    `json:"work_content" gorm:"type:text"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// GetJDText 获取完整的 JD 文本（serviceCondition + workContent）
func (p *Position) GetJDText() string {
	jd := ""
	if p.ServiceCondition != "" {
		jd = p.ServiceCondition
	}
	if p.WorkContent != "" {
		if jd != "" {
			jd += "\n\n"
		}
		jd += p.WorkContent
	}
	return jd
}
