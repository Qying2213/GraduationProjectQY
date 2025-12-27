package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Username   string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email      string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password   string         `gorm:"size:255;not null" json:"-"`
	Role       string         `gorm:"size:20;default:'candidate'" json:"role"` // admin, hr_manager, recruiter, interviewer, viewer
	Avatar     string         `gorm:"size:255" json:"avatar"`
	Phone      string         `gorm:"size:20" json:"phone"`
	Department string         `gorm:"size:50" json:"department"`
	Position   string         `gorm:"size:50" json:"position"`
	RealName   string         `gorm:"size:50" json:"real_name"`
	Status     string         `gorm:"size:20;default:'active'" json:"status"` // active, inactive, suspended
}

// HashPassword 密码加密
func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
