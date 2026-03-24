package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"resume-service/ai"
	"resume-service/models"
	"resume-service/ocr"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AIParseHandler AI解析处理器
type AIParseHandler struct {
	DB         *gorm.DB
	CozeClient *ai.CozeClient
}

// NewAIParseHandler 创建AI解析处理器
func NewAIParseHandler(db *gorm.DB) *AIParseHandler {
	return &AIParseHandler{
		DB:         db,
		CozeClient: ai.NewCozeClient(),
	}
}

// AIParseRequest AI解析请求
type AIParseRequest struct {
	ResumeID uint   `json:"resume_id" form:"resume_id"`
	JDText   string `json:"jd_text" form:"jd_text"`
}

type parseSource struct {
	Resume   *models.Resume
	FileName string
	FilePath string
	FileData []byte
	JDText   string
	Cleanup  func()
}

// AIParseResume AI智能解析简历
func (h *AIParseHandler) AIParseResume(c *gin.Context) {
	source, ok := h.resolveParseSource(c)
	if !ok {
		return
	}
	defer cleanupSource(source)

	if source.Resume != nil {
		log.Printf("[AI解析] 开始解析简历: ID=%d, 文件=%s", source.Resume.ID, source.FileName)
	} else {
		log.Printf("[AI解析] 开始解析上传文件: 文件=%s", source.FileName)
	}

	// 1. 使用OCR提取文本
	ocrResult, err := ocr.ExtractTextFromFile(source.FilePath)
	if err != nil {
		log.Printf("[AI解析] OCR提取失败: %v", err)
		// OCR失败不阻断流程，继续尝试AI解析
	} else {
		log.Printf("[AI解析] OCR提取成功: 页数=%d, 文本长度=%d", ocrResult.Pages, len(ocrResult.Text))
	}

	// 2. 调用Coze AI解析
	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()

	var result interface{}

	if h.CozeClient.IsConfigured() {
		log.Printf("[AI解析] 调用Coze AI解析...")
		aiResult, err := h.CozeClient.ParseResumeWithAI(ctx, source.FileName, source.FileData, source.JDText)
		if err != nil {
			log.Printf("[AI解析] Coze解析失败: %v, 降级到本地解析", err)
			// 降级到本地解析
			result = h.localParse(ocrResult)
		} else {
			log.Printf("[AI解析] Coze解析成功")
			result = aiResult
		}
	} else {
		log.Printf("[AI解析] Coze未配置，使用本地解析")
		result = h.localParse(ocrResult)
	}

	// 3. 更新简历状态
	ocrText := ""
	ocrPages := 0
	ocrConfidence := 0.0
	if ocrResult != nil {
		ocrText = ocrResult.Text
		ocrPages = ocrResult.Pages
		ocrConfidence = ocrResult.Confidence
	}

	if source.Resume != nil {
		source.Resume.Status = "parsed"
		if resultJSON, err := json.Marshal(result); err == nil {
			source.Resume.ParsedData = string(resultJSON)
		}
		if ocrText != "" {
			source.Resume.ExtractedText = ocrText
		}
		if err := h.DB.Save(source.Resume).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "保存解析结果失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "解析成功",
		"data": gin.H{
			"resume_id":      source.resumeID(),
			"parsed_result":  result,
			"ocr_text":       truncateText(ocrText, 500),
			"ocr_pages":      ocrPages,
			"ocr_confidence": ocrConfidence,
			"file_name":      source.FileName,
		},
	})
}

// localParse 本地解析（降级方案）
func (h *AIParseHandler) localParse(ocrResult *ocr.OCRResult) map[string]interface{} {
	if ocrResult == nil || ocrResult.Text == "" {
		return map[string]interface{}{
			"error": "无法提取文本内容",
		}
	}

	// 使用正则提取基本信息
	text := ocrResult.Text

	return map[string]interface{}{
		"name":       extractName(text),
		"phone":      extractPhone(text),
		"email":      extractEmail(text),
		"education":  extractEducation(text),
		"skills":     extractSkills(text),
		"experience": extractExperience(text),
		"location":   extractLocation(text),
		"raw_text":   truncateText(text, 2000),
	}
}

// OCRExtract 单独的OCR提取接口
func (h *AIParseHandler) OCRExtract(c *gin.Context) {
	source, ok := h.resolveParseSource(c)
	if !ok {
		return
	}
	defer cleanupSource(source)

	// OCR提取
	result, err := ocr.ExtractTextFromFile(source.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "OCR提取失败: " + err.Error()})
		return
	}

	if source.Resume != nil && result.Text != "" {
		source.Resume.ExtractedText = result.Text
		_ = h.DB.Save(source.Resume).Error
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "提取成功",
		"data": gin.H{
			"resume_id":  source.resumeID(),
			"file_name":  source.FileName,
			"text":       result.Text,
			"pages":      result.Pages,
			"confidence": result.Confidence,
		},
	})
}

func cleanupSource(source *parseSource) {
	if source != nil && source.Cleanup != nil {
		source.Cleanup()
	}
}

func (s *parseSource) resumeID() uint {
	if s != nil && s.Resume != nil {
		return s.Resume.ID
	}
	return 0
}

func (h *AIParseHandler) resolveParseSource(c *gin.Context) (*parseSource, bool) {
	contentType := c.GetHeader("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		return h.resolveMultipartSource(c)
	}
	return h.resolveStoredResumeSource(c)
}

func (h *AIParseHandler) resolveStoredResumeSource(c *gin.Context) (*parseSource, bool) {
	var req AIParseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "参数错误: " + err.Error()})
		return nil, false
	}

	if req.ResumeID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "resume_id 不能为空"})
		return nil, false
	}

	resume, ok := h.loadResume(c, req.ResumeID)
	if !ok {
		return nil, false
	}

	fileData, err := os.ReadFile(resume.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "读取文件失败"})
		return nil, false
	}

	return &parseSource{
		Resume:   resume,
		FileName: resume.FileName,
		FilePath: resume.FilePath,
		FileData: fileData,
		JDText:   req.JDText,
	}, true
}

func (h *AIParseHandler) resolveMultipartSource(c *gin.Context) (*parseSource, bool) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "请选择要解析的文件"})
		return nil, false
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "读取上传文件失败"})
		return nil, false
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "读取上传文件失败"})
		return nil, false
	}

	tempFile, err := os.CreateTemp("", "resume-ai-*"+filepath.Ext(fileHeader.Filename))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "创建临时文件失败"})
		return nil, false
	}

	if _, err := tempFile.Write(fileData); err != nil {
		tempFile.Close()
		_ = os.Remove(tempFile.Name())
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "保存临时文件失败"})
		return nil, false
	}
	if err := tempFile.Close(); err != nil {
		_ = os.Remove(tempFile.Name())
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "保存临时文件失败"})
		return nil, false
	}

	var resume *models.Resume
	if resumeIDStr := c.PostForm("resume_id"); resumeIDStr != "" {
		if resumeID, err := strconv.Atoi(resumeIDStr); err == nil && resumeID > 0 {
			loadedResume, ok := h.loadResume(c, uint(resumeID))
			if !ok {
				_ = os.Remove(tempFile.Name())
				return nil, false
			}
			resume = loadedResume
		}
	}

	return &parseSource{
		Resume:   resume,
		FileName: fileHeader.Filename,
		FilePath: tempFile.Name(),
		FileData: fileData,
		JDText:   c.PostForm("jd_text"),
		Cleanup: func() {
			_ = os.Remove(tempFile.Name())
		},
	}, true
}

func (h *AIParseHandler) loadResume(c *gin.Context, resumeID uint) (*models.Resume, bool) {
	var resume models.Resume
	if err := h.DB.First(&resume, resumeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "简历不存在"})
		return nil, false
	}

	if resume.FilePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "简历文件不存在"})
		return nil, false
	}

	if _, err := os.Stat(resume.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "简历文件已被删除"})
		return nil, false
	}

	return &resume, true
}

// 辅助函数
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}

func extractName(text string) string {
	// 简单的姓名提取逻辑
	lines := splitLines(text)
	for _, line := range lines[:min(5, len(lines))] {
		line = trimSpace(line)
		if len(line) >= 2 && len(line) <= 12 {
			// 检查是否可能是中文姓名
			if isChineseName(line) {
				return line
			}
		}
	}
	return ""
}

func extractPhone(text string) string {
	// 匹配手机号
	for i := 0; i < len(text)-10; i++ {
		if text[i] == '1' && text[i+1] >= '3' && text[i+1] <= '9' {
			phone := text[i : i+11]
			if isAllDigits(phone) {
				return phone
			}
		}
	}
	return ""
}

func extractEmail(text string) string {
	// 简单的邮箱提取
	atIdx := -1
	for i, c := range text {
		if c == '@' {
			atIdx = i
			break
		}
	}
	if atIdx == -1 {
		return ""
	}

	// 向前找用户名
	start := atIdx
	for start > 0 && isEmailChar(text[start-1]) {
		start--
	}

	// 向后找域名
	end := atIdx + 1
	for end < len(text) && isEmailChar(text[end]) {
		end++
	}

	if end > atIdx+1 && start < atIdx {
		return text[start:end]
	}
	return ""
}

func extractEducation(text string) string {
	levels := []string{"博士", "硕士", "研究生", "本科", "学士", "大专", "专科"}
	for _, level := range levels {
		if containsString(text, level) {
			return level
		}
	}
	return ""
}

func extractSkills(text string) []string {
	keywords := []string{
		"Go", "Golang", "Python", "Java", "JavaScript", "TypeScript",
		"Vue", "React", "Docker", "Kubernetes", "MySQL", "PostgreSQL",
		"Redis", "MongoDB", "Linux", "Git", "微服务", "分布式",
	}

	var skills []string
	textLower := toLower(text)
	for _, skill := range keywords {
		if containsString(textLower, toLower(skill)) {
			skills = append(skills, skill)
		}
	}
	return skills
}

func extractExperience(text string) string {
	// 匹配 "X年经验" 模式
	runes := []rune(text)
	for i := 0; i < len(runes)-3; i++ {
		if runes[i] >= '0' && runes[i] <= '9' {
			if i+1 < len(runes) && runes[i+1] == '年' {
				return string(runes[i]) + "年"
			}
			if i+2 < len(runes) && runes[i+1] >= '0' && runes[i+1] <= '9' && runes[i+2] == '年' {
				return string(runes[i:i+2]) + "年"
			}
		}
	}
	return ""
}

func extractLocation(text string) string {
	cities := []string{
		"北京", "上海", "深圳", "广州", "杭州", "成都", "南京", "武汉", "西安", "苏州",
	}
	for _, city := range cities {
		if containsString(text, city) {
			return city
		}
	}
	return ""
}

// 辅助函数
func splitLines(text string) []string {
	var lines []string
	var current []byte
	for i := 0; i < len(text); i++ {
		if text[i] == '\n' {
			lines = append(lines, string(current))
			current = nil
		} else {
			current = append(current, text[i])
		}
	}
	if len(current) > 0 {
		lines = append(lines, string(current))
	}
	return lines
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

func isChineseName(s string) bool {
	// 简单检查是否包含中文字符
	for _, r := range s {
		if r >= 0x4e00 && r <= 0x9fa5 {
			return true
		}
	}
	return false
}

func isAllDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func isEmailChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') || c == '.' || c == '_' || c == '-' || c == '+'
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && findString(s, substr) >= 0
}

func findString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
