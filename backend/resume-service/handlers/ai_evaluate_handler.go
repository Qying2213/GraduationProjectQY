package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"resume-service/evaluator"
	"resume-service/models"
	"resume-service/ocr"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AIEvaluateHandler AI 评估处理器
type AIEvaluateHandler struct {
	DB              *gorm.DB
	Evaluator       *evaluator.CozeEvaluator
	EmbeddingClient *EmbeddingClient
	RAGEnabled      bool
}

// EmbeddingClient 嵌入向量客户端
type EmbeddingClient struct {
	Endpoint string
	APIKey   string
	ModelID  string
}

// NewAIEvaluateHandler 创建 AI 评估处理器
func NewAIEvaluateHandler(db *gorm.DB) *AIEvaluateHandler {
	// 初始化 Embedding 客户端
	embeddingClient := &EmbeddingClient{
		Endpoint: getEnvDefault("VOLC_ENDPOINT", "https://ark.cn-beijing.volces.com/api/v3/embeddings/multimodal"),
		APIKey:   os.Getenv("ARK_API_KEY"),
		ModelID:  getEnvDefault("VOLC_MODEL_ID", "doubao-embedding-vision-251215"),
	}

	ragEnabled := embeddingClient.APIKey != ""

	return &AIEvaluateHandler{
		DB:              db,
		Evaluator:       evaluator.NewCozeEvaluator(),
		EmbeddingClient: embeddingClient,
		RAGEnabled:      ragEnabled,
	}
}

func getEnvDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
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

// GetCurrentTask 获取当前正在处理的任务（支持多个）
func (h *AIEvaluateHandler) GetCurrentTask(c *gin.Context) {
	// 查询所有状态为 processing 的简历
	var resumes []struct {
		ID       uint   `json:"id"`
		FileName string `json:"file_name"`
		JobID    *uint  `json:"job_id"`
		Status   string `json:"status"`
	}

	err := h.DB.Table("resumes").Where("status = ?", "processing").Find(&resumes).Error
	if err != nil || len(resumes) == 0 {
		// 没有正在处理的任务
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "no task",
			"data":    nil,
		})
		return
	}

	// 构建任务列表
	tasks := make([]gin.H, 0, len(resumes))
	for _, resume := range resumes {
		// 获取职位信息
		var jobTitle string
		if resume.JobID != nil {
			var job struct {
				Title string `json:"title"`
			}
			h.DB.Table("jobs").Where("id = ?", *resume.JobID).First(&job)
			jobTitle = job.Title
		}

		tasks = append(tasks, gin.H{
			"resumeId":    resume.ID,
			"fileName":    resume.FileName,
			"jobTitle":    jobTitle,
			"currentStep": 4, // AI评估阶段
			"stepDetail":  "正在调用Coze AI进行智能评估...",
		})
	}

	// 返回所有任务
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"tasks": tasks,
			"count": len(tasks),
		},
	})
}

// EvaluateByResumeID 根据简历ID和JD进行AI评估
// 完整流程: PDF → OCR提取文本 → Embedding向量化 → RAG检索相似人才 → Coze AI评估
func (h *AIEvaluateHandler) EvaluateByResumeID(c *gin.Context) {
	var req struct {
		ResumeID      uint   `json:"resume_id" binding:"required"`
		JDText        string `json:"jd_text"`
		JobID         uint   `json:"job_id"`
		CandidateName string `json:"candidate_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "请提供简历ID"})
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
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "简历不存在"})
		return
	}

	// 更新简历状态为处理中
	h.DB.Model(&resume).Update("status", "processing")

	// 获取JD文本
	jdText := req.JDText
	jobID := req.JobID

	// 如果没有传JD文本，尝试从job_id获取
	if jdText == "" && jobID == 0 && resume.JobID != nil {
		jobID = *resume.JobID
	}

	// 如果有job_id，从数据库获取职位信息生成JD
	if jdText == "" && jobID > 0 {
		var job struct {
			Title        string `json:"title"`
			Department   string `json:"department"`
			Location     string `json:"location"`
			Salary       string `json:"salary"`
			Description  string `json:"description"`
			Requirements string `json:"requirements"`
			Skills       string `json:"skills"`
		}
		if err := h.DB.Table("jobs").Where("id = ?", jobID).First(&job).Error; err == nil {
			// 组合职位信息生成JD
			parts := []string{}
			if job.Title != "" {
				parts = append(parts, "职位名称："+job.Title)
			}
			if job.Department != "" {
				parts = append(parts, "所属部门："+job.Department)
			}
			if job.Location != "" {
				parts = append(parts, "工作地点："+job.Location)
			}
			if job.Salary != "" {
				parts = append(parts, "薪资范围："+job.Salary)
			}
			if job.Description != "" {
				parts = append(parts, "\n岗位职责：\n"+job.Description)
			}
			if job.Requirements != "" {
				parts = append(parts, "\n任职要求：\n"+job.Requirements)
			}
			if job.Skills != "" {
				parts = append(parts, "\n技能要求："+job.Skills)
			}
			jdText = ""
			for _, p := range parts {
				jdText += p + "\n"
			}
		}
	}

	// 如果还是没有JD，返回错误
	if jdText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "请提供职位描述或关联职位"})
		return
	}

	candidateName := req.CandidateName
	if candidateName == "" {
		candidateName = resume.FileName
	}

	// ========== 步骤1: OCR 提取文本 ==========
	fmt.Println("\n========== [AI评估] 步骤1: OCR文本提取 ==========")
	ocrResult, err := ocr.ExtractTextFromFile(resume.FilePath)
	var resumeText string
	if err != nil {
		fmt.Printf("[OCR] 提取失败: %v，将使用PDF直接上传\n", err)
		resumeText = ""
	} else {
		resumeText = ocrResult.Text
		fmt.Printf("[OCR] 成功提取文本，长度: %d 字符，页数: %d，置信度: %.2f\n",
			len(resumeText), ocrResult.Pages, ocrResult.Confidence)
		// 打印前500字符预览
		preview := resumeText
		if len(preview) > 500 {
			preview = preview[:500] + "..."
		}
		fmt.Printf("[OCR] 文本预览:\n%s\n", preview)
	}

	// ========== 步骤2: Embedding 向量化 ==========
	var resumeEmbedding []float64
	var ragContext string
	if h.RAGEnabled && resumeText != "" {
		fmt.Println("\n========== [AI评估] 步骤2: Embedding向量化 ==========")
		embedding, err := h.getEmbedding(resumeText)
		if err != nil {
			fmt.Printf("[Embedding] 向量化失败: %v\n", err)
		} else {
			resumeEmbedding = embedding
			fmt.Printf("[Embedding] 成功生成向量，维度: %d\n", len(embedding))

			// ========== 步骤3: RAG 检索相似人才 ==========
			fmt.Println("\n========== [AI评估] 步骤3: RAG检索相似人才 ==========")
			ragContext = h.queryRAG(resumeText, resumeEmbedding)
			if ragContext != "" {
				fmt.Printf("[RAG] 检索到相似人才信息，长度: %d 字符\n", len(ragContext))
			}
		}
	} else {
		fmt.Println("\n========== [AI评估] 步骤2&3: 跳过Embedding和RAG (未配置或无文本) ==========")
	}

	// ========== 步骤4: Coze AI 评估 ==========
	fmt.Println("\n========== [AI评估] 步骤4: Coze AI评估 ==========")

	// 读取简历文件（Coze需要原始PDF）
	pdfBytes, err := os.ReadFile(resume.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "读取简历文件失败"})
		return
	}

	// 如果有RAG上下文，将其添加到JD中增强评估
	enhancedJD := jdText
	if ragContext != "" {
		enhancedJD = jdText + "\n\n【参考信息-相似人才画像】\n" + ragContext
		fmt.Printf("[Coze] JD已增强RAG上下文，总长度: %d 字符\n", len(enhancedJD))
	}

	// 调用 AI 评估
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
	defer cancel()

	result, err := h.Evaluator.EvaluateResume(ctx, candidateName, enhancedJD, pdfBytes)
	if err != nil {
		h.DB.Model(&resume).Update("status", "failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "AI 评估失败: " + err.Error()})
		return
	}

	fmt.Printf("[Coze] 评估完成! 总分: %.1f, 等级: %s\n", result.TotalScore, result.Grade)

	// 更新简历的匹配分数
	resume.MatchScore = int(result.TotalScore)
	resume.Status = "parsed"
	if resultJSON, err := json.Marshal(result); err == nil {
		resume.ParsedData = string(resultJSON)
	}
	// 保存OCR提取的文本
	if resumeText != "" {
		resume.ExtractedText = resumeText
	}
	h.DB.Save(&resume)

	// 如果有embedding，保存到向量数据库
	if len(resumeEmbedding) > 0 {
		h.saveResumeEmbedding(resume.ID, resumeEmbedding)
	}

	// 保存评估结果到 EvaluationResult 表
	evalResult := h.saveEvaluationResult(&resume, result, candidateName, "ai_evaluate")

	fmt.Println("\n========== [AI评估] 完成! ==========")

	// 返回评估结果
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "评估成功",
		"data": gin.H{
			"evaluation_id":    evalResult.ID,
			"resume_id":        resume.ID,
			"candidate_name":   candidateName,
			"total_score":      result.TotalScore,
			"grade":            result.Grade,
			"jd_match_score":   result.JDMatchScore,
			"age_score":        result.AgeScore,
			"experience_score": result.ExperienceScore,
			"education_score":  result.EducationScore,
			"company_score":    result.CompanyScore,
			"tech_score":       result.TechScore,
			"project_score":    result.ProjectScore,
			"recommendation":   result.Recommendation,
			"matched_skills":   result.MatchedSkills,
			"missing_skills":   result.MissingSkills,
			"summary":          result.Summary,
			"ocr_extracted":    resumeText != "",
			"embedding_used":   len(resumeEmbedding) > 0,
			"rag_enhanced":     ragContext != "",
		},
	})
}

// getEmbedding 调用 Volcengine Doubao 获取文本向量
func (h *AIEvaluateHandler) getEmbedding(text string) ([]float64, error) {
	if h.EmbeddingClient.APIKey == "" {
		return nil, fmt.Errorf("ARK_API_KEY 未配置")
	}

	// 截断文本避免超长
	if len(text) > 8000 {
		text = text[:8000]
	}

	// 构建请求体 - Volcengine multimodal embedding API
	reqBody := map[string]interface{}{
		"model": h.EmbeddingClient.ModelID,
		"input": []map[string]interface{}{
			{
				"type": "text",
				"text": text,
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", h.EmbeddingClient.Endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.EmbeddingClient.APIKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("embedding API error: %d - %s", resp.StatusCode, string(body))
	}

	// 解析响应 - 火山引擎多模态API返回的是单个对象，不是数组
	var result struct {
		Data struct {
			Embedding []float64 `json:"embedding"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w, body: %s", err, string(body))
	}

	if len(result.Data.Embedding) == 0 {
		return nil, fmt.Errorf("empty embedding response: %s", string(body))
	}

	return result.Data.Embedding, nil
}

// queryRAG 查询RAG获取相似人才信息
func (h *AIEvaluateHandler) queryRAG(text string, embedding []float64) string {
	// 调用 recommendation-service 的 RAG 接口
	ragURL := "http://localhost:8087/api/v1/recommendations/rag/query"

	reqBody := map[string]interface{}{
		"query": text,
		"top_k": 3,
		"type":  "talent",
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("[RAG] 请求构建失败: %v\n", err)
		return ""
	}

	req, err := http.NewRequest("POST", ragURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("[RAG] 请求创建失败: %v\n", err)
		return ""
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[RAG] 请求失败: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Printf("[RAG] 响应错误: %d - %s\n", resp.StatusCode, string(body))
		return ""
	}

	// 解析响应
	var result struct {
		Code int `json:"code"`
		Data struct {
			Results []struct {
				Content    string  `json:"content"`
				Similarity float64 `json:"similarity"`
			} `json:"results"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("[RAG] 响应解析失败: %v\n", err)
		return ""
	}

	if result.Code != 0 || len(result.Data.Results) == 0 {
		return ""
	}

	// 组合相似人才信息
	var context string
	for i, r := range result.Data.Results {
		context += fmt.Sprintf("【相似人才%d】(相似度: %.1f%%)\n%s\n\n", i+1, r.Similarity*100, r.Content)
	}

	return context
}

// saveResumeEmbedding 保存简历向量到数据库
func (h *AIEvaluateHandler) saveResumeEmbedding(resumeID uint, embedding []float64) {
	// 调用 recommendation-service 索引接口
	indexURL := "http://localhost:8087/api/v1/rag/index-resume"

	reqBody := map[string]interface{}{
		"resume_id": resumeID,
		"embedding": embedding,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", indexURL, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[Embedding] 保存向量失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Printf("[Embedding] 简历向量已保存到数据库\n")
	}
}

// EvaluateUploadedFile 上传文件并进行AI评估
// 完整流程: PDF → OCR提取文本 → Embedding向量化 → RAG检索 → Coze AI评估
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

	// 保存临时文件用于OCR
	tempFile, err := os.CreateTemp("", "resume_*.pdf")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建临时文件失败"})
		return
	}
	defer os.Remove(tempFile.Name())
	tempFile.Write(pdfBytes)
	tempFile.Close()

	// ========== 步骤1: OCR 提取文本 ==========
	fmt.Println("\n========== [AI评估-上传] 步骤1: OCR文本提取 ==========")
	ocrResult, err := ocr.ExtractTextFromFile(tempFile.Name())
	var resumeText string
	if err != nil {
		fmt.Printf("[OCR] 提取失败: %v\n", err)
	} else {
		resumeText = ocrResult.Text
		fmt.Printf("[OCR] 成功提取文本，长度: %d 字符，页数: %d\n", len(resumeText), ocrResult.Pages)
	}

	// ========== 步骤2: Embedding 向量化 ==========
	var ragContext string
	if h.RAGEnabled && resumeText != "" {
		fmt.Println("\n========== [AI评估-上传] 步骤2: Embedding向量化 ==========")
		embedding, err := h.getEmbedding(resumeText)
		if err != nil {
			fmt.Printf("[Embedding] 向量化失败: %v\n", err)
		} else {
			fmt.Printf("[Embedding] 成功生成向量，维度: %d\n", len(embedding))

			// ========== 步骤3: RAG 检索 ==========
			fmt.Println("\n========== [AI评估-上传] 步骤3: RAG检索 ==========")
			ragContext = h.queryRAG(resumeText, embedding)
		}
	}

	// ========== 步骤4: Coze AI 评估 ==========
	fmt.Println("\n========== [AI评估-上传] 步骤4: Coze AI评估 ==========")

	// 增强JD
	enhancedJD := jdText
	if ragContext != "" {
		enhancedJD = jdText + "\n\n【参考信息-相似人才画像】\n" + ragContext
	}

	// 调用 AI 评估
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
	defer cancel()

	result, err := h.Evaluator.EvaluateResume(ctx, candidateName, enhancedJD, pdfBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI 评估失败: " + err.Error()})
		return
	}

	fmt.Println("\n========== [AI评估-上传] 完成! ==========")

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

// saveEvaluationResult 保存评估结果到数据库
func (h *AIEvaluateHandler) saveEvaluationResult(resume *models.Resume, result *evaluator.EvaluationResult, candidateName string, evalType string) *models.EvaluationResult {
	// 计算加权综合分数
	// JD匹配度权重 30%，6维度综合权重 70%
	// 6维度总分60分，换算成100分制
	dimensionTotal := float64(result.AgeScore + result.ExperienceScore + result.EducationScore +
		result.CompanyScore + result.TechScore + result.ProjectScore)
	dimensionScore := dimensionTotal / 60.0 * 100.0 // 换算成100分制

	// 综合得分 = JD匹配度 × 0.3 + 6维度得分 × 0.7
	weightedScore := float64(result.JDMatchScore)*0.3 + dimensionScore*0.7

	// 如果 Coze 返回了 TotalScore，优先使用它
	finalScore := weightedScore
	if result.TotalScore > 0 {
		finalScore = result.TotalScore
	}

	// 确定匹配等级
	matchLevel := "low"
	if finalScore >= 80 {
		matchLevel = "high"
	} else if finalScore >= 60 {
		matchLevel = "medium"
	}

	// 序列化技能列表
	matchedSkillsJSON, _ := json.Marshal(result.MatchedSkills)
	missingSkillsJSON, _ := json.Marshal(result.MissingSkills)

	// 构建维度数据 - 分数转换成百分比显示
	// JD匹配度已经是100分制，其他维度是10分制需要乘以10
	dimensions := []map[string]interface{}{
		{"name": "JD匹配度", "score": result.JDMatchScore, "max": 100},
		{"name": "年龄适配", "score": result.AgeScore * 10, "max": 100},
		{"name": "工作经验", "score": result.ExperienceScore * 10, "max": 100},
		{"name": "学历背景", "score": result.EducationScore * 10, "max": 100},
		{"name": "公司背景", "score": result.CompanyScore * 10, "max": 100},
		{"name": "技术能力", "score": result.TechScore * 10, "max": 100},
		{"name": "项目经验", "score": result.ProjectScore * 10, "max": 100},
	}
	dimensionsJSON, _ := json.Marshal(dimensions)

	evalResult := &models.EvaluationResult{
		ResumeID:             resume.ID,
		ResumeName:           resume.FileName,
		ResumeFile:           resume.FilePath,
		ParsedName:           candidateName,
		ParsedSkills:         string(matchedSkillsJSON),
		MatchScore:           finalScore,
		MatchLevel:           matchLevel,
		MatchDetails:         string(missingSkillsJSON),
		Status:               "completed",
		EvalType:             evalType,
		ReportSummary:        result.Summary,
		ReportRecommendation: result.Recommendation,
		ReportStrengths:      string(matchedSkillsJSON),
		ReportGaps:           string(missingSkillsJSON),
		ReportDimensions:     string(dimensionsJSON),
	}

	h.DB.Create(evalResult)
	return evalResult
}
