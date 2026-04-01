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
	"strings"
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

// ProcessTrace 评估链路追踪数据，用于页面展示 OCR/Embedding/RAG/LLM 过程
type ProcessTrace struct {
	GeneratedAt time.Time          `json:"generated_at"`
	OCR         OCRStepTrace       `json:"ocr"`
	Embedding   EmbeddingStepTrace `json:"embedding"`
	RAG         RAGStepTrace       `json:"rag"`
	LLM         LLMStepTrace       `json:"llm"`
}

type OCRStepTrace struct {
	Enabled     bool    `json:"enabled"`
	Success     bool    `json:"success"`
	Pages       int     `json:"pages"`
	Confidence  float64 `json:"confidence"`
	TextLength  int     `json:"text_length"`
	TextPreview string  `json:"text_preview"`
	Error       string  `json:"error"`
}

type EmbeddingStepTrace struct {
	Enabled   bool   `json:"enabled"`
	Success   bool   `json:"success"`
	Model     string `json:"model"`
	Dimension int    `json:"dimension"`
	Error     string `json:"error"`
}

type RAGHit struct {
	Content    string  `json:"content"`
	Similarity float64 `json:"similarity"`
}

type RAGStepTrace struct {
	Enabled bool     `json:"enabled"`
	Success bool     `json:"success"`
	TopK    int      `json:"top_k"`
	Hits    []RAGHit `json:"hits"`
	Error   string   `json:"error"`
}

type LLMStepTrace struct {
	Provider string `json:"provider"`
	Success  bool   `json:"success"`
	Error    string `json:"error"`
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

	handler := &AIEvaluateHandler{
		DB:              db,
		Evaluator:       evaluator.NewCozeEvaluator(),
		EmbeddingClient: embeddingClient,
		RAGEnabled:      ragEnabled,
	}
	handler.ensureAIProcessLogTable()
	handler.ensureEvaluationResultColumns()
	return handler
}

func getEnvDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// ensureAIProcessLogTable 确保 AI 过程日志表存在（避免依赖手动迁移）
func (h *AIEvaluateHandler) ensureAIProcessLogTable() {
	sqls := []string{
		`CREATE TABLE IF NOT EXISTS ai_process_logs (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP,
			evaluation_id INTEGER REFERENCES evaluation_results(id),
			resume_id INTEGER REFERENCES resumes(id) NOT NULL,
			status VARCHAR(20) DEFAULT 'completed',
			process_trace TEXT,
			error_msg TEXT
		)`,
		`CREATE INDEX IF NOT EXISTS idx_ai_process_logs_evaluation_id ON ai_process_logs(evaluation_id)`,
		`CREATE INDEX IF NOT EXISTS idx_ai_process_logs_resume_id ON ai_process_logs(resume_id)`,
		`CREATE INDEX IF NOT EXISTS idx_ai_process_logs_status ON ai_process_logs(status)`,
		`CREATE INDEX IF NOT EXISTS idx_ai_process_logs_created_at ON ai_process_logs(created_at DESC)`,
	}

	for _, sql := range sqls {
		if err := h.DB.Exec(sql).Error; err != nil {
			fmt.Printf("[AIProcessLog] ensure table failed: %v\n", err)
			return
		}
	}
}

func (h *AIEvaluateHandler) ensureEvaluationResultColumns() {
	sqls := []string{
		`ALTER TABLE evaluation_results ADD COLUMN IF NOT EXISTS parsed_school VARCHAR(100)`,
		`ALTER TABLE evaluation_results ADD COLUMN IF NOT EXISTS parsed_report TEXT`,
		`ALTER TABLE evaluation_results ADD COLUMN IF NOT EXISTS raw_result TEXT`,
	}

	for _, sql := range sqls {
		if err := h.DB.Exec(sql).Error; err != nil {
			fmt.Printf("[EvaluationResult] ensure columns failed: %v\n", err)
			return
		}
	}
}

func (h *AIEvaluateHandler) saveProcessLog(evaluationID *uint, resumeID uint, status string, trace *ProcessTrace, errorMsg string) {
	if trace == nil || resumeID == 0 {
		return
	}

	traceJSON, err := json.Marshal(trace)
	if err != nil {
		fmt.Printf("[AIProcessLog] marshal failed: %v\n", err)
		return
	}

	logRow := &models.AIProcessLog{
		EvaluationID: evaluationID,
		ResumeID:     resumeID,
		Status:       status,
		ProcessTrace: string(traceJSON),
		ErrorMsg:     errorMsg,
	}

	if err := h.DB.Create(logRow).Error; err != nil {
		fmt.Printf("[AIProcessLog] save failed: %v\n", err)
	}
}

func tracePreview(text string, maxLen int) string {
	if maxLen <= 0 || len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}

// resolveFilePath 解析简历文件路径，处理绝对路径和相对路径两种情况
// 数据库中可能存储的路径格式：
// 1. 绝对路径: /Users/xxx/Desktop/GraduationProjectQY/backend/resume-service/uploads/xxx.pdf
// 2. 相对路径: /uploads/resumes/xxx.pdf (旧版格式)
// 3. 简单文件名: xxx.pdf
func resolveFilePath(storedPath string) string {
	fmt.Printf("[resolveFilePath] 输入路径: %s\n", storedPath)

	// 如果是绝对路径且文件存在，直接返回
	if filepath.IsAbs(storedPath) {
		if _, err := os.Stat(storedPath); err == nil {
			fmt.Printf("[resolveFilePath] 绝对路径存在: %s\n", storedPath)
			return storedPath
		}
	}

	// 获取当前工作目录
	wd, _ := os.Getwd()
	fmt.Printf("[resolveFilePath] 工作目录: %s\n", wd)

	// 尝试不同的路径组合
	possiblePaths := []string{
		storedPath, // 原始路径
		filepath.Join(wd, "uploads", filepath.Base(storedPath)),            // ./uploads/filename
		filepath.Join(wd, storedPath),                                      // 相对于工作目录
		filepath.Join(wd, "uploads", "resumes", filepath.Base(storedPath)), // ./uploads/resumes/filename
	}

	// 如果原始路径包含 /uploads/，提取文件名尝试
	if strings.Contains(storedPath, "/uploads/") {
		parts := strings.Split(storedPath, "/uploads/")
		if len(parts) > 1 {
			possiblePaths = append(possiblePaths, filepath.Join(wd, "uploads", parts[len(parts)-1]))
		}
	}

	for _, p := range possiblePaths {
		if _, err := os.Stat(p); err == nil {
			fmt.Printf("[resolveFilePath] 找到文件: %s\n", p)
			return p
		}
		fmt.Printf("[resolveFilePath] 路径不存在: %s\n", p)
	}

	// 都找不到，返回原始路径（让后续代码报错）
	fmt.Printf("[resolveFilePath] 所有路径都不存在，返回原始路径\n")
	return storedPath
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

	trace := &ProcessTrace{
		GeneratedAt: time.Now(),
		OCR: OCRStepTrace{
			Enabled: true,
		},
		Embedding: EmbeddingStepTrace{
			Enabled: h.RAGEnabled,
			Model:   h.EmbeddingClient.ModelID,
		},
		RAG: RAGStepTrace{
			Enabled: h.RAGEnabled,
			TopK:    3,
			Hits:    []RAGHit{},
		},
		LLM: LLMStepTrace{
			Provider: "coze",
		},
	}

	// ========== 步骤1: OCR 提取文本 ==========
	fmt.Println("\n========== [AI评估] 步骤1: OCR文本提取 ==========")
	ocrResult, err := ocr.ExtractTextFromFile(resolveFilePath(resume.FilePath))
	var resumeText string
	if err != nil {
		fmt.Printf("[OCR] 提取失败: %v，将使用PDF直接上传\n", err)
		resumeText = ""
		trace.OCR.Success = false
		trace.OCR.Error = err.Error()
	} else {
		resumeText = ocrResult.Text
		trace.OCR.Success = true
		trace.OCR.Pages = ocrResult.Pages
		trace.OCR.Confidence = ocrResult.Confidence
		trace.OCR.TextLength = len(resumeText)
		trace.OCR.TextPreview = tracePreview(resumeText, 500)
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
			trace.Embedding.Success = false
			trace.Embedding.Error = err.Error()
		} else {
			resumeEmbedding = embedding
			trace.Embedding.Success = true
			trace.Embedding.Dimension = len(embedding)
			fmt.Printf("[Embedding] 成功生成向量，维度: %d\n", len(embedding))

			// ========== 步骤3: RAG 检索相似人才 ==========
			fmt.Println("\n========== [AI评估] 步骤3: RAG检索相似人才 ==========")
			ragHitsContext, ragHits, ragErr := h.queryRAG(resumeText, resumeEmbedding)
			ragContext = ragHitsContext
			if ragErr != nil {
				trace.RAG.Success = false
				trace.RAG.Error = ragErr.Error()
			} else {
				trace.RAG.Success = len(ragHits) > 0
				trace.RAG.Hits = ragHits
			}
			if ragContext != "" {
				fmt.Printf("[RAG] 检索到相似人才信息，长度: %d 字符\n", len(ragContext))
			}
		}
	} else {
		fmt.Println("\n========== [AI评估] 步骤2&3: 跳过Embedding和RAG (未配置或无文本) ==========")
		if h.RAGEnabled && resumeText == "" {
			trace.Embedding.Error = "OCR未提取到可用文本"
			trace.RAG.Error = "OCR未提取到可用文本"
		}
	}

	// ========== 步骤4: Coze AI 评估 ==========
	fmt.Println("\n========== [AI评估] 步骤4: Coze AI评估 ==========")

	// 读取简历文件（Coze需要原始PDF）
	pdfBytes, err := os.ReadFile(resolveFilePath(resume.FilePath))
	if err != nil {
		h.DB.Model(&resume).Update("status", "failed")
		trace.LLM.Success = false
		trace.LLM.Error = err.Error()
		h.saveProcessLog(nil, resume.ID, "failed", trace, err.Error())
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "简历文件不存在，请重新上传简历"})
			return
		}
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

	result, err := h.Evaluator.EvaluateResume(ctx, candidateName, enhancedJD, resumeText, pdfBytes)
	if err != nil {
		h.DB.Model(&resume).Update("status", "failed")
		trace.LLM.Success = false
		trace.LLM.Error = err.Error()
		h.saveProcessLog(nil, resume.ID, "failed", trace, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "AI 评估失败: " + err.Error()})
		return
	}
	trace.LLM.Success = true

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
		h.saveResumeEmbedding(&resume, resumeEmbedding)
	}

	// 保存评估结果到 EvaluationResult 表
	evalResult := h.saveEvaluationResult(&resume, result, candidateName, "ai_evaluate")
	h.saveProcessLog(&evalResult.ID, resume.ID, "completed", trace, "")

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
			"parsed_report":    result.ParsedReport,
			"raw_result":       result.RawResult,
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
func (h *AIEvaluateHandler) queryRAG(text string, embedding []float64) (string, []RAGHit, error) {
	_ = embedding

	// 调用 recommendation-service 的 RAG 接口
	ragURL := "http://localhost:8087/internal/recommendations/rag/query"

	reqBody := map[string]interface{}{
		"query": text,
		"top_k": 3,
		"type":  "talent",
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", nil, fmt.Errorf("请求构建失败: %w", err)
	}

	req, err := http.NewRequest("POST", ragURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", nil, fmt.Errorf("请求创建失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if internalAPIKey := os.Getenv("INTERNAL_API_KEY"); internalAPIKey != "" {
		req.Header.Set("X-Internal-Token", internalAPIKey)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return "", nil, fmt.Errorf("响应错误: %d - %s", resp.StatusCode, string(body))
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
		return "", nil, fmt.Errorf("响应解析失败: %w", err)
	}

	if result.Code != 0 || len(result.Data.Results) == 0 {
		return "", []RAGHit{}, nil
	}

	// 组合相似人才信息
	var context string
	hits := make([]RAGHit, 0, len(result.Data.Results))
	for i, r := range result.Data.Results {
		context += fmt.Sprintf("【相似人才%d】(相似度: %.1f%%)\n%s\n\n", i+1, r.Similarity*100, r.Content)
		hits = append(hits, RAGHit{
			Content:    tracePreview(r.Content, 200),
			Similarity: r.Similarity,
		})
	}

	return context, hits, nil
}

// saveResumeEmbedding 将简历对应人才写入向量索引
func (h *AIEvaluateHandler) saveResumeEmbedding(resume *models.Resume, embedding []float64) {
	_ = embedding

	if resume == nil || resume.TalentID == nil || *resume.TalentID == 0 {
		fmt.Printf("[Embedding] 跳过索引：resume_id=%d 未关联 talent_id\n", resume.ID)
		return
	}

	// recommendation-service 当前支持 index-talent/index-job 接口
	indexURL := "http://localhost:8087/internal/recommendations/rag/index-talent"

	reqBody := map[string]interface{}{
		"talent_id": *resume.TalentID,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", indexURL, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	if internalAPIKey := os.Getenv("INTERNAL_API_KEY"); internalAPIKey != "" {
		req.Header.Set("X-Internal-Token", internalAPIKey)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[Embedding] 保存向量失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Printf("[Embedding] 人才向量索引完成: talent_id=%d\n", *resume.TalentID)
		return
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("[Embedding] 索引失败: status=%d, body=%s\n", resp.StatusCode, string(body))
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
			ragResultContext, _, ragErr := h.queryRAG(resumeText, embedding)
			if ragErr != nil {
				fmt.Printf("[RAG] 检索失败: %v\n", ragErr)
			} else {
				ragContext = ragResultContext
			}
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

	result, err := h.Evaluator.EvaluateResume(ctx, candidateName, enhancedJD, resumeText, pdfBytes)
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
		pdfBytes, err := os.ReadFile(resolveFilePath(resume.FilePath))
		if err != nil {
			errors = append(errors, "简历 "+strconv.Itoa(int(resume.ID))+" 文件读取失败")
			continue
		}

		// 调用 AI 评估
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
		result, err := h.Evaluator.EvaluateResume(ctx, resume.FileName, req.JDText, resume.ExtractedText, pdfBytes)
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
	parsedReportJSON, _ := json.Marshal(result.ParsedReport)
	rawResultJSON, _ := json.Marshal(result.RawResult)

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

	basicInfo := extractBasicInfo(result.ParsedReport)
	school := toStringValue(basicInfo["学校"])
	phone := toStringValue(basicInfo["手机"])
	email := toStringValue(basicInfo["邮箱"])
	education := toStringValue(basicInfo["学历"])
	experience := toStringValue(basicInfo["工作经验"])
	location := toStringValue(basicInfo["城市"])
	if location == "" {
		location = toStringValue(basicInfo["地点"])
	}
	if extractedName := toStringValue(basicInfo["姓名"]); extractedName != "" {
		candidateName = extractedName
	}

	evalResult := &models.EvaluationResult{
		ResumeID:             resume.ID,
		ResumeName:           truncateForColumn(resume.FileName, 100),
		ResumeFile:           truncateForColumn(resume.FilePath, 500),
		ParsedName:           truncateForColumn(candidateName, 50),
		ParsedPhone:          truncateForColumn(phone, 20),
		ParsedEmail:          truncateForColumn(email, 100),
		ParsedEducation:      truncateForColumn(education, 50),
		ParsedSchool:         truncateForColumn(school, 100),
		ParsedExperience:     truncateForColumn(experience, 50),
		ParsedLocation:       truncateForColumn(location, 50),
		ParsedSkills:         string(matchedSkillsJSON),
		ParsedReport:         string(parsedReportJSON),
		RawResult:            string(rawResultJSON),
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

	if err := h.DB.Create(evalResult).Error; err != nil {
		fmt.Printf("[EvaluationResult] save failed: %v\n", err)
		return &models.EvaluationResult{}
	}
	return evalResult
}

func extractBasicInfo(report map[string]interface{}) map[string]interface{} {
	if report == nil {
		return map[string]interface{}{}
	}
	if basicInfo, ok := report["基本信息"].(map[string]interface{}); ok {
		return basicInfo
	}
	return map[string]interface{}{}
}

func toStringValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	case fmt.Stringer:
		return strings.TrimSpace(v.String())
	default:
		return ""
	}
}

func truncateForColumn(value string, max int) string {
	if max <= 0 {
		return value
	}
	runes := []rune(strings.TrimSpace(value))
	if len(runes) <= max {
		return string(runes)
	}
	return string(runes[:max])
}
