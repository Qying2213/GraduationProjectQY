package handlers

import (
	"bytes"
	"encoding/json"
	"interview-service/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
		api.POST("/:id/reschedule", handler.RescheduleInterview)
	}

	return r
}

func TestCreateInterview(t *testing.T) {
	db := setupTestDB()
	router := setupTestRouter(db)
	futureDate := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	pastDate := time.Now().Add(-24 * time.Hour).Format("2006-01-02")

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
				Date:          futureDate,
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
		{
			name: "创建过去时间面试应失败",
			request: models.InterviewScheduleRequest{
				CandidateID:   2,
				CandidateName: "李四",
				PositionID:    1,
				Position:      "Go开发工程师",
				Type:          "initial",
				Date:          pastDate,
				Time:          "14:00",
				Duration:      60,
				InterviewerID: 2,
				Interviewer:   "王五",
				Method:        "video",
				Location:      "腾讯会议",
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   1005,
		},
		{
			name: "创建无效时间格式面试应失败",
			request: models.InterviewScheduleRequest{
				CandidateID:   3,
				CandidateName: "赵六",
				PositionID:    2,
				Position:      "前端开发工程师",
				Type:          "initial",
				Date:          futureDate,
				Time:          "invalid-time",
				Duration:      45,
				InterviewerID: 3,
				Interviewer:   "钱七",
				Method:        "onsite",
				Location:      "会议室B",
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   1005,
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
		{
			CandidateID:   3,
			CandidateName: "赵六",
			PositionID:    3,
			Position:      "测试开发",
			Type:          models.InterviewTypeInitial,
			Date:          "2024-12-27",
			Time:          "09:30",
			InterviewerID: 1,
			Interviewer:   "李四",
			Status:        models.InterviewStatusScheduled,
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
			expectedCount:  3,
		},
		{
			name:           "按状态筛选",
			query:          "?status=scheduled",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "按日期筛选",
			query:          "?date=2024-12-25",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "按日期范围筛选",
			query:          "?start_date=2024-12-25&end_date=2024-12-26",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "按面试官筛选",
			query:          "?interviewer_id=1",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "按候选人筛选",
			query:          "?candidate_id=2",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "组合筛选",
			query:          "?interviewer_id=1&status=scheduled",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
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

func TestUpdateInterview(t *testing.T) {
	db := setupTestDB()

	interview := models.Interview{
		CandidateID:   1,
		CandidateName: "张三",
		PositionID:    1,
		Position:      "Go开发",
		Type:          models.InterviewTypeInitial,
		Date:          "2026-12-25",
		Time:          "14:00",
		InterviewerID: 1,
		Interviewer:   "李四",
		Status:        models.InterviewStatusScheduled,
	}
	db.Create(&interview)

	router := setupTestRouter(db)

	updateReq := map[string]interface{}{
		"time":           "16:30",
		"interviewer_id": 99,
		"interviewer":    "新面试官",
		"status":         "completed",
		"location":       "线上会议",
		"notes":          "改由技术负责人面试",
		"feedback":       "表现优秀",
		"rating":         5,
	}
	body, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest("PUT", "/api/v1/interviews/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(0), response["code"])

	var updated models.Interview
	db.First(&updated, 1)
	assert.Equal(t, "16:30", updated.Time)
	assert.Equal(t, uint(99), updated.InterviewerID)
	assert.Equal(t, "新面试官", updated.Interviewer)
	assert.Equal(t, models.InterviewStatusCompleted, updated.Status)
	assert.Equal(t, "线上会议", updated.Location)
	assert.Equal(t, "改由技术负责人面试", updated.Notes)
	assert.Equal(t, "表现优秀", updated.Feedback)
	assert.Equal(t, 5, updated.Rating)
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

func TestCompleteInterview(t *testing.T) {
	db := setupTestDB()

	interview := models.Interview{
		CandidateID:   1,
		CandidateName: "张三",
		PositionID:    1,
		Position:      "Go开发",
		Type:          models.InterviewTypeInitial,
		Date:          "2026-12-25",
		Time:          "14:00",
		InterviewerID: 1,
		Interviewer:   "李四",
		Status:        models.InterviewStatusScheduled,
	}
	db.Create(&interview)

	router := setupTestRouter(db)

	completeReq := map[string]interface{}{
		"feedback": "综合表现优秀",
		"rating":   5,
	}
	body, _ := json.Marshal(completeReq)
	req, _ := http.NewRequest("POST", "/api/v1/interviews/1/complete", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updated models.Interview
	db.First(&updated, 1)
	assert.Equal(t, models.InterviewStatusCompleted, updated.Status)
	assert.Equal(t, "综合表现优秀", updated.Feedback)
	assert.Equal(t, 5, updated.Rating)
}

func TestRescheduleInterview(t *testing.T) {
	t.Run("成功改期", func(t *testing.T) {
		db := setupTestDB()
		interview := models.Interview{
			CandidateID:   1,
			CandidateName: "张三",
			PositionID:    1,
			Position:      "Go开发",
			Type:          models.InterviewTypeInitial,
			Date:          "2026-12-25",
			Time:          "14:00",
			InterviewerID: 1,
			Interviewer:   "李四",
			Status:        models.InterviewStatusScheduled,
			Notes:         "原始安排",
		}
		db.Create(&interview)
		router := setupTestRouter(db)

		newDate := time.Now().Add(48 * time.Hour).Format("2006-01-02")
		payload := map[string]interface{}{
			"date":   newDate,
			"time":   "15:30",
			"reason": "面试官出差",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/api/v1/interviews/1/reschedule", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, float64(0), response["code"])

		var updated models.Interview
		db.First(&updated, 1)
		assert.Equal(t, newDate, updated.Date)
		assert.Equal(t, "15:30", updated.Time)
		assert.Equal(t, models.InterviewStatusScheduled, updated.Status)
		assert.Contains(t, updated.Notes, "[改期]")
		assert.Contains(t, updated.Notes, "面试官出差")
	})

	t.Run("改期到过去时间应失败", func(t *testing.T) {
		db := setupTestDB()
		interview := models.Interview{
			CandidateID:   1,
			CandidateName: "张三",
			PositionID:    1,
			Position:      "Go开发",
			Type:          models.InterviewTypeInitial,
			Date:          "2026-12-25",
			Time:          "14:00",
			InterviewerID: 1,
			Interviewer:   "李四",
			Status:        models.InterviewStatusScheduled,
		}
		db.Create(&interview)
		router := setupTestRouter(db)

		pastDate := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
		payload := map[string]interface{}{
			"date": pastDate,
			"time": "10:00",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/api/v1/interviews/1/reschedule", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, float64(1005), response["code"])
	})

	t.Run("改期无效时间格式应失败", func(t *testing.T) {
		db := setupTestDB()
		interview := models.Interview{
			CandidateID:   1,
			CandidateName: "张三",
			PositionID:    1,
			Position:      "Go开发",
			Type:          models.InterviewTypeInitial,
			Date:          "2026-12-25",
			Time:          "14:00",
			InterviewerID: 1,
			Interviewer:   "李四",
			Status:        models.InterviewStatusScheduled,
		}
		db.Create(&interview)
		router := setupTestRouter(db)

		payload := map[string]interface{}{
			"date": "2026-12-30",
			"time": "invalid",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/api/v1/interviews/1/reschedule", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, float64(1005), response["code"])
	})
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
