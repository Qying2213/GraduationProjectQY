package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"common/middleware"
	"talent-service/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTalentTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&models.Talent{})
	return db
}

func setupTalentRouter(handler *TalentHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	authMiddleware := func(c *gin.Context) {
		role := c.GetHeader("X-Test-Role")
		if role == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Set("role", role)
		c.Next()
	}

	api := r.Group("/api/v1/talents")
	api.Use(authMiddleware, middleware.RoleAuth("admin", "hr", "hr_manager", "recruiter"))
	{
		api.POST("", handler.CreateTalent)
		api.GET("", handler.ListTalents)
		api.GET("/:id", handler.GetTalent)
		api.PUT("/:id", handler.UpdateTalent)
		api.DELETE("/:id", handler.DeleteTalent)
	}

	return r
}

func TestTalentAccessControl(t *testing.T) {
	db := setupTalentTestDB()
	handler := NewTalentHandler(db)
	router := setupTalentRouter(handler)

	t.Run("未登录访问应返回401", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/talents", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("candidate 角色访问应返回403", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/talents", nil)
		req.Header.Set("X-Test-Role", "candidate")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("hr 角色访问列表应返回200", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/talents", nil)
		req.Header.Set("X-Test-Role", "hr")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestTalentCRUDFlowByHR(t *testing.T) {
	db := setupTalentTestDB()
	handler := NewTalentHandler(db)
	router := setupTalentRouter(handler)

	createBody := map[string]interface{}{
		"name":       "张三",
		"email":      "zhangsan@example.com",
		"phone":      "13800138000",
		"status":     "active",
		"education":  "本科",
		"location":   "深圳",
		"salary":     "20K-30K",
		"summary":    "Go后端开发工程师",
		"experience": 3,
	}

	body, _ := json.Marshal(createBody)
	createReq, _ := http.NewRequest("POST", "/api/v1/talents", bytes.NewBuffer(body))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("X-Test-Role", "hr")
	createResp := httptest.NewRecorder()
	router.ServeHTTP(createResp, createReq)

	assert.Equal(t, http.StatusCreated, createResp.Code)

	var createResult map[string]interface{}
	_ = json.Unmarshal(createResp.Body.Bytes(), &createResult)
	assert.Equal(t, float64(0), createResult["code"])

	createdTalent := createResult["data"].(map[string]interface{})
	talentID := int(createdTalent["id"].(float64))

	t.Run("列表应包含新建人才", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/talents?page=1&page_size=10", nil)
		req.Header.Set("X-Test-Role", "hr")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].(map[string]interface{})
		talents := data["talents"].([]interface{})
		assert.Len(t, talents, 1)
	})

	t.Run("详情应返回正确人才", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/talents/"+itoa(talentID), nil)
		req.Header.Set("X-Test-Role", "hr")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].(map[string]interface{})
		assert.Equal(t, "张三", data["name"])
	})

	t.Run("更新后应持久化", func(t *testing.T) {
		updateBody, _ := json.Marshal(map[string]interface{}{
			"name":     "张三-更新",
			"location": "广州",
			"summary":  "更新后的简介",
		})
		req, _ := http.NewRequest("PUT", "/api/v1/talents/"+itoa(talentID), bytes.NewBuffer(updateBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Test-Role", "hr_manager")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var saved models.Talent
		err := db.First(&saved, talentID).Error
		assert.NoError(t, err)
		assert.Equal(t, "张三-更新", saved.Name)
		assert.Equal(t, "广州", saved.Location)
		assert.Equal(t, "更新后的简介", saved.Summary)
	})

	t.Run("删除后详情应返回404", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/v1/talents/"+itoa(talentID), nil)
		req.Header.Set("X-Test-Role", "admin")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		getReq, _ := http.NewRequest("GET", "/api/v1/talents/"+itoa(talentID), nil)
		getReq.Header.Set("X-Test-Role", "admin")
		getResp := httptest.NewRecorder()
		router.ServeHTTP(getResp, getReq)

		assert.Equal(t, http.StatusNotFound, getResp.Code)
	})
}

func itoa(v int) string {
	return strconv.Itoa(v)
}
