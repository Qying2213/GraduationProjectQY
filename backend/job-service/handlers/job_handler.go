package handlers

import (
	"job-service/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JobHandler struct {
	DB *gorm.DB
}

func NewJobHandler(db *gorm.DB) *JobHandler {
	return &JobHandler{DB: db}
}

// CreateJob 创建职位
func (h *JobHandler) CreateJob(c *gin.Context) {
	var job models.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userID, exists := c.Get("user_id"); exists {
		job.CreatedBy = userID.(uint)
	}

	if err := h.DB.Create(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "Job created successfully",
		"data":    job,
	})
}

// ListJobs 获取职位列表
func (h *JobHandler) ListJobs(c *gin.Context) {
	var jobs []models.Job

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	jobType := c.Query("type")
	location := c.Query("location")
	search := c.Query("search")
	keyword := c.Query("keyword")
	level := c.Query("level")
	experience := c.Query("experience")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	offset := (page - 1) * pageSize

	query := h.DB.Model(&models.Job{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if jobType != "" {
		query = query.Where("type = ?", jobType)
	}

	if location != "" {
		query = query.Where("location ILIKE ?", "%"+location+"%")
	}

	// 支持 keyword 搜索（标题、描述、技能）
	if keyword != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 支持级别筛选
	if level != "" {
		query = query.Where("level = ?", level)
	}

	// 支持经验筛选（映射到 level）
	if experience != "" {
		switch experience {
		case "0":
			query = query.Where("level = ?", "junior")
		case "1-3":
			query = query.Where("level IN ?", []string{"junior", "mid"})
		case "3-5":
			query = query.Where("level IN ?", []string{"mid", "senior"})
		case "5-10":
			query = query.Where("level IN ?", []string{"senior", "expert"})
		}
	}

	var total int64
	query.Count(&total)

	// 排序
	allowedSortFields := map[string]bool{"created_at": true, "salary": true, "title": true}
	if !allowedSortFields[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}
	orderClause := sortBy + " " + sortOrder

	if err := query.Order(orderClause).Offset(offset).Limit(pageSize).Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"jobs":      jobs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetJob 获取职位详情
func (h *JobHandler) GetJob(c *gin.Context) {
	id := c.Param("id")
	var job models.Job

	if err := h.DB.First(&job, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    job,
	})
}

// UpdateJob 更新职位
func (h *JobHandler) UpdateJob(c *gin.Context) {
	id := c.Param("id")
	var job models.Job

	if err := h.DB.First(&job, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Save(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update job"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Job updated successfully",
		"data":    job,
	})
}

// DeleteJob 删除职位
func (h *JobHandler) DeleteJob(c *gin.Context) {
	id := c.Param("id")

	if err := h.DB.Delete(&models.Job{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete job"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Job deleted successfully",
	})
}

// GetJobStats 获取职位统计
func (h *JobHandler) GetJobStats(c *gin.Context) {
	var stats struct {
		TotalJobs  int64 `json:"total_jobs"`
		OpenJobs   int64 `json:"open_jobs"`
		ClosedJobs int64 `json:"closed_jobs"`
		FilledJobs int64 `json:"filled_jobs"`
	}

	h.DB.Model(&models.Job{}).Count(&stats.TotalJobs)
	h.DB.Model(&models.Job{}).Where("status = ?", "open").Count(&stats.OpenJobs)
	h.DB.Model(&models.Job{}).Where("status = ?", "closed").Count(&stats.ClosedJobs)
	h.DB.Model(&models.Job{}).Where("status = ?", "filled").Count(&stats.FilledJobs)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    stats,
	})
}
