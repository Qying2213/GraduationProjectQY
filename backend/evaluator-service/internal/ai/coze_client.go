package ai

import (
	"encoding/json"
	"fmt"
	"strings"

	"evaluator-service/internal/models"
)

// CozeClient 从 Coze 返回的 JSON 数据中解析评估结果
type CozeClient struct {
	cozeData map[string]interface{}
}

// NewCozeClient 创建一个新的 CozeClient，使用 Coze 返回的数据
func NewCozeClient(cozeData map[string]interface{}) *CozeClient {
	return &CozeClient{cozeData: cozeData}
}

// cleanJSONString 清理 JSON 字符串，去掉 markdown 代码块标记，并尝试修复截断的 JSON
func cleanJSONString(s string) string {
	s = strings.TrimSpace(s)
	// 去掉开头的 ```json 或 ```
	if strings.HasPrefix(s, "```json") {
		s = strings.TrimPrefix(s, "```json")
	} else if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```")
	}
	// 去掉结尾的 ```
	if strings.HasSuffix(s, "```") {
		s = strings.TrimSuffix(s, "```")
	}
	s = strings.TrimSpace(s)

	// 尝试修复截断的 JSON：补全缺失的括号
	if s != "" && !json.Valid([]byte(s)) {
		// 统计括号数量
		openBraces := strings.Count(s, "{")
		closeBraces := strings.Count(s, "}")
		openBrackets := strings.Count(s, "[")
		closeBrackets := strings.Count(s, "]")

		// 补全缺失的括号
		for i := 0; i < openBrackets-closeBrackets; i++ {
			s += "]"
		}
		for i := 0; i < openBraces-closeBraces; i++ {
			s += "}"
		}

		// 如果还是无效，尝试在末尾加上常见的结束符
		if !json.Valid([]byte(s)) {
			// 尝试补全字符串引号
			s = strings.TrimRight(s, ", \t\n\r")
			if !json.Valid([]byte(s)) {
				s += "\"}]}"
			}
		}
	}

	return s
}

func (c *CozeClient) OCR(images [][]byte) (string, error) {
	return "", nil
}

func (c *CozeClient) Structure(text string) (string, error) {
	// 从 Coze 数据中提取简历结构化内容
	// 如果 Coze 返回了结构化简历，可以从基本信息中提取
	if c.cozeData == nil {
		return text, nil
	}

	// 尝试从基本信息中提取简历内容
	if basicInfo, ok := c.getBasicInfo(); ok {
		if name, ok := basicInfo["姓名"].(string); ok && name != "" {
			return fmt.Sprintf("# 简历\n\n姓名: %s\n", name), nil
		}
	}

	return text, nil
}

func (c *CozeClient) EvaluateJD(resumeMD, jd string) (models.JDMatchResult, error) {
	if c.cozeData == nil {
		return models.JDMatchResult{}, fmt.Errorf("coze data is nil")
	}

	// 解析 output 或 result 字段（JSON 字符串）
	resultStr, ok := c.cozeData["output"].(string)
	if !ok {
		resultStr, ok = c.cozeData["result"].(string)
	}
	if !ok {
		return models.JDMatchResult{}, fmt.Errorf("output/result field not found or not a string")
	}

	// 清理 markdown 代码块标记
	resultStr = cleanJSONString(resultStr)

	var resultData map[string]interface{}
	if err := json.Unmarshal([]byte(resultStr), &resultData); err != nil {
		return models.JDMatchResult{}, fmt.Errorf("failed to parse result JSON: %w", err)
	}

	// 提取 JD匹配度
	jdMatchData, ok := resultData["JD匹配度"].(map[string]interface{})
	if !ok {
		return models.JDMatchResult{}, fmt.Errorf("JD匹配度 field not found")
	}

	result := models.JDMatchResult{}

	// 匹配分数
	if score, ok := jdMatchData["匹配分数"].(float64); ok {
		result.Score = int(score)
	}

	// 匹配总结
	if summary, ok := jdMatchData["匹配总结"].(string); ok {
		result.Summary = summary
	}

	// 匹配的技能
	if matchedSkills, ok := jdMatchData["匹配的技能"].([]interface{}); ok {
		result.MatchedSkills = make([]string, 0, len(matchedSkills))
		for _, skill := range matchedSkills {
			if s, ok := skill.(string); ok {
				result.MatchedSkills = append(result.MatchedSkills, s)
			}
		}
	}

	// 缺失的技能
	if missingSkills, ok := jdMatchData["缺失的技能"].([]interface{}); ok {
		result.MissingSkills = make([]string, 0, len(missingSkills))
		for _, skill := range missingSkills {
			if s, ok := skill.(string); ok {
				result.MissingSkills = append(result.MissingSkills, s)
			}
		}
	}

	return result, nil
}

func (c *CozeClient) EvaluateRequirement(resumeMD string) (models.RequirementResult, error) {
	if c.cozeData == nil {
		return models.RequirementResult{}, fmt.Errorf("coze data is nil")
	}

	// 解析 output 或 result 字段
	resultStr, ok := c.cozeData["output"].(string)
	if !ok {
		resultStr, ok = c.cozeData["result"].(string)
	}
	if !ok {
		return models.RequirementResult{}, fmt.Errorf("output/result field not found or not a string")
	}

	// 清理 markdown 代码块标记
	resultStr = cleanJSONString(resultStr)

	var resultData map[string]interface{}
	if err := json.Unmarshal([]byte(resultStr), &resultData); err != nil {
		return models.RequirementResult{}, fmt.Errorf("failed to parse result JSON: %w", err)
	}

	result := models.RequirementResult{
		BlacklistPass: true,
		BlacklistHits: nil,
		OverallPass:   true,
	}

	// 从基本信息中提取学历和经验信息
	if basicInfo, ok := resultData["基本信息"].(map[string]interface{}); ok {
		// 学历
		if education, ok := basicInfo["学历"].(string); ok {
			result.EducationDetail = education
			result.EducationPass = education != ""
		}

		// 工作经验（从工作经验字段中提取年限）
		if experience, ok := basicInfo["工作经验"].(string); ok && experience != "" {
			result.ExperiencePass = true
			// 尝试从工作经验文本中提取年限（简单解析）
			// 这里可以根据实际格式进行更复杂的解析
			years := 0.0
			if experience != "" {
				years = 1.0 // 默认值，实际应该从文本中解析
			}
			result.ExperienceYears = &years
		}
	}

	// 从录用建议中判断是否通过
	if rec, ok := resultData["录用建议"].(map[string]interface{}); ok {
		if conclusion, ok := rec["结论"].(string); ok {
			// 如果结论是"建议录用"或"推荐"，则通过
			result.OverallPass = conclusion == "建议录用" || conclusion == "推荐"
		}
	}

	return result, nil
}

func (c *CozeClient) Score(resumeMD, jd, criteria string) (models.ScoringResult, error) {
	if c.cozeData == nil {
		return models.ScoringResult{}, fmt.Errorf("coze data is nil")
	}

	// 解析 output 或 result 字段
	resultStr, ok := c.cozeData["output"].(string)
	if !ok {
		resultStr, ok = c.cozeData["result"].(string)
	}
	if !ok {
		return models.ScoringResult{}, fmt.Errorf("output/result field not found or not a string")
	}

	// 清理 markdown 代码块标记
	resultStr = cleanJSONString(resultStr)

	var resultData map[string]interface{}
	if err := json.Unmarshal([]byte(resultStr), &resultData); err != nil {
		return models.ScoringResult{}, fmt.Errorf("failed to parse result JSON: %w", err)
	}

	result := models.ScoringResult{}

	// 从基本信息中提取总分和评级
	if basicInfo, ok := resultData["基本信息"].(map[string]interface{}); ok {
		// 最终得分
		if totalScore, ok := basicInfo["最终得分"].(float64); ok {
			result.TotalScore = totalScore
		}

		// 评级
		if grade, ok := basicInfo["评级"].(string); ok {
			result.Grade = grade
		}
	}

	// 从各维度得分中提取分数和说明
	if scores, ok := resultData["各维度得分"].(map[string]interface{}); ok {
		// 公司背景
		if company, ok := scores["公司背景"].(map[string]interface{}); ok {
			if score, ok := company["得分"].(float64); ok {
				result.CompanyScore = int(score)
			}
			if reason, ok := company["说明"].(string); ok {
				result.CompanyReason = reason
			}
		}

		// 学历背景
		if education, ok := scores["学历背景"].(map[string]interface{}); ok {
			if score, ok := education["得分"].(float64); ok {
				result.EducationScore = int(score)
			}
			if reason, ok := education["说明"].(string); ok {
				result.EducationReason = reason
			}
		}

		// 工作经验
		if experience, ok := scores["工作经验"].(map[string]interface{}); ok {
			if score, ok := experience["得分"].(float64); ok {
				result.ExperienceScore = int(score)
			}
			if reason, ok := experience["说明"].(string); ok {
				result.ExperienceReason = reason
			}
		}

		// 年龄
		if age, ok := scores["年龄"].(map[string]interface{}); ok {
			if score, ok := age["得分"].(float64); ok {
				result.AgeScore = int(score)
			}
			if reason, ok := age["说明"].(string); ok {
				result.AgeReason = reason
			}
		}

		// 技术能力
		if tech, ok := scores["技术能力"].(map[string]interface{}); ok {
			if score, ok := tech["得分"].(float64); ok {
				result.TechScore = int(score)
			}
			if reason, ok := tech["说明"].(string); ok {
				result.TechReason = reason
			}
		}

		// 项目经历
		if project, ok := scores["项目经历"].(map[string]interface{}); ok {
			if score, ok := project["得分"].(float64); ok {
				result.ProjectScore = int(score)
			}
			if reason, ok := project["说明"].(string); ok {
				result.ProjectReason = reason
			}
		}
	}

	return result, nil
}

func (c *CozeClient) GenerateInterviewQuestions(resumeMD string, eval models.EvaluationResult) ([]models.InterviewQuestion, error) {
	// 从 Coze 数据中提取面试问题（如果有的话）
	// 目前 Coze 返回的数据中没有面试问题，返回空数组
	return []models.InterviewQuestion{}, nil
}

// getBasicInfo 辅助方法：从 Coze 数据中提取基本信息
func (c *CozeClient) getBasicInfo() (map[string]interface{}, bool) {
	if c.cozeData == nil {
		return nil, false
	}

	// 优先使用 output 字段，兼容 result 字段
	resultStr, ok := c.cozeData["output"].(string)
	if !ok {
		resultStr, ok = c.cozeData["result"].(string)
	}
	if !ok {
		return nil, false
	}

	// 清理 markdown 代码块标记
	resultStr = cleanJSONString(resultStr)

	var resultData map[string]interface{}
	if err := json.Unmarshal([]byte(resultStr), &resultData); err != nil {
		return nil, false
	}

	basicInfo, ok := resultData["基本信息"].(map[string]interface{})
	return basicInfo, ok
}
