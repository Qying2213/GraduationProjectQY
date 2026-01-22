package parser

import (
	"encoding/json"
	"regexp"
	"strings"
)

// ParsedResume 解析后的简历结构
type ParsedResume struct {
	Name        string           `json:"name"`
	Phone       string           `json:"phone"`
	Email       string           `json:"email"`
	Location    string           `json:"location"`
	Education   string           `json:"education"`
	School      string           `json:"school"`
	Major       string           `json:"major"`
	Experience  string           `json:"experience"`
	Skills      []string         `json:"skills"`
	Summary     string           `json:"summary"`
	WorkHistory []WorkExperience `json:"work_history"`
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

// ResumeParser 简历解析器
type ResumeParser struct{}

// NewResumeParser 创建解析器实例
func NewResumeParser() *ResumeParser {
	return &ResumeParser{}
}

// Parse 解析简历文本
func (p *ResumeParser) Parse(text string) (*ParsedResume, error) {
	result := &ParsedResume{}

	// 提取姓名（通常在简历开头）
	result.Name = p.extractName(text)

	// 提取手机号
	result.Phone = p.extractPhone(text)

	// 提取邮箱
	result.Email = p.extractEmail(text)

	// 提取技能
	result.Skills = p.extractSkills(text)

	// 提取教育背景
	result.Education = p.extractEducation(text)
	result.School = p.extractSchool(text)
	result.Major = p.extractMajor(text)

	// 提取工作经验年限
	result.Experience = p.extractExperience(text)

	// 提取地点
	result.Location = p.extractLocation(text)

	// 提取工作经历
	result.WorkHistory = p.extractWorkHistory(text)

	// 执行风控检查
	result.RiskItems = p.checkRisks(result)
	result.RiskScore = p.calculateRiskScore(result.RiskItems)

	return result, nil
}

// checkRisks 检查风险项
func (p *ResumeParser) checkRisks(resume *ParsedResume) []RiskItem {
	var risks []RiskItem

	// 检查必要信息是否缺失
	if resume.Name == "" {
		risks = append(risks, RiskItem{
			Type:    "missing_name",
			Level:   "warning",
			Message: "未能提取到姓名信息",
		})
	}

	if resume.Phone == "" && resume.Email == "" {
		risks = append(risks, RiskItem{
			Type:    "missing_contact",
			Level:   "warning",
			Message: "未提供有效联系方式",
		})
	}

	if resume.Education == "" {
		risks = append(risks, RiskItem{
			Type:    "missing_education",
			Level:   "info",
			Message: "未提供学历信息",
		})
	}

	// 检查工作经历时间冲突
	for i := 0; i < len(resume.WorkHistory)-1; i++ {
		current := resume.WorkHistory[i]
		next := resume.WorkHistory[i+1]

		if current.EndDate != "" && next.StartDate != "" {
			if current.EndDate > next.StartDate {
				risks = append(risks, RiskItem{
					Type:    "time_conflict",
					Level:   "warning",
					Message: "工作经历时间可能存在重叠",
				})
			}
		}
	}

	return risks
}

// calculateRiskScore 计算风险分数
func (p *ResumeParser) calculateRiskScore(risks []RiskItem) int {
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

// extractSchool 提取学校名称
func (p *ResumeParser) extractSchool(text string) string {
	// 匹配常见大学名称模式
	schoolRegex := regexp.MustCompile(`([\x{4e00}-\x{9fa5}]{2,10}(?:大学|学院|学校))`)
	if matches := schoolRegex.FindStringSubmatch(text); len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// extractMajor 提取专业
func (p *ResumeParser) extractMajor(text string) string {
	majorRegex := regexp.MustCompile(`专业[：:]\s*([\x{4e00}-\x{9fa5}a-zA-Z]+)`)
	if matches := majorRegex.FindStringSubmatch(text); len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// extractWorkHistory 提取工作经历
func (p *ResumeParser) extractWorkHistory(text string) []WorkExperience {
	var history []WorkExperience

	// 简单的工作经历提取（实际应用中可能需要更复杂的逻辑）
	// 匹配公司名称模式
	companyRegex := regexp.MustCompile(`([\x{4e00}-\x{9fa5}a-zA-Z]{2,20}(?:公司|集团|科技|网络|信息|软件))`)
	companies := companyRegex.FindAllString(text, -1)

	for _, company := range companies {
		if len(company) > 2 {
			history = append(history, WorkExperience{
				Company: company,
			})
		}
	}

	return history
}

// ParseToJSON 解析并返回JSON字符串
func (p *ResumeParser) ParseToJSON(text string) (string, error) {
	result, err := p.Parse(text)
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// extractName 提取姓名
func (p *ResumeParser) extractName(text string) string {
	// 中文姓名正则：2-4个汉字
	nameRegex := regexp.MustCompile(`姓\s*名[：:]\s*([\x{4e00}-\x{9fa5}]{2,4})`)
	if matches := nameRegex.FindStringSubmatch(text); len(matches) > 1 {
		return matches[1]
	}

	// 尝试从开头提取
	lines := strings.Split(text, "\n")
	for _, line := range lines[:min(5, len(lines))] {
		line = strings.TrimSpace(line)
		if len(line) >= 2 && len(line) <= 12 {
			// 检查是否全是中文
			if matched, _ := regexp.MatchString(`^[\x{4e00}-\x{9fa5}]{2,4}$`, line); matched {
				return line
			}
		}
	}

	return ""
}

// extractPhone 提取手机号
func (p *ResumeParser) extractPhone(text string) string {
	phoneRegex := regexp.MustCompile(`1[3-9]\d{9}`)
	if matches := phoneRegex.FindString(text); matches != "" {
		return matches
	}
	return ""
}

// extractEmail 提取邮箱
func (p *ResumeParser) extractEmail(text string) string {
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	if matches := emailRegex.FindString(text); matches != "" {
		return matches
	}
	return ""
}

// extractSkills 提取技能
func (p *ResumeParser) extractSkills(text string) []string {
	// 常见技术技能关键词
	skillKeywords := []string{
		// 编程语言
		"Java", "Python", "Go", "Golang", "JavaScript", "TypeScript", "C++", "C#", "PHP", "Ruby", "Swift", "Kotlin", "Rust",
		// 前端
		"Vue", "Vue.js", "Vue3", "React", "Angular", "HTML", "CSS", "SCSS", "Less", "Webpack", "Vite", "Node.js", "Next.js",
		// 后端
		"Spring", "Spring Boot", "Django", "Flask", "Express", "Gin", "FastAPI", "Laravel",
		// 数据库
		"MySQL", "PostgreSQL", "MongoDB", "Redis", "Elasticsearch", "Oracle", "SQL Server",
		// 云和DevOps
		"Docker", "Kubernetes", "K8s", "AWS", "Azure", "GCP", "Linux", "Nginx", "Jenkins", "Git", "CI/CD",
		// 大数据和AI
		"Hadoop", "Spark", "Flink", "TensorFlow", "PyTorch", "机器学习", "深度学习",
		// 其他
		"微服务", "分布式", "高并发", "消息队列", "RabbitMQ", "Kafka", "gRPC", "RESTful",
	}

	textLower := strings.ToLower(text)
	var skills []string
	seen := make(map[string]bool)

	for _, skill := range skillKeywords {
		skillLower := strings.ToLower(skill)
		if strings.Contains(textLower, skillLower) && !seen[skillLower] {
			skills = append(skills, skill)
			seen[skillLower] = true
		}
	}

	return skills
}

// extractEducation 提取教育背景
func (p *ResumeParser) extractEducation(text string) string {
	educationLevels := []string{"博士", "硕士", "研究生", "本科", "学士", "大专", "专科"}

	for _, level := range educationLevels {
		if strings.Contains(text, level) {
			return level
		}
	}

	return ""
}

// extractExperience 提取工作经验
func (p *ResumeParser) extractExperience(text string) string {
	// 匹配 "X年工作经验" 或 "X年经验"
	expRegex := regexp.MustCompile(`(\d+)\s*年[工作]*经验`)
	if matches := expRegex.FindStringSubmatch(text); len(matches) > 1 {
		return matches[1] + "年"
	}

	// 匹配 "工作经验：X年"
	expRegex2 := regexp.MustCompile(`工作经验[：:]\s*(\d+)\s*年`)
	if matches := expRegex2.FindStringSubmatch(text); len(matches) > 1 {
		return matches[1] + "年"
	}

	return ""
}

// extractLocation 提取地点
func (p *ResumeParser) extractLocation(text string) string {
	cities := []string{
		"北京", "上海", "深圳", "广州", "杭州", "成都", "南京", "武汉", "西安", "苏州",
		"天津", "重庆", "长沙", "郑州", "青岛", "大连", "宁波", "厦门", "福州", "合肥",
	}

	for _, city := range cities {
		if strings.Contains(text, city) {
			return city
		}
	}

	return ""
}

// CalculateMatchScore 计算简历与职位的匹配度
func (p *ResumeParser) CalculateMatchScore(resume *ParsedResume, jobSkills []string, jobExperience int, jobEducation string) int {
	score := 0
	maxScore := 100

	// 技能匹配 (最高50分)
	if len(jobSkills) > 0 {
		matchedSkills := 0
		for _, skill := range resume.Skills {
			for _, jobSkill := range jobSkills {
				if strings.EqualFold(skill, jobSkill) {
					matchedSkills++
					break
				}
			}
		}
		skillScore := (matchedSkills * 50) / len(jobSkills)
		if skillScore > 50 {
			skillScore = 50
		}
		score += skillScore
	}

	// 学历匹配 (最高25分)
	educationRank := map[string]int{
		"博士": 5, "硕士": 4, "研究生": 4, "本科": 3, "学士": 3, "大专": 2, "专科": 2,
	}
	resumeEduRank := educationRank[resume.Education]
	jobEduRank := educationRank[jobEducation]
	if resumeEduRank >= jobEduRank {
		score += 25
	} else if resumeEduRank == jobEduRank-1 {
		score += 15
	}

	// 经验匹配 (最高25分)
	if resume.Experience != "" {
		expYears := 0
		if _, err := regexp.MatchString(`\d+`, resume.Experience); err == nil {
			expRegex := regexp.MustCompile(`(\d+)`)
			if matches := expRegex.FindStringSubmatch(resume.Experience); len(matches) > 1 {
				expYears, _ = parseInt(matches[1])
			}
		}
		if expYears >= jobExperience {
			score += 25
		} else if expYears >= jobExperience-1 {
			score += 15
		}
	}

	if score > maxScore {
		score = maxScore
	}

	return score
}

func parseInt(s string) (int, error) {
	var result int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		}
	}
	return result, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
