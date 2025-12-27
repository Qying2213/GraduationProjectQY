package models

import "time"

// Credential stores third-party credentials for an organization.
// PasswordCipher is the encrypted form (AES-GCM recommended).
// NOTE: Do not log plaintext.
type Credential struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `json:"user_id" gorm:"index"` // 关联用户，数据隔离
	Org            string    `gorm:"size:100;uniqueIndex:idx_user_org,priority:2" json:"org"`
	Account        string    `gorm:"size:200" json:"account"`
	PasswordCipher string    `gorm:"type:text" json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
