package handlers

import (
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type RecommendationHandler struct {
}

func NewRecommendationHandler() *RecommendationHandler {
	return &RecommendationHandler{}
}

type TalentProfile struct {
	ID         uint     `json:"id"`
	Name       string   `json:"name"`
	Skills     []string `json:"skills"`
	Experience int      `json:"experience"`
	Education  string   `json:"education"`
	Location   string   `json:"location"`
}

type JobProfile struct {
	ID           uint     `json:"id"`
	Title        string   `json:"title"`
	Skills       []string `json:"skills"`
	Location     string   `json:"location"`
	Requirements []string `json:"requirements"`
	Level        string   `json:"level"`
}

type Recommendation struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	Score      float64 `json:"score"`
	Reason     string  `json:"reason"`
	MatchLevel string  `json:"match_level"` // high, medium, low
}

// calculateMatchScore 计算匹配分数
func calculateMatchScore(talentSkills, jobSkills []string, talentExp int, talentLoc, jobLoc string) float64 {
	score := 0.0

	// 技能匹配 (60%)
	skillMatch := 0
	for _, ts := range talentSkills {
		for _, js := range jobSkills {
			if strings.EqualFold(strings.TrimSpace(ts), strings.TrimSpace(js)) {
				skillMatch++
				break
			}
		}
	}

	if len(jobSkills) > 0 {
		score += (float64(skillMatch) / float64(len(jobSkills))) * 60
	}

	// 经验匹配 (20%)
	if talentExp >= 0 && talentExp <= 2 {
		score += 20 * 0.5 // junior
	} else if talentExp >= 3 && talentExp <= 5 {
		score += 20 * 0.8 // mid
	} else {
		score += 20 // senior
	}

	// 地理位置匹配 (20%)
	if strings.Contains(strings.ToLower(talentLoc), strings.ToLower(jobLoc)) ||
		strings.Contains(strings.ToLower(jobLoc), strings.ToLower(talentLoc)) {
		score += 20
	}

	return math.Min(score, 100)
}

// RecommendJobsForTalent 为人才推荐职位
func (h *RecommendationHandler) RecommendJobsForTalent(c *gin.Context) {
	var talent TalentProfile
	if err := c.ShouldBindJSON(&talent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 模拟职位数据（实际应该从job-service获取）
	jobs := []JobProfile{
		{ID: 1, Title: "Senior Go Developer", Skills: []string{"Go", "Docker", "Kubernetes"}, Location: "北京", Level: "senior"},
		{ID: 2, Title: "Frontend Engineer", Skills: []string{"Vue", "TypeScript", "React"}, Location: "上海", Level: "mid"},
		{ID: 3, Title: "Full Stack Developer", Skills: []string{"Go", "Vue", "PostgreSQL"}, Location: "深圳", Level: "mid"},
		{ID: 4, Title: "Backend Developer", Skills: []string{"Go", "Redis", "MySQL"}, Location: "杭州", Level: "junior"},
	}

	var recommendations []Recommendation

	for _, job := range jobs {
		score := calculateMatchScore(talent.Skills, job.Skills, talent.Experience, talent.Location, job.Location)

		matchLevel := "low"
		reason := "部分技能匹配"
		if score >= 80 {
			matchLevel = "high"
			reason = "高度匹配：技能、经验和地理位置都很符合"
		} else if score >= 60 {
			matchLevel = "medium"
			reason = "中等匹配：主要技能相符"
		}

		recommendations = append(recommendations, Recommendation{
			ID:         job.ID,
			Name:       job.Title,
			Score:      score,
			Reason:     reason,
			MatchLevel: matchLevel,
		})
	}

	// 按分数排序
	for i := 0; i < len(recommendations); i++ {
		for j := i + 1; j < len(recommendations); j++ {
			if recommendations[i].Score < recommendations[j].Score {
				recommendations[i], recommendations[j] = recommendations[j], recommendations[i]
			}
		}
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 模拟人才数据（实际应该从talent-service获取）
	talents := []TalentProfile{
		{ID: 1, Name: "张三", Skills: []string{"Go", "Docker", "Kubernetes", "Redis"}, Experience: 5, Location: "北京"},
		{ID: 2, Name: "李四", Skills: []string{"Vue", "TypeScript", "React", "Node.js"}, Experience: 3, Location: "上海"},
		{ID: 3, Name: "王五", Skills: []string{"Go", "Vue", "PostgreSQL", "Docker"}, Experience: 4, Location: "深圳"},
		{ID: 4, Name: "赵六", Skills: []string{"Java", "Spring", "MySQL"}, Experience: 2, Location: "北京"},
	}

	var recommendations []Recommendation

	for _, talent := range talents {
		score := calculateMatchScore(talent.Skills, job.Skills, talent.Experience, talent.Location, job.Location)

		matchLevel := "low"
		reason := "部分技能匹配"
		if score >= 80 {
			matchLevel = "high"
			reason = "高度匹配：技能、经验和地理位置都很符合"
		} else if score >= 60 {
			matchLevel = "medium"
			reason = "中等匹配：主要技能相符"
		}

		recommendations = append(recommendations, Recommendation{
			ID:         talent.ID,
			Name:       talent.Name,
			Score:      score,
			Reason:     reason,
			MatchLevel: matchLevel,
		})
	}

	// 按分数排序
	for i := 0; i < len(recommendations); i++ {
		for j := i + 1; j < len(recommendations); j++ {
			if recommendations[i].Score < recommendations[j].Score {
				recommendations[i], recommendations[j] = recommendations[j], recommendations[i]
			}
		}
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
		"total_recommendations": 150,
		"successful_matches":    45,
		"pending_reviews":       22,
		"success_rate":          30.0,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    stats,
	})
}
