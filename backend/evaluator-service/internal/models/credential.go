package models

import "time"

// Credential 保存某个组织的第三方系统凭据。
// PasswordCipher 是加密后的密码内容，禁止在日志中输出明文。
type Credential struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `json:"user_id" gorm:"index"` // 关联用户，数据隔离
	Org            string    `gorm:"size:100;uniqueIndex:idx_user_org,priority:2" json:"org"`
	Account        string    `gorm:"size:200" json:"account"`
	PasswordCipher string    `gorm:"type:text" json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
