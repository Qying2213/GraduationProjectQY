package handlers

import (
	"net/http"
	"resume-service/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ResumeHandler struct {
	DB *gorm.DB
}

func NewResumeHandler(db *gorm.DB) *ResumeHandler {
	return &ResumeHandler{DB: db}
}

// UploadResume 上传简历
func (h *ResumeHandler) UploadResume(c *gin.Context) {
	var resume models.Resume
	if err := c.ShouldBindJSON(&resume); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Create(&resume).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload resume"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "Resume uploaded successfully",
		"data":    resume,
	})
}

// GetResume 获取简历
func (h *ResumeHandler) GetResume(c *gin.Context) {
	id := c.Param("id")
	var resume models.Resume

	if err := h.DB.First(&resume, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resume not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    resume,
	})
}

// ListResumes 获取简历列表
func (h *ResumeHandler) ListResumes(c *gin.Context) {
	var resumes []models.Resume

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	talentID := c.Query("talent_id")

	offset := (page - 1) * pageSize

	query := h.DB.Model(&models.Resume{})

	if talentID != "" {
		query = query.Where("talent_id = ?", talentID)
	}

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Find(&resumes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resumes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"resumes":   resumes,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// DeleteResume 删除简历
func (h *ResumeHandler) DeleteResume(c *gin.Context) {
	id := c.Param("id")

	if err := h.DB.Delete(&models.Resume{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete resume"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Resume deleted successfully",
	})
}

// Application handlers

// CreateApplication 创建申请
func (h *ResumeHandler) CreateApplication(c *gin.Context) {
	var app models.Application
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Create(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create application"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "Application created successfully",
		"data":    app,
	})
}

// ListApplications 获取申请列表
func (h *ResumeHandler) ListApplications(c *gin.Context) {
	var applications []models.Application

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	jobID := c.Query("job_id")
	talentID := c.Query("talent_id")
	status := c.Query("status")

	offset := (page - 1) * pageSize

	query := h.DB.Model(&models.Application{})

	if jobID != "" {
		query = query.Where("job_id = ?", jobID)
	}

	if talentID != "" {
		query = query.Where("talent_id = ?", talentID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"applications": applications,
			"total":        total,
			"page":         page,
			"page_size":    pageSize,
		},
	})
}

// UpdateApplication 更新申请状态
func (h *ResumeHandler) UpdateApplication(c *gin.Context) {
	id := c.Param("id")
	var app models.Application

	if err := h.DB.First(&app, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	var req struct {
		Status string `json:"status"`
		Notes  string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app.Status = req.Status
	app.Notes = req.Notes

	if err := h.DB.Save(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Application updated successfully",
		"data":    app,
	})
}
