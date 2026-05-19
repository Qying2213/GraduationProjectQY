package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 是系统账号表模型。
// 它既服务登录鉴权，也承载后台角色、部门职位和候选人基础资料。
type User struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Username   string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email      string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password   string         `gorm:"size:255;not null" json:"-"`
	Role       string         `gorm:"size:20;default:'candidate'" json:"role"` // admin, hr_manager, recruiter, interviewer, viewer
	Avatar     string         `gorm:"type:text" json:"avatar"`
	Phone      string         `gorm:"size:20" json:"phone"`
	Department string         `gorm:"size:50" json:"department"`
	Position   string         `gorm:"size:50" json:"position"`
	RealName   string         `gorm:"size:50" json:"real_name"`
	Status     string         `gorm:"size:20;default:'active'" json:"status"` // active, inactive, suspended
}

// HashPassword 使用 bcrypt 存储密码。
// 数据库中不保存明文密码，避免账号数据泄露时直接暴露用户密码。

func (u *User) HashPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

// CheckPassword 校验登录密码。
// 先按 bcrypt 校验，失败后兼容历史明文数据，便于旧数据平滑过渡。
func (u *User) CheckPassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err == nil {
		return true
	}
	return u.Password == password
}
