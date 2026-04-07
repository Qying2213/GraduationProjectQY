package handlers

import (
	"net/http"
	"strconv"
	"strings"
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

// TalentWithScore 带匹配分数的人才结构
type TalentWithScore struct {
	models.Talent
	MatchScore float64 `json:"match_score"`
}

// ListTalents 获取人才列表
func (h *TalentHandler) ListTalents(c *gin.Context) {
	var talents []models.Talent

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	search := c.Query("search")
	experience := c.Query("experience")
	skills := c.QueryArray("skills")
	location := c.Query("location")

	offset := (page - 1) * pageSize

	query := h.DB.Model(&models.Talent{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 增强关键词搜索：搜索姓名、技能、经验描述
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(
			"name ILIKE ? OR email ILIKE ? OR summary ILIKE ? OR current_position ILIKE ? OR ? = ANY(skills)",
			searchPattern, searchPattern, searchPattern, searchPattern, search,
		)
	}

	// 技能筛选：支持多技能筛选
	if len(skills) > 0 {
		for _, skill := range skills {
			query = query.Where("? = ANY(skills)", skill)
		}
	}

	// 地区筛选
	if location != "" {
		query = query.Where("location ILIKE ?", "%"+location+"%")
	}

	// 经验筛选
	if experience != "" {
		switch experience {
		case "0":
			query = query.Where("experience = 0")
		case "1-3":
			query = query.Where("experience >= 1 AND experience <= 3")
		case "3-5":
			query = query.Where("experience >= 3 AND experience <= 5")
		case "5-10":
			query = query.Where("experience >= 5 AND experience <= 10")
		case "10+":
			query = query.Where("experience > 10")
		}
	}

	var total int64
	query.Count(&total)

	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&talents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch talents"})
		return
	}

	// 计算匹配分数
	talentsWithScore := make([]TalentWithScore, len(talents))
	for i, talent := range talents {
		score := h.calculateMatchScore(talent, search, skills)
		if latestEvalScore := h.latestEvaluationScore(talent); latestEvalScore > 0 {
			score = latestEvalScore
		}
		talentsWithScore[i] = TalentWithScore{
			Talent:     talent,
			MatchScore: score,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"talents":   talentsWithScore,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// calculateMatchScore 计算人才匹配分数
func (h *TalentHandler) calculateMatchScore(talent models.Talent, keyword string, filterSkills []string) float64 {
	score := 0.0
	maxScore := 100.0

	// 基础分数：有完整信息的人才得分更高
	if talent.Name != "" {
		score += 10
	}
	if talent.Email != "" {
		score += 5
	}
	if talent.Phone != "" {
		score += 5
	}
	if len(talent.Skills) > 0 {
		score += 15
	}
	if talent.Experience > 0 {
		score += 10
	}
	if talent.Education != "" {
		score += 10
	}
	if talent.Summary != "" {
		score += 10
	}

	// 关键词匹配加分
	if keyword != "" {
		keywordLower := strings.ToLower(keyword)
		// 姓名匹配
		if strings.Contains(strings.ToLower(talent.Name), keywordLower) {
			score += 15
		}
		// 技能匹配
		for _, skill := range talent.Skills {
			if strings.Contains(strings.ToLower(skill), keywordLower) {
				score += 10
				break
			}
		}
		// 简介匹配
		if strings.Contains(strings.ToLower(talent.Summary), keywordLower) {
			score += 5
		}
	}

	// 技能筛选匹配加分
	if len(filterSkills) > 0 {
		matchedSkills := 0
		for _, filterSkill := range filterSkills {
			filterSkillLower := strings.ToLower(filterSkill)
			for _, talentSkill := range talent.Skills {
				if strings.ToLower(talentSkill) == filterSkillLower {
					matchedSkills++
					break
				}
			}
		}
		// 根据匹配的技能数量加分
		if matchedSkills > 0 {
			skillMatchRatio := float64(matchedSkills) / float64(len(filterSkills))
			score += skillMatchRatio * 20
		}
	}

	// 确保分数在0-100之间
	if score > maxScore {
		score = maxScore
	}
	if score < 0 {
		score = 0
	}

	return score
}

func (h *TalentHandler) latestEvaluationScore(talent models.Talent) float64 {
	if talent.ResumeID == nil || *talent.ResumeID == 0 {
		return 0
	}

	var result struct {
		MatchScore float64 `json:"match_score"`
	}
	if err := h.DB.Table("evaluation_results").
		Select("match_score").
		Where("resume_id = ?", *talent.ResumeID).
		Order("created_at DESC").
		Limit(1).
		Scan(&result).Error; err != nil {
		return 0
	}

	return result.MatchScore
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

	// 定义允许更新的字段
	var updateReq struct {
		Name            string   `json:"name"`
		Email           string   `json:"email"`
		Phone           string   `json:"phone"`
		Skills          []string `json:"skills"`
		Experience      int      `json:"experience"`
		Education       string   `json:"education"`
		Location        string   `json:"location"`
		Salary          string   `json:"salary"`
		CurrentCompany  string   `json:"current_company"`
		CurrentPosition string   `json:"current_position"`
		Summary         string   `json:"summary"`
		Status          string   `json:"status"`
	}

	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 只更新允许的字段
	updates := make(map[string]interface{})
	if updateReq.Name != "" {
		updates["name"] = updateReq.Name
	}
	if updateReq.Email != "" {
		updates["email"] = updateReq.Email
	}
	if updateReq.Phone != "" {
		updates["phone"] = updateReq.Phone
	}
	if len(updateReq.Skills) > 0 {
		updates["skills"] = updateReq.Skills
	}
	if updateReq.Experience > 0 {
		updates["experience"] = updateReq.Experience
	}
	if updateReq.Education != "" {
		updates["education"] = updateReq.Education
	}
	if updateReq.Location != "" {
		updates["location"] = updateReq.Location
	}
	if updateReq.Salary != "" {
		updates["salary"] = updateReq.Salary
	}
	if updateReq.CurrentCompany != "" {
		updates["current_company"] = updateReq.CurrentCompany
	}
	if updateReq.CurrentPosition != "" {
		updates["current_position"] = updateReq.CurrentPosition
	}
	if updateReq.Summary != "" {
		updates["summary"] = updateReq.Summary
	}
	if updateReq.Status != "" {
		updates["status"] = updateReq.Status
	}

	if err := h.DB.Model(&talent).Updates(updates).Error; err != nil {
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

// GetTalentStats 获取人才统计
func (h *TalentHandler) GetTalentStats(c *gin.Context) {
	var stats struct {
		TotalTalents  int64 `json:"total_talents"`
		ActiveTalents int64 `json:"active_talents"`
		NewThisMonth  int64 `json:"new_this_month"`
		ByEducation   []struct {
			Education string `json:"education"`
			Count     int64  `json:"count"`
		} `json:"by_education"`
	}

	h.DB.Model(&models.Talent{}).Count(&stats.TotalTalents)
	h.DB.Model(&models.Talent{}).Where("status = ?", "active").Count(&stats.ActiveTalents)
	h.DB.Model(&models.Talent{}).Where("created_at >= date_trunc('month', CURRENT_DATE)").Count(&stats.NewThisMonth)

	// 按学历统计
	h.DB.Model(&models.Talent{}).
		Select("education, count(*) as count").
		Group("education").
		Scan(&stats.ByEducation)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    stats,
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
