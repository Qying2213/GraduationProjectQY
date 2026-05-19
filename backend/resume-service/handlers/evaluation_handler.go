// 本文件负责对外暴露已持久化的 AI 评估结果和过程追踪。
// 昂贵的 AI 执行逻辑在 ai_evaluate_handler.go 中完成；这里主要服务评估列表、
// 详情页以及前端的 OCR/Embedding/RAG/LLM 流程展示。
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"resume-service/models"
	"strconv"
	"strings"

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

func shouldUseLatestOnly(raw string) bool {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "1", "true", "yes", "y":
		return true
	case "0", "false", "no", "n":
		return false
	default:
		return true
	}
}

// latestEvaluationsQuery 让列表页默认只展示每份简历最新一次评估。
// 同一份简历可能因为调整 JD 或 prompt 被重复评估，但日常看板通常更关心最新结论。
func (h *EvaluationHandler) latestEvaluationsQuery() *gorm.DB {
	subQuery := h.DB.Model(&models.EvaluationResult{}).
		Select("DISTINCT ON (resume_id) *").
		Order("resume_id, created_at DESC, id DESC")

	return h.DB.Table("(?) AS latest_evaluations", subQuery)
}

func applyEvaluationFilters(query *gorm.DB, status, evalType, search string) *gorm.DB {
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if evalType != "" {
		query = query.Where("eval_type = ?", evalType)
	}

	if search != "" {
		query = query.Where("resume_name ILIKE ? OR parsed_name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	return query
}

// ListEvaluations 返回 AI 评估结果页使用的分页列表。
// 它既支持审计式历史记录，也支持 HR 日常使用的“只看最新结果”视图。
func (h *EvaluationHandler) ListEvaluations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	evalType := c.Query("eval_type")
	search := c.Query("search")
	latestOnly := shouldUseLatestOnly(c.DefaultQuery("latest_only", "true"))
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	offset := (page - 1) * pageSize

	var query *gorm.DB
	if latestOnly {
		query = h.latestEvaluationsQuery()
	} else {
		query = h.DB.Model(&models.EvaluationResult{})
	}
	query = applyEvaluationFilters(query, status, evalType, search)

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
	if err := query.Order(orderClause).Offset(offset).Limit(pageSize).Scan(&evaluations).Error; err != nil {
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
			"latest_only": latestOnly,
		},
	})
}

// GetEvaluation 返回单条归一化后的评估记录。
// 格式化步骤会把 PostgreSQL 中存储的 JSON 字符串转换为前端更容易消费的结构。
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

// GetEvaluationProcess 获取评估链路详情（OCR/Embedding/RAG/LLM）。
// 这个接口是毕业设计中的可解释性展示入口：它不仅返回最终分数，还展示 AI 子步骤
// 是否执行、是否成功以及失败原因。
func (h *EvaluationHandler) GetEvaluationProcess(c *gin.Context) {
	id := c.Param("id")

	var processLog models.AIProcessLog
	err := h.DB.Where("evaluation_id = ?", id).
		Order("created_at DESC").
		First(&processLog).Error

	// 兼容旧数据：若没有 evaluation_id 关联记录，则按 resume_id 回查最近一条。
	// 这样数据库迁移前产生的评估记录仍能在前端展示过程信息。
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
			"id":              processLog.ID,
			"evaluation_id":   processLog.EvaluationID,
			"resume_id":       processLog.ResumeID,
			"status":          processLog.Status,
			"error_msg":       processLog.ErrorMsg,
			"trace":           trace,
			"created_at":      processLog.CreatedAt,
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
	status := c.Query("status")
	evalType := c.Query("eval_type")
	search := c.Query("search")
	latestOnly := shouldUseLatestOnly(c.DefaultQuery("latest_only", "true"))

	var stats struct {
		Total         int64   `json:"total"`
		Completed     int64   `json:"completed"`
		HighMatch     int64   `json:"high_match"`
		MediumMatch   int64   `json:"medium_match"`
		LowMatch      int64   `json:"low_match"`
		AvgMatchScore float64 `json:"avg_match_score"`
		AvgRiskScore  float64 `json:"avg_risk_score"`
	}

	buildQuery := func() *gorm.DB {
		var query *gorm.DB
		if latestOnly {
			query = h.latestEvaluationsQuery()
		} else {
			query = h.DB.Model(&models.EvaluationResult{})
		}
		return applyEvaluationFilters(query, status, evalType, search)
	}

	buildQuery().Count(&stats.Total)
	buildQuery().Where("status = ?", "completed").Count(&stats.Completed)
	buildQuery().Where("match_level = ?", "high").Count(&stats.HighMatch)
	buildQuery().Where("match_level = ?", "medium").Count(&stats.MediumMatch)
	buildQuery().Where("match_level = ?", "low").Count(&stats.LowMatch)

	var avgScores struct {
		AvgMatch float64
		AvgRisk  float64
	}
	buildQuery().
		Select("AVG(match_score) as avg_match, AVG(risk_score) as avg_risk").
		Where("status = ?", "completed").
		Scan(&avgScores)

	stats.AvgMatchScore = avgScores.AvgMatch
	stats.AvgRiskScore = avgScores.AvgRisk

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"total":             stats.Total,
			"completed":         stats.Completed,
			"high_match":        stats.HighMatch,
			"medium_match":      stats.MediumMatch,
			"low_match":         stats.LowMatch,
			"avg_match_score":   stats.AvgMatchScore,
			"avg_risk_score":    stats.AvgRiskScore,
			"latest_only":       latestOnly,
			"stats_scope_label": map[bool]string{true: "latest_per_resume", false: "all_records"}[latestOnly],
		},
	})
}

// SaveEvaluation 保存评估结果（内部调用）
func (h *EvaluationHandler) SaveEvaluation(eval *models.EvaluationResult) error {
	return h.DB.Create(eval).Error
}

// formatEvaluation 格式化评估结果
func (h *EvaluationHandler) formatEvaluation(eval *models.EvaluationResult) map[string]interface{} {
	var fallbackEvalData map[string]interface{}
	var fallbackExtractedText string
	if eval.ResumeID > 0 && (eval.ParsedReport == "" || eval.RawResult == "" ||
		eval.ParsedName == "" || eval.ParsedPhone == "" || eval.ParsedEmail == "" ||
		eval.ParsedEducation == "" || eval.ParsedExperience == "" || eval.ParsedLocation == "") {
		var resume models.Resume
		if err := h.DB.Select("parsed_data", "extracted_text").First(&resume, eval.ResumeID).Error; err == nil {
			fallbackExtractedText = resume.ExtractedText
			if resume.ParsedData != "" {
				_ = json.Unmarshal([]byte(resume.ParsedData), &fallbackEvalData)
			}
		}
	}

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
		"parsed_school":         eval.ParsedSchool,
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

	if eval.ParsedReport != "" {
		var parsedReport map[string]interface{}
		if err := json.Unmarshal([]byte(eval.ParsedReport), &parsedReport); err == nil {
			result["parsed_report"] = parsedReport
		}
	} else if fallbackParsedReport, ok := fallbackEvalData["parsed_report"].(map[string]interface{}); ok {
		result["parsed_report"] = fallbackParsedReport
	}

	if eval.RawResult != "" {
		var rawResult map[string]interface{}
		if err := json.Unmarshal([]byte(eval.RawResult), &rawResult); err == nil {
			result["raw_result"] = rawResult
		}
	} else if fallbackRawResult, ok := fallbackEvalData["raw_result"].(map[string]interface{}); ok {
		result["raw_result"] = fallbackRawResult
	}

	parsedReport := asMap(result["parsed_report"])
	parsedName, parsedSchool, parsedPhone, parsedEmail, parsedEducation, parsedExperience, parsedLocation :=
		extractEvaluationBasicFields(parsedReport, fallbackExtractedText, eval.ParsedName)

	if result["parsed_name"] == "" && parsedName != "" {
		result["parsed_name"] = parsedName
	}
	if result["parsed_phone"] == "" && parsedPhone != "" {
		result["parsed_phone"] = parsedPhone
	}
	if result["parsed_email"] == "" && parsedEmail != "" {
		result["parsed_email"] = parsedEmail
	}
	if result["parsed_education"] == "" && parsedEducation != "" {
		result["parsed_education"] = parsedEducation
	}
	if result["parsed_school"] == "" && parsedSchool != "" {
		result["parsed_school"] = parsedSchool
	}
	if result["parsed_experience"] == "" && parsedExperience != "" {
		result["parsed_experience"] = parsedExperience
	}
	if result["parsed_location"] == "" && parsedLocation != "" {
		result["parsed_location"] = parsedLocation
	}

	return result
}

func asMap(value interface{}) map[string]interface{} {
	if m, ok := value.(map[string]interface{}); ok {
		return m
	}
	return map[string]interface{}{}
}
