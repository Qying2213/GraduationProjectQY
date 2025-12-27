package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"resume-service/models"
	"resume-service/parser"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 文件存储目录
const UploadDir = "./uploads/resumes"

type ResumeHandler struct {
	DB     *gorm.DB
	Parser *parser.ResumeParser
}

func NewResumeHandler(db *gorm.DB) *ResumeHandler {
	// 确保上传目录存在
	os.MkdirAll(UploadDir, 0755)
	return &ResumeHandler{
		DB:     db,
		Parser: parser.NewResumeParser(),
	}
}

// UploadResumeFile 上传简历文件
func (h *ResumeHandler) UploadResumeFile(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的文件"})
		return
	}

	// 检查文件类型
	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{".pdf": true, ".doc": true, ".docx": true}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只支持 PDF、DOC、DOCX 格式"})
		return
	}

	// 检查文件大小（最大10MB）
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过10MB"})
		return
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	filePath := filepath.Join(UploadDir, filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败"})
		return
	}

	// 获取其他表单数据
	talentID, _ := strconv.Atoi(c.PostForm("talent_id"))
	jobID, _ := strconv.Atoi(c.PostForm("job_id"))

	// 创建简历记录（talent_id 和 job_id 为 0 时设为 nil）
	resume := models.Resume{
		FilePath: filePath,
		FileName: file.Filename,
		FileSize: file.Size,
		FileType: ext,
		Status:   "pending",
	}
	if talentID > 0 {
		tid := uint(talentID)
		resume.TalentID = &tid
	}
	if jobID > 0 {
		jid := uint(jobID)
		resume.JobID = &jid
	}

	if err := h.DB.Create(&resume).Error; err != nil {
		// 删除已上传的文件
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "简历记录创建失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "简历上传成功",
		"data":    resume,
	})
}

// DownloadResume 下载简历文件
func (h *ResumeHandler) DownloadResume(c *gin.Context) {
	id := c.Param("id")
	var resume models.Resume

	if err := h.DB.First(&resume, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "简历不存在"})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(resume.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", resume.FileName))
	c.File(resume.FilePath)
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
	status := c.Query("status")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	offset := (page - 1) * pageSize

	query := h.DB.Model(&models.Resume{})

	if talentID != "" {
		query = query.Where("talent_id = ?", talentID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	// 构建排序
	allowedSortFields := map[string]bool{"created_at": true, "status": true, "file_name": true}
	if !allowedSortFields[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}
	orderClause := sortBy + " " + sortOrder

	if err := query.Order(orderClause).Offset(offset).Limit(pageSize).Find(&resumes).Error; err != nil {
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

// ParseResume 解析简历文本
func (h *ResumeHandler) ParseResume(c *gin.Context) {
	var req struct {
		Text string `json:"text" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供简历文本内容"})
		return
	}

	// 解析简历
	result, err := h.Parser.Parse(req.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "简历解析失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "解析成功",
		"data":    result,
	})
}

// MatchResumeToJob 计算简历与职位的匹配度
func (h *ResumeHandler) MatchResumeToJob(c *gin.Context) {
	var req struct {
		ResumeText    string   `json:"resume_text" binding:"required"`
		JobSkills     []string `json:"job_skills"`
		JobExperience int      `json:"job_experience"`
		JobEducation  string   `json:"job_education"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 解析简历
	parsedResume, err := h.Parser.Parse(req.ResumeText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "简历解析失败"})
		return
	}

	// 计算匹配度
	score := h.Parser.CalculateMatchScore(parsedResume, req.JobSkills, req.JobExperience, req.JobEducation)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "匹配计算成功",
		"data": gin.H{
			"parsed_resume": parsedResume,
			"match_score":   score,
		},
	})
}

// ListResumesForEvaluation 获取简历列表（用于自动评估系统）
// 返回包含简历文件内容的完整数据
func (h *ResumeHandler) ListResumesForEvaluation(c *gin.Context) {
	var resumes []models.Resume

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "100"))
	status := c.DefaultQuery("status", "pending") // 默认只获取待评估的简历

	offset := (page - 1) * pageSize

	query := h.DB.Model(&models.Resume{})

	// 只获取有文件的简历
	query = query.Where("file_path != '' AND file_path IS NOT NULL")

	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&resumes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resumes"})
		return
	}

	// 构建返回数据，包含文件内容
	type ResumeWithFile struct {
		ID         uint   `json:"id"`
		TalentID   *uint  `json:"talent_id"`
		JobID      *uint  `json:"job_id"`
		FileName   string `json:"file_name"`
		FileType   string `json:"file_type"`
		Status     string `json:"status"`
		HasFile    bool   `json:"has_file"`
		FileBase64 string `json:"file_base64,omitempty"`
	}

	result := make([]ResumeWithFile, 0, len(resumes))
	for _, resume := range resumes {
		item := ResumeWithFile{
			ID:       resume.ID,
			TalentID: resume.TalentID,
			JobID:    resume.JobID,
			FileName: resume.FileName,
			FileType: resume.FileType,
			Status:   resume.Status,
			HasFile:  false,
		}

		// 读取文件内容并转为 base64
		if resume.FilePath != "" {
			if fileBytes, err := os.ReadFile(resume.FilePath); err == nil {
				item.HasFile = true
				item.FileBase64 = base64.StdEncoding.EncodeToString(fileBytes)
			}
		}

		result = append(result, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"resumes":   result,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// UpdateResumeStatus 更新简历状态（用于自动评估系统标记已评估）
func (h *ResumeHandler) UpdateResumeStatus(c *gin.Context) {
	id := c.Param("id")
	var resume models.Resume

	if err := h.DB.First(&resume, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resume not found"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resume.Status = req.Status
	if err := h.DB.Save(&resume).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update resume status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Resume status updated successfully",
		"data":    resume,
	})
}
