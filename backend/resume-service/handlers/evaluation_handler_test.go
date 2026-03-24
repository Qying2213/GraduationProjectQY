package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"resume-service/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupEvaluationTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.EvaluationResult{}, &models.AIProcessLog{}))
	return db
}

func TestGetEvaluationProcessReturnsEmptyTraceWhenMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupEvaluationTestDB(t)
	eval := models.EvaluationResult{
		ResumeID:   25,
		ResumeName: "test.pdf",
		Status:     "completed",
		EvalType:   "ai_evaluate",
	}
	require.NoError(t, db.Create(&eval).Error)

	handler := NewEvaluationHandler(db)
	router := gin.New()
	router.GET("/evaluations/:id/process", handler.GetEvaluationProcess)

	req := httptest.NewRequest(http.MethodGet, "/evaluations/1/process", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	require.Equal(t, float64(0), response["code"])
	require.Equal(t, "该记录暂无链路详情", response["message"])

	data := response["data"].(map[string]interface{})
	require.Equal(t, false, data["trace_available"])
	require.Nil(t, data["trace"])
	require.Equal(t, "missing", data["status"])
}
