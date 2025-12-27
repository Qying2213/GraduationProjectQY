package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"evaluator-service/internal/logging"
	"evaluator-service/internal/models"
	"evaluator-service/internal/repository"
	"evaluator-service/internal/script"
)

// PositionService 岗位服务
type PositionService struct {
	repo *repository.PositionRepository
	log  *logging.Logger
}

// NewPositionService 创建岗位服务
func NewPositionService(repo *repository.PositionRepository, log *logging.Logger) *PositionService {
	return &PositionService{repo: repo, log: log}
}

// PositionItem Python 脚本返回的岗位数据结构
type PositionItem struct {
	PostID           string `json:"post_id"`
	PostName         string `json:"post_name"`
	RecruitType      string `json:"recruit_type"`
	ServiceCondition string `json:"service_condition"`
	WorkContent      string `json:"work_content"`
}

// positionsPayload Python 脚本返回的岗位列表
type positionsPayload struct {
	Total     int            `json:"total"`
	Positions []PositionItem `json:"positions"`
}

// SyncPositions 同步岗位数据
func (s *PositionService) SyncPositions(ctx context.Context, userID uint, env map[string]string) (int, error) {
	s.log.Info("Starting position sync", logging.KV("user_id", userID))

	// 调用 Python 脚本获取岗位列表
	positions, err := s.fetchPositionsFromScript(ctx, env)
	if err != nil {
		s.log.Error("Failed to fetch positions from script", logging.Err(err))
		return 0, err
	}

	s.log.Info("Fetched positions from script", logging.KV("count", len(positions)))

	// 批量保存到数据库
	var positionModels []models.Position
	for _, p := range positions {
		if p.PostID == "" {
			continue
		}
		positionModels = append(positionModels, models.Position{
			UserID:           userID,
			PostID:           p.PostID,
			PostName:         p.PostName,
			RecruitType:      p.RecruitType,
			ServiceCondition: p.ServiceCondition,
			WorkContent:      p.WorkContent,
		})
	}

	if len(positionModels) == 0 {
		s.log.Info("No positions to sync")
		return 0, nil
	}

	// 批量 Upsert
	if err := s.repo.BatchUpsert(positionModels); err != nil {
		s.log.Error("Failed to batch upsert positions", logging.Err(err))
		return 0, err
	}

	s.log.Info("Position sync completed", logging.KV("synced_count", len(positionModels)))
	return len(positionModels), nil
}

// fetchPositionsFromScript 调用 Python 脚本获取岗位列表
func (s *PositionService) fetchPositionsFromScript(ctx context.Context, env map[string]string) ([]PositionItem, error) {
	scriptPath := "internal/script/wintalent_fetch.py"
	if !filepath.IsAbs(scriptPath) {
		abs, err := filepath.Abs(scriptPath)
		if err == nil {
			scriptPath = abs
		}
	}

	// 获取 Python 解释器路径
	pythonPath := getPythonPath()

	cmd := exec.CommandContext(ctx, pythonPath, scriptPath, "--positions-only")
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("script failed: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("run script failed: %w", err)
	}

	var payload positionsPayload
	if err := json.Unmarshal(out, &payload); err != nil {
		return nil, fmt.Errorf("parse positions json failed: %w", err)
	}

	return payload.Positions, nil
}

// getPythonPath 返回 Python 解释器路径（使用 script 包的统一函数）
func getPythonPath() string {
	return script.GetPythonPath()
}

// GetAllByUser 获取用户的所有岗位
func (s *PositionService) GetAllByUser(userID uint) ([]models.Position, error) {
	return s.repo.FindAllByUser(userID)
}

// GetByID 根据 ID 获取岗位
func (s *PositionService) GetByID(id uint) (*models.Position, error) {
	return s.repo.FindByID(id)
}

// GetByIDAndUser 根据 ID 和用户 ID 获取岗位
func (s *PositionService) GetByIDAndUser(id, userID uint) (*models.Position, error) {
	return s.repo.FindByIDAndUser(id, userID)
}
