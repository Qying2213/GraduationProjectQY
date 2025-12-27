package repository

import (
	"errors"
	"time"

	"evaluator-service/internal/models"

	"gorm.io/gorm"
)

type CandidateRepository struct{ db *gorm.DB }

func NewCandidateRepository(db *gorm.DB) *CandidateRepository { return &CandidateRepository{db: db} }

func (r *CandidateRepository) DB() *gorm.DB { return r.db }

// ==================== 基础 CRUD（带用户隔离）====================

func (r *CandidateRepository) Create(c *models.Candidate) error { return r.db.Create(c).Error }

// CreateWithUser 创建候选人并关联用户
func (r *CandidateRepository) CreateWithUser(c *models.Candidate, userID uint) error {
	c.UserID = userID
	return r.db.Create(c).Error
}

func (r *CandidateRepository) Update(c *models.Candidate) error { return r.db.Save(c).Error }

// UpdateByUser 更新候选人（验证用户所有权）
func (r *CandidateRepository) UpdateByUser(c *models.Candidate, userID uint) error {
	// 先验证所有权
	var existing models.Candidate
	if err := r.db.Where("id = ? AND user_id = ?", c.ID, userID).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("无权访问该资源")
		}
		return err
	}
	return r.db.Save(c).Error
}

func (r *CandidateRepository) Delete(id uint) error {
	return r.db.Delete(&models.Candidate{}, id).Error
}

// DeleteByUser 删除候选人（验证用户所有权）
func (r *CandidateRepository) DeleteByUser(id, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Candidate{})
	if result.RowsAffected == 0 {
		return errors.New("无权访问该资源")
	}
	return result.Error
}

func (r *CandidateRepository) Get(id uint) (*models.Candidate, error) {
	var c models.Candidate
	if err := r.db.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

// GetByUser 根据ID和用户ID查询候选人
func (r *CandidateRepository) GetByUser(id, userID uint) (*models.Candidate, error) {
	var c models.Candidate
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("无权访问该资源")
		}
		return nil, err
	}
	return &c, nil
}

// FindByID 根据ID查询候选人（别名方法，与Get功能相同）
func (r *CandidateRepository) FindByID(id uint) (*models.Candidate, error) {
	return r.Get(id)
}

// FindByIDAndUser 根据ID和用户ID查询候选人
func (r *CandidateRepository) FindByIDAndUser(id, userID uint) (*models.Candidate, error) {
	return r.GetByUser(id, userID)
}

// ==================== 列表查询（带用户隔离）====================

type ListFilter struct {
	Status   string
	Grade    string
	MinScore *float64
	Search   string
}

func (r *CandidateRepository) List(f ListFilter) ([]models.Candidate, error) {
	q := r.db.Model(&models.Candidate{})
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.Grade != "" {
		q = q.Where("grade = ?", f.Grade)
	}
	if f.MinScore != nil {
		q = q.Where("total_score >= ?", *f.MinScore)
	}
	if f.Search != "" {
		q = q.Where("name LIKE ?", "%"+f.Search+"%")
	}
	var list []models.Candidate
	if err := q.Order("created_at desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// ListByUser 按用户ID查询候选人列表
func (r *CandidateRepository) ListByUser(userID uint, f ListFilter) ([]models.Candidate, error) {
	q := r.db.Model(&models.Candidate{}).Where("user_id = ?", userID)
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.Grade != "" {
		q = q.Where("grade = ?", f.Grade)
	}
	if f.MinScore != nil {
		q = q.Where("total_score >= ?", *f.MinScore)
	}
	if f.Search != "" {
		q = q.Where("name LIKE ?", "%"+f.Search+"%")
	}
	var list []models.Candidate
	if err := q.Order("created_at desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *CandidateRepository) All() ([]models.Candidate, error) {
	var list []models.Candidate
	if err := r.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// AllByUser 获取用户的所有候选人
func (r *CandidateRepository) AllByUser(userID uint) ([]models.Candidate, error) {
	var list []models.Candidate
	if err := r.db.Where("user_id = ?", userID).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *CandidateRepository) GetIn(ids []uint) ([]models.Candidate, error) {
	var list []models.Candidate
	if len(ids) == 0 {
		return list, nil
	}
	if err := r.db.Where("id IN ?", ids).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// GetInByUser 获取指定ID列表中属于用户的候选人
func (r *CandidateRepository) GetInByUser(ids []uint, userID uint) ([]models.Candidate, error) {
	var list []models.Candidate
	if len(ids) == 0 {
		return list, nil
	}
	if err := r.db.Where("id IN ? AND user_id = ?", ids, userID).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *CandidateRepository) DeleteAll() (int64, error) {
	res := r.db.Where("1=1").Delete(&models.Candidate{})
	return res.RowsAffected, res.Error
}

// DeleteAllByUser 删除用户的所有候选人
func (r *CandidateRepository) DeleteAllByUser(userID uint) (int64, error) {
	res := r.db.Where("user_id = ?", userID).Delete(&models.Candidate{})
	return res.RowsAffected, res.Error
}

// ==================== 统计查询（带用户隔离）====================

func (r *CandidateRepository) Count() (int64, error) {
	var n int64
	if err := r.db.Model(&models.Candidate{}).Count(&n).Error; err != nil {
		return 0, err
	}
	return n, nil
}

// CountByUser 统计用户的候选人数量
func (r *CandidateRepository) CountByUser(userID uint) (int64, error) {
	var n int64
	if err := r.db.Model(&models.Candidate{}).Where("user_id = ?", userID).Count(&n).Error; err != nil {
		return 0, err
	}
	return n, nil
}

func (r *CandidateRepository) CountByStatus(status string) (int64, error) {
	var n int64
	if err := r.db.Model(&models.Candidate{}).Where("status = ?", status).Count(&n).Error; err != nil {
		return 0, err
	}
	return n, nil
}

// CountByStatusAndUser 按状态和用户统计候选人数量
func (r *CandidateRepository) CountByStatusAndUser(status string, userID uint) (int64, error) {
	var n int64
	if err := r.db.Model(&models.Candidate{}).Where("status = ? AND user_id = ?", status, userID).Count(&n).Error; err != nil {
		return 0, err
	}
	return n, nil
}

func (r *CandidateRepository) RecentCounts(days int) (map[string]int, error) {
	res := map[string]int{}
	var list []models.Candidate
	from := time.Now().AddDate(0, 0, -days)
	if err := r.db.Where("created_at >= ?", from).Find(&list).Error; err != nil {
		return nil, err
	}
	for _, c := range list {
		key := c.CreatedAt.Format("2006-01-02")
		res[key] = res[key] + 1
	}
	return res, nil
}

// RecentCountsByUser 按用户统计最近几天的候选人数量
func (r *CandidateRepository) RecentCountsByUser(days int, userID uint) (map[string]int, error) {
	res := map[string]int{}
	var list []models.Candidate
	from := time.Now().AddDate(0, 0, -days)
	if err := r.db.Where("created_at >= ? AND user_id = ?", from, userID).Find(&list).Error; err != nil {
		return nil, err
	}
	for _, c := range list {
		key := c.CreatedAt.Format("2006-01-02")
		res[key] = res[key] + 1
	}
	return res, nil
}

// ==================== 钉钉相关查询（带用户隔离）====================

// FindUnnotified 查询未通知的候选人（通知次数为0），按分数降序
func (r *CandidateRepository) FindUnnotified(limit int) ([]models.Candidate, error) {
	var candidates []models.Candidate
	err := r.db.Where("notify_count = ?", 0).
		Order("total_score DESC").
		Limit(limit).
		Find(&candidates).Error
	return candidates, err
}

// FindUnnotifiedByUser 查询用户未通知的候选人
func (r *CandidateRepository) FindUnnotifiedByUser(limit int, userID uint) ([]models.Candidate, error) {
	var candidates []models.Candidate
	err := r.db.Where("notify_count = ? AND user_id = ?", 0, userID).
		Order("total_score DESC").
		Limit(limit).
		Find(&candidates).Error
	return candidates, err
}

// FindRecentlyNotified 查询最近通知的候选人（最近24小时内通知的），按分数降序
func (r *CandidateRepository) FindRecentlyNotified(limit int) ([]models.Candidate, error) {
	var candidates []models.Candidate
	yesterday := time.Now().Add(-24 * time.Hour)
	err := r.db.Where("last_notify_at > ?", yesterday).
		Order("total_score DESC").
		Limit(limit).
		Find(&candidates).Error
	return candidates, err
}

// FindRecentlyNotifiedByUser 查询用户最近通知的候选人
func (r *CandidateRepository) FindRecentlyNotifiedByUser(limit int, userID uint) ([]models.Candidate, error) {
	var candidates []models.Candidate
	yesterday := time.Now().Add(-24 * time.Hour)
	err := r.db.Where("last_notify_at > ? AND user_id = ?", yesterday, userID).
		Order("total_score DESC").
		Limit(limit).
		Find(&candidates).Error
	return candidates, err
}

// ==================== ApplyID 相关查询（重复检测）====================

// FindByApplyIDsAndUser 根据 ApplyID 列表和用户ID查询已存在的候选人
func (r *CandidateRepository) FindByApplyIDsAndUser(applyIDs []string, userID uint) ([]models.Candidate, error) {
	var candidates []models.Candidate
	if len(applyIDs) == 0 {
		return candidates, nil
	}
	// 过滤空字符串
	validIDs := make([]string, 0, len(applyIDs))
	for _, id := range applyIDs {
		if id != "" {
			validIDs = append(validIDs, id)
		}
	}
	if len(validIDs) == 0 {
		return candidates, nil
	}
	err := r.db.Where("apply_id IN ? AND user_id = ?", validIDs, userID).Find(&candidates).Error
	return candidates, err
}

// FindByApplyIDAndUser 根据单个 ApplyID 和用户ID查询候选人
func (r *CandidateRepository) FindByApplyIDAndUser(applyID string, userID uint) (*models.Candidate, error) {
	if applyID == "" {
		return nil, nil
	}
	var c models.Candidate
	err := r.db.Where("apply_id = ? AND user_id = ?", applyID, userID).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}
