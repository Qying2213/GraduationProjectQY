package coze

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"evaluator-service/internal/config"
	"evaluator-service/internal/logging"
)

var httpClient = &http.Client{Timeout: 300 * time.Second}

// Client wraps Coze API configuration
type Client struct {
	cfg *config.Config
	log *logging.Logger
}

// NewClient creates a new Coze client with configuration
func NewClient(cfg *config.Config) *Client {
	return &Client{cfg: cfg, log: logging.New()}
}

// uploadFile uploads raw bytes to Coze and returns file_id
func (c *Client) uploadFile(ctx context.Context, filename string, data []byte) (string, error) {
	url := fmt.Sprintf("%s/v1/files/upload", c.cfg.Coze.BaseURL)
	c.log.Info("Uploading file to Coze", logging.KV("url", url), logging.KV("filename", filename), logging.KV("size", len(data)))

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, err := w.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("create form file: %w", err)
	}
	if _, err := part.Write(data); err != nil {
		return "", fmt.Errorf("write file data: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.cfg.Coze.Token)
	req.Header.Set("Content-Type", w.FormDataContentType())

	start := time.Now()
	resp, err := httpClient.Do(req)
	dur := time.Since(start)
	if err != nil {
		c.log.Error("Upload request failed", logging.Err(err), logging.KV("duration", dur.String()))
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	prev := string(body)
	if len(prev) > 200 {
		prev = prev[:200]
	}
	c.log.Info("Upload response", logging.KV("status", resp.StatusCode), logging.KV("body_preview", prev))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("coze upload http %d: %s", resp.StatusCode, string(body))
	}

	var out struct {
		Code int             `json:"code"`
		Msg  string          `json:"msg"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("parse upload response: %w", err)
	}
	if out.Code != 0 {
		if out.Msg == "" {
			out.Msg = "upload failed"
		}
		return "", fmt.Errorf("coze upload error: %s", out.Msg)
	}

	var dataMap map[string]any
	if err := json.Unmarshal(out.Data, &dataMap); err != nil {
		return "", fmt.Errorf("parse upload data: %w", err)
	}
	id, _ := dataMap["id"].(string)
	if id == "" {
		return "", fmt.Errorf("coze upload response missing file id")
	}
	return id, nil
}

// RunWorkflow uploads the resume to Coze to get file_id, then runs the workflow with that file_id
// Endpoint: POST {baseURL}/v1/workflow/run (workflow_id in body)
func (c *Client) RunWorkflow(ctx context.Context, name string, jdText string, resume []byte) (map[string]any, error) {
	c.log.Info("Coze RunWorkflow started",
		logging.KV("name", name),
		logging.KV("jd_text_len", len(jdText)),
		logging.KV("resume_size", len(resume)),
		logging.KV("workflow_id", c.cfg.Coze.WorkflowID))

	// Ensure filename ends with .pdf
	filename := name
	if filename == "" {
		filename = "resume.pdf"
	}
	if !strings.HasSuffix(strings.ToLower(filename), ".pdf") {
		base := strings.TrimSuffix(filename, filepath.Ext(filename))
		filename = base + ".pdf"
	}

	// 1) Upload to get file_id
	fileID, err := c.uploadFile(ctx, filename, resume)
	if err != nil {
		c.log.Error("File upload failed", logging.Err(err))
		return nil, err
	}
	c.log.Info("File uploaded", logging.KV("file_id", fileID))

	// 2) Run workflow with file_id
	requestBody := map[string]interface{}{
		"workflow_id":   c.cfg.Coze.WorkflowID,
		"response_mode": "blocking",
		"parameters": map[string]interface{}{
			"name":        name,
			"jd_text":     jdText,
			"resume_file": map[string]any{"file_id": fileID},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		c.log.Error("Failed to marshal request body", logging.Err(err))
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	c.log.Info("Request body prepared", logging.KV("body_size", len(jsonBody)))

	url := fmt.Sprintf("%s/v1/workflow/run", c.cfg.Coze.BaseURL)
	c.log.Info("Sending request to Coze", logging.KV("url", url))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		c.log.Error("Failed to create request", logging.Err(err))
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.cfg.Coze.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	startTime := time.Now()
	resp, err := httpClient.Do(req)
	duration := time.Since(startTime)
	if err != nil {
		c.log.Error("HTTP request failed",
			logging.Err(err),
			logging.KV("duration", duration.String()),
			logging.KV("url", url))
		return nil, err
	}
	defer resp.Body.Close()

	c.log.Info("Received response from Coze",
		logging.KV("status_code", resp.StatusCode),
		logging.KV("duration", duration.String()),
		logging.KV("content_type", resp.Header.Get("Content-Type")))

	b, _ := io.ReadAll(resp.Body)
	previewLen := 200
	if len(b) < previewLen {
		previewLen = len(b)
	}
	c.log.Info("Response body read", logging.KV("body_len", len(b)), logging.KV("body_preview", string(b[:previewLen])))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		c.log.Error("Coze API returned error status",
			logging.KV("status_code", resp.StatusCode),
			logging.KV("response_body", string(b)))
		return nil, fmt.Errorf("coze http %d: %s", resp.StatusCode, string(b))
	}

	// Coze commonly returns envelope { code, msg, data }
	var envelope map[string]any
	if err := json.Unmarshal(b, &envelope); err != nil {
		c.log.Error("Failed to unmarshal response", logging.Err(err))
		return map[string]any{"raw": string(b)}, nil
	}

	// If code exists and is non-zero, return as error with msg
	if v, ok := envelope["code"]; ok {
		if f, ok2 := v.(float64); ok2 && int(f) != 0 {
			msg, _ := envelope["msg"].(string)
			return nil, fmt.Errorf("coze error code=%d: %s", int(f), msg)
		}
	}

	// Normalize nested data to expose top-level fields for downstream logic
	// Coze 工作流返回格式: {"code": 0, "data": {"output": "{...JSON...}", "reasoning_content": ""}}
	if dataVal, ok := envelope["data"]; ok {
		switch dv := dataVal.(type) {
		case string:
			// data is a JSON string. Try to parse it and extract result
			c.log.Info("Parsing data string", logging.KV("data_len", len(dv)), logging.KV("data_preview", truncateString(dv, 500)))
			var dataObj map[string]any
			if err := json.Unmarshal([]byte(dv), &dataObj); err == nil {
				// 提取 output 字段（Coze 工作流返回格式）
				if output, ok := dataObj["output"].(string); ok && output != "" {
					envelope["output"] = output
				}
				// 兼容旧格式 result 字段
				if res, ok := dataObj["result"].(string); ok && res != "" {
					envelope["result"] = res
					c.log.Info("Extracted result from data", logging.KV("result_len", len(res)))
				} else {
					c.log.Warn("No result field found in data object", logging.KV("data_keys", getMapKeys(dataObj)))
				}
			} else {
				c.log.Error("Failed to parse data JSON string",
					logging.Err(err),
					logging.KV("data_len", len(dv)),
					logging.KV("data_content", truncateString(dv, 1000)))
				return nil, fmt.Errorf("failed to parse data JSON: %w", err)
			}
		case map[string]any:
			// 提取 output 字段（Coze 工作流返回格式）
			if output, ok := dv["output"].(string); ok && output != "" {
				envelope["output"] = output
			}
			// 兼容旧格式 result 字段
			if res, ok := dv["result"].(string); ok && res != "" {
				envelope["result"] = res
			}
		default:
			// leave as-is
		}
	}

	c.log.Info("Coze RunWorkflow completed successfully", logging.KV("response_keys", getMapKeys(envelope)))
	return envelope, nil
}

func getMapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "... (truncated)"
}

// Global client instance (for backward compatibility)
var globalClient *Client

// Init initializes the global Coze client
func Init(cfg *config.Config) {
	globalClient = NewClient(cfg)
}

// RunWorkflow uses the global client (for backward compatibility)
func RunWorkflow(ctx context.Context, name string, jdText string, resume []byte) (map[string]any, error) {
	if globalClient == nil {
		return nil, fmt.Errorf("coze client not initialized, call coze.Init() first")
	}
	return globalClient.RunWorkflow(ctx, name, jdText, resume)
}
