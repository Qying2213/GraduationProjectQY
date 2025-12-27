package repository

import (
	"errors"

	"evaluator-service/internal/models"

	"gorm.io/gorm"
)

type CredentialRepository struct{ db *gorm.DB }

func NewCredentialRepository(db *gorm.DB) *CredentialRepository { return &CredentialRepository{db: db} }

func (r *CredentialRepository) GetByOrg(org string) (*models.Credential, error) {
	var c models.Credential
	if err := r.db.Where("org = ?", org).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

// GetByOrgAndUser 根据组织和用户ID获取凭据
func (r *CredentialRepository) GetByOrgAndUser(org string, userID uint) (*models.Credential, error) {
	var c models.Credential
	if err := r.db.Where("org = ? AND user_id = ?", org, userID).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

// ListByUser 获取用户的所有凭据
func (r *CredentialRepository) ListByUser(userID uint) ([]models.Credential, error) {
	var list []models.Credential
	if err := r.db.Where("user_id = ?", userID).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *CredentialRepository) Upsert(org, account, passwordCipher string) (*models.Credential, error) {
	var c models.Credential
	err := r.db.Where("org = ?", org).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c = models.Credential{Org: org, Account: account, PasswordCipher: passwordCipher}
		if err := r.db.Create(&c).Error; err != nil {
			return nil, err
		}
		return &c, nil
	}
	if err != nil {
		return nil, err
	}
	c.Account = account
	c.PasswordCipher = passwordCipher
	if err := r.db.Save(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

// UpsertByUser 创建或更新用户的凭据
func (r *CredentialRepository) UpsertByUser(org, account, passwordCipher string, userID uint) (*models.Credential, error) {
	var c models.Credential
	err := r.db.Where("org = ? AND user_id = ?", org, userID).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c = models.Credential{
			UserID:         userID,
			Org:            org,
			Account:        account,
			PasswordCipher: passwordCipher,
		}
		if err := r.db.Create(&c).Error; err != nil {
			return nil, err
		}
		return &c, nil
	}
	if err != nil {
		return nil, err
	}
	c.Account = account
	c.PasswordCipher = passwordCipher
	if err := r.db.Save(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

// DeleteByUser 删除用户的凭据
func (r *CredentialRepository) DeleteByUser(org string, userID uint) error {
	result := r.db.Where("org = ? AND user_id = ?", org, userID).Delete(&models.Credential{})
	return result.Error
}
