package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"common/cache"
	"common/database"
	"job-service/models"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type JobHandler struct {
	DB *gorm.DB
}

func NewJobHandler(db *gorm.DB) *JobHandler {
	return &JobHandler{DB: db}
}

const (
	jobDetailCachePrefix = "job:detail"
	jobHotZSetKey        = "job:hot:zset"
)

type HotJob struct {
	models.Job
	ViewCount int64 `json:"view_count"`
}

func buildJobDetailCacheKey(jobID interface{}) string {
	return cache.BuildKey(jobDetailCachePrefix, jobID)
}

func (h *JobHandler) incrementJobView(ctx context.Context, jobID uint) {
	client := database.GetRedis()
	if client == nil {
		return
	}

	member := strconv.FormatUint(uint64(jobID), 10)
	if err := client.ZIncrBy(ctx, jobHotZSetKey, 1, member).Err(); err != nil {
		log.Printf("increment job view failed: job_id=%d err=%v", jobID, err)
	}
}

func (h *JobHandler) invalidateJobDetailCache(ctx context.Context, jobID uint) {
	client := database.GetRedis()
	if client == nil {
		return
	}

	key := buildJobDetailCacheKey(jobID)
	if err := client.Del(ctx, key).Err(); err != nil {
		log.Printf("invalidate job cache failed: job_id=%d err=%v", jobID, err)
	}
}

func (h *JobHandler) removeJobFromHotRanking(ctx context.Context, jobID uint) {
	client := database.GetRedis()
	if client == nil {
		return
	}

	member := strconv.FormatUint(uint64(jobID), 10)
	if err := client.ZRem(ctx, jobHotZSetKey, member).Err(); err != nil {
		log.Printf("remove job from hot ranking failed: job_id=%d err=%v", jobID, err)
	}
}

func (h *JobHandler) respondWithFallbackHotJobs(c *gin.Context) {
	var jobs []models.Job
	if err := h.DB.Where("status = ?", "open").Order("created_at DESC").Limit(10).Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hot jobs"})
		return
	}
	hotJobs := make([]HotJob, 0, len(jobs))
	for _, job := range jobs {
		hotJobs = append(hotJobs, HotJob{
			Job:       job,
			ViewCount: 0,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    hotJobs,
	})
}

// CreateJob 创建职位
func (h *JobHandler) CreateJob(c *gin.Context) {
	var job models.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 优先使用认证上下文中的 user_id；未启用认证时回退到请求体中的 created_by
	if userID, exists := c.Get("user_id"); exists {
		switch uid := userID.(type) {
		case uint:
			job.CreatedBy = uid
		case int:
			if uid > 0 {
				job.CreatedBy = uint(uid)
			}
		case float64:
			if uid > 0 {
				job.CreatedBy = uint(uid)
			}
		}
	}

	// 兼容无鉴权调用：自动绑定一个有效用户，避免 created_by 外键失败
	if job.CreatedBy == 0 {
		var creator struct {
			ID uint `gorm:"column:id"`
		}
		if err := h.DB.Table("users").Select("id").Order("id ASC").Limit(1).First(&creator).Error; err == nil && creator.ID > 0 {
			job.CreatedBy = creator.ID
		}
	}

	if job.CreatedBy == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "created_by is required and no valid user exists"})
		return
	}

	if job.Education == "" {
		job.Education = "不限"
	}
	if job.Headcount <= 0 {
		job.Headcount = 1
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
	education := c.Query("education")
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

	if education != "" {
		query = query.Where("education = ?", education)
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
	allowedSortFields := map[string]bool{"created_at": true, "salary": true, "title": true, "headcount": true}
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

	// 查询每个职位的申请人数
	type JobWithApplicants struct {
		models.Job
		Applicants int64 `json:"applicants"`
	}
	
	type ApplicantCount struct {
		JobID      uint  `json:"job_id"`
		Applicants int64 `json:"applicants"`
	}
	
	jobIDs := make([]uint, 0, len(jobs))
	for _, job := range jobs {
		jobIDs = append(jobIDs, job.ID)
	}
	
	countMap := make(map[uint]int64)
	if len(jobIDs) > 0 {
		var counts []ApplicantCount
		err := h.DB.Table("applications").
			Select("job_id, COUNT(*) as applicants").
			Where("job_id IN ?", jobIDs).
			Where("deleted_at IS NULL").
			Group("job_id").
			Scan(&counts).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applicant counts"})
			return
		}
	
		for _, item := range counts {
			countMap[item.JobID] = item.Applicants
		}
	}
	
	jobsWithApplicants := make([]JobWithApplicants, 0, len(jobs))
	for _, job := range jobs {
		jobsWithApplicants = append(jobsWithApplicants, JobWithApplicants{
			Job:        job,
			Applicants: countMap[job.ID],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"jobs":      jobsWithApplicants,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetJob 获取职位详情
func (h *JobHandler) GetJob(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()
	cacheKey := buildJobDetailCacheKey(id)

	var job models.Job

	if err := cache.GetJSON(ctx, cacheKey, &job); err == nil {
		h.incrementJobView(ctx, job.ID)

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data":    job,
		})
		return
	} else if err != redis.Nil {
		log.Printf("get job cache failed: key=%s err=%v", cacheKey, err)
	}

	if err := h.DB.First(&job, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	if err := cache.Set(ctx, cacheKey, job, cache.DefaultTTL); err != nil {
		log.Printf("set job cache failed: key=%s err=%v", cacheKey, err)
	}

	h.incrementJobView(ctx, job.ID)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    job,
	})
}

func (h *JobHandler) GetHotJobs(c *gin.Context) {
	ctx := c.Request.Context()
	client := database.GetRedis()
	if client == nil {
		h.respondWithFallbackHotJobs(c)
		return
	}

	results, err := client.ZRevRangeWithScores(ctx, jobHotZSetKey, 0, 9).Result()
	if err != nil {
		log.Printf("get hot jobs failed: err=%v", err)
		h.respondWithFallbackHotJobs(c)
		return
	}

	if len(results) == 0 {
		h.respondWithFallbackHotJobs(c)
		return
	}

	ids := make([]uint, 0, len(results))
	scoreMap := make(map[uint]int64, len(results))

	for _, item := range results {
		memberStr, ok := item.Member.(string)
		if !ok {
			continue
		}

		parsedID, err := strconv.ParseUint(memberStr, 10, 64)
		if err != nil {
			continue
		}

		jobID := uint(parsedID)
		ids = append(ids, jobID)
		scoreMap[jobID] = int64(item.Score)
	}

	if len(ids) == 0 {
		h.respondWithFallbackHotJobs(c)
		return
	}

	var jobs []models.Job
	if err := h.DB.Where("id IN ?", ids).Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hot jobs"})
		return
	}

	jobMap := make(map[uint]models.Job, len(jobs))
	for _, job := range jobs {
		jobMap[job.ID] = job
	}

	hotJobs := make([]HotJob, 0, len(ids))
	for _, id := range ids {
		job, ok := jobMap[id]
		if !ok {
			continue
		}

		hotJobs = append(hotJobs, HotJob{
			Job:       job,
			ViewCount: scoreMap[id],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    hotJobs,
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
	h.invalidateJobDetailCache(c.Request.Context(), job.ID)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Job updated successfully",
		"data":    job,
	})
}

// DeleteJob 删除职位
func (h *JobHandler) DeleteJob(c *gin.Context) {
	id := c.Param("id")

	var job models.Job
	if err := h.DB.First(&job, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	if err := h.DB.Delete(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete job"})
		return
	}

	ctx := c.Request.Context()
	h.invalidateJobDetailCache(ctx, job.ID)
	h.removeJobFromHotRanking(ctx, job.ID)

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

// GetJobApplications 获取职位申请列表（HR视角）
// GET /jobs/:id/applications
// 支持状态筛选和分页
// Requirements: 6.1, 6.2
func (h *JobHandler) GetJobApplications(c *gin.Context) {
	// 获取职位ID
	jobID := c.Param("id")

	// 验证职位是否存在
	var job models.Job
	if err := h.DB.First(&job, jobID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "职位不存在",
		})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 构建查询
	query := h.DB.Table("applications").
		Select(`
			applications.id,
			applications.job_id,
			applications.talent_id,
			applications.resume_id,
			applications.status,
			applications.cover_letter,
			applications.notes,
			applications.created_at,
			applications.updated_at,
			talents.name as candidate_name,
			talents.email as candidate_email,
			talents.phone as candidate_phone,
			talents.summary as resume_summary
		`).
		Joins("LEFT JOIN talents ON applications.talent_id = talents.id").
		Where("applications.job_id = ?", jobID).
		Where("applications.deleted_at IS NULL")

	// 状态筛选
	if status != "" {
		query = query.Where("applications.status = ?", status)
	}

	// 获取总数
	var total int64
	countQuery := h.DB.Table("applications").
		Where("job_id = ?", jobID).
		Where("deleted_at IS NULL")
	if status != "" {
		countQuery = countQuery.Where("status = ?", status)
	}
	countQuery.Count(&total)

	// 获取申请列表
	var applications []models.ApplicationWithCandidate
	if err := query.Order("applications.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取申请列表失败",
			"error":   err.Error(),
		})
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
