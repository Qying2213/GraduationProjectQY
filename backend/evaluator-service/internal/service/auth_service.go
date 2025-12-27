package service

import (
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"path/filepath"
	"time"

	"evaluator-service/internal/config"
	"evaluator-service/internal/logging"
	"evaluator-service/internal/models"
	"evaluator-service/internal/repository"
	"evaluator-service/internal/script"
	"evaluator-service/internal/utils"

	"gorm.io/gorm"
)

// LoginResult 登录结果
type LoginResult struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      UserInfo  `json:"user"`
}

// UserInfo 用户基本信息（不含敏感数据）
type UserInfo struct {
	ID       uint   `json:"id"`
	CorpCode string `json:"corp_code"`
	Username string `json:"username"`
}

// AuthService 认证服务
type AuthService struct {
	cfg      *config.Config
	log      *logging.Logger
	userRepo *repository.UserRepository
	credRepo *repository.CredentialRepository
	jwtSvc   *JWTService
}

// NewAuthService 创建认证服务实例
func NewAuthService(cfg *config.Config, log *logging.Logger, userRepo *repository.UserRepository, db *gorm.DB) *AuthService {
	return &AuthService{
		cfg:      cfg,
		log:      log,
		userRepo: userRepo,
		credRepo: repository.NewCredentialRepository(db),
		jwtSvc:   NewJWTService(cfg),
	}
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, corpCode, username, password string) (*LoginResult, error) {
	// 毕业设计模式：跳过第三方验证，直接使用本地账号
	if corpCode == "graduate" {
		return s.loginGraduate(ctx, username, password)
	}

	// 1. 调用 Python 脚本验证凭据
	if err := s.verifyCredentials(ctx, corpCode, username, password); err != nil {
		return nil, err
	}

	// 2. 加密密码
	encKey := s.getEncryptionKey()
	passwordCipher, err := utils.EncryptAESGCM(encKey, password)
	if err != nil {
		s.log.Error("encrypt password failed", logging.KV("error", err.Error()))
		return nil, errors.New("凭据加密失败")
	}

	// 3. 查找或创建用户
	user, isNew, err := s.userRepo.FindOrCreate(corpCode, username, passwordCipher)
	if err != nil {
		s.log.Error("find or create user failed", logging.KV("error", err.Error()))
		return nil, errors.New("用户创建失败")
	}

	if isNew {
		s.log.Info("new user created", logging.KV("user_id", user.ID), logging.KV("username", username))
	}

	// 4. 保存凭据到 credentials 表（供批量评估使用）
	_, err = s.credRepo.UpsertByUser("motern", username, passwordCipher, user.ID)
	if err != nil {
		s.log.Error("save credentials failed", logging.KV("error", err.Error()))
		// 不阻断登录流程，只记录错误
	}

	// 5. 生成 JWT token
	token, expiresAt, err := s.jwtSvc.GenerateToken(user)
	if err != nil {
		s.log.Error("generate token failed", logging.KV("error", err.Error()))
		return nil, errors.New("token 生成失败")
	}

	return &LoginResult{
		Token:     token,
		ExpiresAt: expiresAt,
		User: UserInfo{
			ID:       user.ID,
			CorpCode: user.CorpCode,
			Username: user.Username,
		},
	}, nil
}

// loginGraduate 毕业设计模式登录（简单验证）
func (s *AuthService) loginGraduate(ctx context.Context, username, password string) (*LoginResult, error) {
	// 简单验证：admin/admin123 或任意用户名密码
	// 生产环境应该对接真实的用户系统
	validUsers := map[string]string{
		"admin":   "admin123",
		"hr_li":   "123456",
		"hr_wang": "123456",
	}

	// 检查是否是预设用户
	if expectedPwd, ok := validUsers[username]; ok {
		if password != expectedPwd {
			return nil, errors.New("用户名或密码错误")
		}
	}
	// 如果不是预设用户，也允许登录（演示用）

	// 加密密码
	encKey := s.getEncryptionKey()
	passwordCipher, err := utils.EncryptAESGCM(encKey, password)
	if err != nil {
		return nil, errors.New("凭据加密失败")
	}

	// 查找或创建用户
	user, _, err := s.userRepo.FindOrCreate("graduate", username, passwordCipher)
	if err != nil {
		return nil, errors.New("用户创建失败")
	}

	// 生成 JWT token
	token, expiresAt, err := s.jwtSvc.GenerateToken(user)
	if err != nil {
		return nil, errors.New("token 生成失败")
	}

	return &LoginResult{
		Token:     token,
		ExpiresAt: expiresAt,
		User: UserInfo{
			ID:       user.ID,
			CorpCode: "graduate",
			Username: user.Username,
		},
	}, nil
}

// ValidateToken 验证 token
func (s *AuthService) ValidateToken(token string) (*UserClaims, error) {
	return s.jwtSvc.ValidateToken(token)
}

// GetUserByID 根据 ID 获取用户
func (s *AuthService) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	return s.userRepo.GetByID(userID)
}

// GetUserCredentials 获取用户的解密凭据（用于调用 Python 脚本）
func (s *AuthService) GetUserCredentials(ctx context.Context, userID uint) (corpCode, username, password string, err error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return "", "", "", err
	}

	encKey := s.getEncryptionKey()
	password, err = utils.DecryptAESGCM(encKey, user.PasswordCipher)
	if err != nil {
		return "", "", "", errors.New("凭据解密失败")
	}

	return user.CorpCode, user.Username, password, nil
}

// verifyCredentials 调用 Python 脚本验证凭据
func (s *AuthService) verifyCredentials(ctx context.Context, corpCode, username, password string) error {
	scriptPath := filepath.Join("internal", "script", "wintalent_fetch.py")

	// 获取 Python 解释器路径
	pythonPath := s.getPythonPath()

	cmd := exec.CommandContext(ctx, pythonPath, scriptPath,
		"--login-only",
		"--corp-code", corpCode,
		"--username", username,
		"--password", password,
	)

	output, err := cmd.Output()
	if err != nil {
		// 尝试解析错误输出
		if exitErr, ok := err.(*exec.ExitError); ok {
			var result struct {
				Status  string `json:"status"`
				Message string `json:"message"`
			}
			if json.Unmarshal(output, &result) == nil && result.Message != "" {
				return errors.New(result.Message)
			}
			if json.Unmarshal(exitErr.Stderr, &result) == nil && result.Message != "" {
				return errors.New(result.Message)
			}
		}
		return errors.New("第三方认证服务暂时不可用")
	}

	// 解析成功响应
	var result struct {
		Status string `json:"status"`
		User   struct {
			CorpCode string `json:"corp_code"`
			Username string `json:"username"`
		} `json:"user"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return errors.New("认证响应解析失败")
	}

	if result.Status != "success" {
		if result.Message != "" {
			return errors.New(result.Message)
		}
		return errors.New("用户名或密码错误")
	}

	return nil
}

// getPythonPath 获取 Python 解释器路径（使用 script 包的统一函数）
func (s *AuthService) getPythonPath() string {
	return script.GetPythonPath()
}

// getEncryptionKey 获取加密密钥
func (s *AuthService) getEncryptionKey() []byte {
	key := s.cfg.Credentials.EncKey
	if key == "" {
		key = "resume-evaluator-default-enc-key"
	}
	// 确保密钥长度为 32 字节（AES-256）
	keyBytes := []byte(key)
	if len(keyBytes) < 32 {
		padded := make([]byte, 32)
		copy(padded, keyBytes)
		return padded
	}
	return keyBytes[:32]
}
