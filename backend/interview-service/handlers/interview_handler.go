package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"interview-service/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InterviewHandler 负责面试流程接口。
// 它连接面试安排、候选人通知、面试反馈和统计数据，是招聘流程后半段的核心 handler。
type InterviewHandler struct {
	DB *gorm.DB
}

func NewInterviewHandler(db *gorm.DB) *InterviewHandler {
	return &InterviewHandler{DB: db}
}

// CreateInterview 创建面试安排。
// 它会校验面试时间必须在未来，落库后异步向候选人发送面试通知。
func (h *InterviewHandler) CreateInterview(c *gin.Context) {
	var req struct {
		models.InterviewScheduleRequest
		CreatedBy     uint `json:"created_by"`
		ApplicationID uint `json:"application_id"` // 关联的申请ID
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	// 验证面试日期必须是未来时间，避免创建已经过期的面试。
	interviewDateTime, err := time.ParseInLocation("2006-01-02 15:04", req.Date+" "+req.Time, time.Local)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1005,
			"message": "日期时间格式无效",
		})
		return
	}

	if interviewDateTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1005,
			"message": "面试日期必须是未来时间",
		})
		return
	}

	// 获取创建人ID（优先从请求体获取，其次从context获取）
	var createdBy uint
	if req.CreatedBy > 0 {
		createdBy = req.CreatedBy
	} else if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(uint); ok {
			createdBy = uid
		}
	}

	interview := models.Interview{
		CandidateID:   req.CandidateID,
		CandidateName: req.CandidateName,
		PositionID:    req.PositionID,
		Position:      req.Position,
		Type:          models.InterviewType(req.Type),
		Date:          req.Date,
		Time:          req.Time,
		Duration:      req.Duration,
		InterviewerID: req.InterviewerID,
		Interviewer:   req.Interviewer,
		Method:        models.InterviewMethod(req.Method),
		Location:      req.Location,
		Notes:         req.Notes,
		Status:        models.InterviewStatusScheduled,
		CreatedBy:     createdBy,
	}

	if err := h.DB.Create(&interview).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "Failed to create interview: " + err.Error(),
		})
		return
	}

	// 面试安排创建成功后，异步发送通知给候选人。
	go h.sendInterviewNotification(&interview, req.ApplicationID)

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "Interview scheduled successfully",
		"data":    interview,
	})
}

// sendInterviewNotification 发送面试通知给候选人。
// 通知通过 message-service 的内部接口落库，后续可由站内信或实时消息展示。
func (h *InterviewHandler) sendInterviewNotification(interview *models.Interview, applicationID uint) {
	// 获取候选人的用户ID
	var userID uint

	// 首先尝试从 talents 表获取用户ID
	var talent struct {
		UserID uint `gorm:"column:user_id"`
	}
	if err := h.DB.Table("talents").Select("user_id").Where("id = ?", interview.CandidateID).First(&talent).Error; err == nil && talent.UserID > 0 {
		userID = talent.UserID
	} else {
		// 如果 talents 表没有 user_id，尝试从 applications 表关联查找
		if applicationID > 0 {
			var app struct {
				TalentID uint `gorm:"column:talent_id"`
			}
			if err := h.DB.Table("applications").Select("talent_id").Where("id = ?", applicationID).First(&app).Error; err == nil {
				// 使用 talent_id 作为 user_id（某些系统中可能相同）
				userID = app.TalentID
			}
		}
	}

	// 如果仍然没有找到用户ID，使用候选人ID
	if userID == 0 {
		userID = interview.CandidateID
	}

	// 构建通知消息
	methodText := map[string]string{
		"onsite": "现场面试",
		"video":  "视频面试",
		"phone":  "电话面试",
	}[string(interview.Method)]
	if methodText == "" {
		methodText = "面试"
	}

	typeText := map[string]string{
		"initial": "初试",
		"second":  "复试",
		"final":   "终面",
		"hr":      "HR面试",
	}[string(interview.Type)]
	if typeText == "" {
		typeText = "面试"
	}

	title := fmt.Sprintf("面试邀请 - %s", interview.Position)
	content := fmt.Sprintf(
		"您好，%s！\n\n您已被邀请参加【%s】职位的%s。\n\n面试详情：\n- 面试类型：%s\n- 面试方式：%s\n- 面试时间：%s %s\n- 时长：%d分钟\n- 面试官：%s\n- 地点/链接：%s\n\n请准时参加面试，祝您面试顺利！",
		interview.CandidateName,
		interview.Position,
		typeText,
		typeText,
		methodText,
		interview.Date,
		interview.Time,
		interview.Duration,
		interview.Interviewer,
		interview.Location,
	)

	// 发送通知到 message-service
	messageServiceURL := os.Getenv("MESSAGE_SERVICE_URL")
	if messageServiceURL == "" {
		messageServiceURL = "http://localhost:8085"
	}

	notificationPayload := map[string]interface{}{
		"receiver_id": userID,
		"title":       title,
		"content":     content,
		"type":        "interview",
	}
	h.sendInternalMessage(messageServiceURL, notificationPayload)
}

// ListInterviews 获取面试列表。
// 支持按状态、日期范围、面试官和候选人筛选，用于面试日程页面。
func (h *InterviewHandler) ListInterviews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	date := c.Query("date")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	interviewerID := c.Query("interviewer_id")
	candidateID := c.Query("candidate_id")

	offset := (page - 1) * pageSize

	query := h.DB.Model(&models.Interview{})

	// 过滤条件
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if date != "" {
		query = query.Where("date = ?", date)
	}
	if startDate != "" && endDate != "" {
		query = query.Where("date >= ? AND date <= ?", startDate, endDate)
	}
	if interviewerID != "" {
		query = query.Where("interviewer_id = ?", interviewerID)
	}
	if candidateID != "" {
		query = query.Where("candidate_id = ?", candidateID)
	}

	var total int64
	query.Count(&total)

	var interviews []models.Interview
	if err := query.Order("date ASC, time ASC").Offset(offset).Limit(pageSize).Find(&interviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "Failed to fetch interviews: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": models.InterviewListResponse{
			Interviews: interviews,
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
		},
	})
}

// GetInterview 获取单个面试详情
func (h *InterviewHandler) GetInterview(c *gin.Context) {
	id := c.Param("id")
	var interview models.Interview

	if err := h.DB.First(&interview, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    1,
			"message": "Interview not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    interview,
	})
}

// UpdateInterview 更新面试信息
func (h *InterviewHandler) UpdateInterview(c *gin.Context) {
	id := c.Param("id")
	var interview models.Interview

	if err := h.DB.First(&interview, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    1,
			"message": "Interview not found",
		})
		return
	}

	var req models.InterviewUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Date != "" {
		updates["date"] = req.Date
	}
	if req.Time != "" {
		updates["time"] = req.Time
	}
	if req.Duration > 0 {
		updates["duration"] = req.Duration
	}
	if req.InterviewerID > 0 {
		updates["interviewer_id"] = req.InterviewerID
		updates["interviewer"] = req.Interviewer
	}
	if req.Method != "" {
		updates["method"] = req.Method
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Notes != "" {
		updates["notes"] = req.Notes
	}
	if req.Feedback != "" {
		updates["feedback"] = req.Feedback
	}
	if req.Rating > 0 {
		updates["rating"] = req.Rating
	}

	if err := h.DB.Model(&interview).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "Failed to update interview: " + err.Error(),
		})
		return
	}

	// 重新获取更新后的数据
	if err := h.DB.First(&interview, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "Failed to fetch updated interview: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Interview updated successfully",
		"data":    interview,
	})
}

// DeleteInterview 删除面试
func (h *InterviewHandler) DeleteInterview(c *gin.Context) {
	id := c.Param("id")

	if err := h.DB.Delete(&models.Interview{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "Failed to delete interview: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Interview deleted successfully",
	})
}

// CancelInterview 取消面试
func (h *InterviewHandler) CancelInterview(c *gin.Context) {
	id := c.Param("id")
	var interview models.Interview

	if err := h.DB.First(&interview, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    1,
			"message": "Interview not found",
		})
		return
	}

	if err := h.DB.Model(&interview).Update("status", models.InterviewStatusCancelled).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "Failed to cancel interview: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Interview cancelled successfully",
	})
}

// CompleteInterview 完成面试
func (h *InterviewHandler) CompleteInterview(c *gin.Context) {
	id := c.Param("id")
	var interview models.Interview

	if err := h.DB.First(&interview, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    1,
			"message": "Interview not found",
		})
		return
	}

	// 可选：接收反馈
	var feedback struct {
		Feedback string `json:"feedback"`
		Rating   int    `json:"rating"`
	}
	if err := c.ShouldBindJSON(&feedback); err != nil {
		// 忽略绑定错误，因为反馈是可选的
		feedback.Feedback = ""
		feedback.Rating = 0
	}

	updates := map[string]interface{}{
		"status": models.InterviewStatusCompleted,
	}
	if feedback.Feedback != "" {
		updates["feedback"] = feedback.Feedback
	}
	if feedback.Rating > 0 {
		updates["rating"] = feedback.Rating
	}

	if err := h.DB.Model(&interview).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "Failed to complete interview: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Interview marked as completed",
	})
}

// GetInterviewStats 获取面试统计
func (h *InterviewHandler) GetInterviewStats(c *gin.Context) {
	var stats models.InterviewStats

	// 总面试数
	h.DB.Model(&models.Interview{}).Count(&stats.TotalInterviews)

	// 各状态统计
	h.DB.Model(&models.Interview{}).Where("status = ?", models.InterviewStatusScheduled).Count(&stats.ScheduledInterviews)
	h.DB.Model(&models.Interview{}).Where("status = ?", models.InterviewStatusCompleted).Count(&stats.CompletedInterviews)
	h.DB.Model(&models.Interview{}).Where("status = ?", models.InterviewStatusCancelled).Count(&stats.CancelledInterviews)

	// 今日面试
	today := time.Now().Format("2006-01-02")
	h.DB.Model(&models.Interview{}).Where("date = ?", today).Count(&stats.TodayInterviews)

	// 本周面试
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekStart := now.AddDate(0, 0, -weekday+1).Format("2006-01-02")
	weekEnd := now.AddDate(0, 0, 7-weekday).Format("2006-01-02")
	h.DB.Model(&models.Interview{}).Where("date >= ? AND date <= ?", weekStart, weekEnd).Count(&stats.WeekInterviews)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    stats,
	})
}

// GetTodayInterviews 获取今日面试
func (h *InterviewHandler) GetTodayInterviews(c *gin.Context) {
	today := time.Now().Format("2006-01-02")

	var interviews []models.Interview
	if err := h.DB.Where("date = ? AND status = ?", today, models.InterviewStatusScheduled).
		Order("time ASC").
		Find(&interviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "Failed to fetch today's interviews: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    interviews,
	})
}

// GetInterviewerSchedule 获取面试官日程
func (h *InterviewHandler) GetInterviewerSchedule(c *gin.Context) {
	interviewerID := c.Param("interviewer_id")
	startDate := c.DefaultQuery("start_date", time.Now().Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().AddDate(0, 0, 7).Format("2006-01-02"))

	var interviews []models.Interview
	if err := h.DB.Where("interviewer_id = ? AND date >= ? AND date <= ? AND status = ?",
		interviewerID, startDate, endDate, models.InterviewStatusScheduled).
		Order("date ASC, time ASC").
		Find(&interviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "Failed to fetch interviewer schedule: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    interviews,
	})
}

// SubmitFeedback 提交面试反馈
func (h *InterviewHandler) SubmitFeedback(c *gin.Context) {
	id := c.Param("id")
	var interview models.Interview

	if err := h.DB.First(&interview, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    1,
			"message": "Interview not found",
		})
		return
	}

	var req struct {
		Rating         int    `json:"rating" binding:"required,min=1,max=5"`
		Strengths      string `json:"strengths"`
		Weaknesses     string `json:"weaknesses"`
		Comments       string `json:"comments"`
		Recommendation string `json:"recommendation"` // pass, fail, pending
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	// 获取面试官ID
	var interviewerID uint
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(uint); ok {
			interviewerID = uid
		} else {
			interviewerID = interview.InterviewerID
		}
	} else {
		interviewerID = interview.InterviewerID
	}

	// 创建反馈记录
	feedback := models.InterviewFeedback{
		InterviewID:    interview.ID,
		InterviewerID:  interviewerID,
		Rating:         req.Rating,
		Strengths:      req.Strengths,
		Weaknesses:     req.Weaknesses,
		Comments:       req.Comments,
		Recommendation: req.Recommendation,
	}

	if err := h.DB.Create(&feedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "Failed to submit feedback: " + err.Error(),
		})
		return
	}

	// 更新面试记录
	updates := map[string]interface{}{
		"status":   models.InterviewStatusCompleted,
		"rating":   req.Rating,
		"feedback": req.Comments,
	}
	if err := h.DB.Model(&interview).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "Failed to update interview: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Feedback submitted successfully",
		"data":    feedback,
	})
}

// GetFeedback 获取面试反馈
func (h *InterviewHandler) GetFeedback(c *gin.Context) {
	id := c.Param("id")

	var feedbacks []models.InterviewFeedback
	if err := h.DB.Where("interview_id = ?", id).Find(&feedbacks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "Failed to fetch feedback: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    feedbacks,
	})
}

// GetCandidateInterviews 获取候选人的所有面试
func (h *InterviewHandler) GetCandidateInterviews(c *gin.Context) {
	candidateID := c.Param("candidate_id")

	var interviews []models.Interview
	if err := h.DB.Where("candidate_id = ?", candidateID).
		Order("date DESC, time DESC").
		Find(&interviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "Failed to fetch candidate interviews: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    interviews,
	})
}

// RescheduleInterview 重新安排面试
func (h *InterviewHandler) RescheduleInterview(c *gin.Context) {
	id := c.Param("id")
	var interview models.Interview

	if err := h.DB.First(&interview, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    1,
			"message": "Interview not found",
		})
		return
	}

	var req struct {
		Date   string `json:"date" binding:"required"`
		Time   string `json:"time" binding:"required"`
		Reason string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	// Requirement 7.5: 验证面试日期必须是未来时间
	newDateTime, err := time.ParseInLocation("2006-01-02 15:04", req.Date+" "+req.Time, time.Local)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1005,
			"message": "日期时间格式无效",
		})
		return
	}

	if newDateTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1005,
			"message": "面试日期必须是未来时间",
		})
		return
	}

	oldDate := interview.Date
	oldTime := interview.Time

	updates := map[string]interface{}{
		"date":   req.Date,
		"time":   req.Time,
		"status": models.InterviewStatusScheduled,
		"notes":  interview.Notes + "\n[改期] 从 " + oldDate + " " + oldTime + " 改至 " + req.Date + " " + req.Time + "。原因: " + req.Reason,
	}

	if err := h.DB.Model(&interview).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "Failed to reschedule interview: " + err.Error(),
		})
		return
	}

	h.DB.First(&interview, id)

	// 发送改期通知给候选人
	go h.sendRescheduleNotification(&interview, oldDate, oldTime, req.Reason)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Interview rescheduled successfully",
		"data":    interview,
	})
}

// sendRescheduleNotification 发送面试改期通知
func (h *InterviewHandler) sendRescheduleNotification(interview *models.Interview, oldDate, oldTime, reason string) {
	// 获取候选人的用户ID
	var userID uint
	var talent struct {
		UserID uint `gorm:"column:user_id"`
	}
	if err := h.DB.Table("talents").Select("user_id").Where("id = ?", interview.CandidateID).First(&talent).Error; err == nil && talent.UserID > 0 {
		userID = talent.UserID
	} else {
		userID = interview.CandidateID
	}

	title := fmt.Sprintf("面试改期通知 - %s", interview.Position)
	content := fmt.Sprintf(
		"您好，%s！\n\n您的【%s】职位面试时间已调整。\n\n原时间：%s %s\n新时间：%s %s\n改期原因：%s\n\n请按新时间准时参加面试，祝您面试顺利！",
		interview.CandidateName,
		interview.Position,
		oldDate,
		oldTime,
		interview.Date,
		interview.Time,
		reason,
	)

	messageServiceURL := os.Getenv("MESSAGE_SERVICE_URL")
	if messageServiceURL == "" {
		messageServiceURL = "http://localhost:8085"
	}

	notificationPayload := map[string]interface{}{
		"receiver_id": userID,
		"title":       title,
		"content":     content,
		"type":        "interview",
	}
	h.sendInternalMessage(messageServiceURL, notificationPayload)
}

// sendInternalMessage 发送系统内部消息到 message-service
func (h *InterviewHandler) sendInternalMessage(messageServiceURL string, payload map[string]interface{}) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Failed to marshal notification: %v\n", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, messageServiceURL+"/internal/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Failed to create notification request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// 服务间认证（可选）
	if token := os.Getenv("INTERNAL_API_KEY"); token != "" {
		req.Header.Set("X-Internal-Token", token)
	}

	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send notification: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		fmt.Printf("Notification service returned status: %d\n", resp.StatusCode)
	}
}
