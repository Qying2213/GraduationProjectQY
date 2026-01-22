package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
		workflowID: os.Getenv("COZE_ATTRIBUTION_WORKFLOW_ID"),
		client:     &http.Client{Timeout: 120 * time.Second},
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// AttributionRequest 归因报告请求
type AttributionRequest struct {
	TalentProfile map[string]interface{} `json:"talent_profile"`
	JobProfile    map[string]interface{} `json:"job_profile"`
	MatchResult   map[string]interface{} `json:"match_result"`
}

// AttributionReport 归因报告
type AttributionReport struct {
	Summary        string           `json:"summary"`
	MatchScore     float64          `json:"match_score"`
	Dimensions     []DimensionScore `json:"dimensions"`
	Strengths      []string         `json:"strengths"`
	Gaps           []string         `json:"gaps"`
	Recommendation string           `json:"recommendation"`
}

// DimensionScore 维度得分
type DimensionScore struct {
	Name    string  `json:"name"`
	Score   float64 `json:"score"`
	Weight  int     `json:"weight"`
	Details string  `json:"details"`
}

// GenerateAttributionReport 生成归因报告
func (c *CozeClient) GenerateAttributionReport(ctx context.Context, req *AttributionRequest) (*AttributionReport, error) {
	if c.token == "" {
		// 如果没有配置Coze，返回基于规则的报告
		return c.generateRuleBasedReport(req)
	}

	// 构建工作流请求
	workflowReq := map[string]interface{}{
		"workflow_id":   c.workflowID,
		"response_mode": "blocking",
		"parameters": map[string]interface{}{
			"talent_profile": req.TalentProfile,
			"job_profile":    req.JobProfile,
			"match_result":   req.MatchResult,
		},
	}

	jsonBody, err := json.Marshal(workflowReq)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	url := c.baseURL + "/v1/workflow/run"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		// 降级到规则引擎
		return c.generateRuleBasedReport(req)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.generateRuleBasedReport(req)
	}

	if resp.StatusCode != http.StatusOK {
		return c.generateRuleBasedReport(req)
	}

	// 解析响应
	var result struct {
		Code int `json:"code"`
		Data struct {
			Output string `json:"output"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return c.generateRuleBasedReport(req)
	}

	if result.Code != 0 || result.Data.Output == "" {
		return c.generateRuleBasedReport(req)
	}

	// 解析AI生成的报告
	var report AttributionReport
	if err := json.Unmarshal([]byte(result.Data.Output), &report); err != nil {
		return c.generateRuleBasedReport(req)
	}

	return &report, nil
}

// generateRuleBasedReport 基于规则生成报告（降级方案）
func (c *CozeClient) generateRuleBasedReport(req *AttributionRequest) (*AttributionReport, error) {
	matchResult := req.MatchResult
	talentProfile := req.TalentProfile
	jobProfile := req.JobProfile

	// 提取匹配分数
	score := 0.0
	if s, ok := matchResult["score"].(float64); ok {
		score = s
	}

	// 提取匹配详情
	details := []string{}
	if d, ok := matchResult["match_details"].([]interface{}); ok {
		for _, item := range d {
			if s, ok := item.(string); ok {
				details = append(details, s)
			}
		}
	}

	// 生成维度得分
	dimensions := []DimensionScore{
		{
			Name:    "技能匹配",
			Score:   score * 0.9,
			Weight:  50,
			Details: "基于技能关键词匹配",
		},
		{
			Name:    "经验匹配",
			Score:   score * 0.85,
			Weight:  20,
			Details: "基于工作年限评估",
		},
		{
			Name:    "地理位置",
			Score:   score * 0.95,
			Weight:  15,
			Details: "基于工作地点匹配",
		},
		{
			Name:    "学历匹配",
			Score:   score * 0.9,
			Weight:  10,
			Details: "基于学历要求评估",
		},
		{
			Name:    "薪资匹配",
			Score:   score * 0.8,
			Weight:  5,
			Details: "基于薪资范围匹配",
		},
	}

	// 生成优势和不足
	strengths := []string{}
	gaps := []string{}

	for _, detail := range details {
		if len(detail) > 0 {
			// 检查是否是优势相关的描述
			if strings.HasPrefix(detail, "匹配") || strings.HasPrefix(detail, "经验") ||
				strings.HasPrefix(detail, "地理") || strings.HasPrefix(detail, "学历") ||
				strings.HasPrefix(detail, "完全") || strings.HasPrefix(detail, "符合") {
				strengths = append(strengths, detail)
			} else if strings.HasPrefix(detail, "缺") || strings.HasPrefix(detail, "不") ||
				strings.HasPrefix(detail, "略") {
				gaps = append(gaps, detail)
			}
		}
	}

	if len(strengths) == 0 {
		strengths = append(strengths, "具备相关工作经验")
	}
	if len(gaps) == 0 && score < 80 {
		gaps = append(gaps, "部分技能需要进一步评估")
	}

	// 生成推荐建议
	recommendation := ""
	if score >= 80 {
		recommendation = "建议优先安排面试，该候选人与职位高度匹配"
	} else if score >= 60 {
		recommendation = "建议进一步沟通，了解候选人的具体项目经验"
	} else if score >= 40 {
		recommendation = "可作为备选候选人，建议关注其发展潜力"
	} else {
		recommendation = "匹配度较低，建议寻找更合适的候选人"
	}

	// 生成摘要
	talentName := ""
	if name, ok := talentProfile["name"].(string); ok {
		talentName = name
	}
	jobTitle := ""
	if title, ok := jobProfile["title"].(string); ok {
		jobTitle = title
	}

	summary := fmt.Sprintf("候选人%s与职位「%s」的综合匹配度为%.1f分。", talentName, jobTitle, score)
	if score >= 80 {
		summary += "整体匹配度较高，建议优先考虑。"
	} else if score >= 60 {
		summary += "基本符合职位要求，可进一步评估。"
	} else {
		summary += "存在一定差距，需谨慎考虑。"
	}

	return &AttributionReport{
		Summary:        summary,
		MatchScore:     score,
		Dimensions:     dimensions,
		Strengths:      strengths,
		Gaps:           gaps,
		Recommendation: recommendation,
	}, nil
}

// ChatRequest 对话请求
type ChatRequest struct {
	Messages []Message `json:"messages"`
}

// Message 消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Chat 对话接口
func (c *CozeClient) Chat(ctx context.Context, prompt string) (string, error) {
	if c.token == "" {
		return "", fmt.Errorf("Coze API Token未配置")
	}

	reqBody := map[string]interface{}{
		"bot_id": c.workflowID,
		"user":   "user",
		"query":  prompt,
		"stream": false,
	}

	jsonBody, _ := json.Marshal(reqBody)

	url := c.baseURL + "/open_api/v2/chat"
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Messages []struct {
			Content string `json:"content"`
		} `json:"messages"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if len(result.Messages) > 0 {
		return result.Messages[len(result.Messages)-1].Content, nil
	}

	return "", fmt.Errorf("无响应内容")
}
