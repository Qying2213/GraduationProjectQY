package script

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

// configuredPythonPath 存储从配置中读取的 Python 路径
var configuredPythonPath string

// SetPythonPath 设置 Python 解释器路径（由 main 函数调用）
func SetPythonPath(path string) {
	configuredPythonPath = path
}

// GetPythonPath 返回 Python 解释器路径，优先使用配置值，其次使用项目虚拟环境
func GetPythonPath() string {
	// 优先使用配置的路径
	if configuredPythonPath != "" {
		return configuredPythonPath
	}
	// 尝试使用项目根目录的 venv
	venvPython := "venv/bin/python3"
	if abs, err := filepath.Abs(venvPython); err == nil {
		if _, err := os.Stat(abs); err == nil {
			return abs
		}
	}
	// 回退到系统 python3
	return "python3"
}

// WintalentItem mirrors items produced by wintalent_fetch.py --json-out
type WintalentItem struct {
	Name         string                 `json:"name"`
	ApplyID      string                 `json:"apply_id"`  // 招聘系统申请ID
	ResumeID     string                 `json:"resume_id"` // 简历ID
	JD           map[string]any         `json:"jd"`
	ResumePDFB64 string                 `json:"resume_pdf_b64"`
	Extra        map[string]interface{} `json:"-"`
}

// wintalentPagePayload is the new paged JSON payload: { total: number, items: [...] }
// When python is invoked with --json-out --page-no, it returns this structure.
type wintalentPagePayload struct {
	Total int             `json:"total"`
	Items []WintalentItem `json:"items"`
}

// RunWintalentFetch executes the python script and parses its JSON output (single shot, no paging params).
// scriptRelPath: relative path from repo root to the script (e.g., internal/script/wintalent_fetch.py)
func RunWintalentFetch(ctx context.Context, scriptRelPath string) ([]WintalentItem, error) {
	return RunWintalentFetchWithEnv(ctx, scriptRelPath, nil)
}

// RunWintalentFetchWithEnv executes python with optional env vars (e.g., WT_USERNAME, WT_PASSWORD).
func RunWintalentFetchWithEnv(ctx context.Context, scriptRelPath string, env map[string]string) ([]WintalentItem, error) {
	// Resolve possible relative path
	sp := scriptRelPath
	if !filepath.IsAbs(sp) {
		abs, err := filepath.Abs(sp)
		if err == nil {
			sp = abs
		}
	}
	// 使用配置的 Python 路径
	pythonPath := GetPythonPath()
	cmd := exec.CommandContext(ctx, pythonPath, sp, "--json-out")
	// 继承当前进程的环境变量
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("run wintalent_fetch.py failed: %w", err)
	}
	// Prefer new payload shape { total, items }
	var payload wintalentPagePayload
	if err := json.Unmarshal(out, &payload); err == nil && payload.Items != nil {
		return payload.Items, nil
	}
	// Fallback to legacy pure array
	var items []WintalentItem
	if err := json.Unmarshal(out, &items); err != nil {
		return nil, fmt.Errorf("parse wintalent json failed: %w", err)
	}
	return items, nil
}

// RunWintalentFetchPage executes the python script for a specific page and page size.
// NOTE: The underlying python currently may fetch page_no and page_no+1 together.
// To mitigate duplication at caller side, consider slicing to pageSize items.
func RunWintalentFetchPage(ctx context.Context, scriptRelPath string, pageNo int, pageSize int) ([]WintalentItem, error) {
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	sp := scriptRelPath
	if !filepath.IsAbs(sp) {
		abs, err := filepath.Abs(sp)
		if err == nil {
			sp = abs
		}
	}
	cmd := exec.CommandContext(
		ctx,
		GetPythonPath(),
		sp,
		"--json-out",
		"--page-no", strconv.Itoa(pageNo),
		"--page-size", strconv.Itoa(pageSize),
	)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("run wintalent_fetch.py page failed (page=%d,size=%d): %w", pageNo, pageSize, err)
	}
	var items []WintalentItem
	if err := json.Unmarshal(out, &items); err != nil {
		return nil, fmt.Errorf("parse wintalent json failed (page=%d,size=%d): %w", pageNo, pageSize, err)
	}
	return items, nil
}

// WithTimeout helper returns a child context with timeout suitable for script runtime.
func WithTimeout(parent context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, d)
}

// RunGraduateFetch executes the graduate_fetch.py script to fetch resumes from graduation project backend.
// scriptRelPath: relative path from repo root to the script (e.g., internal/script/graduate_fetch.py)
func RunGraduateFetch(ctx context.Context, scriptRelPath string) ([]WintalentItem, error) {
	return RunGraduateFetchWithEnv(ctx, scriptRelPath, nil)
}

// RunGraduateFetchWithEnv executes graduate_fetch.py with optional env vars (e.g., GRADUATE_API_URL).
func RunGraduateFetchWithEnv(ctx context.Context, scriptRelPath string, env map[string]string) ([]WintalentItem, error) {
	sp := scriptRelPath
	if !filepath.IsAbs(sp) {
		abs, err := filepath.Abs(sp)
		if err == nil {
			sp = abs
		}
	}
	pythonPath := GetPythonPath()
	cmd := exec.CommandContext(ctx, pythonPath, sp, "--json-out")
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("run graduate_fetch.py failed: %w", err)
	}
	var payload wintalentPagePayload
	if err := json.Unmarshal(out, &payload); err == nil && payload.Items != nil {
		return payload.Items, nil
	}
	var items []WintalentItem
	if err := json.Unmarshal(out, &items); err != nil {
		return nil, fmt.Errorf("parse graduate json failed: %w", err)
	}
	return items, nil
}

// RunGraduateFetchPage executes graduate_fetch.py for a specific page.
func RunGraduateFetchPage(ctx context.Context, scriptRelPath string, pageNo int, pageSize int) ([]WintalentItem, int, error) {
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	sp := scriptRelPath
	if !filepath.IsAbs(sp) {
		abs, err := filepath.Abs(sp)
		if err == nil {
			sp = abs
		}
	}
	cmd := exec.CommandContext(
		ctx,
		GetPythonPath(),
		sp,
		"--json-out",
		"--page-no", strconv.Itoa(pageNo),
		"--page-size", strconv.Itoa(pageSize),
	)
	out, err := cmd.Output()
	if err != nil {
		return nil, 0, fmt.Errorf("run graduate_fetch.py page failed (page=%d,size=%d): %w", pageNo, pageSize, err)
	}
	var payload wintalentPagePayload
	if err := json.Unmarshal(out, &payload); err != nil {
		return nil, 0, fmt.Errorf("parse graduate json failed (page=%d,size=%d): %w", pageNo, pageSize, err)
	}
	return payload.Items, payload.Total, nil
}
