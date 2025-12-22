package handlers

import (
	"bytes"
	"encoding/json"
	"interview-service/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Interview{}, &models.InterviewFeedback{})
	return db
}

func setupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	handler := NewInterviewHandler(db)

	api := r.Group("/api/v1/interviews")
	{
		api.POST("", handler.CreateInterview)
		api.GET("", handler.ListInterviews)
		api.GET("/stats", handler.GetInterviewStats)
		api.GET("/today", handler.GetTodayInterviews)
		api.GET("/:id", handler.GetInterview)
		api.PUT("/:id", handler.UpdateInterview)
		api.DELETE("/:id", handler.DeleteInterview)
		api.POST("/:id/cancel", handler.CancelInterview)
		api.POST("/:id/complete", handler.CompleteInterview)
		api.POST("/:id/feedback", handler.SubmitFeedback)
		api.GET("/:id/feedback", handler.GetFeedback)
	}

	return r
}

func TestCreateInterview(t *testing.T) {
	db := setupTestDB()
	router := setupTestRouter(db)

	tests := []struct {
		name           string
		request        models.InterviewScheduleRequest
		expectedStatus int
		expectedCode   float64
	}{
		{
			name: "成功创建面试",
			request: models.InterviewScheduleRequest{
				CandidateID:   1,
				CandidateName: "张三",
				PositionID:    1,
				Position:      "Go开发工程师",
				Type:          "initial",
				Date:          "2024-12-25",
				Time:          "14:00",
				Duration:      60,
				InterviewerID: 1,
				Interviewer:   "李四",
				Method:        "onsite",
				Location:      "会议室A",
			},
			expectedStatus: http.StatusCreated,
			expectedCode:   0,
		},
		{
			name: "缺少必填字段",
			request: models.InterviewScheduleRequest{
				CandidateID: 1,
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/api/v1/interviews", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Equal(t, tt.expectedCode, response["code"])
		})
	}
}

func TestListInterviews(t *testing.T) {
	db := setupTestDB()

	// 创建测试数据
	interviews := []models.Interview{
		{
			CandidateID:   1,
			CandidateName: "张三",
			PositionID:    1,
			Position:      "Go开发",
			Type:          models.InterviewTypeInitial,
			Date:          "2024-12-25",
			Time:          "14:00",
			InterviewerID: 1,
			Interviewer:   "李四",
			Status:        models.InterviewStatusScheduled,
		},
		{
			CandidateID:   2,
			CandidateName: "王五",
			PositionID:    2,
			Position:      "前端开发",
			Type:          models.InterviewTypeSecond,
			Date:          "2024-12-26",
			Time:          "10:00",
			InterviewerID: 2,
			Interviewer:   "赵六",
			Status:        models.InterviewStatusCompleted,
		},
	}
	for _, i := range interviews {
		db.Create(&i)
	}

	router := setupTestRouter(db)

	tests := []struct {
		name           string
		query          string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "获取所有面试",
			query:          "",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "按状态筛选",
			query:          "?status=scheduled",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "按日期筛选",
			query:          "?date=2024-12-25",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/v1/interviews"+tt.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			data := response["data"].(map[string]interface{})
			interviews := data["interviews"].([]interface{})
			assert.Equal(t, tt.expectedCount, len(interviews))
		})
	}
}

func TestGetInterview(t *testing.T) {
	db := setupTestDB()

	interview := models.Interview{
		CandidateID:   1,
		CandidateName: "张三",
		PositionID:    1,
		Position:      "Go开发",
		Type:          models.InterviewTypeInitial,
		Date:          "2024-12-25",
		Time:          "14:00",
		InterviewerID: 1,
		Interviewer:   "李四",
		Status:        models.InterviewStatusScheduled,
	}
	db.Create(&interview)

	router := setupTestRouter(db)

	tests := []struct {
		name           string
		id             string
		expectedStatus int
		expectedCode   float64
	}{
		{
			name:           "获取存在的面试",
			id:             "1",
			expectedStatus: http.StatusOK,
			expectedCode:   0,
		},
		{
			name:           "获取不存在的面试",
			id:             "999",
			expectedStatus: http.StatusNotFound,
			expectedCode:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/v1/interviews/"+tt.id, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Equal(t, tt.expectedCode, response["code"])
		})
	}
}

func TestCancelInterview(t *testing.T) {
	db := setupTestDB()

	interview := models.Interview{
		CandidateID:   1,
		CandidateName: "张三",
		PositionID:    1,
		Position:      "Go开发",
		Type:          models.InterviewTypeInitial,
		Date:          "2024-12-25",
		Time:          "14:00",
		InterviewerID: 1,
		Interviewer:   "李四",
		Status:        models.InterviewStatusScheduled,
	}
	db.Create(&interview)

	router := setupTestRouter(db)

	req, _ := http.NewRequest("POST", "/api/v1/interviews/1/cancel", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 验证状态已更新
	var updated models.Interview
	db.First(&updated, 1)
	assert.Equal(t, models.InterviewStatusCancelled, updated.Status)
}

func TestSubmitFeedback(t *testing.T) {
	db := setupTestDB()

	interview := models.Interview{
		CandidateID:   1,
		CandidateName: "张三",
		PositionID:    1,
		Position:      "Go开发",
		Type:          models.InterviewTypeInitial,
		Date:          "2024-12-25",
		Time:          "14:00",
		InterviewerID: 1,
		Interviewer:   "李四",
		Status:        models.InterviewStatusScheduled,
	}
	db.Create(&interview)

	router := setupTestRouter(db)

	feedback := map[string]interface{}{
		"rating":         4,
		"strengths":      "技术扎实，沟通能力强",
		"weaknesses":     "项目经验略少",
		"comments":       "整体表现良好",
		"recommendation": "pass",
	}

	body, _ := json.Marshal(feedback)
	req, _ := http.NewRequest("POST", "/api/v1/interviews/1/feedback", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 验证反馈已创建
	var fb models.InterviewFeedback
	db.Where("interview_id = ?", 1).First(&fb)
	assert.Equal(t, 4, fb.Rating)
	assert.Equal(t, "pass", fb.Recommendation)

	// 验证面试状态已更新
	var updated models.Interview
	db.First(&updated, 1)
	assert.Equal(t, models.InterviewStatusCompleted, updated.Status)
}

func TestGetInterviewStats(t *testing.T) {
	db := setupTestDB()

	// 创建测试数据
	interviews := []models.Interview{
		{CandidateID: 1, CandidateName: "A", PositionID: 1, Position: "P1", Type: "initial", Date: "2024-12-25", Time: "10:00", InterviewerID: 1, Interviewer: "I1", Status: models.InterviewStatusScheduled},
		{CandidateID: 2, CandidateName: "B", PositionID: 1, Position: "P1", Type: "initial", Date: "2024-12-25", Time: "14:00", InterviewerID: 1, Interviewer: "I1", Status: models.InterviewStatusCompleted},
		{CandidateID: 3, CandidateName: "C", PositionID: 1, Position: "P1", Type: "initial", Date: "2024-12-26", Time: "10:00", InterviewerID: 1, Interviewer: "I1", Status: models.InterviewStatusCancelled},
	}
	for _, i := range interviews {
		db.Create(&i)
	}

	router := setupTestRouter(db)

	req, _ := http.NewRequest("GET", "/api/v1/interviews/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(3), data["total_interviews"])
	assert.Equal(t, float64(1), data["scheduled_interviews"])
	assert.Equal(t, float64(1), data["completed_interviews"])
	assert.Equal(t, float64(1), data["cancelled_interviews"])
}
