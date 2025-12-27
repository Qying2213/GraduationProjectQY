package repository

import (
	"time"

	"evaluator-service/internal/models"

	"gorm.io/gorm"
)

// UserRepository 用户数据访问层
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建新用户
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByID 根据 ID 获取用户
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByCorpCodeAndUsername 根据组织代码和用户名获取用户
func (r *UserRepository) GetByCorpCodeAndUsername(corpCode, username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("corp_code = ? AND username = ?", corpCode, username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateLastLogin 更新最后登录时间
func (r *UserRepository) UpdateLastLogin(id uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("last_login_at", time.Now()).Error
}

// UpdatePassword 更新用户密码（加密后的）
func (r *UserRepository) UpdatePassword(id uint, passwordCipher string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("password_cipher", passwordCipher).Error
}

// FindOrCreate 查找或创建用户，返回用户和是否新创建
func (r *UserRepository) FindOrCreate(corpCode, username, passwordCipher string) (*models.User, bool, error) {
	var user models.User
	err := r.db.Where("corp_code = ? AND username = ?", corpCode, username).First(&user).Error
	if err == nil {
		// 用户已存在，更新密码和登录时间
		user.PasswordCipher = passwordCipher
		user.LastLoginAt = time.Now()
		if err := r.db.Save(&user).Error; err != nil {
			return nil, false, err
		}
		return &user, false, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, false, err
	}
	// 创建新用户
	user = models.User{
		CorpCode:       corpCode,
		Username:       username,
		PasswordCipher: passwordCipher,
		LastLoginAt:    time.Now(),
	}
	if err := r.db.Create(&user).Error; err != nil {
		return nil, false, err
	}
	return &user, true, nil
}
