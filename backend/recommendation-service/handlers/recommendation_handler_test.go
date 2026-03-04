package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	handler := NewRecommendationHandler(nil) // 测试时不使用数据库

	api := r.Group("/api/v1/recommendations")
	{
		api.POST("/jobs-for-talent", handler.RecommendJobsForTalent)
		api.POST("/talents-for-job", handler.RecommendTalentsForJob)
		api.GET("/stats", handler.GetRecommendationStats)
	}

	return r
}

func TestRecommendJobsForTalent(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name           string
		talent         TalentProfile
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name: "Go开发者应该匹配Go相关职位",
			talent: TalentProfile{
				ID:         1,
				Name:       "张三",
				Skills:     []string{"Go", "Docker", "Kubernetes"},
				Experience: 5,
				Education:  "本科",
				Location:   "北京",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, float64(0), body["code"])
				data := body["data"].([]interface{})
				assert.Greater(t, len(data), 0)

				// 第一个推荐应该是高匹配度
				firstRec := data[0].(map[string]interface{})
				assert.Greater(t, firstRec["score"].(float64), 60.0)
			},
		},
		{
			name: "前端开发者应该匹配前端职位",
			talent: TalentProfile{
				ID:         2,
				Name:       "李四",
				Skills:     []string{"Vue", "TypeScript", "React"},
				Experience: 3,
				Education:  "硕士",
				Location:   "上海",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, float64(0), body["code"])
				data := body["data"].([]interface{})
				assert.Greater(t, len(data), 0)
			},
		},
		{
			name: "无技能人才应该返回低匹配度",
			talent: TalentProfile{
				ID:         3,
				Name:       "王五",
				Skills:     []string{},
				Experience: 0,
				Location:   "其他城市",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, float64(0), body["code"])
				data := body["data"].([]interface{})
				if len(data) > 0 {
					firstRec := data[0].(map[string]interface{})
					assert.Less(t, firstRec["score"].(float64), 50.0)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.talent)
			req, _ := http.NewRequest("POST", "/api/v1/recommendations/jobs-for-talent", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			tt.checkResponse(t, response)
		})
	}
}

func TestRecommendTalentsForJob(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name           string
		job            JobProfile
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name: "Go职位应该匹配Go开发者",
			job: JobProfile{
				ID:       1,
				Title:    "高级Go开发工程师",
				Skills:   []string{"Go", "Docker", "Kubernetes"},
				Location: "北京",
				Level:    "senior",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, float64(0), body["code"])
				data := body["data"].([]interface{})
				assert.Greater(t, len(data), 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.job)
			req, _ := http.NewRequest("POST", "/api/v1/recommendations/talents-for-job", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			tt.checkResponse(t, response)
		})
	}
}

func TestRecommendJobsForTalentInvalidJSON(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("POST", "/api/v1/recommendations/jobs-for-talent", strings.NewReader("{"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(1), response["code"])
}

func TestRecommendTalentsForJobInvalidJSON(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("POST", "/api/v1/recommendations/talents-for-job", strings.NewReader("{"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(1), response["code"])
}

func TestRecommendJobsForTalentScoreSortedDesc(t *testing.T) {
	router := setupRouter()

	talent := TalentProfile{
		ID:         9,
		Name:       "测试候选人",
		Skills:     []string{"Go", "Docker"},
		Experience: 4,
		Education:  "本科",
		Location:   "北京",
	}

	body, _ := json.Marshal(talent)
	req, _ := http.NewRequest("POST", "/api/v1/recommendations/jobs-for-talent", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].([]interface{})
	assert.Greater(t, len(data), 1)

	lastScore := data[0].(map[string]interface{})["score"].(float64)
	for i := 1; i < len(data); i++ {
		score := data[i].(map[string]interface{})["score"].(float64)
		assert.LessOrEqual(t, score, lastScore, "推荐结果应按分数降序排序")
		lastScore = score
	}
}

func TestGetRecommendationStats(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/v1/recommendations/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, float64(0), response["code"])
	data := response["data"].(map[string]interface{})
	assert.Contains(t, data, "total_recommendations")
	assert.Contains(t, data, "successful_matches")
	assert.Contains(t, data, "success_rate")
}

func TestCalculateAdvancedMatchScoreDeterministic(t *testing.T) {
	talent := TalentProfile{
		Skills:     []string{"Go", "Docker", "Kubernetes"},
		Experience: 5,
		Education:  "本科",
		Location:   "北京",
		Salary:     "30-40K",
	}
	job := JobProfile{
		Skills:   []string{"Go", "Docker"},
		Location: "北京",
		Level:    "senior",
		Salary:   "30-50K",
	}

	score1, details1 := calculateAdvancedMatchScore(talent, job)
	score2, details2 := calculateAdvancedMatchScore(talent, job)

	assert.Equal(t, score1, score2)
	assert.Equal(t, details1, details2)
}

func TestCalculateAdvancedMatchScore(t *testing.T) {
	tests := []struct {
		name     string
		talent   TalentProfile
		job      JobProfile
		minScore float64
		maxScore float64
	}{
		{
			name: "完美匹配",
			talent: TalentProfile{
				Skills:     []string{"Go", "Docker", "Kubernetes"},
				Experience: 7,
				Education:  "硕士",
				Location:   "北京",
			},
			job: JobProfile{
				Skills:   []string{"Go", "Docker", "Kubernetes"},
				Location: "北京",
				Level:    "senior",
			},
			minScore: 80,
			maxScore: 100,
		},
		{
			name: "部分匹配",
			talent: TalentProfile{
				Skills:     []string{"Go", "MySQL"},
				Experience: 3,
				Education:  "本科",
				Location:   "上海",
			},
			job: JobProfile{
				Skills:   []string{"Go", "Docker", "Kubernetes"},
				Location: "北京",
				Level:    "mid",
			},
			minScore: 30,
			maxScore: 60,
		},
		{
			name: "不匹配",
			talent: TalentProfile{
				Skills:     []string{"Java", "Spring"},
				Experience: 1,
				Education:  "大专",
				Location:   "成都",
			},
			job: JobProfile{
				Skills:   []string{"Go", "Docker", "Kubernetes"},
				Location: "北京",
				Level:    "senior",
			},
			minScore: 0,
			maxScore: 40,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, details := calculateAdvancedMatchScore(tt.talent, tt.job)
			assert.GreaterOrEqual(t, score, tt.minScore, "分数应该大于等于最小值")
			assert.LessOrEqual(t, score, tt.maxScore, "分数应该小于等于最大值")
			assert.NotEmpty(t, details, "应该返回匹配详情")
		})
	}
}

func TestSkillMatch(t *testing.T) {
	tests := []struct {
		name         string
		talentSkills []string
		jobSkills    []string
		minScore     float64
	}{
		{
			name:         "完全匹配",
			talentSkills: []string{"Go", "Docker", "Kubernetes"},
			jobSkills:    []string{"Go", "Docker", "Kubernetes"},
			minScore:     0.9,
		},
		{
			name:         "部分匹配",
			talentSkills: []string{"Go", "MySQL"},
			jobSkills:    []string{"Go", "Docker", "Kubernetes"},
			minScore:     0.3,
		},
		{
			name:         "大小写不敏感",
			talentSkills: []string{"go", "DOCKER"},
			jobSkills:    []string{"Go", "Docker"},
			minScore:     0.9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, _ := calculateSkillMatch(tt.talentSkills, tt.jobSkills)
			assert.GreaterOrEqual(t, score, tt.minScore)
		})
	}
}

func TestSkillMatchWithEmptyJobSkills(t *testing.T) {
	score, details := calculateSkillMatch([]string{"Go"}, []string{})
	assert.Equal(t, 0.5, score)
	assert.Contains(t, details[0], "职位未指定技能要求")
}

func TestLocationMatch(t *testing.T) {
	tests := []struct {
		name      string
		talentLoc string
		jobLoc    string
		minScore  float64
	}{
		{
			name:      "完全匹配",
			talentLoc: "北京",
			jobLoc:    "北京",
			minScore:  1.0,
		},
		{
			name:      "同城市群",
			talentLoc: "上海",
			jobLoc:    "杭州",
			minScore:  0.6,
		},
		{
			name:      "不匹配",
			talentLoc: "北京",
			jobLoc:    "深圳",
			minScore:  0.2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, _ := calculateLocationMatch(tt.talentLoc, tt.jobLoc)
			assert.GreaterOrEqual(t, score, tt.minScore)
		})
	}
}

func TestParseSkills(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "PostgreSQL数组格式",
			input:    "{Go,Docker, Kubernetes}",
			expected: []string{"Go", "Docker", "Kubernetes"},
		},
		{
			name:     "空字符串",
			input:    "",
			expected: []string{},
		},
		{
			name:     "仅大括号",
			input:    "{}",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseSkills(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
