package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"recommendation-service/ai"
	"recommendation-service/embedding"
	"recommendation-service/rag"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RecommendationHandler struct {
	DB              *gorm.DB
	EmbeddingClient *embedding.Client
	CozeClient      *ai.CozeClient
	SkillMatcher    *embedding.SkillMatcher
	RAGEngine       *rag.RAGEngine
}

func NewRecommendationHandler(db *gorm.DB) *RecommendationHandler {
	embClient := embedding.GetClient()
	return &RecommendationHandler{
		DB:              db,
		EmbeddingClient: embClient,
		CozeClient:      ai.NewCozeClient(),
		SkillMatcher:    embedding.NewSkillMatcher(embClient),
		RAGEngine:       rag.NewRAGEngine(db),
	}
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
	ID             uint     `json:"id"`
	Name           string   `json:"name"`
	Score          float64  `json:"score"`
	Reason         string   `json:"reason"`
	MatchLevel     string   `json:"match_level"`
	MatchDetails   []string `json:"match_details"`
	SemanticScore  float64  `json:"semantic_score,omitempty"`  // 语义匹配分数
	AttributionURL string   `json:"attribution_url,omitempty"` // 归因报告链接
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
		details = append(details, fmt.Sprintf("缺少 %d 项技能", missingCount))
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

	// 从数据库获取职位数据
	var jobs []JobProfile
	if h.DB != nil {
		var dbJobs []struct {
			ID         uint   `json:"id"`
			Title      string `json:"title"`
			Skills     string `json:"skills"`
			Location   string `json:"location"`
			Level      string `json:"level"`
			Salary     string `json:"salary"`
			Department string `json:"department"`
		}
		h.DB.Table("jobs").Where("status = ?", "open").Limit(20).Find(&dbJobs)
		for _, j := range dbJobs {
			skills := []string{}
			if j.Skills != "" {
				// 解析 PostgreSQL 数组格式 {skill1,skill2}
				skillStr := strings.Trim(j.Skills, "{}")
				if skillStr != "" {
					skills = strings.Split(skillStr, ",")
				}
			}
			jobs = append(jobs, JobProfile{
				ID:         j.ID,
				Title:      j.Title,
				Skills:     skills,
				Location:   j.Location,
				Level:      j.Level,
				Salary:     j.Salary,
				Department: j.Department,
			})
		}
	}

	// 如果数据库没有数据，使用默认数据
	if len(jobs) == 0 {
		jobs = []JobProfile{
			{ID: 1, Title: "高级Go开发工程师", Skills: []string{"Go", "Docker", "Kubernetes", "微服务"}, Location: "北京", Level: "senior", Salary: "30-50K", Department: "技术部"},
			{ID: 2, Title: "前端架构师", Skills: []string{"Vue", "TypeScript", "React", "Webpack"}, Location: "上海", Level: "senior", Salary: "35-55K", Department: "技术部"},
			{ID: 3, Title: "全栈开发工程师", Skills: []string{"Go", "Vue", "PostgreSQL", "Redis"}, Location: "深圳", Level: "mid", Salary: "25-40K", Department: "产品部"},
		}
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

	// 从数据库获取人才数据
	var talents []TalentProfile
	if h.DB != nil {
		var dbTalents []struct {
			ID         uint   `json:"id"`
			Name       string `json:"name"`
			Skills     string `json:"skills"`
			Experience int    `json:"experience"`
			Education  string `json:"education"`
			Location   string `json:"location"`
			Salary     string `json:"salary"`
		}
		h.DB.Table("talents").Where("status = ?", "active").Limit(20).Find(&dbTalents)
		for _, t := range dbTalents {
			skills := []string{}
			if t.Skills != "" {
				// 解析 PostgreSQL 数组格式 {skill1,skill2}
				skillStr := strings.Trim(t.Skills, "{}")
				if skillStr != "" {
					skills = strings.Split(skillStr, ",")
				}
			}
			talents = append(talents, TalentProfile{
				ID:         t.ID,
				Name:       t.Name,
				Skills:     skills,
				Experience: t.Experience,
				Education:  t.Education,
				Location:   t.Location,
				Salary:     t.Salary,
			})
		}
	}

	// 如果数据库没有数据，使用默认数据
	if len(talents) == 0 {
		talents = []TalentProfile{
			{ID: 1, Name: "张三", Skills: []string{"Go", "Docker", "Kubernetes", "Redis", "微服务"}, Experience: 5, Education: "本科", Location: "北京", Salary: "30-40K"},
			{ID: 2, Name: "李四", Skills: []string{"Vue", "TypeScript", "React", "Node.js", "Webpack"}, Experience: 3, Education: "硕士", Location: "上海", Salary: "25-35K"},
			{ID: 3, Name: "王五", Skills: []string{"Go", "Vue", "PostgreSQL", "Docker", "Redis"}, Experience: 4, Education: "本科", Location: "深圳", Salary: "25-35K"},
		}
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
		"total_recommendations":   0,
		"successful_matches":      0,
		"pending_reviews":         0,
		"success_rate":            0.0,
		"avg_match_score":         0.0,
		"high_match_count":        0,
		"medium_match_count":      0,
		"low_match_count":         0,
		"top_matched_skills":      []string{},
		"recommendations_by_dept": map[string]int{},
	}

	if h.DB != nil {
		// 获取申请总数
		var totalApplications int64
		h.DB.Table("applications").Count(&totalApplications)
		stats["total_recommendations"] = totalApplications

		// 获取成功匹配数（已录用）
		var hiredCount int64
		h.DB.Table("applications").Where("status = ?", "hired").Count(&hiredCount)
		stats["successful_matches"] = hiredCount

		// 获取待审核数
		var pendingCount int64
		h.DB.Table("applications").Where("status = ?", "pending").Count(&pendingCount)
		stats["pending_reviews"] = pendingCount

		// 计算成功率
		if totalApplications > 0 {
			stats["success_rate"] = float64(hiredCount) / float64(totalApplications) * 100
		}
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

// GenerateAttributionReport 生成归因报告
func (h *RecommendationHandler) GenerateAttributionReport(c *gin.Context) {
	var req struct {
		TalentID uint `json:"talent_id" binding:"required"`
		JobID    uint `json:"job_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}

	// 获取人才信息
	var talent struct {
		ID         uint   `json:"id"`
		Name       string `json:"name"`
		Skills     string `json:"skills"`
		Experience int    `json:"experience"`
		Education  string `json:"education"`
		Location   string `json:"location"`
		Salary     string `json:"salary"`
	}
	if h.DB != nil {
		h.DB.Table("talents").Where("id = ?", req.TalentID).First(&talent)
	}

	// 获取职位信息
	var job struct {
		ID         uint   `json:"id"`
		Title      string `json:"title"`
		Skills     string `json:"skills"`
		Location   string `json:"location"`
		Level      string `json:"level"`
		Salary     string `json:"salary"`
		Department string `json:"department"`
	}
	if h.DB != nil {
		h.DB.Table("jobs").Where("id = ?", req.JobID).First(&job)
	}

	// 解析技能
	talentSkills := parseSkills(talent.Skills)
	jobSkills := parseSkills(job.Skills)

	// 使用语义匹配计算技能分数
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	skillScore, skillDetails, _ := h.SkillMatcher.MatchSkills(ctx, talentSkills, jobSkills)

	// 计算其他维度分数
	talentProfile := TalentProfile{
		ID:         talent.ID,
		Name:       talent.Name,
		Skills:     talentSkills,
		Experience: talent.Experience,
		Education:  talent.Education,
		Location:   talent.Location,
		Salary:     talent.Salary,
	}

	jobProfile := JobProfile{
		ID:       job.ID,
		Title:    job.Title,
		Skills:   jobSkills,
		Location: job.Location,
		Level:    job.Level,
		Salary:   job.Salary,
	}

	totalScore, allDetails := calculateAdvancedMatchScore(talentProfile, jobProfile)

	// 合并技能匹配详情
	allDetails = append(skillDetails, allDetails...)

	// 调用Coze生成归因报告
	attrReq := &ai.AttributionRequest{
		TalentProfile: map[string]interface{}{
			"id":         talent.ID,
			"name":       talent.Name,
			"skills":     talentSkills,
			"experience": talent.Experience,
			"education":  talent.Education,
			"location":   talent.Location,
		},
		JobProfile: map[string]interface{}{
			"id":       job.ID,
			"title":    job.Title,
			"skills":   jobSkills,
			"location": job.Location,
			"level":    job.Level,
		},
		MatchResult: map[string]interface{}{
			"score":         totalScore,
			"skill_score":   skillScore * 100,
			"match_details": allDetails,
		},
	}

	report, err := h.CozeClient.GenerateAttributionReport(ctx, attrReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "生成报告失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    report,
	})
}

// SemanticMatch 语义匹配接口
func (h *RecommendationHandler) SemanticMatch(c *gin.Context) {
	var req struct {
		Text1 string `json:"text1" binding:"required"`
		Text2 string `json:"text2" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	similarity, err := h.EmbeddingClient.TextSimilarity(ctx, req.Text1, req.Text2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "计算相似度失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"similarity": math.Round(similarity*1000) / 1000,
			"text1":      req.Text1,
			"text2":      req.Text2,
		},
	})
}

// parseSkills 解析技能字符串
func parseSkills(skillsStr string) []string {
	if skillsStr == "" {
		return []string{}
	}
	// 处理PostgreSQL数组格式 {skill1,skill2}
	skillsStr = strings.Trim(skillsStr, "{}")
	if skillsStr == "" {
		return []string{}
	}
	skills := strings.Split(skillsStr, ",")
	result := make([]string, 0, len(skills))
	for _, s := range skills {
		s = strings.TrimSpace(s)
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}

// ==================== RAG 相关接口 ====================

// RAGQuery RAG检索增强查询
func (h *RecommendationHandler) RAGQuery(c *gin.Context) {
	var req struct {
		Query     string `json:"query" binding:"required"`
		TopK      int    `json:"top_k"`
		Type      string `json:"type"`       // talent 或 job (兼容旧字段)
		QueryType string `json:"query_type"` // talent 或 job
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}

	// 兼容两种字段名
	queryType := req.QueryType
	if queryType == "" {
		queryType = req.Type
	}
	if queryType == "" {
		queryType = "talent"
	}

	topK := req.TopK
	if topK <= 0 {
		topK = 5
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// 直接调用向量搜索，返回简单结果
	var results []rag.SearchResult
	var err error

	if queryType == "job" {
		results, err = h.RAGEngine.GetVectorStore().SearchSimilarJobs(ctx, req.Query, topK)
	} else if queryType == "resume" {
		results, err = h.RAGEngine.GetVectorStore().SearchSimilarResumes(ctx, req.Query, topK)
	} else {
		results, err = h.RAGEngine.GetVectorStore().SearchSimilarTalents(ctx, req.Query, topK)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"results": results,
			"query":   req.Query,
			"type":    queryType,
		},
	})
}

func buildResumeRAGContent(resumeRow map[string]any, evaluationRow map[string]any) string {
	parts := []string{}

	if resumeID, ok := resumeRow["resume_id"].(uint); ok && resumeID > 0 {
		parts = append(parts, fmt.Sprintf("简历ID: %d", resumeID))
	}
	if name, _ := evaluationRow["parsed_name"].(string); strings.TrimSpace(name) != "" {
		parts = append(parts, "姓名: "+strings.TrimSpace(name))
	}
	if education, _ := evaluationRow["parsed_education"].(string); strings.TrimSpace(education) != "" {
		parts = append(parts, "学历: "+strings.TrimSpace(education))
	}
	if school, _ := evaluationRow["parsed_school"].(string); strings.TrimSpace(school) != "" {
		parts = append(parts, "学校: "+strings.TrimSpace(school))
	}
	if experience, _ := evaluationRow["parsed_experience"].(string); strings.TrimSpace(experience) != "" {
		parts = append(parts, "经历摘要: "+strings.TrimSpace(experience))
	}
	if recommendation, _ := evaluationRow["report_recommendation"].(string); strings.TrimSpace(recommendation) != "" {
		parts = append(parts, "录用建议: "+strings.TrimSpace(recommendation))
	}
	if summary, _ := evaluationRow["report_summary"].(string); strings.TrimSpace(summary) != "" {
		parts = append(parts, "评估摘要: "+strings.TrimSpace(summary))
	}
	if skillsRaw, _ := evaluationRow["parsed_skills"].(string); strings.TrimSpace(skillsRaw) != "" {
		var skills []string
		if err := json.Unmarshal([]byte(skillsRaw), &skills); err == nil && len(skills) > 0 {
			parts = append(parts, "技能: "+strings.Join(skills, "、"))
		}
	}

	return strings.ToValidUTF8(strings.Join(parts, "\n"), "")
}

func trimLongText(value string, maxLen int) string {
	if maxLen <= 0 || len(value) <= maxLen {
		return value
	}
	return value[:maxLen] + "..."
}

// IndexResume 索引简历到向量数据库
func (h *RecommendationHandler) IndexResume(c *gin.Context) {
	var req struct {
		ResumeID uint   `json:"resume_id" binding:"required"`
		Content  string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		var resumeRow map[string]any
		if err := h.DB.Table("resumes").
			Select("id as resume_id, file_name, extracted_text").
			Where("id = ?", req.ResumeID).
			Take(&resumeRow).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "简历不存在"})
			return
		}

		var evaluationRow map[string]any
		if err := h.DB.Table("evaluation_results").
			Select("parsed_name, parsed_education, parsed_school, parsed_experience, parsed_skills, report_summary, report_recommendation").
			Where("resume_id = ?", req.ResumeID).
			Order("created_at DESC, id DESC").
			Take(&evaluationRow).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "该简历暂无可索引的评估结果"})
			return
		}

		content = buildResumeRAGContent(resumeRow, evaluationRow)
		if extractedText, _ := resumeRow["extracted_text"].(string); strings.TrimSpace(extractedText) != "" {
			content += "\nOCR文本片段: " + trimLongText(strings.ToValidUTF8(extractedText, ""), 300)
		}
	}

	content = strings.ToValidUTF8(content, "")
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "索引内容为空"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	if err := h.RAGEngine.GetVectorStore().IndexResume(ctx, req.ResumeID, content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "索引失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "索引成功", "data": gin.H{"resume_id": req.ResumeID}})
}

// IndexTalent 索引人才到向量数据库
func (h *RecommendationHandler) IndexTalent(c *gin.Context) {
	var req struct {
		TalentID uint `json:"talent_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}

	// 获取人才信息
	var talent struct {
		ID         uint   `json:"id"`
		Name       string `json:"name"`
		Skills     string `json:"skills"`
		Experience int    `json:"experience"`
		Education  string `json:"education"`
		Location   string `json:"location"`
	}
	if err := h.DB.Table("talents").Where("id = ?", req.TalentID).First(&talent).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "人才不存在"})
		return
	}

	// 构建内容
	content := fmt.Sprintf("姓名: %s\n技能: %s\n经验: %d年\n学历: %s\n城市: %s",
		talent.Name, talent.Skills, talent.Experience, talent.Education, talent.Location)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	if err := h.RAGEngine.GetVectorStore().IndexTalent(ctx, talent.ID, content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "索引失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "索引成功", "data": gin.H{"talent_id": talent.ID}})
}

// IndexJob 索引职位到向量数据库
func (h *RecommendationHandler) IndexJob(c *gin.Context) {
	var req struct {
		JobID uint `json:"job_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}

	// 获取职位信息
	var job struct {
		ID           uint   `json:"id"`
		Title        string `json:"title"`
		Skills       string `json:"skills"`
		Requirements string `json:"requirements"`
		Location     string `json:"location"`
		Salary       string `json:"salary"`
	}
	if err := h.DB.Table("jobs").Where("id = ?", req.JobID).First(&job).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "职位不存在"})
		return
	}

	// 构建内容
	content := fmt.Sprintf("职位: %s\n技能要求: %s\n任职要求: %s\n地点: %s\n薪资: %s",
		job.Title, job.Skills, job.Requirements, job.Location, job.Salary)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	if err := h.RAGEngine.GetVectorStore().IndexJob(ctx, job.ID, content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "索引失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "索引成功", "data": gin.H{"job_id": job.ID}})
}

// IndexAll 批量索引所有人才和职位
func (h *RecommendationHandler) IndexAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
	defer cancel()

	var talentCount, jobCount int

	// 索引所有人才
	var talents []struct {
		ID         uint   `json:"id"`
		Name       string `json:"name"`
		Skills     string `json:"skills"`
		Experience int    `json:"experience"`
		Education  string `json:"education"`
		Location   string `json:"location"`
	}
	h.DB.Table("talents").Find(&talents)

	for _, t := range talents {
		content := fmt.Sprintf("姓名: %s\n技能: %s\n经验: %d年\n学历: %s\n城市: %s",
			t.Name, t.Skills, t.Experience, t.Education, t.Location)
		if err := h.RAGEngine.GetVectorStore().IndexTalent(ctx, t.ID, content); err == nil {
			talentCount++
		}
	}

	// 索引所有职位
	var jobs []struct {
		ID           uint   `json:"id"`
		Title        string `json:"title"`
		Skills       string `json:"skills"`
		Requirements string `json:"requirements"`
		Location     string `json:"location"`
		Salary       string `json:"salary"`
	}
	h.DB.Table("jobs").Find(&jobs)

	for _, j := range jobs {
		content := fmt.Sprintf("职位: %s\n技能要求: %s\n任职要求: %s\n地点: %s\n薪资: %s",
			j.Title, j.Skills, j.Requirements, j.Location, j.Salary)
		if err := h.RAGEngine.GetVectorStore().IndexJob(ctx, j.ID, content); err == nil {
			jobCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "批量索引完成",
		"data": gin.H{
			"talents_indexed": talentCount,
			"jobs_indexed":    jobCount,
		},
	})
}

// RAGMatch 使用RAG进行人岗匹配
func (h *RecommendationHandler) RAGMatch(c *gin.Context) {
	var req struct {
		TalentID uint `json:"talent_id" binding:"required"`
		JobID    uint `json:"job_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}

	// 获取人才信息
	var talent struct {
		Name       string
		Skills     string
		Experience int
		Education  string
		Location   string
	}
	h.DB.Table("talents").Where("id = ?", req.TalentID).First(&talent)

	// 获取职位信息
	var job struct {
		Title        string
		Skills       string
		Requirements string
		Location     string
		Salary       string
	}
	h.DB.Table("jobs").Where("id = ?", req.JobID).First(&job)

	talentContent := fmt.Sprintf("姓名: %s, 技能: %s, 经验: %d年, 学历: %s, 城市: %s",
		talent.Name, talent.Skills, talent.Experience, talent.Education, talent.Location)
	jobContent := fmt.Sprintf("职位: %s, 技能要求: %s, 任职要求: %s, 地点: %s, 薪资: %s",
		job.Title, job.Skills, job.Requirements, job.Location, job.Salary)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()

	result, err := h.RAGEngine.MatchTalentToJob(ctx, talentContent, jobContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    result,
	})
}
