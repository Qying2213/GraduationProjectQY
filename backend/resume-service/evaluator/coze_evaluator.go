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
	"regexp"
	"strconv"
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
	ParsedReport    map[string]interface{} `json:"parsed_report,omitempty"`
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
func (e *CozeEvaluator) EvaluateResume(ctx context.Context, name string, jdText string, resumeText string, resumePDF []byte) (*EvaluationResult, error) {
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
			"resume_text": truncateString(strings.TrimSpace(resumeText), 12000),
			// Coze 工作流的 file 类型入参要求传序列化后的 JSON 字符串，而不是对象本身。
			"resume_file": fmt.Sprintf("{\"file_id\":\"%s\"}", fileID),
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

	// 打印完整的原始响应用于调试
	fmt.Println("\n========== [Coze Debug] START ==========")
	fmt.Printf("[Coze Debug] HTTP Status: %d\n", resp.StatusCode)
	fmt.Println("[Coze Debug] Raw Response (FULL):")
	fmt.Println(string(body))
	fmt.Println("========== [Coze Debug] END RAW ==========")

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

	// 提取 output - Coze 返回格式: {"data": "{\"result\": \"...\"}"}
	var outputStr string

	// 情况1: data 是 JSON 字符串，需要先解析
	if dataStr, ok := envelope["data"].(string); ok {
		fmt.Printf("[Coze Debug] data is string, parsing...\n")
		var dataObj map[string]interface{}
		if err := json.Unmarshal([]byte(dataStr), &dataObj); err == nil {
			// 从解析后的 data 中提取 result
			if resultVal, ok := dataObj["result"].(string); ok {
				outputStr = resultVal
				fmt.Printf("[Coze Debug] Found data.result: %s\n", truncateString(outputStr, 500))
			} else if outputVal, ok := dataObj["output"].(string); ok {
				outputStr = outputVal
				fmt.Printf("[Coze Debug] Found data.output: %s\n", truncateString(outputStr, 500))
			}
		}
	}

	// 情况2: data 是对象
	if outputStr == "" {
		if data, ok := envelope["data"].(map[string]interface{}); ok {
			fmt.Printf("[Coze Debug] data is object, keys: %v\n", getMapKeys(data))
			if resultVal, ok := data["result"].(string); ok {
				outputStr = resultVal
			} else if output, ok := data["output"].(string); ok {
				outputStr = output
			}
		}
	}

	// 情况3: 直接在 envelope 中
	if outputStr == "" {
		if output, ok := envelope["output"].(string); ok {
			outputStr = output
			fmt.Printf("[Coze Debug] Found envelope.output: %s\n", truncateString(outputStr, 500))
		}
	}

	if outputStr == "" {
		fmt.Printf("[Coze Debug] WARNING: output is empty! Full envelope: %v\n", envelope)
		return result, nil
	}

	fmt.Printf("[Coze Debug] Successfully extracted output, length: %d\n", len(outputStr))
	fmt.Println("\n========== [Coze Debug] EXTRACTED OUTPUT (FULL) ==========")
	fmt.Println(outputStr)
	fmt.Println("========== [Coze Debug] END EXTRACTED ==========")

	rawOutputStr := outputStr

	// 清理 JSON 字符串（移除 markdown 代码块标记）
	outputStr = cleanJSONString(outputStr)
	fmt.Printf("[Coze Debug] After cleaning, length: %d\n", len(outputStr))
	fmt.Println("\n========== [Coze Debug] CLEANED JSON (FULL) ==========")
	fmt.Println(outputStr)
	fmt.Println("========== [Coze Debug] END CLEANED ==========")

	var resultData map[string]interface{}
	if err := json.Unmarshal([]byte(outputStr), &resultData); err != nil {
		fmt.Printf("[Coze Debug] JSON parse error: %v\n", err)
		fmt.Printf("[Coze Debug] Attempting regex extraction as fallback...\n")

		// JSON 解析失败时，使用正则提取关键数据
		extractedResult := extractScoresFromText(outputStr)
		if extractedResult != nil {
			extractedResult.Name = name
			extractedResult.RawResult = envelope
			fmt.Printf("[Coze Debug] Regex extraction successful! TotalScore: %.1f, Grade: %s\n",
				extractedResult.TotalScore, extractedResult.Grade)
			return extractedResult, nil
		}

		return result, nil // 解析失败返回基本结果
	}

	fmt.Printf("[Coze Debug] JSON parsed successfully! Keys: %v\n", getMapKeys(resultData))
	result.ParsedReport = resultData
	if recovered := restoreInterviewQuestions(resultData, rawOutputStr); recovered > 0 {
		fmt.Printf("[Coze Debug] Recovered %d interview questions from truncated output\n", recovered)
	}

	// 打印解析后的各个部分
	if basicInfo, ok := resultData["基本信息"].(map[string]interface{}); ok {
		fmt.Printf("[Coze Debug] 基本信息: %v\n", basicInfo)
	}
	if scores, ok := resultData["各维度得分"].(map[string]interface{}); ok {
		fmt.Printf("[Coze Debug] 各维度得分: %v\n", scores)
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

// restoreInterviewQuestions 尝试从原始输出里恢复「面试题目」数组。
// 背景：当 Coze 输出在题目区域被截断时，cleanJSONString 会把面试题置空以保证 JSON 可解析。
// 这里额外做一次「已完整对象」恢复，避免前端完全看不到题目。
func restoreInterviewQuestions(parsed map[string]interface{}, rawOutput string) int {
	if parsed == nil || rawOutput == "" {
		return 0
	}

	rec, ok := parsed["录用建议"].(map[string]interface{})
	if !ok || rec == nil {
		rec = map[string]interface{}{}
		parsed["录用建议"] = rec
	}

	if existing, ok := rec["面试题目"].([]interface{}); ok && len(existing) > 0 {
		return 0
	}

	recovered := extractInterviewQuestionObjects(rawOutput)
	if len(recovered) == 0 {
		return 0
	}

	items := make([]interface{}, 0, len(recovered))
	for _, question := range recovered {
		items = append(items, question)
	}
	rec["面试题目"] = items
	return len(items)
}

func extractInterviewQuestionObjects(rawOutput string) []map[string]interface{} {
	keyIdx := strings.Index(rawOutput, `"面试题目"`)
	if keyIdx < 0 {
		return nil
	}

	arrayStartRel := strings.Index(rawOutput[keyIdx:], "[")
	if arrayStartRel < 0 {
		return nil
	}

	arrayContent := rawOutput[keyIdx+arrayStartRel+1:]
	return parseCompleteJSONObjectsFromArray(arrayContent)
}

func parseCompleteJSONObjectsFromArray(input string) []map[string]interface{} {
	if input == "" {
		return nil
	}

	results := make([]map[string]interface{}, 0)
	inString := false
	escaped := false
	depth := 0
	objStart := -1

	for i := 0; i < len(input); i++ {
		ch := input[i]

		if inString {
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == '"' {
				inString = false
			}
			continue
		}

		switch ch {
		case '"':
			inString = true
		case '{':
			if depth == 0 {
				objStart = i
			}
			depth++
		case '}':
			if depth == 0 {
				continue
			}
			depth--
			if depth == 0 && objStart >= 0 {
				rawObj := input[objStart : i+1]
				var obj map[string]interface{}
				if err := json.Unmarshal([]byte(rawObj), &obj); err == nil {
					results = append(results, obj)
				}
				objStart = -1
			}
		case ']':
			if depth == 0 {
				return results
			}
		}
	}

	return results
}

// cleanJSONString 清理 JSON 字符串
func cleanJSONString(s string) string {
	s = strings.TrimSpace(s)

	// 移除开头的 ```json 或 ```
	if strings.HasPrefix(s, "```json") {
		s = strings.TrimPrefix(s, "```json")
	} else if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```")
	}

	// 移除结尾的 ```
	if strings.HasSuffix(s, "```") {
		s = strings.TrimSuffix(s, "```")
	}

	s = strings.TrimSpace(s)

	// 尝试修复不完整的 JSON
	// 计算括号平衡
	openBraces := strings.Count(s, "{")
	closeBraces := strings.Count(s, "}")
	openBrackets := strings.Count(s, "[")
	closeBrackets := strings.Count(s, "]")

	// 如果缺少闭合括号，尝试修复
	if openBraces > closeBraces || openBrackets > closeBrackets {
		// 找到最后一个完整的对象结束位置
		// 策略：找到 "录用建议" 部分的结束，截断面试题目
		if idx := strings.Index(s, `"面试题目"`); idx > 0 {
			// 找到面试题目之前的位置，截断
			// 往前找到 "薪资建议" 的值结束位置
			salaryIdx := strings.Index(s, `"薪资建议"`)
			if salaryIdx > 0 {
				// 找到薪资建议值的结束引号
				afterSalary := s[salaryIdx+len(`"薪资建议"`):]
				// 跳过 ": "
				colonIdx := strings.Index(afterSalary, `"`)
				if colonIdx >= 0 {
					afterSalary = afterSalary[colonIdx+1:]
					// 找到值的结束引号
					endQuoteIdx := strings.Index(afterSalary, `"`)
					if endQuoteIdx >= 0 {
						// 截断到薪资建议结束，然后补齐括号
						cutPoint := salaryIdx + len(`"薪资建议"`) + colonIdx + 1 + endQuoteIdx + 1
						s = s[:cutPoint] + `,"面试题目":[]}}`
						return strings.TrimSpace(s)
					}
				}
			}
		}

		// 通用修复：补齐缺失的括号
		// 先补齐方括号
		missingBrackets := openBrackets - closeBrackets
		for i := 0; i < missingBrackets; i++ {
			s += "]"
		}
		// 再补齐花括号
		missingBraces := openBraces - closeBraces
		for i := 0; i < missingBraces; i++ {
			s += "}"
		}
	}

	return strings.TrimSpace(s)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// extractScoresFromText 从文本中使用正则提取分数（JSON解析失败时的备用方案）
func extractScoresFromText(text string) *EvaluationResult {
	result := &EvaluationResult{}

	// 提取最终得分
	if match := regexp.MustCompile(`"最终得分"\s*:\s*(\d+\.?\d*)`).FindStringSubmatch(text); len(match) > 1 {
		if score, err := strconv.ParseFloat(match[1], 64); err == nil {
			result.TotalScore = score
		}
	}

	// 提取评级
	if match := regexp.MustCompile(`"评级"\s*:\s*"([A-Z])"`).FindStringSubmatch(text); len(match) > 1 {
		result.Grade = match[1]
	}

	// 提取 JD 匹配分数
	if match := regexp.MustCompile(`"匹配分数"\s*:\s*(\d+)`).FindStringSubmatch(text); len(match) > 1 {
		if score, err := strconv.Atoi(match[1]); err == nil {
			result.JDMatchScore = score
		}
	}

	// 提取各维度得分 - 年龄
	if match := regexp.MustCompile(`"年龄"\s*:\s*\{[^}]*"得分"\s*:\s*(\d+)`).FindStringSubmatch(text); len(match) > 1 {
		if score, err := strconv.Atoi(match[1]); err == nil {
			result.AgeScore = score
		}
	}

	// 提取各维度得分 - 工作经验
	if match := regexp.MustCompile(`"工作经验"\s*:\s*\{[^}]*"得分"\s*:\s*(\d+)`).FindStringSubmatch(text); len(match) > 1 {
		if score, err := strconv.Atoi(match[1]); err == nil {
			result.ExperienceScore = score
		}
	}

	// 提取各维度得分 - 学历背景
	if match := regexp.MustCompile(`"学历背景"\s*:\s*\{[^}]*"得分"\s*:\s*(\d+)`).FindStringSubmatch(text); len(match) > 1 {
		if score, err := strconv.Atoi(match[1]); err == nil {
			result.EducationScore = score
		}
	}

	// 提取各维度得分 - 公司背景
	if match := regexp.MustCompile(`"公司背景"\s*:\s*\{[^}]*"得分"\s*:\s*(\d+)`).FindStringSubmatch(text); len(match) > 1 {
		if score, err := strconv.Atoi(match[1]); err == nil {
			result.CompanyScore = score
		}
	}

	// 提取各维度得分 - 技术能力
	if match := regexp.MustCompile(`"技术能力"\s*:\s*\{[^}]*"得分"\s*:\s*(\d+)`).FindStringSubmatch(text); len(match) > 1 {
		if score, err := strconv.Atoi(match[1]); err == nil {
			result.TechScore = score
		}
	}

	// 提取各维度得分 - 项目经历
	if match := regexp.MustCompile(`"项目经历"\s*:\s*\{[^}]*"得分"\s*:\s*(\d+)`).FindStringSubmatch(text); len(match) > 1 {
		if score, err := strconv.Atoi(match[1]); err == nil {
			result.ProjectScore = score
		}
	}

	// 提取录用建议
	if match := regexp.MustCompile(`"结论"\s*:\s*"([^"]+)"`).FindStringSubmatch(text); len(match) > 1 {
		result.Recommendation = match[1]
	}

	// 提取匹配总结
	if match := regexp.MustCompile(`"匹配总结"\s*:\s*"([^"]+)"`).FindStringSubmatch(text); len(match) > 1 {
		result.Summary = match[1]
	}

	// 如果提取到了关键数据，返回结果
	if result.TotalScore > 0 || result.Grade != "" {
		return result
	}

	return nil
}

// getMapKeys 获取 map 的所有 key
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// truncateString 截断字符串
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// max 返回两个整数中的较大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
