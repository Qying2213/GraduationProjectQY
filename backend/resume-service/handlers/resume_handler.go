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
	"github.com/lib/pq"
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

	// 防止路径遍历攻击
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "非法文件名"})
		return
	}

	filePath := filepath.Join(UploadDir, filename)

	// 确保文件路径在上传目录内
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "无效的文件路径"})
		return
	}
	absUploadDir, err := filepath.Abs(UploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "服务器配置错误"})
		return
	}
	if !strings.HasPrefix(absFilePath, absUploadDir) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "非法文件路径"})
		return
	}

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

	type ResumeWithDetails struct {
		models.Resume
		TalentName string `json:"talent_name"`
		JobTitle   string `json:"job_title"`
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

	result := make([]ResumeWithDetails, len(resumes))
	for i, resume := range resumes {
		result[i].Resume = resume
		// 获取人才名称
		if resume.TalentID != nil {
			var talent struct {
				Name string `json:"name"`
			}
			h.DB.Table("talents").Where("id = ?", *resume.TalentID).First(&talent)
			result[i].TalentName = talent.Name
		}
		// 获取职位名称
		if resume.JobID != nil {
			var job struct {
				Title string `json:"title"`
			}
			h.DB.Table("jobs").Where("id = ?", *resume.JobID).First(&job)
			result[i].JobTitle = job.Title
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

	// 删除文件（确保文件路径在上传目录内）
	if resume.FilePath != "" {
		absFilePath, err := filepath.Abs(resume.FilePath)
		if err == nil {
			absUploadDir, err := filepath.Abs(UploadDir)
			if err == nil && strings.HasPrefix(absFilePath, absUploadDir) {
				os.Remove(resume.FilePath)
			}
		}
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
// 增强逻辑：
// - 检查重复申请（同一求职者不能重复投递同一职位）
// - 验证求职者是否有简历
// - 申请创建后自动发送通知给HR
// Requirements: 2.2, 2.3, 2.4, 2.6
func (h *ResumeHandler) CreateApplication(c *gin.Context) {
	var app models.Application
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}

	// 验证必要字段
	if app.JobID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "job_id is required"})
		return
	}
	if app.TalentID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "talent_id is required"})
		return
	}

	// Requirement 2.6: 验证求职者是否有简历
	var resumeCount int64
	h.DB.Table("resumes").Where("talent_id = ?", app.TalentID).Count(&resumeCount)
	if resumeCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1002,
			"message": "请先上传简历",
		})
		return
	}

	// Requirement 2.4: 检查重复申请（同一求职者不能重复投递同一职位）
	var existingApp models.Application
	result := h.DB.Where("job_id = ? AND talent_id = ?", app.JobID, app.TalentID).First(&existingApp)
	if result.Error == nil {
		// 已存在申请记录
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "您已投递过该职位",
		})
		return
	}

	// Requirement 2.2: 确保申请状态为 "pending"
	app.Status = "pending"

	// 创建申请记录
	if err := h.DB.Create(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "Failed to create application"})
		return
	}

	// Requirement 2.3: 申请创建后自动发送通知给HR
	// 获取职位信息，找到创建该职位的HR
	var job struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		CreatedBy uint   `json:"created_by"`
	}
	if err := h.DB.Table("jobs").Where("id = ?", app.JobID).First(&job).Error; err == nil && job.CreatedBy > 0 {
		// 获取求职者信息
		var talent struct {
			Name string `json:"name"`
		}
		h.DB.Table("talents").Where("id = ?", app.TalentID).First(&talent)

		// 创建通知消息给HR
		notification := map[string]interface{}{
			"sender_id":   nil, // 系统消息
			"receiver_id": job.CreatedBy,
			"title":       "新的职位申请",
			"content":     fmt.Sprintf("求职者 %s 投递了您发布的职位「%s」，请及时查看。", talent.Name, job.Title),
			"type":        "application",
			"is_read":     false,
			"created_at":  time.Now(),
		}
		h.DB.Table("messages").Create(notification)
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

	// 处理 talent_id=me 的情况，表示查询当前登录用户的申请
	if talentID == "me" {
		// 从 JWT 中获取当前用户 ID
		jwtUserID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "error": "未授权，请先登录"})
			return
		}
		userID, ok := jwtUserID.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "invalid user_id type"})
			return
		}

		// 根据 user_id 查找对应的 talent_id
		var talent struct {
			ID uint `json:"id"`
		}
		if err := h.DB.Table("talents").Where("user_id = ?", userID).First(&talent).Error; err != nil {
			// 如果没有找到对应的 talent，返回空列表
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "success",
				"data": gin.H{
					"applications": []interface{}{},
					"total":        0,
					"page":         page,
					"page_size":    pageSize,
				},
			})
			return
		}
		talentID = strconv.FormatUint(uint64(talent.ID), 10)
	}

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
			Name       string         `json:"name"`
			Location   string         `json:"location"`
			Experience int            `json:"experience"`
			Salary     string         `json:"salary"`
			Skills     pq.StringArray `gorm:"type:text[]" json:"skills"`
		}
		h.DB.Table("talents").Where("id = ?", app.TalentID).First(&talent)
		result[i].TalentName = talent.Name
		result[i].Location = talent.Location
		result[i].Experience = talent.Experience
		result[i].Salary = talent.Salary
		result[i].Skills = []string(talent.Skills)

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
// 增强逻辑：
// - 状态更新时发送通知给求职者
// - 记录状态变更历史
// Requirements: 3.2, 6.4
func (h *ResumeHandler) UpdateApplication(c *gin.Context) {
	id := c.Param("id")
	var app models.Application

	if err := h.DB.First(&app, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "Application not found"})
		return
	}

	// 保存旧状态用于比较和记录历史
	oldStatus := app.Status

	var req struct {
		Status string `json:"status"`
		Notes  string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}

	// 验证状态值是否有效
	validStatuses := map[string]bool{
		"pending":   true,
		"viewed":    true,
		"interview": true,
		"offer":     true,
		"rejected":  true,
	}
	if req.Status != "" && !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的状态值，有效值为: pending, viewed, interview, offer, rejected"})
		return
	}

	// 记录状态变更历史
	statusChanged := req.Status != "" && req.Status != oldStatus
	if statusChanged {
		// 构建状态变更历史记录
		historyEntry := fmt.Sprintf("[%s] 状态从 %s 变更为 %s", time.Now().Format("2006-01-02 15:04:05"), getStatusDisplayName(oldStatus), getStatusDisplayName(req.Status))
		if req.Notes != "" {
			historyEntry += fmt.Sprintf(" - 备注: %s", req.Notes)
		}

		// 追加到现有notes中
		if app.Notes != "" {
			app.Notes = app.Notes + "\n" + historyEntry
		} else {
			app.Notes = historyEntry
		}

		app.Status = req.Status
	} else if req.Notes != "" {
		// 如果只是更新备注，追加备注
		noteEntry := fmt.Sprintf("[%s] 备注: %s", time.Now().Format("2006-01-02 15:04:05"), req.Notes)
		if app.Notes != "" {
			app.Notes = app.Notes + "\n" + noteEntry
		} else {
			app.Notes = noteEntry
		}
	}

	if err := h.DB.Save(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "Failed to update application"})
		return
	}

	// Requirement 3.2, 6.4: 状态更新时发送通知给求职者
	if statusChanged {
		h.sendApplicationStatusNotification(&app, oldStatus, req.Status)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Application updated successfully",
		"data":    app,
	})
}

// sendApplicationStatusNotification 发送申请状态变更通知给求职者
// Requirements: 3.2, 6.4
func (h *ResumeHandler) sendApplicationStatusNotification(app *models.Application, oldStatus, newStatus string) {
	// 获取求职者信息（包括关联的用户ID）
	var talent struct {
		Name   string `json:"name"`
		UserID *uint  `json:"user_id"`
	}
	if err := h.DB.Table("talents").Where("id = ?", app.TalentID).First(&talent).Error; err != nil {
		log.Printf("[通知] 获取求职者信息失败: %v", err)
		return
	}

	// 如果求职者没有关联用户账号，无法发送通知
	if talent.UserID == nil || *talent.UserID == 0 {
		log.Printf("[通知] 求职者 %s (ID: %d) 没有关联用户账号，跳过通知", talent.Name, app.TalentID)
		return
	}

	// 获取职位信息
	var job struct {
		Title string `json:"title"`
	}
	if err := h.DB.Table("jobs").Where("id = ?", app.JobID).First(&job).Error; err != nil {
		log.Printf("[通知] 获取职位信息失败: %v", err)
		return
	}

	// 构建通知内容
	title := "申请状态更新"
	content := fmt.Sprintf("您投递的职位「%s」申请状态已更新为「%s」。", job.Title, getStatusDisplayName(newStatus))

	// 根据不同状态添加额外提示
	switch newStatus {
	case "viewed":
		content += " HR已查看您的简历，请耐心等待进一步通知。"
	case "interview":
		content += " 恭喜您进入面试环节！请关注后续面试安排通知。"
	case "offer":
		content += " 恭喜您获得录用！请及时查看offer详情。"
	case "rejected":
		content += " 很遗憾本次未能通过，祝您求职顺利！"
	}

	// 创建通知消息
	notification := map[string]interface{}{
		"sender_id":   nil, // 系统消息
		"receiver_id": *talent.UserID,
		"title":       title,
		"content":     content,
		"type":        "application_status",
		"is_read":     false,
		"created_at":  time.Now(),
	}

	if err := h.DB.Table("messages").Create(notification).Error; err != nil {
		log.Printf("[通知] 创建通知消息失败: %v", err)
		return
	}

	log.Printf("[通知] ✓ 已发送申请状态变更通知给用户 %d: %s -> %s", *talent.UserID, oldStatus, newStatus)
}

// getStatusDisplayName 获取状态的中文显示名称
func getStatusDisplayName(status string) string {
	statusNames := map[string]string{
		"pending":   "待处理",
		"viewed":    "已查看",
		"interview": "面试中",
		"offer":     "已录用",
		"rejected":  "已拒绝",
	}
	if name, ok := statusNames[status]; ok {
		return name
	}
	return status
}

// DeleteApplication 删除/撤回申请
// 求职者可以撤回自己的申请
func (h *ResumeHandler) DeleteApplication(c *gin.Context) {
	id := c.Param("id")

	// 获取当前用户ID
	jwtUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "error": "未授权，请先登录"})
		return
	}
	userID, ok := jwtUserID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "invalid user_id type"})
		return
	}

	var app models.Application
	if err := h.DB.First(&app, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "申请不存在"})
		return
	}

	// 验证是否是申请人本人（通过talent关联的user_id）
	var talent struct {
		UserID *uint `json:"user_id"`
	}
	if err := h.DB.Table("talents").Where("id = ?", app.TalentID).First(&talent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取申请人信息失败"})
		return
	}

	// 检查权限：只有申请人本人可以撤回
	if talent.UserID == nil || *talent.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权操作此申请"})
		return
	}

	// 只有待处理和已查看状态可以撤回
	if app.Status != "pending" && app.Status != "viewed" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "当前状态不允许撤回"})
		return
	}

	// 软删除申请
	if err := h.DB.Delete(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "撤回失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "申请已撤回",
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
