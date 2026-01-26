package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CozeClient Coze API客户端
type CozeClient struct {
	baseURL    string
	token      string
	workflowID string
	client     *http.Client
}

// NewCozeClient 创建Coze客户端
func NewCozeClient() *CozeClient {
	return &CozeClient{
		baseURL:    getEnv("COZE_BASE_URL", "https://api.coze.cn"),
		token:      os.Getenv("COZE_API_TOKEN"),
		workflowID: os.Getenv("COZE_RESUME_WORKFLOW_ID"),
		client:     &http.Client{Timeout: 300 * time.Second},
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// IsConfigured 检查是否已配置
func (c *CozeClient) IsConfigured() bool {
	return c.token != "" && c.workflowID != ""
}

// ParsedResumeAI AI解析的简历结构
type ParsedResumeAI struct {
	Name        string           `json:"name"`
	Phone       string           `json:"phone"`
	Email       string           `json:"email"`
	Location    string           `json:"location"`
	Education   string           `json:"education"`
	School      string           `json:"school"`
	Major       string           `json:"major"`
	Experience  string           `json:"experience"`
	Skills      []string         `json:"skills"`
	WorkHistory []WorkExperience `json:"work_history"`
	Summary     string           `json:"summary"`
	RiskItems   []RiskItem       `json:"risk_items"`
	RiskScore   int              `json:"risk_score"`
}

// WorkExperience 工作经历
type WorkExperience struct {
	Company   string `json:"company"`
	Position  string `json:"position"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Desc      string `json:"description"`
}

// RiskItem 风险项
type RiskItem struct {
	Type    string `json:"type"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

// ParseResumeWithAI 使用AI解析简历
func (c *CozeClient) ParseResumeWithAI(ctx context.Context, filename string, fileData []byte, jdText string) (*ParsedResumeAI, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("Coze API未配置")
	}

	// 1. 上传文件获取file_id
	fileID, err := c.uploadFile(ctx, filename, fileData)
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}

	// 2. 调用工作流
	result, err := c.runWorkflow(ctx, filename, jdText, fileID)
	if err != nil {
		return nil, fmt.Errorf("调用工作流失败: %w", err)
	}

	return result, nil
}

// uploadFile 上传文件到Coze
func (c *CozeClient) uploadFile(ctx context.Context, filename string, data []byte) (string, error) {
	url := c.baseURL + "/v1/files/upload"

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	// 确保文件名有正确的扩展名
	if !strings.HasSuffix(strings.ToLower(filename), ".pdf") {
		base := strings.TrimSuffix(filename, filepath.Ext(filename))
		filename = base + ".pdf"
	}

	part, err := w.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	if _, err := part.Write(data); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("上传失败 HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Code int `json:"code"`
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.Code != 0 {
		return "", fmt.Errorf("上传失败: %s", result.Msg)
	}

	return result.Data.ID, nil
}

// runWorkflow 运行Coze工作流
func (c *CozeClient) runWorkflow(ctx context.Context, name, jdText, fileID string) (*ParsedResumeAI, error) {
	url := c.baseURL + "/v1/workflow/run"

	reqBody := map[string]interface{}{
		"workflow_id":   c.workflowID,
		"response_mode": "blocking",
		"parameters": map[string]interface{}{
			"name":        name,
			"jd_text":     jdText,
			"resume_file": map[string]interface{}{"file_id": fileID},
		},
	}

	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("工作流调用失败 HTTP %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var result struct {
		Code int `json:"code"`
		Data struct {
			Output string `json:"output"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("工作流执行失败: %s", result.Msg)
	}

	// 解析output中的JSON
	var parsed ParsedResumeAI
	if err := json.Unmarshal([]byte(result.Data.Output), &parsed); err != nil {
		// 尝试直接返回原始输出
		return &ParsedResumeAI{
			Summary: result.Data.Output,
		}, nil
	}

	// 执行风控检查
	parsed.RiskItems = c.checkRisks(&parsed)
	parsed.RiskScore = c.calculateRiskScore(parsed.RiskItems)

	return &parsed, nil
}

// checkRisks 检查风险项
func (c *CozeClient) checkRisks(resume *ParsedResumeAI) []RiskItem {
	var risks []RiskItem

	// 检查工作经历时间冲突
	for i := 0; i < len(resume.WorkHistory)-1; i++ {
		current := resume.WorkHistory[i]
		next := resume.WorkHistory[i+1]

		// 简单的时间冲突检测
		if current.EndDate != "" && next.StartDate != "" {
			if current.EndDate > next.StartDate {
				risks = append(risks, RiskItem{
					Type:    "time_conflict",
					Level:   "warning",
					Message: fmt.Sprintf("工作经历时间可能存在重叠: %s 与 %s", current.Company, next.Company),
				})
			}
		}
	}

	// 检查学历信息完整性
	if resume.Education == "" {
		risks = append(risks, RiskItem{
			Type:    "missing_info",
			Level:   "info",
			Message: "未提供学历信息",
		})
	}

	// 检查联系方式
	if resume.Phone == "" && resume.Email == "" {
		risks = append(risks, RiskItem{
			Type:    "missing_contact",
			Level:   "warning",
			Message: "未提供有效联系方式",
		})
	}

	// 检查工作经历空白期
	for i := 0; i < len(resume.WorkHistory)-1; i++ {
		current := resume.WorkHistory[i]
		next := resume.WorkHistory[i+1]

		if current.EndDate != "" && next.StartDate != "" {
			// 简单检查是否有超过6个月的空白期
			// 这里简化处理，实际应该解析日期
			if len(current.EndDate) >= 7 && len(next.StartDate) >= 7 {
				endYear := current.EndDate[:4]
				startYear := next.StartDate[:4]
				if endYear < startYear {
					risks = append(risks, RiskItem{
						Type:    "gap_period",
						Level:   "info",
						Message: fmt.Sprintf("存在工作空白期: %s 至 %s", current.EndDate, next.StartDate),
					})
				}
			}
		}
	}

	return risks
}

// calculateRiskScore 计算风险分数
func (c *CozeClient) calculateRiskScore(risks []RiskItem) int {
	score := 0
	for _, risk := range risks {
		switch risk.Level {
		case "high":
			score += 30
		case "warning":
			score += 15
		case "info":
			score += 5
		}
	}
	if score > 100 {
		score = 100
	}
	return score
}
