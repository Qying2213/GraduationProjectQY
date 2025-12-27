package handlers

import (
	"net/http"
	"strconv"
	"talent-service/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TalentHandler struct {
	DB *gorm.DB
}

func NewTalentHandler(db *gorm.DB) *TalentHandler {
	return &TalentHandler{DB: db}
}

// CreateTalent 创建人才
func (h *TalentHandler) CreateTalent(c *gin.Context) {
	var talent models.Talent
	if err := c.ShouldBindJSON(&talent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Create(&talent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create talent: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "Talent created successfully",
		"data":    talent,
	})
}

// ListTalents 获取人才列表
func (h *TalentHandler) ListTalents(c *gin.Context) {
	var talents []models.Talent

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	search := c.Query("search")

	offset := (page - 1) * pageSize

	query := h.DB.Model(&models.Talent{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Find(&talents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch talents"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"talents":   talents,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetTalent 获取单个人才详情
func (h *TalentHandler) GetTalent(c *gin.Context) {
	id := c.Param("id")
	var talent models.Talent

	if err := h.DB.First(&talent, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Talent not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    talent,
	})
}

// UpdateTalent 更新人才信息
func (h *TalentHandler) UpdateTalent(c *gin.Context) {
	id := c.Param("id")
	var talent models.Talent

	if err := h.DB.First(&talent, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Talent not found"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Model(&talent).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update talent: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Talent updated successfully",
		"data":    talent,
	})
}

// DeleteTalent 删除人才
func (h *TalentHandler) DeleteTalent(c *gin.Context) {
	id := c.Param("id")

	if err := h.DB.Delete(&models.Talent{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete talent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Talent deleted successfully",
	})
}

// SearchTalents 搜索人才（高级搜索）
func (h *TalentHandler) SearchTalents(c *gin.Context) {
	var talents []models.Talent

	keyword := c.Query("keyword")
	skills := c.QueryArray("skills")
	minExp, _ := strconv.Atoi(c.DefaultQuery("min_experience", "0"))
	maxExp, _ := strconv.Atoi(c.DefaultQuery("max_experience", "100"))
	education := c.Query("education")
	location := c.Query("location")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	query := h.DB.Model(&models.Talent{})

	// 关键词搜索（搜索姓名、技能、职位等）
	if keyword != "" {
		query = query.Where("name ILIKE ? OR current_position ILIKE ? OR summary ILIKE ? OR ? = ANY(skills)",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", keyword)
	}

	if len(skills) > 0 {
		query = query.Where("skills && ?", skills)
	}

	query = query.Where("experience >= ? AND experience <= ?", minExp, maxExp)

	if education != "" {
		query = query.Where("education = ?", education)
	}

	if location != "" {
		query = query.Where("location ILIKE ?", "%"+location+"%")
	}

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Find(&talents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search talents"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"talents":   talents,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
