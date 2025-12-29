package handlers

import (
	"net/http"
	"strconv"
	"time"

	"common/elasticsearch"

	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	logService *elasticsearch.LogService
}

func NewLogHandler() *LogHandler {
	return &LogHandler{
		logService: elasticsearch.NewLogService("log-service"),
	}
}

// QueryLogs 查询日志
func (h *LogHandler) QueryLogs(c *gin.Context) {
	params := &elasticsearch.QueryParams{}

	// 解析分页参数
	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		params.Page = page
	}
	if pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "20")); err == nil {
		params.PageSize = pageSize
	}

	// 解析筛选参数
	params.Service = c.Query("service")
	params.Username = c.Query("username")
	params.Method = c.Query("method")
	params.Path = c.Query("path")
	params.Action = c.Query("action")
	params.Module = c.Query("module")
	params.Level = c.Query("level")
	params.Keyword = c.Query("keyword")

	if userID, err := strconv.ParseUint(c.Query("user_id"), 10, 32); err == nil {
		params.UserID = uint(userID)
	}

	// 解析时间范围
	if startTime := c.Query("start_time"); startTime != "" {
		if t, err := time.Parse(time.RFC3339, startTime); err == nil {
			params.StartTime = t
		} else if t, err := time.Parse("2006-01-02", startTime); err == nil {
			params.StartTime = t
		}
	}
	if endTime := c.Query("end_time"); endTime != "" {
		if t, err := time.Parse(time.RFC3339, endTime); err == nil {
			params.EndTime = t
		} else if t, err := time.Parse("2006-01-02", endTime); err == nil {
			params.EndTime = t.Add(24*time.Hour - time.Second)
		}
	}

	result, err := h.logService.Query(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询日志失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"total":     result.Total,
			"logs":      result.Logs,
			"page":      params.Page,
			"page_size": params.PageSize,
		},
	})
}

// GetStats 获取统计信息
func (h *LogHandler) GetStats(c *gin.Context) {
	// 默认统计最近24小时
	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour)

	if start := c.Query("start_time"); start != "" {
		if t, err := time.Parse(time.RFC3339, start); err == nil {
			startTime = t
		} else if t, err := time.Parse("2006-01-02", start); err == nil {
			startTime = t
		}
	}
	if end := c.Query("end_time"); end != "" {
		if t, err := time.Parse(time.RFC3339, end); err == nil {
			endTime = t
		} else if t, err := time.Parse("2006-01-02", end); err == nil {
			endTime = t.Add(24*time.Hour - time.Second)
		}
	}

	stats, err := h.logService.GetStats(startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取统计失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    stats,
	})
}

// GetServices 获取服务列表
func (h *LogHandler) GetServices(c *gin.Context) {
	services := []string{
		"user-service",
		"job-service",
		"resume-service",
		"interview-service",
		"message-service",
		"talent-service",
		"log-service",
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    services,
	})
}

// GetActions 获取操作类型列表
func (h *LogHandler) GetActions(c *gin.Context) {
	actions := []map[string]string{
		{"value": "登录", "label": "登录"},
		{"value": "注册", "label": "注册"},
		{"value": "查询", "label": "查询"},
		{"value": "查看", "label": "查看"},
		{"value": "创建", "label": "创建"},
		{"value": "更新", "label": "更新"},
		{"value": "删除", "label": "删除"},
		{"value": "上传", "label": "上传"},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    actions,
	})
}

// Cleanup 清理旧日志
func (h *LogHandler) Cleanup(c *gin.Context) {
	// 默认清理30天前的日志
	days := 30
	if d, err := strconv.Atoi(c.DefaultQuery("days", "30")); err == nil {
		days = d
	}

	before := time.Now().AddDate(0, 0, -days)
	if err := h.logService.DeleteOldLogs(before); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "清理日志失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "清理成功",
	})
}
