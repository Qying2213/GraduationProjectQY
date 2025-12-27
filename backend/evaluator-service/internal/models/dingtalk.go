package models

import "time"

// DingTalkConfig 钉钉机器人配置
type DingTalkConfig struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `json:"user_id" gorm:"index"` // 关联用户，数据隔离
	Name         string    `gorm:"size:100;uniqueIndex:idx_user_name,priority:2" json:"name"` // 机器人名称，用于区分多个机器人
	ClientID     string    `gorm:"size:200" json:"client_id"`
	ClientSecret string    `gorm:"size:200" json:"client_secret"`
	Webhook      string    `gorm:"type:text" json:"webhook"`
	Secret       string    `gorm:"size:200" json:"secret"`
	PushTime     string    `gorm:"size:10;default:'09:00'" json:"push_time"` // 格式: HH:MM
	PushLimit    int       `gorm:"default:10" json:"push_limit"`             // 每次推送候选人数量
	AtUserIDs          string `gorm:"type:text" json:"at_user_ids"`                    // 逗号分隔的钉钉UserID
	Enabled            bool   `gorm:"default:false" json:"enabled"`                    // 是否启用钉钉机器人
	AutoPushOnComplete bool   `gorm:"default:false" json:"auto_push_on_complete"`      // 评估完成后自动推送
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定表名为 ding_talk_configs
func (DingTalkConfig) TableName() string {
	return "ding_talk_configs"
}
