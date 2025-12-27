package models

import "time"

// User 系统用户，使用 Wintalent 账号登录
type User struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CorpCode       string    `gorm:"size:100;index" json:"corp_code"`
	Username       string    `gorm:"size:200;uniqueIndex:idx_corp_user" json:"username"`
	PasswordCipher string    `gorm:"type:text" json:"-"` // AES-256 加密存储的 Wintalent 密码
	LastLoginAt    time.Time `json:"last_login_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
