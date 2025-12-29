package elasticsearch

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	// 操作日志索引名
	OperationLogIndex = "operation_logs"
)

// OperationLog 操作日志结构
type OperationLog struct {
	ID           string                 `json:"id"`
	Timestamp    time.Time              `json:"timestamp"`
	Service      string                 `json:"service"`       // 服务名称
	UserID       uint                   `json:"user_id"`       // 用户ID
	Username     string                 `json:"username"`      // 用户名
	IP           string                 `json:"ip"`            // 客户端IP
	Method       string                 `json:"method"`        // HTTP方法
	Path         string                 `json:"path"`          // 请求路径
	Query        string                 `json:"query"`         // 查询参数
	StatusCode   int                    `json:"status_code"`   // 响应状态码
	Duration     int64                  `json:"duration"`      // 请求耗时(ms)
	RequestBody  string                 `json:"request_body"`  // 请求体
	ResponseBody string                 `json:"response_body"` // 响应体
	UserAgent    string                 `json:"user_agent"`    // 用户代理
	Action       string                 `json:"action"`        // 操作类型
	Module       string                 `json:"module"`        // 模块名称
	Description  string                 `json:"description"`   // 操作描述
	Extra        map[string]interface{} `json:"extra"`         // 额外信息
	Level        string                 `json:"level"`         // 日志级别: info, warn, error
	ErrorMsg     string                 `json:"error_msg"`     // 错误信息
}

// OperationLogMapping ES索引映射
var OperationLogMapping = map[string]interface{}{
	"settings": map[string]interface{}{
		"number_of_shards":   1,
		"number_of_replicas": 0,
		"index": map[string]interface{}{
			"lifecycle": map[string]interface{}{
				"name": "logs_policy",
			},
		},
	},
	"mappings": map[string]interface{}{
		"properties": map[string]interface{}{
			"id":            map[string]string{"type": "keyword"},
			"timestamp":     map[string]string{"type": "date"},
			"service":       map[string]string{"type": "keyword"},
			"user_id":       map[string]string{"type": "integer"},
			"username":      map[string]string{"type": "keyword"},
			"ip":            map[string]string{"type": "ip"},
			"method":        map[string]string{"type": "keyword"},
			"path":          map[string]string{"type": "keyword"},
			"query":         map[string]string{"type": "text"},
			"status_code":   map[string]string{"type": "integer"},
			"duration":      map[string]string{"type": "long"},
			"request_body":  map[string]string{"type": "text"},
			"response_body": map[string]string{"type": "text"},
			"user_agent":    map[string]string{"type": "text"},
			"action":        map[string]string{"type": "keyword"},
			"module":        map[string]string{"type": "keyword"},
			"description":   map[string]string{"type": "text"},
			"level":         map[string]string{"type": "keyword"},
			"error_msg":     map[string]string{"type": "text"},
		},
	},
}

// LogService 日志服务
type LogService struct {
	serviceName string
}

// NewLogService 创建日志服务
func NewLogService(serviceName string) *LogService {
	// 确保索引存在
	CreateIndex(OperationLogIndex, OperationLogMapping)
	return &LogService{serviceName: serviceName}
}

// Log 记录操作日志
func (s *LogService) Log(log *OperationLog) error {
	if log.ID == "" {
		log.ID = uuid.New().String()
	}
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}
	if log.Service == "" {
		log.Service = s.serviceName
	}
	if log.Level == "" {
		log.Level = "info"
	}

	return IndexDocument(OperationLogIndex, log.ID, log)
}

// LogAsync 异步记录日志
func (s *LogService) LogAsync(log *OperationLog) {
	go func() {
		if err := s.Log(log); err != nil {
			fmt.Printf("Failed to log operation: %v\n", err)
		}
	}()
}

// QueryParams 查询参数
type QueryParams struct {
	Page      int       `form:"page" json:"page"`
	PageSize  int       `form:"page_size" json:"page_size"`
	Service   string    `form:"service" json:"service"`
	UserID    uint      `form:"user_id" json:"user_id"`
	Username  string    `form:"username" json:"username"`
	Method    string    `form:"method" json:"method"`
	Path      string    `form:"path" json:"path"`
	Action    string    `form:"action" json:"action"`
	Module    string    `form:"module" json:"module"`
	Level     string    `form:"level" json:"level"`
	StartTime time.Time `form:"start_time" json:"start_time"`
	EndTime   time.Time `form:"end_time" json:"end_time"`
	Keyword   string    `form:"keyword" json:"keyword"`
}

// QueryResult 查询结果
type QueryResult struct {
	Total int64          `json:"total"`
	Logs  []OperationLog `json:"logs"`
}

// Query 查询日志
func (s *LogService) Query(params *QueryParams) (*QueryResult, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}

	// 构建查询
	must := []map[string]interface{}{}

	if params.Service != "" {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{"service": params.Service},
		})
	}
	if params.UserID > 0 {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{"user_id": params.UserID},
		})
	}
	if params.Username != "" {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{"username": params.Username},
		})
	}
	if params.Method != "" {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{"method": params.Method},
		})
	}
	if params.Path != "" {
		must = append(must, map[string]interface{}{
			"wildcard": map[string]interface{}{"path": "*" + params.Path + "*"},
		})
	}
	if params.Action != "" {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{"action": params.Action},
		})
	}
	if params.Module != "" {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{"module": params.Module},
		})
	}
	if params.Level != "" {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{"level": params.Level},
		})
	}
	if params.Keyword != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  params.Keyword,
				"fields": []string{"path", "description", "request_body", "response_body"},
			},
		})
	}

	// 时间范围
	if !params.StartTime.IsZero() || !params.EndTime.IsZero() {
		rangeQuery := map[string]interface{}{}
		if !params.StartTime.IsZero() {
			rangeQuery["gte"] = params.StartTime.Format(time.RFC3339)
		}
		if !params.EndTime.IsZero() {
			rangeQuery["lte"] = params.EndTime.Format(time.RFC3339)
		}
		must = append(must, map[string]interface{}{
			"range": map[string]interface{}{"timestamp": rangeQuery},
		})
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
		"sort": []map[string]interface{}{
			{"timestamp": map[string]string{"order": "desc"}},
		},
		"from": (params.Page - 1) * params.PageSize,
		"size": params.PageSize,
	}

	if len(must) == 0 {
		query["query"] = map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
	}

	docs, total, err := Search(OperationLogIndex, query)
	if err != nil {
		return nil, err
	}

	logs := make([]OperationLog, 0, len(docs))
	for _, doc := range docs {
		log := OperationLog{}
		if v, ok := doc["id"].(string); ok {
			log.ID = v
		}
		if v, ok := doc["_id"].(string); ok && log.ID == "" {
			log.ID = v
		}
		if v, ok := doc["timestamp"].(string); ok {
			log.Timestamp, _ = time.Parse(time.RFC3339, v)
		}
		if v, ok := doc["service"].(string); ok {
			log.Service = v
		}
		if v, ok := doc["user_id"].(float64); ok {
			log.UserID = uint(v)
		}
		if v, ok := doc["username"].(string); ok {
			log.Username = v
		}
		if v, ok := doc["ip"].(string); ok {
			log.IP = v
		}
		if v, ok := doc["method"].(string); ok {
			log.Method = v
		}
		if v, ok := doc["path"].(string); ok {
			log.Path = v
		}
		if v, ok := doc["query"].(string); ok {
			log.Query = v
		}
		if v, ok := doc["status_code"].(float64); ok {
			log.StatusCode = int(v)
		}
		if v, ok := doc["duration"].(float64); ok {
			log.Duration = int64(v)
		}
		if v, ok := doc["action"].(string); ok {
			log.Action = v
		}
		if v, ok := doc["module"].(string); ok {
			log.Module = v
		}
		if v, ok := doc["description"].(string); ok {
			log.Description = v
		}
		if v, ok := doc["level"].(string); ok {
			log.Level = v
		}
		if v, ok := doc["error_msg"].(string); ok {
			log.ErrorMsg = v
		}
		if v, ok := doc["user_agent"].(string); ok {
			log.UserAgent = v
		}
		logs = append(logs, log)
	}

	return &QueryResult{
		Total: total,
		Logs:  logs,
	}, nil
}

// GetStats 获取统计信息
func (s *LogService) GetStats(startTime, endTime time.Time) (map[string]interface{}, error) {
	query := map[string]interface{}{
		"size": 0,
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				"timestamp": map[string]interface{}{
					"gte": startTime.Format(time.RFC3339),
					"lte": endTime.Format(time.RFC3339),
				},
			},
		},
		"aggs": map[string]interface{}{
			"by_service": map[string]interface{}{
				"terms": map[string]interface{}{"field": "service"},
			},
			"by_method": map[string]interface{}{
				"terms": map[string]interface{}{"field": "method"},
			},
			"by_level": map[string]interface{}{
				"terms": map[string]interface{}{"field": "level"},
			},
			"by_action": map[string]interface{}{
				"terms": map[string]interface{}{"field": "action"},
			},
			"avg_duration": map[string]interface{}{
				"avg": map[string]interface{}{"field": "duration"},
			},
			"by_hour": map[string]interface{}{
				"date_histogram": map[string]interface{}{
					"field":             "timestamp",
					"calendar_interval": "hour",
				},
			},
		},
	}

	es := GetClient()
	if es == nil {
		return nil, fmt.Errorf("ES client not initialized")
	}

	docs, total, err := Search(OperationLogIndex, query)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total": total,
		"docs":  docs,
	}, nil
}

// DeleteOldLogs 删除旧日志
func (s *LogService) DeleteOldLogs(before time.Time) error {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				"timestamp": map[string]interface{}{
					"lt": before.Format(time.RFC3339),
				},
			},
		},
	}
	return DeleteByQuery(OperationLogIndex, query)
}
