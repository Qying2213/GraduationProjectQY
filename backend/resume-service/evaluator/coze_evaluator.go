package evaluator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

// CozeConfig Coze API 配置
type CozeConfig struct {
	BaseURL    string
	Token      string
	WorkflowID string
}

// CozeEvaluator Coze AI 简历评估器
type CozeEvaluator struct {
	config     CozeConfig
	httpClient *http.Client
}

// EvaluationResult AI 评估结果
type EvaluationResult struct {
	Name            string                 `json:"name"`
	TotalScore      float64                `json:"total_score"`
	Grade           string                 `json:"grade"`
	JDMatchScore    int                    `json:"jd_match_score"`
	AgeScore        int                    `json:"age_score"`
	ExperienceScore int                    `json:"experience_score"`
	EducationScore  int                    `json:"education_score"`
	CompanyScore    int                    `json:"company_score"`
	TechScore       int                    `json:"tech_score"`
	ProjectScore    int                    `json:"project_score"`
	Recommendation  string                 `json:"recommendation"`
	MatchedSkills   []string               `json:"matched_skills"`
	MissingSkills   []string               `json:"missing_skills"`
	Summary         string                 `json:"summary"`
	RawResult       map[string]interface{} `json:"raw_result"`
}

// NewCozeEvaluator 创建 Coze 评估器
func NewCozeEvaluator() *CozeEvaluator {
	config := CozeConfig{
		BaseURL:    getEnv("COZE_BASE_URL", "https://api.coze.cn"),
		Token:      getEnv("COZE_TOKEN", ""),
		WorkflowID: getEnv("COZE_WORKFLOW_ID", ""),
	}

	return &CozeEvaluator{
		config:     config,
		httpClient: &http.Client{Timeout: 300 * time.Second},
	}
}

// NewCozeEvaluatorWithConfig 使用指定配置创建评估器
func NewCozeEvaluatorWithConfig(config CozeConfig) *CozeEvaluator {
	return &CozeEvaluator{
		config:     config,
		httpClient: &http.Client{Timeout: 300 * time.Second},
	}
}

// IsConfigured 检查是否已配置
func (e *CozeEvaluator) IsConfigured() bool {
	return e.config.Token != "" && e.config.WorkflowID != ""
}

// uploadFile 上传文件到 Coze
func (e *CozeEvaluator) uploadFile(ctx context.Context, filename string, data []byte) (string, error) {
	url := fmt.Sprintf("%s/v1/files/upload", e.config.BaseURL)

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, err := w.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("create form file: %w", err)
	}
	if _, err := part.Write(data); err != nil {
		return "", fmt.Errorf("write file data: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+e.config.Token)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("coze upload http %d: %s", resp.StatusCode, string(body))
	}

	var out struct {
		Code int             `json:"code"`
		Msg  string          `json:"msg"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("parse upload response: %w", err)
	}
	if out.Code != 0 {
		return "", fmt.Errorf("coze upload error: %s", out.Msg)
	}

	var dataMap map[string]interface{}
	if err := json.Unmarshal(out.Data, &dataMap); err != nil {
		return "", fmt.Errorf("parse upload data: %w", err)
	}
	id, _ := dataMap["id"].(string)
	if id == "" {
		return "", fmt.Errorf("coze upload response missing file id")
	}
	return id, nil
}

// EvaluateResume 评估简历
func (e *CozeEvaluator) EvaluateResume(ctx context.Context, name string, jdText string, resumePDF []byte) (*EvaluationResult, error) {
	if !e.IsConfigured() {
		return nil, fmt.Errorf("Coze 未配置，请设置 COZE_TOKEN 和 COZE_WORKFLOW_ID 环境变量")
	}

	// 确保文件名以 .pdf 结尾
	filename := name
	if filename == "" {
		filename = "resume.pdf"
	}
	if !strings.HasSuffix(strings.ToLower(filename), ".pdf") {
		filename = strings.TrimSuffix(filename, "."+strings.Split(filename, ".")[len(strings.Split(filename, "."))-1]) + ".pdf"
	}

	// 1. 上传文件获取 file_id
	fileID, err := e.uploadFile(ctx, filename, resumePDF)
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}

	// 2. 调用工作流
	requestBody := map[string]interface{}{
		"workflow_id":   e.config.WorkflowID,
		"response_mode": "blocking",
		"parameters": map[string]interface{}{
			"name":        name,
			"jd_text":     jdText,
			"resume_file": map[string]interface{}{"file_id": fileID},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	url := fmt.Sprintf("%s/v1/workflow/run", e.config.BaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+e.config.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("coze http %d: %s", resp.StatusCode, string(body))
	}

	// 3. 解析响应
	var envelope map[string]interface{}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查错误码
	if code, ok := envelope["code"].(float64); ok && int(code) != 0 {
		msg, _ := envelope["msg"].(string)
		return nil, fmt.Errorf("coze error: %s", msg)
	}

	// 提取结果
	return e.parseResult(envelope, name)
}

// parseResult 解析 Coze 返回的结果
func (e *CozeEvaluator) parseResult(envelope map[string]interface{}, name string) (*EvaluationResult, error) {
	result := &EvaluationResult{
		Name:      name,
		RawResult: envelope,
	}

	// 提取 output 或 data.output
	var outputStr string
	if data, ok := envelope["data"].(map[string]interface{}); ok {
		if output, ok := data["output"].(string); ok {
			outputStr = output
		}
	}
	if outputStr == "" {
		if output, ok := envelope["output"].(string); ok {
			outputStr = output
		}
	}

	if outputStr == "" {
		return result, nil
	}

	// 清理 JSON 字符串
	outputStr = cleanJSONString(outputStr)

	var resultData map[string]interface{}
	if err := json.Unmarshal([]byte(outputStr), &resultData); err != nil {
		return result, nil // 解析失败返回基本结果
	}

	// 提取基本信息
	if basicInfo, ok := resultData["基本信息"].(map[string]interface{}); ok {
		if score, ok := basicInfo["最终得分"].(float64); ok {
			result.TotalScore = score
		}
		if grade, ok := basicInfo["评级"].(string); ok {
			result.Grade = grade
		}
	}

	// 提取各维度得分
	if scores, ok := resultData["各维度得分"].(map[string]interface{}); ok {
		if age, ok := scores["年龄"].(map[string]interface{}); ok {
			if score, ok := age["得分"].(float64); ok {
				result.AgeScore = int(score)
			}
		}
		if exp, ok := scores["工作经验"].(map[string]interface{}); ok {
			if score, ok := exp["得分"].(float64); ok {
				result.ExperienceScore = int(score)
			}
		}
		if edu, ok := scores["学历背景"].(map[string]interface{}); ok {
			if score, ok := edu["得分"].(float64); ok {
				result.EducationScore = int(score)
			}
		}
		if company, ok := scores["公司背景"].(map[string]interface{}); ok {
			if score, ok := company["得分"].(float64); ok {
				result.CompanyScore = int(score)
			}
		}
		if tech, ok := scores["技术能力"].(map[string]interface{}); ok {
			if score, ok := tech["得分"].(float64); ok {
				result.TechScore = int(score)
			}
		}
		if project, ok := scores["项目经历"].(map[string]interface{}); ok {
			if score, ok := project["得分"].(float64); ok {
				result.ProjectScore = int(score)
			}
		}
	}

	// 提取 JD 匹配度
	if jdMatch, ok := resultData["JD匹配度"].(map[string]interface{}); ok {
		if score, ok := jdMatch["匹配分数"].(float64); ok {
			result.JDMatchScore = int(score)
		}
		if summary, ok := jdMatch["匹配总结"].(string); ok {
			result.Summary = summary
		}
		if matched, ok := jdMatch["匹配的技能"].([]interface{}); ok {
			for _, s := range matched {
				if str, ok := s.(string); ok {
					result.MatchedSkills = append(result.MatchedSkills, str)
				}
			}
		}
		if missing, ok := jdMatch["缺失的技能"].([]interface{}); ok {
			for _, s := range missing {
				if str, ok := s.(string); ok {
					result.MissingSkills = append(result.MissingSkills, str)
				}
			}
		}
	}

	// 提取录用建议
	if rec, ok := resultData["录用建议"].(map[string]interface{}); ok {
		if conclusion, ok := rec["结论"].(string); ok {
			result.Recommendation = conclusion
		}
	}

	return result, nil
}

// cleanJSONString 清理 JSON 字符串
func cleanJSONString(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```json") {
		s = strings.TrimPrefix(s, "```json")
	} else if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```")
	}
	if strings.HasSuffix(s, "```") {
		s = strings.TrimSuffix(s, "```")
	}
	return strings.TrimSpace(s)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
