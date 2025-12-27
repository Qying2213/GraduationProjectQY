package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"resume-service/evaluator"
	"resume-service/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AIEvaluateHandler AI 评估处理器
type AIEvaluateHandler struct {
	DB        *gorm.DB
	Evaluator *evaluator.CozeEvaluator
}

// NewAIEvaluateHandler 创建 AI 评估处理器
func NewAIEvaluateHandler(db *gorm.DB) *AIEvaluateHandler {
	return &AIEvaluateHandler{
		DB:        db,
		Evaluator: evaluator.NewCozeEvaluator(),
	}
}

// AIEvaluateRequest AI 评估请求
type AIEvaluateRequest struct {
	ResumeID      uint   `json:"resume_id"`      // 简历ID（从数据库获取）
	JDText        string `json:"jd_text"`        // 职位描述
	CandidateName string `json:"candidate_name"` // 候选人姓名
}

// AIEvaluateResponse AI 评估响应
type AIEvaluateResponse struct {
	ResumeID        uint     `json:"resume_id"`
	CandidateName   string   `json:"candidate_name"`
	TotalScore      float64  `json:"total_score"`
	Grade           string   `json:"grade"`
	JDMatchScore    int      `json:"jd_match_score"`
	AgeScore        int      `json:"age_score"`
	ExperienceScore int      `json:"experience_score"`
	EducationScore  int      `json:"education_score"`
	CompanyScore    int      `json:"company_score"`
	TechScore       int      `json:"tech_score"`
	ProjectScore    int      `json:"project_score"`
	Recommendation  string   `json:"recommendation"`
	MatchedSkills   []string `json:"matched_skills"`
	MissingSkills   []string `json:"missing_skills"`
	Summary         string   `json:"summary"`
}

// CheckAIConfig 检查 AI 配置状态
func (h *AIEvaluateHandler) CheckAIConfig(c *gin.Context) {
	configured := h.Evaluator.IsConfigured()
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"configured": configured,
			"provider":   "coze",
		},
	})
}

// EvaluateByResumeID 根据简历ID进行AI评估
func (h *AIEvaluateHandler) EvaluateByResumeID(c *gin.Context) {
	var req AIEvaluateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查 AI 是否配置
	if !h.Evaluator.IsConfigured() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"code":    503,
			"message": "AI 服务未配置，请设置 COZE_TOKEN 和 COZE_WORKFLOW_ID 环境变量",
		})
		return
	}

	// 获取简历信息
	var resume models.Resume
	if err := h.DB.First(&resume, req.ResumeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "简历不存在"})
		return
	}

	// 读取简历文件
	pdfBytes, err := os.ReadFile(resume.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取简历文件失败"})
		return
	}

	// 调用 AI 评估
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
	defer cancel()

	candidateName := req.CandidateName
	if candidateName == "" {
		candidateName = resume.FileName
	}

	result, err := h.Evaluator.EvaluateResume(ctx, candidateName, req.JDText, pdfBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI 评估失败: " + err.Error()})
		return
	}

	// 更新简历的匹配分数
	resume.MatchScore = int(result.TotalScore)
	resume.Status = "evaluated"
	if resultJSON, err := json.Marshal(result); err == nil {
		resume.ParsedData = string(resultJSON)
	}
	h.DB.Save(&resume)

	// 返回评估结果
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "评估成功",
		"data": AIEvaluateResponse{
			ResumeID:        resume.ID,
			CandidateName:   candidateName,
			TotalScore:      result.TotalScore,
			Grade:           result.Grade,
			JDMatchScore:    result.JDMatchScore,
			AgeScore:        result.AgeScore,
			ExperienceScore: result.ExperienceScore,
			EducationScore:  result.EducationScore,
			CompanyScore:    result.CompanyScore,
			TechScore:       result.TechScore,
			ProjectScore:    result.ProjectScore,
			Recommendation:  result.Recommendation,
			MatchedSkills:   result.MatchedSkills,
			MissingSkills:   result.MissingSkills,
			Summary:         result.Summary,
		},
	})
}

// EvaluateUploadedFile 上传文件并进行AI评估
func (h *AIEvaluateHandler) EvaluateUploadedFile(c *gin.Context) {
	// 检查 AI 是否配置
	if !h.Evaluator.IsConfigured() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"code":    503,
			"message": "AI 服务未配置，请设置 COZE_TOKEN 和 COZE_WORKFLOW_ID 环境变量",
		})
		return
	}

	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的文件"})
		return
	}
	defer file.Close()

	// 检查文件类型
	ext := filepath.Ext(header.Filename)
	if ext != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "目前只支持 PDF 格式的简历"})
		return
	}

	// 读取文件内容
	pdfBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取文件失败"})
		return
	}

	// 获取其他参数
	jdText := c.PostForm("jd_text")
	candidateName := c.PostForm("candidate_name")
	if candidateName == "" {
		candidateName = header.Filename
	}

	// 调用 AI 评估
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
	defer cancel()

	result, err := h.Evaluator.EvaluateResume(ctx, candidateName, jdText, pdfBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI 评估失败: " + err.Error()})
		return
	}

	// 返回评估结果
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "评估成功",
		"data": AIEvaluateResponse{
			CandidateName:   candidateName,
			TotalScore:      result.TotalScore,
			Grade:           result.Grade,
			JDMatchScore:    result.JDMatchScore,
			AgeScore:        result.AgeScore,
			ExperienceScore: result.ExperienceScore,
			EducationScore:  result.EducationScore,
			CompanyScore:    result.CompanyScore,
			TechScore:       result.TechScore,
			ProjectScore:    result.ProjectScore,
			Recommendation:  result.Recommendation,
			MatchedSkills:   result.MatchedSkills,
			MissingSkills:   result.MissingSkills,
			Summary:         result.Summary,
		},
	})
}

// BatchEvaluate 批量评估简历
func (h *AIEvaluateHandler) BatchEvaluate(c *gin.Context) {
	var req struct {
		ResumeIDs []uint `json:"resume_ids" binding:"required"`
		JDText    string `json:"jd_text"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查 AI 是否配置
	if !h.Evaluator.IsConfigured() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"code":    503,
			"message": "AI 服务未配置",
		})
		return
	}

	// 获取所有简历
	var resumes []models.Resume
	if err := h.DB.Where("id IN ?", req.ResumeIDs).Find(&resumes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取简历失败"})
		return
	}

	results := make([]AIEvaluateResponse, 0)
	errors := make([]string, 0)

	for _, resume := range resumes {
		// 读取简历文件
		pdfBytes, err := os.ReadFile(resume.FilePath)
		if err != nil {
			errors = append(errors, "简历 "+strconv.Itoa(int(resume.ID))+" 文件读取失败")
			continue
		}

		// 调用 AI 评估
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
		result, err := h.Evaluator.EvaluateResume(ctx, resume.FileName, req.JDText, pdfBytes)
		cancel()

		if err != nil {
			errors = append(errors, "简历 "+strconv.Itoa(int(resume.ID))+" 评估失败: "+err.Error())
			continue
		}

		// 更新简历
		resume.MatchScore = int(result.TotalScore)
		resume.Status = "evaluated"
		if resultJSON, err := json.Marshal(result); err == nil {
			resume.ParsedData = string(resultJSON)
		}
		h.DB.Save(&resume)

		results = append(results, AIEvaluateResponse{
			ResumeID:       resume.ID,
			CandidateName:  resume.FileName,
			TotalScore:     result.TotalScore,
			Grade:          result.Grade,
			JDMatchScore:   result.JDMatchScore,
			Recommendation: result.Recommendation,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "批量评估完成",
		"data": gin.H{
			"results": results,
			"errors":  errors,
			"total":   len(req.ResumeIDs),
			"success": len(results),
			"failed":  len(errors),
		},
	})
}

// GetEvaluationResult 获取评估结果
func (h *AIEvaluateHandler) GetEvaluationResult(c *gin.Context) {
	id := c.Param("id")
	var resume models.Resume

	if err := h.DB.First(&resume, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "简历不存在"})
		return
	}

	if resume.ParsedData == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "该简历尚未进行 AI 评估"})
		return
	}

	var result evaluator.EvaluationResult
	if err := json.Unmarshal([]byte(resume.ParsedData), &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析评估结果失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    result,
	})
}
