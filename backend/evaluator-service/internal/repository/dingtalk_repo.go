package repository

import (
	"errors"

	"evaluator-service/internal/models"

	"gorm.io/gorm"
)

type DingTalkRepository struct {
	db *gorm.DB
}

func NewDingTalkRepository(db *gorm.DB) *DingTalkRepository {
	return &DingTalkRepository{db: db}
}

// Get 获取单个配置（兼容旧接口，返回第一个启用的配置）
func (r *DingTalkRepository) Get() (*models.DingTalkConfig, error) {
	var configs []models.DingTalkConfig
	err := r.db.Find(&configs).Error
	if err != nil {
		return nil, err
	}

	for _, cfg := range configs {
		if cfg.Enabled {
			return &cfg, nil
		}
	}

	if len(configs) > 0 {
		return &configs[0], nil
	}

	return nil, nil
}

// GetByUser 获取用户的第一个启用配置
func (r *DingTalkRepository) GetByUser(userID uint) (*models.DingTalkConfig, error) {
	var configs []models.DingTalkConfig
	err := r.db.Where("user_id = ?", userID).Find(&configs).Error
	if err != nil {
		return nil, err
	}

	for _, cfg := range configs {
		if cfg.Enabled {
			return &cfg, nil
		}
	}

	if len(configs) > 0 {
		return &configs[0], nil
	}

	return nil, nil
}

// GetByID 根据ID获取配置
func (r *DingTalkRepository) GetByID(id uint) (*models.DingTalkConfig, error) {
	var config models.DingTalkConfig
	err := r.db.First(&config, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &config, err
}

// GetByIDAndUser 根据ID和用户ID获取配置
func (r *DingTalkRepository) GetByIDAndUser(id, userID uint) (*models.DingTalkConfig, error) {
	var config models.DingTalkConfig
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("无权访问该资源")
	}
	return &config, err
}

// GetByName 根据名称获取配置
func (r *DingTalkRepository) GetByName(name string) (*models.DingTalkConfig, error) {
	var config models.DingTalkConfig
	err := r.db.Where("name = ?", name).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &config, err
}

// GetByNameAndUser 根据名称和用户ID获取配置
func (r *DingTalkRepository) GetByNameAndUser(name string, userID uint) (*models.DingTalkConfig, error) {
	var config models.DingTalkConfig
	err := r.db.Where("name = ? AND user_id = ?", name, userID).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &config, err
}

// List 获取所有配置
func (r *DingTalkRepository) List() ([]models.DingTalkConfig, error) {
	var configs []models.DingTalkConfig
	err := r.db.Order("created_at DESC").Find(&configs).Error
	return configs, err
}

// ListByUser 获取用户的所有配置
func (r *DingTalkRepository) ListByUser(userID uint) ([]models.DingTalkConfig, error) {
	var configs []models.DingTalkConfig
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&configs).Error
	return configs, err
}

// ListEnabled 获取所有启用的配置
func (r *DingTalkRepository) ListEnabled() ([]models.DingTalkConfig, error) {
	var configs []models.DingTalkConfig
	var allConfigs []models.DingTalkConfig
	err := r.db.Order("created_at DESC").Find(&allConfigs).Error
	if err != nil {
		return nil, err
	}

	for _, cfg := range allConfigs {
		if cfg.Enabled {
			configs = append(configs, cfg)
		}
	}

	return configs, nil
}

// ListEnabledByUser 获取用户所有启用的配置
func (r *DingTalkRepository) ListEnabledByUser(userID uint) ([]models.DingTalkConfig, error) {
	var configs []models.DingTalkConfig
	var allConfigs []models.DingTalkConfig
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&allConfigs).Error
	if err != nil {
		return nil, err
	}

	for _, cfg := range allConfigs {
		if cfg.Enabled {
			configs = append(configs, cfg)
		}
	}

	return configs, nil
}

// Upsert 创建或更新配置
func (r *DingTalkRepository) Upsert(config *models.DingTalkConfig) error {
	if config.ID > 0 {
		return r.db.Save(config).Error
	}

	if config.Name != "" {
		var existing models.DingTalkConfig
		err := r.db.Where("name = ?", config.Name).First(&existing).Error

		if err == nil {
			config.ID = existing.ID
			config.CreatedAt = existing.CreatedAt
			return r.db.Save(config).Error
		}

		if err != gorm.ErrRecordNotFound {
			return err
		}
	}

	return r.db.Create(config).Error
}

// UpsertByUser 创建或更新用户的配置
func (r *DingTalkRepository) UpsertByUser(config *models.DingTalkConfig, userID uint) error {
	config.UserID = userID

	if config.ID > 0 {
		// 验证所有权
		var existing models.DingTalkConfig
		err := r.db.Where("id = ? AND user_id = ?", config.ID, userID).First(&existing).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("无权访问该资源")
			}
			return err
		}
		return r.db.Save(config).Error
	}

	if config.Name != "" {
		var existing models.DingTalkConfig
		err := r.db.Where("name = ? AND user_id = ?", config.Name, userID).First(&existing).Error

		if err == nil {
			config.ID = existing.ID
			config.CreatedAt = existing.CreatedAt
			return r.db.Save(config).Error
		}

		if err != gorm.ErrRecordNotFound {
			return err
		}
	}

	return r.db.Create(config).Error
}

// Delete 删除配置
func (r *DingTalkRepository) Delete(id uint) error {
	return r.db.Delete(&models.DingTalkConfig{}, id).Error
}

// DeleteByUser 删除用户的配置
func (r *DingTalkRepository) DeleteByUser(id, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.DingTalkConfig{})
	if result.RowsAffected == 0 {
		return errors.New("无权访问该资源")
	}
	return result.Error
}
