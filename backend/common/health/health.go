package health

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthStatus 健康状态
type HealthStatus struct {
	Status    string           `json:"status"`
	Service   string           `json:"service"`
	Version   string           `json:"version"`
	Uptime    string           `json:"uptime"`
	Timestamp time.Time        `json:"timestamp"`
	Checks    map[string]Check `json:"checks"`
	System    SystemInfo       `json:"system"`
}

// Check 检查项
type Check struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// SystemInfo 系统信息
type SystemInfo struct {
	GoVersion    string `json:"go_version"`
	NumGoroutine int    `json:"num_goroutine"`
	NumCPU       int    `json:"num_cpu"`
	MemAlloc     uint64 `json:"mem_alloc_mb"`
}

var startTime = time.Now()

// Handler 健康检查处理器
type Handler struct {
	serviceName string
	version     string
	db          *gorm.DB
}

// NewHandler 创建健康检查处理器
func NewHandler(serviceName, version string, db *gorm.DB) *Handler {
	return &Handler{
		serviceName: serviceName,
		version:     version,
		db:          db,
	}
}

// HealthCheck 健康检查端点
func (h *Handler) HealthCheck(c *gin.Context) {
	checks := make(map[string]Check)

	// 数据库检查
	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err != nil {
			checks["database"] = Check{Status: "unhealthy", Message: err.Error()}
		} else if err := sqlDB.Ping(); err != nil {
			checks["database"] = Check{Status: "unhealthy", Message: err.Error()}
		} else {
			checks["database"] = Check{Status: "healthy"}
		}
	}

	// 系统信息
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	status := "healthy"
	for _, check := range checks {
		if check.Status != "healthy" {
			status = "unhealthy"
			break
		}
	}

	health := HealthStatus{
		Status:    status,
		Service:   h.serviceName,
		Version:   h.version,
		Uptime:    time.Since(startTime).String(),
		Timestamp: time.Now(),
		Checks:    checks,
		System: SystemInfo{
			GoVersion:    runtime.Version(),
			NumGoroutine: runtime.NumGoroutine(),
			NumCPU:       runtime.NumCPU(),
			MemAlloc:     m.Alloc / 1024 / 1024,
		},
	}

	statusCode := http.StatusOK
	if status != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, health)
}

// LivenessCheck 存活检查
func (h *Handler) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}

// ReadinessCheck 就绪检查
func (h *Handler) ReadinessCheck(c *gin.Context) {
	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err != nil || sqlDB.Ping() != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"reason": "database connection failed",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}
