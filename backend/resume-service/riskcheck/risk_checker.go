package riskcheck

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// RiskLevel 风险等级
type RiskLevel string

const (
	RiskLevelHigh   RiskLevel = "high"
	RiskLevelMedium RiskLevel = "medium"
	RiskLevelLow    RiskLevel = "low"
)

// RiskItem 风险项
type RiskItem struct {
	Type        string    `json:"type"`        // 风险类型
	Level       RiskLevel `json:"level"`       // 风险等级
	Description string    `json:"description"` // 风险描述
	Detail      string    `json:"detail"`      // 详细信息
}

// RiskCheckResult 风控检查结果
type RiskCheckResult struct {
	HasRisk   bool       `json:"has_risk"`   // 是否有风险
	RiskItems []RiskItem `json:"risk_items"` // 风险项列表
	Score     int        `json:"score"`      // 风险评分 (0-100, 100为无风险)
}

// ResumeData 简历数据（用于风控检查）
type ResumeData struct {
	Name       string           `json:"name"`
	Age        int              `json:"age"`
	Education  []EducationItem  `json:"education"`
	Experience []ExperienceItem `json:"experience"`
	Skills     []string         `json:"skills"`
}

// EducationItem 教育经历
type EducationItem struct {
	School    string `json:"school"`
	Degree    string `json:"degree"` // 本科、硕士、博士
	Major     string `json:"major"`
	StartYear int    `json:"start_year"`
	EndYear   int    `json:"end_year"`
}

// ExperienceItem 工作经历
type ExperienceItem struct {
	Company   string `json:"company"`
	Position  string `json:"position"`
	StartDate string `json:"start_date"` // 格式: 2020-01
	EndDate   string `json:"end_date"`   // 格式: 2023-05 或 "至今"
}

// RiskChecker 风控检查器
type RiskChecker struct{}

// NewRiskChecker 创建风控检查器
func NewRiskChecker() *RiskChecker {
	return &RiskChecker{}
}

// Check 执行完整风控检查
func (r *RiskChecker) Check(resume *ResumeData) *RiskCheckResult {
	result := &RiskCheckResult{
		RiskItems: []RiskItem{},
		Score:     100,
	}

	// 1. 检查时间冲突
	timeConflicts := r.CheckTimeConflict(resume)
	result.RiskItems = append(result.RiskItems, timeConflicts...)

	// 2. 检查学历造假
	eduFrauds := r.CheckEducationFraud(resume)
	result.RiskItems = append(result.RiskItems, eduFrauds...)

	// 3. 检查逻辑一致性
	inconsistencies := r.CheckLogicalConsistency(resume)
	result.RiskItems = append(result.RiskItems, inconsistencies...)

	// 计算风险评分
	for _, item := range result.RiskItems {
		switch item.Level {
		case RiskLevelHigh:
			result.Score -= 30
		case RiskLevelMedium:
			result.Score -= 15
		case RiskLevelLow:
			result.Score -= 5
		}
	}

	if result.Score < 0 {
		result.Score = 0
	}

	result.HasRisk = len(result.RiskItems) > 0

	return result
}

// CheckTimeConflict 检查时间冲突（工作/教育经历时间段重叠）
func (r *RiskChecker) CheckTimeConflict(resume *ResumeData) []RiskItem {
	var risks []RiskItem

	// 检查工作经历时间重叠
	for i := 0; i < len(resume.Experience); i++ {
		for j := i + 1; j < len(resume.Experience); j++ {
			exp1 := resume.Experience[i]
			exp2 := resume.Experience[j]

			start1, end1 := parseDate(exp1.StartDate), parseDate(exp1.EndDate)
			start2, end2 := parseDate(exp2.StartDate), parseDate(exp2.EndDate)

			if datesOverlap(start1, end1, start2, end2) {
				risks = append(risks, RiskItem{
					Type:        "time_conflict",
					Level:       RiskLevelHigh,
					Description: "工作经历时间冲突",
					Detail: fmt.Sprintf("「%s」(%s~%s) 与 「%s」(%s~%s) 时间重叠",
						exp1.Company, exp1.StartDate, exp1.EndDate,
						exp2.Company, exp2.StartDate, exp2.EndDate),
				})
			}
		}
	}

	// 检查教育经历时间重叠（同时读两个全日制学位可疑）
	for i := 0; i < len(resume.Education); i++ {
		for j := i + 1; j < len(resume.Education); j++ {
			edu1 := resume.Education[i]
			edu2 := resume.Education[j]

			if yearsOverlap(edu1.StartYear, edu1.EndYear, edu2.StartYear, edu2.EndYear) {
				risks = append(risks, RiskItem{
					Type:        "time_conflict",
					Level:       RiskLevelMedium,
					Description: "教育经历时间重叠",
					Detail: fmt.Sprintf("「%s」(%d~%d) 与 「%s」(%d~%d) 时间重叠",
						edu1.School, edu1.StartYear, edu1.EndYear,
						edu2.School, edu2.StartYear, edu2.EndYear),
				})
			}
		}
	}

	return risks
}

// CheckEducationFraud 检查学历造假（年龄与学历是否匹配）
func (r *RiskChecker) CheckEducationFraud(resume *ResumeData) []RiskItem {
	var risks []RiskItem
	currentYear := time.Now().Year()

	for _, edu := range resume.Education {
		// 计算入学时的年龄
		ageAtEnroll := edu.StartYear - (currentYear - resume.Age)

		// 本科入学年龄通常 17-22 岁
		if edu.Degree == "本科" && (ageAtEnroll < 15 || ageAtEnroll > 25) {
			risks = append(risks, RiskItem{
				Type:        "education_fraud",
				Level:       RiskLevelHigh,
				Description: "学历年龄不匹配",
				Detail: fmt.Sprintf("本科入学年龄异常：%d岁（%s，%d年入学）",
					ageAtEnroll, edu.School, edu.StartYear),
			})
		}

		// 硕士入学年龄通常 21-30 岁
		if edu.Degree == "硕士" && (ageAtEnroll < 19 || ageAtEnroll > 35) {
			risks = append(risks, RiskItem{
				Type:        "education_fraud",
				Level:       RiskLevelMedium,
				Description: "学历年龄不匹配",
				Detail: fmt.Sprintf("硕士入学年龄异常：%d岁（%s，%d年入学）",
					ageAtEnroll, edu.School, edu.StartYear),
			})
		}

		// 检查毕业年份是否在未来
		if edu.EndYear > currentYear {
			risks = append(risks, RiskItem{
				Type:        "education_fraud",
				Level:       RiskLevelHigh,
				Description: "毕业年份异常",
				Detail:      fmt.Sprintf("毕业年份在未来：%s %d年毕业", edu.School, edu.EndYear),
			})
		}

		// 检查学制是否合理
		duration := edu.EndYear - edu.StartYear
		if edu.Degree == "本科" && (duration < 3 || duration > 6) {
			risks = append(risks, RiskItem{
				Type:        "education_fraud",
				Level:       RiskLevelMedium,
				Description: "学制异常",
				Detail:      fmt.Sprintf("本科学制%d年异常（%s）", duration, edu.School),
			})
		}
	}

	return risks
}

// CheckLogicalConsistency 检查逻辑一致性（技能与经历是否匹配）
func (r *RiskChecker) CheckLogicalConsistency(resume *ResumeData) []RiskItem {
	var risks []RiskItem

	// 检查高级技能但经验过短
	totalYears := calculateTotalExperience(resume.Experience)

	// 高级技能关键词
	advancedSkills := []string{"架构设计", "技术总监", "CTO", "首席", "资深", "专家"}
	hasAdvancedSkill := false
	for _, skill := range resume.Skills {
		for _, adv := range advancedSkills {
			if strings.Contains(skill, adv) {
				hasAdvancedSkill = true
				break
			}
		}
	}

	if hasAdvancedSkill && totalYears < 5 {
		risks = append(risks, RiskItem{
			Type:        "logical_inconsistency",
			Level:       RiskLevelMedium,
			Description: "技能与经验不匹配",
			Detail:      fmt.Sprintf("声称具备高级技能，但总工作经验仅%d年", totalYears),
		})
	}

	// 检查工作经历空白期
	gaps := findExperienceGaps(resume.Experience)
	for _, gap := range gaps {
		if gap > 12 { // 超过12个月的空白期
			risks = append(risks, RiskItem{
				Type:        "experience_gap",
				Level:       RiskLevelLow,
				Description: "工作经历存在空白期",
				Detail:      fmt.Sprintf("存在约%d个月的工作空白期", gap),
			})
		}
	}

	return risks
}

// 辅助函数

func parseDate(dateStr string) time.Time {
	if dateStr == "至今" || dateStr == "" {
		return time.Now()
	}
	// 尝试解析 2020-01 格式
	t, err := time.Parse("2006-01", dateStr)
	if err != nil {
		// 尝试解析 2020年1月 格式
		re := regexp.MustCompile(`(\d{4})年(\d{1,2})月?`)
		matches := re.FindStringSubmatch(dateStr)
		if len(matches) >= 3 {
			year, _ := strconv.Atoi(matches[1])
			month, _ := strconv.Atoi(matches[2])
			return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
		}
		return time.Time{}
	}
	return t
}

func datesOverlap(start1, end1, start2, end2 time.Time) bool {
	if start1.IsZero() || start2.IsZero() {
		return false
	}
	return start1.Before(end2) && start2.Before(end1)
}

func yearsOverlap(start1, end1, start2, end2 int) bool {
	return start1 < end2 && start2 < end1
}

func calculateTotalExperience(experiences []ExperienceItem) int {
	totalMonths := 0
	for _, exp := range experiences {
		start := parseDate(exp.StartDate)
		end := parseDate(exp.EndDate)
		if !start.IsZero() && !end.IsZero() {
			months := int(end.Sub(start).Hours() / 24 / 30)
			if months > 0 {
				totalMonths += months
			}
		}
	}
	return totalMonths / 12
}

func findExperienceGaps(experiences []ExperienceItem) []int {
	if len(experiences) < 2 {
		return nil
	}

	// 按开始时间排序（简化处理，假设已排序）
	var gaps []int
	for i := 0; i < len(experiences)-1; i++ {
		end := parseDate(experiences[i].EndDate)
		start := parseDate(experiences[i+1].StartDate)

		if !end.IsZero() && !start.IsZero() && start.After(end) {
			gapMonths := int(start.Sub(end).Hours() / 24 / 30)
			if gapMonths > 0 {
				gaps = append(gaps, gapMonths)
			}
		}
	}
	return gaps
}
