package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"resume-service/models"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type resumeTestJob struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `gorm:"size:200" json:"title"`
}

func (resumeTestJob) TableName() string {
	return "jobs"
}

func setupResumeTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Resume{}, &resumeTestJob{})
	return db
}

func setupResumeRouter(handler *ResumeHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.PUT("/resumes/:id/job", handler.UpdateResumeJob)
	return r
}

func TestUpdateResumeJob_Success(t *testing.T) {
	db := setupResumeTestDB()
	handler := NewResumeHandler(db)
	router := setupResumeRouter(handler)

	db.Create(&resumeTestJob{ID: 2, Title: "高级 Go 开发工程师"})
	db.Create(&models.Resume{ID: 1, FileName: "resume.pdf", Status: "pending"})

	body, _ := json.Marshal(map[string]interface{}{"job_id": 2})
	req, _ := http.NewRequest("PUT", "/resumes/1/job", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "高级 Go 开发工程师", data["job_title"])

	var saved models.Resume
	err := db.First(&saved, 1).Error
	assert.NoError(t, err)
	if assert.NotNil(t, saved.JobID) {
		assert.Equal(t, uint(2), *saved.JobID)
	}
}

func TestUpdateResumeJob_JobNotFound(t *testing.T) {
	db := setupResumeTestDB()
	handler := NewResumeHandler(db)
	router := setupResumeRouter(handler)

	db.Create(&models.Resume{ID: 1, FileName: "resume.pdf", Status: "pending"})

	body, _ := json.Marshal(map[string]interface{}{"job_id": 999})
	req, _ := http.NewRequest("PUT", "/resumes/1/job", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
