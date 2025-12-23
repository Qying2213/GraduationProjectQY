package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-service/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 创建测试数据库
func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{})
	return db
}

// 创建测试路由
func setupRouter(handler *UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
	return r
}

func TestRegister(t *testing.T) {
	db := setupTestDB()
	handler := NewUserHandler(db)
	router := setupRouter(handler)

	t.Run("成功注册新用户", func(t *testing.T) {
		body := RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, float64(0), response["code"])
	})

	t.Run("用户名已存在应该返回错误", func(t *testing.T) {
		// 先创建一个用户
		user := models.User{
			Username: "existuser",
			Email:    "exist@example.com",
			Status:   "active",
		}
		user.HashPassword("password123")
		db.Create(&user)

		body := RegisterRequest{
			Username: "existuser",
			Email:    "new@example.com",
			Password: "password123",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("缺少必填字段应该返回错误", func(t *testing.T) {
		body := map[string]string{
			"username": "testuser2",
			// 缺少 email 和 password
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLogin(t *testing.T) {
	db := setupTestDB()
	handler := NewUserHandler(db)
	router := setupRouter(handler)

	// 创建测试用户
	user := models.User{
		Username: "loginuser",
		Email:    "login@example.com",
		Role:     "admin",
		Status:   "active",
	}
	user.HashPassword("password123")
	db.Create(&user)

	t.Run("正确的凭据应该登录成功", func(t *testing.T) {
		body := LoginRequest{
			Username: "loginuser",
			Password: "password123",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, float64(0), response["code"])

		data := response["data"].(map[string]interface{})
		assert.NotEmpty(t, data["token"])
		assert.NotNil(t, data["user"])
	})

	t.Run("使用邮箱登录应该成功", func(t *testing.T) {
		body := LoginRequest{
			Username: "login@example.com",
			Password: "password123",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("错误的密码应该返回401", func(t *testing.T) {
		body := LoginRequest{
			Username: "loginuser",
			Password: "wrongpassword",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("不存在的用户应该返回401", func(t *testing.T) {
		body := LoginRequest{
			Username: "nonexistent",
			Password: "password123",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestGenerateToken(t *testing.T) {
	t.Run("应该成功生成JWT token", func(t *testing.T) {
		token, err := generateToken(1, "testuser", "admin")

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}
