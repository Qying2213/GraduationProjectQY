package repository

import (
	"errors"

	"evaluator-service/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PositionRepository struct {
	db *gorm.DB
}

func NewPositionRepository(db *gorm.DB) *PositionRepository {
	return &PositionRepository{db: db}
}

func (r *PositionRepository) DB() *gorm.DB {
	return r.db
}

// UpsertByPostID 根据 PostID 和 UserID 插入或更新岗位
func (r *PositionRepository) UpsertByPostID(p *models.Position) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "post_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"post_name", "recruit_type", "service_condition", "work_content", "updated_at"}),
	}).Create(p).Error
}

// BatchUpsert 批量插入或更新岗位
func (r *PositionRepository) BatchUpsert(positions []models.Position) error {
	if len(positions) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "post_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"post_name", "recruit_type", "service_condition", "work_content", "updated_at"}),
	}).Create(&positions).Error
}

// FindAllByUser 获取用户的所有岗位，按名称排序
func (r *PositionRepository) FindAllByUser(userID uint) ([]models.Position, error) {
	var positions []models.Position
	err := r.db.Where("user_id = ?", userID).Order("post_name ASC").Find(&positions).Error
	return positions, err
}

// FindByPostIDAndUser 根据 PostID 和 UserID 查询岗位
func (r *PositionRepository) FindByPostIDAndUser(postID string, userID uint) (*models.Position, error) {
	if postID == "" {
		return nil, nil
	}
	var p models.Position
	err := r.db.Where("post_id = ? AND user_id = ?", postID, userID).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// FindByID 根据 ID 查询岗位
func (r *PositionRepository) FindByID(id uint) (*models.Position, error) {
	var p models.Position
	err := r.db.First(&p, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// FindByIDAndUser 根据 ID 和 UserID 查询岗位
func (r *PositionRepository) FindByIDAndUser(id, userID uint) (*models.Position, error) {
	var p models.Position
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// CountByUser 统计用户的岗位数量
func (r *PositionRepository) CountByUser(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Position{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// DeleteByUser 删除用户的所有岗位
func (r *PositionRepository) DeleteByUser(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.Position{}).Error
}
