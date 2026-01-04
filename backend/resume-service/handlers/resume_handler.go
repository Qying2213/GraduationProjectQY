package handlers

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"resume-service/models"
	"resume-service/parser"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 文件存储目录
var UploadDir string

type ResumeHandler struct {
	DB     *gorm.DB
	Parser *parser.ResumeParser
}

func NewResumeHandler(db *gorm.DB) *ResumeHandler {
	// 获取当前工作目录
	wd, _ := os.Getwd()
	UploadDir = filepath.Join(wd, "uploads")

	// 确保上传目录存在
	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		log.Printf("Warning: Failed to create upload dir: %v", err)
	}
	log.Printf("Upload directory: %s", UploadDir)

	return &ResumeHandler{
		DB:     db,
		Parser: parser.NewResumeParser(),
	}
}

// UploadResumeFile 上传简历文件
func (h *ResumeHandler) UploadResumeFile(c *gin.Context) {
	log.Println("========== UploadResumeFile START ==========")
	log.Printf("[上传] 请求方法: %s", c.Request.Method)
	log.Printf("[上传] Content-Type: %s", c.GetHeader("Content-Type"))
	log.Printf("[上传] Content-Length: %s", c.GetHeader("Content-Length"))

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("[上传] ❌ FormFile 错误: %v", err)
		log.Printf("[上传] 可能原因: Content-Type 不是 multipart/form-data 或 file 字段不存在")
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "请选择要上传的文件: " + err.Error()})
		return
	}

	log.Printf("[上传] ✓ 文件接收成功: 文件名=%s, 大小=%d bytes", file.Filename, file.Size)

	// 检查文件类型（不区分大小写）
	ext := strings.ToLower(filepath.Ext(file.Filename))
	log.Printf("[上传] 文件扩展名: %s", ext)

	allowedExts := map[string]bool{".pdf": true, ".doc": true, ".docx": true}
	if !allowedExts[ext] {
		log.Printf("[上传] ❌ 文件类型不支持: %s", ext)
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "只支持 PDF、DOC、DOCX 格式，当前格式: " + ext})
		return
	}
	log.Printf("[上传] ✓ 文件类型检查通过")

	// 检查文件大小（最大10MB）
	if file.Size > 10*1024*1024 {
		log.Printf("[上传] ❌ 文件太大: %d bytes", file.Size)
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "文件大小不能超过10MB"})
		return
	}
	log.Printf("[上传] ✓ 文件大小检查通过")

	// 生成唯一文件名
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	filePath := filepath.Join(UploadDir, filename)
	log.Printf("[上传] 目标路径: %s", filePath)
	log.Printf("[上传] UploadDir: %s", UploadDir)

	// 检查目录是否存在
	if _, err := os.Stat(UploadDir); os.IsNotExist(err) {
		log.Printf("[上传] 目录不存在，创建目录: %s", UploadDir)
		if err := os.MkdirAll(UploadDir, 0755); err != nil {
			log.Printf("[上传] ❌ 创建目录失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "创建上传目录失败"})
			return
		}
	}

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Printf("[上传] ❌ 保存文件失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "文件保存失败: " + err.Error()})
		return
	}
	log.Printf("[上传] ✓ 文件保存成功")

	// 验证文件是否真的保存了
	if info, err := os.Stat(filePath); err != nil {
		log.Printf("[上传] ❌ 文件验证失败: %v", err)
	} else {
		log.Printf("[上传] ✓ 文件验证成功: 大小=%d bytes", info.Size())
	}

	// 获取其他表单数据
	talentIDStr := c.PostForm("talent_id")
	jobIDStr := c.PostForm("job_id")
	log.Printf("[上传] 表单数据: talent_id=%s, job_id=%s", talentIDStr, jobIDStr)

	talentID, _ := strconv.Atoi(talentIDStr)
	jobID, _ := strconv.Atoi(jobIDStr)

	// 生成访问URL
	fileURL := "/api/v1/resumes/file/" + filename
	log.Printf("[上传] 文件访问URL: %s", fileURL)

	// 创建简历记录
	resume := models.Resume{
		FilePath: filePath,
		FileName: file.Filename,
		FileURL:  fileURL,
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

	log.Printf("[上传] 准备写入数据库...")
	if err := h.DB.Create(&resume).Error; err != nil {
		log.Printf("[上传] ❌ 数据库写入失败: %v", err)
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "简历记录创建失败: " + err.Error()})
		return
	}

	log.Printf("[上传] ✓ 数据库写入成功, ID=%d", resume.ID)
	log.Println("========== UploadResumeFile SUCCESS ==========")

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "简历上传成功",
		"data":    resume,
	})
}

// ServeResumeFile 提供简历文件访问
func (h *ResumeHandler) ServeResumeFile(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join(UploadDir, filename)

	log.Printf("ServeResumeFile: filename=%s, UploadDir=%s, filePath=%s", filename, UploadDir, filePath)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("File not found: %s", filePath)
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "文件不存在"})
		return
	}

	// 设置响应头
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".pdf":
		c.Header("Content-Type", "application/pdf")
	case ".doc":
		c.Header("Content-Type", "application/msword")
	case ".docx":
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	}

	c.File(filePath)
}

// DownloadResume 下载简历文件
func (h *ResumeHandler) DownloadResume(c *gin.Context) {
	id := c.Param("id")
	var resume models.Resume

	if err := h.DB.First(&resume, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "简历不存在"})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(resume.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "文件不存在"})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", resume.FileName))
	c.File(resume.FilePath)
}

// UploadResume 上传简历（JSON方式）
func (h *ResumeHandler) UploadResume(c *gin.Context) {
	var resume models.Resume
	if err := c.ShouldBindJSON(&resume); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}

	if err := h.DB.Create(&resume).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "Failed to upload resume"})
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
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "Resume not found"})
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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	talentID := c.Query("talent_id")
	status := c.Query("status")
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	offset := (page - 1) * pageSize

	type ResumeWithTalent struct {
		models.Resume
		TalentName string `json:"talent_name"`
	}

	query := h.DB.Model(&models.Resume{})

	if talentID != "" {
		query = query.Where("talent_id = ?", talentID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		query = query.Where("file_name ILIKE ? OR talent_id IN (SELECT id FROM talents WHERE name ILIKE ?)", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	allowedSortFields := map[string]bool{"created_at": true, "status": true, "file_name": true}
	if !allowedSortFields[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}
	orderClause := sortBy + " " + sortOrder

	var resumes []models.Resume
	if err := query.Order(orderClause).Offset(offset).Limit(pageSize).Find(&resumes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "Failed to fetch resumes"})
		return
	}

	result := make([]ResumeWithTalent, len(resumes))
	for i, resume := range resumes {
		result[i].Resume = resume
		if resume.TalentID != nil {
			var talent struct {
				Name string `json:"name"`
			}
			h.DB.Table("talents").Where("id = ?", *resume.TalentID).First(&talent)
			result[i].TalentName = talent.Name
		}
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

// DeleteResume 删除简历
func (h *ResumeHandler) DeleteResume(c *gin.Context) {
	id := c.Param("id")
	var resume models.Resume

	if err := h.DB.First(&resume, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "Resume not found"})
		return
	}

	// 删除文件
	if resume.FilePath != "" {
		os.Remove(resume.FilePath)
	}

	if err := h.DB.Delete(&resume).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "Failed to delete resume"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Resume deleted successfully",
	})
}

// CreateApplication 创建申请
func (h *ResumeHandler) CreateApplication(c *gin.Context) {
	var app models.Application
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}

	if err := h.DB.Create(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "Failed to create application"})
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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	jobID := c.Query("job_id")
	talentID := c.Query("talent_id")
	status := c.Query("status")

	offset := (page - 1) * pageSize

	type ApplicationWithDetails struct {
		models.Application
		TalentName string   `json:"talent_name"`
		JobTitle   string   `json:"job_title"`
		Location   string   `json:"location"`
		Experience int      `json:"experience"`
		Salary     string   `json:"salary"`
		Skills     []string `json:"skills"`
		MatchScore int      `json:"match_score"`
	}

	query := h.DB.Model(&models.Application{})

	if jobID != "" {
		query = query.Where("applications.job_id = ?", jobID)
	}
	if talentID != "" {
		query = query.Where("applications.talent_id = ?", talentID)
	}
	if status != "" {
		query = query.Where("applications.status = ?", status)
	}

	var total int64
	query.Count(&total)

	var applications []models.Application
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "Failed to fetch applications"})
		return
	}

	result := make([]ApplicationWithDetails, len(applications))
	for i, app := range applications {
		result[i].Application = app
		result[i].MatchScore = 75

		var talent struct {
			Name       string `json:"name"`
			Location   string `json:"location"`
			Experience int    `json:"experience"`
			Salary     string `json:"salary"`
		}
		h.DB.Table("talents").Where("id = ?", app.TalentID).First(&talent)
		result[i].TalentName = talent.Name
		result[i].Location = talent.Location
		result[i].Experience = talent.Experience
		result[i].Salary = talent.Salary

		var job struct {
			Title string `json:"title"`
		}
		h.DB.Table("jobs").Where("id = ?", app.JobID).First(&job)
		result[i].JobTitle = job.Title
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"applications": result,
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
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "Application not found"})
		return
	}

	var req struct {
		Status string `json:"status"`
		Notes  string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}

	app.Status = req.Status
	app.Notes = req.Notes

	if err := h.DB.Save(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "Failed to update application"})
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
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "请提供简历文本内容"})
		return
	}

	result, err := h.Parser.Parse(req.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "简历解析失败"})
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
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}

	parsedResume, err := h.Parser.Parse(req.ResumeText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "简历解析失败"})
		return
	}

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
func (h *ResumeHandler) ListResumesForEvaluation(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "100"))
	status := c.DefaultQuery("status", "pending")

	offset := (page - 1) * pageSize

	query := h.DB.Model(&models.Resume{})
	query = query.Where("file_path != '' AND file_path IS NOT NULL")

	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var resumes []models.Resume
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&resumes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "Failed to fetch resumes"})
		return
	}

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

// UpdateResumeStatus 更新简历状态
func (h *ResumeHandler) UpdateResumeStatus(c *gin.Context) {
	id := c.Param("id")
	var resume models.Resume

	if err := h.DB.First(&resume, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "Resume not found"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}

	resume.Status = req.Status
	if err := h.DB.Save(&resume).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "Failed to update resume status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Resume status updated successfully",
		"data":    resume,
	})
}
