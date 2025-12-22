package handlers

import (
	"math"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

type RecommendationHandler struct{}

func NewRecommendationHandler() *RecommendationHandler {
	return &RecommendationHandler{}
}

type TalentProfile struct {
	ID              uint     `json:"id"`
	Name            string   `json:"name"`
	Skills          []string `json:"skills"`
	Experience      int      `json:"experience"`
	Education       string   `json:"education"`
	Location        string   `json:"location"`
	Salary          string   `json:"salary"`
	CurrentCompany  string   `json:"current_company"`
	CurrentPosition string   `json:"current_position"`
}

type JobProfile struct {
	ID           uint     `json:"id"`
	Title        string   `json:"title"`
	Skills       []string `json:"skills"`
	Location     string   `json:"location"`
	Requirements []string `json:"requirements"`
	Level        string   `json:"level"`
	Salary       string   `json:"salary"`
	Department   string   `json:"department"`
}

type Recommendation struct {
	ID           uint     `json:"id"`
	Name         string   `json:"name"`
	Score        float64  `json:"score"`
	Reason       string   `json:"reason"`
	MatchLevel   string   `json:"match_level"`
	MatchDetails []string `json:"match_details"`
}

// SkillWeight 技能权重配置
var skillWeights = map[string]float64{
	"go":         1.2,
	"python":     1.1,
	"java":       1.1,
	"kubernetes": 1.3,
	"docker":     1.2,
	"react":      1.1,
	"vue":        1.1,
	"typescript": 1.1,
	"postgresql": 1.0,
	"mysql":      1.0,
	"redis":      1.0,
	"aws":        1.2,
	"机器学习":       1.3,
	"深度学习":       1.3,
}

// EducationScore 学历分数
var educationScores = map[string]float64{
	"博士": 1.0,
	"硕士": 0.9,
	"本科": 0.8,
	"大专": 0.6,
	"高中": 0.4,
}

// calculateAdvancedMatchScore 增强版匹配算法
func calculateAdvancedMatchScore(talent TalentProfile, job JobProfile) (float64, []string) {
	var details []string
	totalScore := 0.0

	// 1. 技能匹配 (50%)
	skillScore, skillDetails := calculateSkillMatch(talent.Skills, job.Skills)
	totalScore += skillScore * 0.5
	details = append(details, skillDetails...)

	// 2. 经验匹配 (20%)
	expScore, expDetail := calculateExperienceMatch(talent.Experience, job.Level)
	totalScore += expScore * 0.2
	details = append(details, expDetail)

	// 3. 地理位置匹配 (15%)
	locScore, locDetail := calculateLocationMatch(talent.Location, job.Location)
	totalScore += locScore * 0.15
	details = append(details, locDetail)

	// 4. 学历匹配 (10%)
	eduScore, eduDetail := calculateEducationMatch(talent.Education, job.Level)
	totalScore += eduScore * 0.1
	details = append(details, eduDetail)

	// 5. 薪资匹配 (5%)
	salaryScore, salaryDetail := calculateSalaryMatch(talent.Salary, job.Salary)
	totalScore += salaryScore * 0.05
	if salaryDetail != "" {
		details = append(details, salaryDetail)
	}

	return math.Min(totalScore*100, 100), details
}

// calculateSkillMatch 计算技能匹配度
func calculateSkillMatch(talentSkills, jobSkills []string) (float64, []string) {
	if len(jobSkills) == 0 {
		return 0.5, []string{"职位未指定技能要求"}
	}

	matchedSkills := []string{}
	totalWeight := 0.0
	matchedWeight := 0.0

	for _, js := range jobSkills {
		jsLower := strings.ToLower(strings.TrimSpace(js))
		weight := skillWeights[jsLower]
		if weight == 0 {
			weight = 1.0
		}
		totalWeight += weight

		for _, ts := range talentSkills {
			tsLower := strings.ToLower(strings.TrimSpace(ts))
			if jsLower == tsLower || strings.Contains(tsLower, jsLower) || strings.Contains(jsLower, tsLower) {
				matchedWeight += weight
				matchedSkills = append(matchedSkills, js)
				break
			}
		}
	}

	score := matchedWeight / totalWeight
	var details []string
	if len(matchedSkills) > 0 {
		details = append(details, "匹配技能: "+strings.Join(matchedSkills, ", "))
	}
	if len(matchedSkills) < len(jobSkills) {
		missingCount := len(jobSkills) - len(matchedSkills)
		details = append(details, "缺少 "+string(rune('0'+missingCount))+" 项技能")
	}

	return score, details
}

// calculateExperienceMatch 计算经验匹配度
func calculateExperienceMatch(experience int, level string) (float64, string) {
	levelRequirements := map[string]struct{ min, ideal, max int }{
		"junior":     {0, 1, 2},
		"mid":        {2, 4, 6},
		"senior":     {5, 7, 10},
		"expert":     {8, 10, 15},
		"management": {5, 8, 15},
	}

	req, ok := levelRequirements[strings.ToLower(level)]
	if !ok {
		req = levelRequirements["mid"]
	}

	var score float64
	var detail string

	if experience >= req.min && experience <= req.max {
		if experience >= req.ideal {
			score = 1.0
			detail = "经验完全匹配"
		} else {
			score = 0.8
			detail = "经验基本匹配"
		}
	} else if experience < req.min {
		score = float64(experience) / float64(req.min) * 0.6
		detail = "经验略显不足"
	} else {
		score = 0.7
		detail = "经验超出要求，可能期望更高"
	}

	return score, detail
}

// calculateLocationMatch 计算地理位置匹配度
func calculateLocationMatch(talentLoc, jobLoc string) (float64, string) {
	if talentLoc == "" || jobLoc == "" {
		return 0.5, "位置信息不完整"
	}

	talentLower := strings.ToLower(talentLoc)
	jobLower := strings.ToLower(jobLoc)

	// 完全匹配
	if strings.Contains(talentLower, jobLower) || strings.Contains(jobLower, talentLower) {
		return 1.0, "地理位置匹配"
	}

	// 同城市群
	cityGroups := [][]string{
		{"北京", "天津", "河北"},
		{"上海", "苏州", "杭州", "南京"},
		{"广州", "深圳", "东莞", "佛山"},
		{"成都", "重庆"},
	}

	for _, group := range cityGroups {
		talentInGroup := false
		jobInGroup := false
		for _, city := range group {
			if strings.Contains(talentLower, city) {
				talentInGroup = true
			}
			if strings.Contains(jobLower, city) {
				jobInGroup = true
			}
		}
		if talentInGroup && jobInGroup {
			return 0.7, "同城市群，可考虑"
		}
	}

	return 0.3, "地理位置不匹配"
}

// calculateEducationMatch 计算学历匹配度
func calculateEducationMatch(education, level string) (float64, string) {
	eduScore, ok := educationScores[education]
	if !ok {
		eduScore = 0.7
	}

	levelEduRequirements := map[string]float64{
		"junior":     0.6,
		"mid":        0.7,
		"senior":     0.8,
		"expert":     0.9,
		"management": 0.8,
	}

	required, ok := levelEduRequirements[strings.ToLower(level)]
	if !ok {
		required = 0.7
	}

	if eduScore >= required {
		return 1.0, "学历符合要求"
	}
	return eduScore / required, "学历略低于要求"
}

// calculateSalaryMatch 计算薪资匹配度
func calculateSalaryMatch(talentSalary, jobSalary string) (float64, string) {
	if talentSalary == "" || jobSalary == "" {
		return 0.5, ""
	}
	// 简化处理，实际应解析薪资范围
	return 0.8, "薪资范围基本匹配"
}

// RecommendJobsForTalent 为人才推荐职位
func (h *RecommendationHandler) RecommendJobsForTalent(c *gin.Context) {
	var talent TalentProfile
	if err := c.ShouldBindJSON(&talent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}

	// 模拟职位数据（实际应该从job-service获取）
	jobs := []JobProfile{
		{ID: 1, Title: "高级Go开发工程师", Skills: []string{"Go", "Docker", "Kubernetes", "微服务"}, Location: "北京", Level: "senior", Salary: "30-50K", Department: "技术部"},
		{ID: 2, Title: "前端架构师", Skills: []string{"Vue", "TypeScript", "React", "Webpack"}, Location: "上海", Level: "senior", Salary: "35-55K", Department: "技术部"},
		{ID: 3, Title: "全栈开发工程师", Skills: []string{"Go", "Vue", "PostgreSQL", "Redis"}, Location: "深圳", Level: "mid", Salary: "25-40K", Department: "产品部"},
		{ID: 4, Title: "后端开发工程师", Skills: []string{"Go", "Redis", "MySQL", "消息队列"}, Location: "杭州", Level: "mid", Salary: "20-35K", Department: "技术部"},
		{ID: 5, Title: "AI算法工程师", Skills: []string{"Python", "机器学习", "深度学习", "TensorFlow"}, Location: "北京", Level: "senior", Salary: "40-70K", Department: "AI部"},
		{ID: 6, Title: "DevOps工程师", Skills: []string{"Docker", "Kubernetes", "Jenkins", "AWS"}, Location: "广州", Level: "mid", Salary: "25-40K", Department: "运维部"},
	}

	var recommendations []Recommendation

	for _, job := range jobs {
		score, details := calculateAdvancedMatchScore(talent, job)

		matchLevel := "low"
		reason := "匹配度较低"
		if score >= 80 {
			matchLevel = "high"
			reason = "高度匹配"
		} else if score >= 60 {
			matchLevel = "medium"
			reason = "中等匹配"
		} else if score >= 40 {
			matchLevel = "low"
			reason = "基本匹配"
		}

		recommendations = append(recommendations, Recommendation{
			ID:           job.ID,
			Name:         job.Title,
			Score:        math.Round(score*10) / 10,
			Reason:       reason,
			MatchLevel:   matchLevel,
			MatchDetails: details,
		})
	}

	// 按分数排序
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// 只返回前10个
	if len(recommendations) > 10 {
		recommendations = recommendations[:10]
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    recommendations,
	})
}

// RecommendTalentsForJob 为职位推荐人才
func (h *RecommendationHandler) RecommendTalentsForJob(c *gin.Context) {
	var job JobProfile
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}

	// 模拟人才数据（实际应该从talent-service获取）
	talents := []TalentProfile{
		{ID: 1, Name: "张三", Skills: []string{"Go", "Docker", "Kubernetes", "Redis", "微服务"}, Experience: 5, Education: "本科", Location: "北京", Salary: "30-40K"},
		{ID: 2, Name: "李四", Skills: []string{"Vue", "TypeScript", "React", "Node.js", "Webpack"}, Experience: 3, Education: "硕士", Location: "上海", Salary: "25-35K"},
		{ID: 3, Name: "王五", Skills: []string{"Go", "Vue", "PostgreSQL", "Docker", "Redis"}, Experience: 4, Education: "本科", Location: "深圳", Salary: "25-35K"},
		{ID: 4, Name: "赵六", Skills: []string{"Java", "Spring", "MySQL", "Redis"}, Experience: 2, Education: "本科", Location: "北京", Salary: "15-25K"},
		{ID: 5, Name: "钱七", Skills: []string{"Python", "机器学习", "TensorFlow", "PyTorch"}, Experience: 6, Education: "博士", Location: "北京", Salary: "50-70K"},
		{ID: 6, Name: "孙八", Skills: []string{"Go", "Kubernetes", "AWS", "Terraform"}, Experience: 7, Education: "硕士", Location: "杭州", Salary: "40-55K"},
	}

	var recommendations []Recommendation

	for _, talent := range talents {
		score, details := calculateAdvancedMatchScore(talent, job)

		matchLevel := "low"
		reason := "匹配度较低"
		if score >= 80 {
			matchLevel = "high"
			reason = "高度匹配"
		} else if score >= 60 {
			matchLevel = "medium"
			reason = "中等匹配"
		} else if score >= 40 {
			matchLevel = "low"
			reason = "基本匹配"
		}

		recommendations = append(recommendations, Recommendation{
			ID:           talent.ID,
			Name:         talent.Name,
			Score:        math.Round(score*10) / 10,
			Reason:       reason,
			MatchLevel:   matchLevel,
			MatchDetails: details,
		})
	}

	// 按分数排序
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// 只返回前10个
	if len(recommendations) > 10 {
		recommendations = recommendations[:10]
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    recommendations,
	})
}

// GetRecommendationStats 获取推荐统计
func (h *RecommendationHandler) GetRecommendationStats(c *gin.Context) {
	stats := gin.H{
		"total_recommendations":   256,
		"successful_matches":      78,
		"pending_reviews":         34,
		"success_rate":            30.5,
		"avg_match_score":         72.3,
		"high_match_count":        45,
		"medium_match_count":      89,
		"low_match_count":         122,
		"top_matched_skills":      []string{"Go", "Vue", "Docker", "Kubernetes", "Python"},
		"recommendations_by_dept": map[string]int{"技术部": 120, "产品部": 56, "运维部": 45, "AI部": 35},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    stats,
	})
}

// BatchRecommend 批量推荐
func (h *RecommendationHandler) BatchRecommend(c *gin.Context) {
	var req struct {
		TalentIDs []uint `json:"talent_ids"`
		JobIDs    []uint `json:"job_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}

	// 返回批量推荐结果
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"processed": len(req.TalentIDs) * len(req.JobIDs),
			"matches":   15,
		},
	})
}
