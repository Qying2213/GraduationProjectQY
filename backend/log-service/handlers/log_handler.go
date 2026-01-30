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

// LogViewerPage 日志查看页面（无需登录）
func (h *LogHandler) LogViewerPage(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>系统日志查看器</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f7fa; min-height: 100vh; }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 20px; }
        .header h1 { font-size: 24px; margin-bottom: 5px; }
        .header p { opacity: 0.8; font-size: 14px; }
        .container { max-width: 1400px; margin: 0 auto; padding: 20px; }
        .filters { background: white; border-radius: 8px; padding: 20px; margin-bottom: 20px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
        .filter-row { display: flex; gap: 15px; flex-wrap: wrap; align-items: center; }
        .filter-item { display: flex; flex-direction: column; gap: 5px; }
        .filter-item label { font-size: 12px; color: #666; font-weight: 500; }
        .filter-item select, .filter-item input { padding: 8px 12px; border: 1px solid #ddd; border-radius: 6px; font-size: 14px; min-width: 150px; }
        .filter-item select:focus, .filter-item input:focus { outline: none; border-color: #667eea; }
        .btn { padding: 10px 20px; border: none; border-radius: 6px; cursor: pointer; font-size: 14px; font-weight: 500; }
        .btn-primary { background: #667eea; color: white; }
        .btn-primary:hover { background: #5a6fd6; }
        .btn-secondary { background: #e0e0e0; color: #333; }
        .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap: 15px; margin-bottom: 20px; }
        .stat-card { background: white; border-radius: 8px; padding: 15px; text-align: center; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
        .stat-card .value { font-size: 28px; font-weight: bold; color: #333; }
        .stat-card .label { font-size: 12px; color: #666; margin-top: 5px; }
        .stat-card.info .value { color: #667eea; }
        .stat-card.warn .value { color: #f59e0b; }
        .stat-card.error .value { color: #ef4444; }
        .log-table { background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
        table { width: 100%; border-collapse: collapse; }
        th { background: #f8f9fa; padding: 12px; text-align: left; font-size: 13px; color: #666; border-bottom: 2px solid #eee; }
        td { padding: 12px; border-bottom: 1px solid #eee; font-size: 13px; }
        tr:hover { background: #f8f9fa; }
        .level { padding: 3px 8px; border-radius: 4px; font-size: 11px; font-weight: 600; }
        .level-info { background: #e0f2fe; color: #0369a1; }
        .level-warn { background: #fef3c7; color: #b45309; }
        .level-error { background: #fee2e2; color: #dc2626; }
        .method { padding: 3px 8px; border-radius: 4px; font-size: 11px; font-weight: 600; }
        .method-GET { background: #d1fae5; color: #059669; }
        .method-POST { background: #dbeafe; color: #2563eb; }
        .method-PUT { background: #fef3c7; color: #d97706; }
        .method-DELETE { background: #fee2e2; color: #dc2626; }
        .status { font-weight: 600; }
        .status-2xx { color: #059669; }
        .status-4xx { color: #d97706; }
        .status-5xx { color: #dc2626; }
        .pagination { display: flex; justify-content: center; gap: 10px; margin-top: 20px; }
        .page-btn { padding: 8px 15px; border: 1px solid #ddd; border-radius: 6px; background: white; cursor: pointer; }
        .page-btn:hover { background: #f5f5f5; }
        .page-btn.active { background: #667eea; color: white; border-color: #667eea; }
        .empty { text-align: center; padding: 60px; color: #999; }
        .loading { text-align: center; padding: 40px; color: #666; }
        .path { max-width: 300px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
        .time { color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>🔍 系统日志查看器</h1>
        <p>实时查看系统操作日志，无需登录</p>
    </div>
    <div class="container">
        <div class="stats" id="stats">
            <div class="stat-card info"><div class="value" id="total">-</div><div class="label">总日志数</div></div>
            <div class="stat-card info"><div class="value" id="info-count">-</div><div class="label">INFO</div></div>
            <div class="stat-card warn"><div class="value" id="warn-count">-</div><div class="label">WARN</div></div>
            <div class="stat-card error"><div class="value" id="error-count">-</div><div class="label">ERROR</div></div>
        </div>
        <div class="filters">
            <div class="filter-row">
                <div class="filter-item">
                    <label>服务</label>
                    <select id="service">
                        <option value="">全部服务</option>
                        <option value="user-service">用户服务</option>
                        <option value="job-service">职位服务</option>
                        <option value="resume-service">简历服务</option>
                        <option value="interview-service">面试服务</option>
                        <option value="message-service">消息服务</option>
                        <option value="talent-service">人才服务</option>
                        <option value="recommendation-service">推荐服务</option>
                    </select>
                </div>
                <div class="filter-item">
                    <label>级别</label>
                    <select id="level">
                        <option value="">全部级别</option>
                        <option value="info">INFO</option>
                        <option value="warn">WARN</option>
                        <option value="error">ERROR</option>
                    </select>
                </div>
                <div class="filter-item">
                    <label>方法</label>
                    <select id="method">
                        <option value="">全部方法</option>
                        <option value="GET">GET</option>
                        <option value="POST">POST</option>
                        <option value="PUT">PUT</option>
                        <option value="DELETE">DELETE</option>
                    </select>
                </div>
                <div class="filter-item">
                    <label>关键词</label>
                    <input type="text" id="keyword" placeholder="搜索路径或内容">
                </div>
                <div class="filter-item" style="justify-content: flex-end;">
                    <label>&nbsp;</label>
                    <div style="display:flex;gap:10px;">
                        <button class="btn btn-primary" onclick="loadLogs()">🔍 查询</button>
                        <button class="btn btn-secondary" onclick="resetFilters()">重置</button>
                        <button class="btn btn-secondary" onclick="loadLogs()">🔄 刷新</button>
                    </div>
                </div>
            </div>
        </div>
        <div class="log-table">
            <table>
                <thead>
                    <tr>
                        <th>时间</th>
                        <th>级别</th>
                        <th>服务</th>
                        <th>方法</th>
                        <th>路径</th>
                        <th>状态</th>
                        <th>耗时</th>
                        <th>操作</th>
                    </tr>
                </thead>
                <tbody id="log-body">
                    <tr><td colspan="8" class="loading">加载中...</td></tr>
                </tbody>
            </table>
        </div>
        <div class="pagination" id="pagination"></div>
    </div>
    <script>
        let currentPage = 1;
        const pageSize = 20;

        async function loadLogs() {
            const params = new URLSearchParams({
                page: currentPage,
                page_size: pageSize,
                service: document.getElementById('service').value,
                level: document.getElementById('level').value,
                method: document.getElementById('method').value,
                keyword: document.getElementById('keyword').value
            });
            
            try {
                const res = await fetch('/api/v1/logs?' + params);
                const data = await res.json();
                
                if (data.code === 0) {
                    renderLogs(data.data.logs || []);
                    renderPagination(data.data.total);
                    updateStats(data.data.logs || []);
                }
            } catch (e) {
                document.getElementById('log-body').innerHTML = '<tr><td colspan="8" class="empty">加载失败: ' + e.message + '</td></tr>';
            }
        }

        function renderLogs(logs) {
            const tbody = document.getElementById('log-body');
            if (!logs.length) {
                tbody.innerHTML = '<tr><td colspan="8" class="empty">暂无日志数据</td></tr>';
                return;
            }
            
            tbody.innerHTML = logs.map(log => {
                const time = new Date(log.timestamp).toLocaleString('zh-CN');
                const levelClass = 'level-' + (log.level || 'info');
                const methodClass = 'method-' + log.method;
                const statusClass = log.status_code < 300 ? 'status-2xx' : (log.status_code < 500 ? 'status-4xx' : 'status-5xx');
                
                return '<tr>' +
                    '<td class="time">' + time + '</td>' +
                    '<td><span class="level ' + levelClass + '">' + (log.level || 'info').toUpperCase() + '</span></td>' +
                    '<td>' + (log.service || '-') + '</td>' +
                    '<td><span class="method ' + methodClass + '">' + log.method + '</span></td>' +
                    '<td class="path" title="' + log.path + '">' + log.path + '</td>' +
                    '<td class="status ' + statusClass + '">' + log.status_code + '</td>' +
                    '<td>' + log.duration + 'ms</td>' +
                    '<td>' + (log.action || '-') + '</td>' +
                '</tr>';
            }).join('');
        }

        function renderPagination(total) {
            const totalPages = Math.ceil(total / pageSize);
            const pagination = document.getElementById('pagination');
            
            if (totalPages <= 1) {
                pagination.innerHTML = '';
                return;
            }
            
            let html = '';
            if (currentPage > 1) html += '<button class="page-btn" onclick="goPage(' + (currentPage-1) + ')">上一页</button>';
            
            for (let i = Math.max(1, currentPage-2); i <= Math.min(totalPages, currentPage+2); i++) {
                html += '<button class="page-btn' + (i === currentPage ? ' active' : '') + '" onclick="goPage(' + i + ')">' + i + '</button>';
            }
            
            if (currentPage < totalPages) html += '<button class="page-btn" onclick="goPage(' + (currentPage+1) + ')">下一页</button>';
            
            pagination.innerHTML = html;
        }

        function updateStats(logs) {
            document.getElementById('total').textContent = logs.length;
            document.getElementById('info-count').textContent = logs.filter(l => l.level === 'info').length;
            document.getElementById('warn-count').textContent = logs.filter(l => l.level === 'warn').length;
            document.getElementById('error-count').textContent = logs.filter(l => l.level === 'error').length;
        }

        function goPage(page) {
            currentPage = page;
            loadLogs();
        }

        function resetFilters() {
            document.getElementById('service').value = '';
            document.getElementById('level').value = '';
            document.getElementById('method').value = '';
            document.getElementById('keyword').value = '';
            currentPage = 1;
            loadLogs();
        }

        // 初始加载
        loadLogs();
        // 每30秒自动刷新
        setInterval(loadLogs, 30000);
    </script>
</body>
</html>`

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}
