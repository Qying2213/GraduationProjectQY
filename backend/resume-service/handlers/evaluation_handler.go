package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"resume-service/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EvaluationHandler 评估结果处理器
type EvaluationHandler struct {
	DB *gorm.DB
}

// NewEvaluationHandler 创建评估结果处理器
func NewEvaluationHandler(db *gorm.DB) *EvaluationHandler {
	return &EvaluationHandler{DB: db}
}

// ListEvaluations 获取评估结果列表
func (h *EvaluationHandler) ListEvaluations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	evalType := c.Query("eval_type")
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	offset := (page - 1) * pageSize

	query := h.DB.Model(&models.EvaluationResult{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if evalType != "" {
		query = query.Where("eval_type = ?", evalType)
	}

	if search != "" {
		query = query.Where("resume_name ILIKE ? OR parsed_name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	// 排序
	allowedSortFields := map[string]bool{
		"created_at":  true,
		"match_score": true,
		"risk_score":  true,
		"status":      true,
	}
	if !allowedSortFields[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}
	orderClause := sortBy + " " + sortOrder

	var evaluations []models.EvaluationResult
	if err := query.Order(orderClause).Offset(offset).Limit(pageSize).Find(&evaluations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "查询失败"})
		return
	}

	// 转换为前端友好的格式
	results := make([]map[string]interface{}, len(evaluations))
	for i, eval := range evaluations {
		results[i] = h.formatEvaluation(&eval)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"evaluations": results,
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
		},
	})
}

// GetEvaluation 获取单个评估结果详情
func (h *EvaluationHandler) GetEvaluation(c *gin.Context) {
	id := c.Param("id")

	var eval models.EvaluationResult
	if err := h.DB.First(&eval, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "评估结果不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    h.formatEvaluation(&eval),
	})
}

// GetEvaluationProcess 获取评估链路详情（OCR/Embedding/RAG/LLM）
func (h *EvaluationHandler) GetEvaluationProcess(c *gin.Context) {
	id := c.Param("id")

	var processLog models.AIProcessLog
	err := h.DB.Where("evaluation_id = ?", id).
		Order("created_at DESC").
		First(&processLog).Error

	// 兼容：若没有 evaluation_id 关联记录，则按 resume_id 回查最近一条
	if err != nil {
		var eval models.EvaluationResult
		if e := h.DB.Select("resume_id").First(&eval, id).Error; e == nil && eval.ResumeID > 0 {
			err = h.DB.Where("resume_id = ?", eval.ResumeID).
				Order("created_at DESC").
				First(&processLog).Error
		}
	}

	if err != nil {
		log.Printf("[EvaluationProcess] trace not found for evaluation_id=%s", id)
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "该记录暂无链路详情",
			"data": gin.H{
				"id":              nil,
				"evaluation_id":   id,
				"resume_id":       nil,
				"status":          "missing",
				"error_msg":       "",
				"trace":           nil,
				"created_at":      nil,
				"trace_available": false,
			},
		})
		return
	}

	trace := map[string]interface{}{}
	if processLog.ProcessTrace != "" {
		_ = json.Unmarshal([]byte(processLog.ProcessTrace), &trace)
	}

	log.Printf(
		"[EvaluationProcess] trace loaded for evaluation_id=%s process_log_id=%d status=%s",
		id,
		processLog.ID,
		processLog.Status,
	)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id":            processLog.ID,
			"evaluation_id": processLog.EvaluationID,
			"resume_id":     processLog.ResumeID,
			"status":        processLog.Status,
			"error_msg":     processLog.ErrorMsg,
			"trace":         trace,
			"created_at":    processLog.CreatedAt,
			"trace_available": true,
		},
	})
}

// DeleteEvaluation 删除评估结果
func (h *EvaluationHandler) DeleteEvaluation(c *gin.Context) {
	id := c.Param("id")

	if err := h.DB.Delete(&models.EvaluationResult{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// GetEvaluationStats 获取评估统计
func (h *EvaluationHandler) GetEvaluationStats(c *gin.Context) {
	var stats struct {
		Total         int64   `json:"total"`
		Completed     int64   `json:"completed"`
		HighMatch     int64   `json:"high_match"`
		MediumMatch   int64   `json:"medium_match"`
		LowMatch      int64   `json:"low_match"`
		AvgMatchScore float64 `json:"avg_match_score"`
		AvgRiskScore  float64 `json:"avg_risk_score"`
	}

	h.DB.Model(&models.EvaluationResult{}).Count(&stats.Total)
	h.DB.Model(&models.EvaluationResult{}).Where("status = ?", "completed").Count(&stats.Completed)
	h.DB.Model(&models.EvaluationResult{}).Where("match_level = ?", "high").Count(&stats.HighMatch)
	h.DB.Model(&models.EvaluationResult{}).Where("match_level = ?", "medium").Count(&stats.MediumMatch)
	h.DB.Model(&models.EvaluationResult{}).Where("match_level = ?", "low").Count(&stats.LowMatch)

	var avgScores struct {
		AvgMatch float64
		AvgRisk  float64
	}
	h.DB.Model(&models.EvaluationResult{}).
		Select("AVG(match_score) as avg_match, AVG(risk_score) as avg_risk").
		Where("status = ?", "completed").
		Scan(&avgScores)

	stats.AvgMatchScore = avgScores.AvgMatch
	stats.AvgRiskScore = avgScores.AvgRisk

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    stats,
	})
}

// SaveEvaluation 保存评估结果（内部调用）
func (h *EvaluationHandler) SaveEvaluation(eval *models.EvaluationResult) error {
	return h.DB.Create(eval).Error
}

// formatEvaluation 格式化评估结果
func (h *EvaluationHandler) formatEvaluation(eval *models.EvaluationResult) map[string]interface{} {
	result := map[string]interface{}{
		"id":                    eval.ID,
		"created_at":            eval.CreatedAt,
		"resume_id":             eval.ResumeID,
		"talent_id":             eval.TalentID,
		"job_id":                eval.JobID,
		"resume_name":           eval.ResumeName,
		"resume_file":           eval.ResumeFile,
		"parsed_name":           eval.ParsedName,
		"parsed_phone":          eval.ParsedPhone,
		"parsed_email":          eval.ParsedEmail,
		"parsed_education":      eval.ParsedEducation,
		"parsed_experience":     eval.ParsedExperience,
		"parsed_location":       eval.ParsedLocation,
		"match_score":           eval.MatchScore,
		"match_level":           eval.MatchLevel,
		"risk_score":            eval.RiskScore,
		"status":                eval.Status,
		"eval_type":             eval.EvalType,
		"report_summary":        eval.ReportSummary,
		"report_recommendation": eval.ReportRecommendation,
	}

	// 解析JSON字段
	if eval.ParsedSkills != "" {
		var skills []string
		json.Unmarshal([]byte(eval.ParsedSkills), &skills)
		result["parsed_skills"] = skills
	}

	if eval.MatchDetails != "" {
		var details []string
		json.Unmarshal([]byte(eval.MatchDetails), &details)
		result["match_details"] = details
	}

	if eval.ReportStrengths != "" {
		var strengths []string
		json.Unmarshal([]byte(eval.ReportStrengths), &strengths)
		result["report_strengths"] = strengths
	}

	if eval.ReportGaps != "" {
		var gaps []string
		json.Unmarshal([]byte(eval.ReportGaps), &gaps)
		result["report_gaps"] = gaps
	}

	if eval.ReportDimensions != "" {
		var dimensions []map[string]interface{}
		json.Unmarshal([]byte(eval.ReportDimensions), &dimensions)
		result["report_dimensions"] = dimensions
	}

	if eval.RiskItems != "" {
		var risks []map[string]interface{}
		json.Unmarshal([]byte(eval.RiskItems), &risks)
		result["risk_items"] = risks
	}

	return result
}
