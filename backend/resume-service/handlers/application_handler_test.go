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

// Job model for testing (simplified)
type Job struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `gorm:"size:200" json:"title"`
	CreatedBy uint      `json:"created_by"`
}

// Talent model for testing (simplified)
type Talent struct {
	ID   uint   `gorm:"primarykey" json:"id"`
	Name string `gorm:"size:100" json:"name"`
}

// Message model for testing
type Message struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	SenderID   *uint     `json:"sender_id"`
	ReceiverID uint      `json:"receiver_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Type       string    `json:"type"`
	IsRead     bool      `json:"is_read"`
}

// 创建测试数据库
func setupApplicationTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	// 创建所有需要的表
	db.AutoMigrate(&models.Resume{}, &models.Application{}, &Job{}, &Talent{}, &Message{})
	return db
}

// 创建测试路由
func setupApplicationRouter(handler *ResumeHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/applications", handler.CreateApplication)
	r.GET("/applications", handler.ListApplications)
	r.PUT("/applications/:id", handler.UpdateApplication)
	return r
}

// TestCreateApplication_Success 测试成功创建申请
// **Validates: Requirements 2.2**
func TestCreateApplication_Success(t *testing.T) {
	db := setupApplicationTestDB()
	handler := NewResumeHandler(db)
	router := setupApplicationRouter(handler)

	// 创建测试数据：HR用户、职位、求职者、简历
	hrUserID := uint(1)
	job := Job{ID: 1, Title: "Go Developer", CreatedBy: hrUserID}
	db.Create(&job)

	talent := Talent{ID: 1, Name: "张三"}
	db.Create(&talent)

	// 创建简历（满足 Requirement 2.6）
	talentID := uint(1)
	resume := models.Resume{ID: 1, TalentID: &talentID, FileName: "resume.pdf"}
	db.Create(&resume)

	t.Run("成功创建申请，状态应为pending", func(t *testing.T) {
		body := map[string]interface{}{
			"job_id":       1,
			"talent_id":    1,
			"resume_id":    1,
			"cover_letter": "I am interested in this position.",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/applications", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, float64(0), response["code"])

		// 验证申请状态为 "pending" (Requirement 2.2)
		data := response["data"].(map[string]interface{})
		assert.Equal(t, "pending", data["status"])
	})
}

// TestCreateApplication_DuplicatePrevention 测试重复申请检查
// **Validates: Requirements 2.4**
// **Property 3: Duplicate Application Prevention**
func TestCreateApplication_DuplicatePrevention(t *testing.T) {
	db := setupApplicationTestDB()
	handler := NewResumeHandler(db)
	router := setupApplicationRouter(handler)

	// 创建测试数据
	job := Job{ID: 1, Title: "Go Developer", CreatedBy: 1}
	db.Create(&job)

	talent := Talent{ID: 1, Name: "张三"}
	db.Create(&talent)

	talentID := uint(1)
	resume := models.Resume{ID: 1, TalentID: &talentID, FileName: "resume.pdf"}
	db.Create(&resume)

	// 创建第一个申请
	existingApp := models.Application{
		JobID:    1,
		TalentID: 1,
		ResumeID: 1,
		Status:   "pending",
	}
	db.Create(&existingApp)

	t.Run("重复申请应该被拒绝并返回错误码1001", func(t *testing.T) {
		body := map[string]interface{}{
			"job_id":    1,
			"talent_id": 1,
			"resume_id": 1,
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/applications", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, float64(1001), response["code"])
		assert.Equal(t, "您已投递过该职位", response["message"])
	})
}

// TestCreateApplication_NoResume 测试无简历时的验证
// **Validates: Requirements 2.6**
// **Property 4: Resume Requirement for Application**
func TestCreateApplication_NoResume(t *testing.T) {
	db := setupApplicationTestDB()
	handler := NewResumeHandler(db)
	router := setupApplicationRouter(handler)

	// 创建测试数据：职位和求职者，但不创建简历
	job := Job{ID: 1, Title: "Go Developer", CreatedBy: 1}
	db.Create(&job)

	talent := Talent{ID: 1, Name: "张三"}
	db.Create(&talent)

	t.Run("没有简历的求职者申请应该被拒绝并返回错误码1002", func(t *testing.T) {
		body := map[string]interface{}{
			"job_id":    1,
			"talent_id": 1,
			"resume_id": 0,
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/applications", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, float64(1002), response["code"])
		assert.Equal(t, "请先上传简历", response["message"])
	})
}

// TestCreateApplication_NotificationSent 测试HR通知发送
// **Validates: Requirements 2.3**
func TestCreateApplication_NotificationSent(t *testing.T) {
	db := setupApplicationTestDB()
	handler := NewResumeHandler(db)
	router := setupApplicationRouter(handler)

	// 创建测试数据
	hrUserID := uint(2)
	job := Job{ID: 1, Title: "Go Developer", CreatedBy: hrUserID}
	db.Create(&job)

	talent := Talent{ID: 1, Name: "李四"}
	db.Create(&talent)

	talentID := uint(1)
	resume := models.Resume{ID: 1, TalentID: &talentID, FileName: "resume.pdf"}
	db.Create(&resume)

	t.Run("创建申请后应该发送通知给HR", func(t *testing.T) {
		body := map[string]interface{}{
			"job_id":    1,
			"talent_id": 1,
			"resume_id": 1,
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/applications", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// 验证通知已发送给HR
		var message Message
		result := db.Where("receiver_id = ? AND type = ?", hrUserID, "application").First(&message)
		assert.NoError(t, result.Error)
		assert.Equal(t, hrUserID, message.ReceiverID)
		assert.Equal(t, "新的职位申请", message.Title)
		assert.Contains(t, message.Content, "李四")
		assert.Contains(t, message.Content, "Go Developer")
	})
}

// TestCreateApplication_MissingJobID 测试缺少job_id
func TestCreateApplication_MissingJobID(t *testing.T) {
	db := setupApplicationTestDB()
	handler := NewResumeHandler(db)
	router := setupApplicationRouter(handler)

	t.Run("缺少job_id应该返回错误", func(t *testing.T) {
		body := map[string]interface{}{
			"talent_id": 1,
			"resume_id": 1,
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/applications", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestCreateApplication_MissingTalentID 测试缺少talent_id
func TestCreateApplication_MissingTalentID(t *testing.T) {
	db := setupApplicationTestDB()
	handler := NewResumeHandler(db)
	router := setupApplicationRouter(handler)

	t.Run("缺少talent_id应该返回错误", func(t *testing.T) {
		body := map[string]interface{}{
			"job_id":    1,
			"resume_id": 1,
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/applications", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestCreateApplication_DifferentJobsSameCandidate 测试同一求职者可以申请不同职位
func TestCreateApplication_DifferentJobsSameCandidate(t *testing.T) {
	db := setupApplicationTestDB()
	handler := NewResumeHandler(db)
	router := setupApplicationRouter(handler)

	// 创建测试数据
	job1 := Job{ID: 1, Title: "Go Developer", CreatedBy: 1}
	job2 := Job{ID: 2, Title: "Python Developer", CreatedBy: 1}
	db.Create(&job1)
	db.Create(&job2)

	talent := Talent{ID: 1, Name: "张三"}
	db.Create(&talent)

	talentID := uint(1)
	resume := models.Resume{ID: 1, TalentID: &talentID, FileName: "resume.pdf"}
	db.Create(&resume)

	// 创建第一个申请
	existingApp := models.Application{
		JobID:    1,
		TalentID: 1,
		ResumeID: 1,
		Status:   "pending",
	}
	db.Create(&existingApp)

	t.Run("同一求职者可以申请不同的职位", func(t *testing.T) {
		body := map[string]interface{}{
			"job_id":    2, // 不同的职位
			"talent_id": 1,
			"resume_id": 1,
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/applications", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, float64(0), response["code"])
	})
}

// TalentWithUser model for testing (with user_id)
type TalentWithUser struct {
	ID     uint   `gorm:"primarykey" json:"id"`
	Name   string `gorm:"size:100" json:"name"`
	UserID *uint  `json:"user_id"`
}

// TableName specifies the table name for TalentWithUser
func (TalentWithUser) TableName() string {
	return "talents"
}

// TestUpdateApplication_StatusChangeNotification 测试状态更新时发送通知给求职者
// **Validates: Requirements 3.2, 6.4**
// **Property 13: Status Update Notification**
func TestUpdateApplication_StatusChangeNotification(t *testing.T) {
	db := setupApplicationTestDB()
	// 重新创建talents表以包含user_id字段
	db.Migrator().DropTable(&Talent{})
	db.AutoMigrate(&TalentWithUser{})

	handler := NewResumeHandler(db)
	router := setupApplicationRouter(handler)

	// 创建测试数据
	job := Job{ID: 1, Title: "Go Developer", CreatedBy: 1}
	db.Create(&job)

	candidateUserID := uint(10)
	talent := TalentWithUser{ID: 1, Name: "王五", UserID: &candidateUserID}
	db.Create(&talent)

	talentID := uint(1)
	resume := models.Resume{ID: 1, TalentID: &talentID, FileName: "resume.pdf"}
	db.Create(&resume)

	// 创建申请
	app := models.Application{
		ID:       1,
		JobID:    1,
		TalentID: 1,
		ResumeID: 1,
		Status:   "pending",
	}
	db.Create(&app)

	t.Run("状态更新时应该发送通知给求职者", func(t *testing.T) {
		body := map[string]interface{}{
			"status": "viewed",
			"notes":  "简历已查看",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("PUT", "/applications/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, float64(0), response["code"])

		// 验证通知已发送给求职者
		var message Message
		result := db.Where("receiver_id = ? AND type = ?", candidateUserID, "application_status").First(&message)
		assert.NoError(t, result.Error)
		assert.Equal(t, candidateUserID, message.ReceiverID)
		assert.Equal(t, "申请状态更新", message.Title)
		assert.Contains(t, message.Content, "Go Developer")
		assert.Contains(t, message.Content, "已查看")
	})
}

// TestUpdateApplication_StatusChangeHistory 测试状态变更历史记录
// **Validates: Requirements 3.2, 6.4**
func TestUpdateApplication_StatusChangeHistory(t *testing.T) {
	db := setupApplicationTestDB()
	// 重新创建talents表以包含user_id字段
	db.Migrator().DropTable(&Talent{})
	db.AutoMigrate(&TalentWithUser{})

	handler := NewResumeHandler(db)
	router := setupApplicationRouter(handler)

	// 创建测试数据
	job := Job{ID: 1, Title: "Go Developer", CreatedBy: 1}
	db.Create(&job)

	candidateUserID := uint(10)
	talent := TalentWithUser{ID: 1, Name: "赵六", UserID: &candidateUserID}
	db.Create(&talent)

	talentID := uint(1)
	resume := models.Resume{ID: 1, TalentID: &talentID, FileName: "resume.pdf"}
	db.Create(&resume)

	// 创建申请
	app := models.Application{
		ID:       1,
		JobID:    1,
		TalentID: 1,
		ResumeID: 1,
		Status:   "pending",
		Notes:    "",
	}
	db.Create(&app)

	t.Run("状态变更应该记录在notes中", func(t *testing.T) {
		body := map[string]interface{}{
			"status": "interview",
			"notes":  "安排面试",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("PUT", "/applications/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 验证notes中包含状态变更历史
		var updatedApp models.Application
		db.First(&updatedApp, 1)
		assert.Contains(t, updatedApp.Notes, "待处理")
		assert.Contains(t, updatedApp.Notes, "面试中")
		assert.Contains(t, updatedApp.Notes, "安排面试")
	})
}

// TestUpdateApplication_MultipleStatusChanges 测试多次状态变更历史累积
// **Validates: Requirements 3.2, 6.4**
func TestUpdateApplication_MultipleStatusChanges(t *testing.T) {
	db := setupApplicationTestDB()
	// 重新创建talents表以包含user_id字段
	db.Migrator().DropTable(&Talent{})
	db.AutoMigrate(&TalentWithUser{})

	handler := NewResumeHandler(db)
	router := setupApplicationRouter(handler)

	// 创建测试数据
	job := Job{ID: 1, Title: "Go Developer", CreatedBy: 1}
	db.Create(&job)

	candidateUserID := uint(10)
	talent := TalentWithUser{ID: 1, Name: "钱七", UserID: &candidateUserID}
	db.Create(&talent)

	talentID := uint(1)
	resume := models.Resume{ID: 1, TalentID: &talentID, FileName: "resume.pdf"}
	db.Create(&resume)

	// 创建申请
	app := models.Application{
		ID:       1,
		JobID:    1,
		TalentID: 1,
		ResumeID: 1,
		Status:   "pending",
		Notes:    "",
	}
	db.Create(&app)

	t.Run("多次状态变更应该累积记录", func(t *testing.T) {
		// 第一次状态变更: pending -> viewed
		body1 := map[string]interface{}{
			"status": "viewed",
		}
		jsonBody1, _ := json.Marshal(body1)
		req1, _ := http.NewRequest("PUT", "/applications/1", bytes.NewBuffer(jsonBody1))
		req1.Header.Set("Content-Type", "application/json")
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// 第二次状态变更: viewed -> interview
		body2 := map[string]interface{}{
			"status": "interview",
			"notes":  "技术面试",
		}
		jsonBody2, _ := json.Marshal(body2)
		req2, _ := http.NewRequest("PUT", "/applications/1", bytes.NewBuffer(jsonBody2))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)

		// 验证notes中包含所有状态变更历史
		var updatedApp models.Application
		db.First(&updatedApp, 1)
		assert.Contains(t, updatedApp.Notes, "待处理")
		assert.Contains(t, updatedApp.Notes, "已查看")
		assert.Contains(t, updatedApp.Notes, "面试中")
		assert.Contains(t, updatedApp.Notes, "技术面试")

		// 验证发送了两条通知
		var messages []Message
		db.Where("receiver_id = ? AND type = ?", candidateUserID, "application_status").Find(&messages)
		assert.Equal(t, 2, len(messages))
	})
}

// TestUpdateApplication_InvalidStatus 测试无效状态值
// **Validates: Requirements 3.3**
func TestUpdateApplication_InvalidStatus(t *testing.T) {
	db := setupApplicationTestDB()
	handler := NewResumeHandler(db)
	router := setupApplicationRouter(handler)

	// 创建测试数据
	job := Job{ID: 1, Title: "Go Developer", CreatedBy: 1}
	db.Create(&job)

	talent := Talent{ID: 1, Name: "测试用户"}
	db.Create(&talent)

	// 创建申请
	app := models.Application{
		ID:       1,
		JobID:    1,
		TalentID: 1,
		ResumeID: 1,
		Status:   "pending",
	}
	db.Create(&app)

	t.Run("无效状态值应该被拒绝", func(t *testing.T) {
		body := map[string]interface{}{
			"status": "invalid_status",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("PUT", "/applications/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, float64(400), response["code"])
	})
}

// TestUpdateApplication_SameStatusNoNotification 测试相同状态不发送通知
// **Validates: Requirements 3.2, 6.4**
func TestUpdateApplication_SameStatusNoNotification(t *testing.T) {
	db := setupApplicationTestDB()
	// 重新创建talents表以包含user_id字段
	db.Migrator().DropTable(&Talent{})
	db.AutoMigrate(&TalentWithUser{})

	handler := NewResumeHandler(db)
	router := setupApplicationRouter(handler)

	// 创建测试数据
	job := Job{ID: 1, Title: "Go Developer", CreatedBy: 1}
	db.Create(&job)

	candidateUserID := uint(10)
	talent := TalentWithUser{ID: 1, Name: "孙八", UserID: &candidateUserID}
	db.Create(&talent)

	talentID := uint(1)
	resume := models.Resume{ID: 1, TalentID: &talentID, FileName: "resume.pdf"}
	db.Create(&resume)

	// 创建申请
	app := models.Application{
		ID:       1,
		JobID:    1,
		TalentID: 1,
		ResumeID: 1,
		Status:   "viewed",
	}
	db.Create(&app)

	t.Run("相同状态更新不应该发送通知", func(t *testing.T) {
		body := map[string]interface{}{
			"status": "viewed", // 相同状态
			"notes":  "再次查看",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("PUT", "/applications/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 验证没有发送通知
		var count int64
		db.Model(&Message{}).Where("receiver_id = ? AND type = ?", candidateUserID, "application_status").Count(&count)
		assert.Equal(t, int64(0), count)
	})
}

// TestUpdateApplication_NotFound 测试申请不存在
func TestUpdateApplication_NotFound(t *testing.T) {
	db := setupApplicationTestDB()
	handler := NewResumeHandler(db)
	router := setupApplicationRouter(handler)

	t.Run("申请不存在应该返回404", func(t *testing.T) {
		body := map[string]interface{}{
			"status": "viewed",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("PUT", "/applications/999", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
